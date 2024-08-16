package mailer

import (
	"bytes"
	"fmt"
	"github.com/ainsleyclark/go-mail/drivers"
	"github.com/ainsleyclark/go-mail/mail"
	"github.com/vanng822/go-premailer/premailer"
	smtpmail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Mail holds the information necessary to connect to an SMTP server
type Mail struct {
	Domain      string
	Templates   string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
	Jobs        chan Message
	Results     chan Result
	API         string
	APIKey      string
	APIUrl      string
}

// Message is the type for an email message
type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Template    string
	Attachments []string
	Data        interface{}
}

// Result contains information regarding the status of the email message
type Result struct {
	Success bool
	Error   error
}

// ListenForMail listens to the mail channel and sends mail when it receives a payload. It runs
// continually in the background, and sends error/success messages back on the Results channel.
// Note that if api and api key are set, it will prefer using an api to send mail
func (m *Mail) ListenForMail() {
	for {
		msg := <-m.Jobs
		err := m.Send(msg)
		if err != nil {
			m.Results <- Result{false, err}
		} else {
			m.Results <- Result{true, nil}
		}
	}
}

func (m *Mail) Send(msg Message) error {
	if len(m.API) > 0 && len(m.APIKey) > 0 && len(m.APIUrl) > 0 && strings.ToLower(m.API) != "smtp" {
		return m.ChooseAPI(msg)
	}
	return m.SendSMTPMessage(msg)
}

func (m *Mail) ChooseAPI(msg Message) error {
	switch m.API {
	case "mailgun", "sparkpost", "sendgrid", "postal", "postmark":
		return m.SendUsingAPI(msg, m.API)
	default:
		return fmt.Errorf("unknown api %s; only mailgun, sparkpost, sendgrid,postal or postmark accepted", m.API)
	}
}

func (m *Mail) SendUsingAPI(msg Message, transport string) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	cfg := mail.Config{
		URL:         m.APIUrl,
		APIKey:      m.APIKey,
		Domain:      m.Domain,
		FromAddress: msg.From,
		FromName:    msg.FromName,
	}
	driver, err := drivers.NewSMTP(cfg)

	switch transport {
	case "mailgun":
		driver, err = drivers.NewMailgun(cfg)
		if err != nil {
			return err
		}
	case "sparkpost":
		driver, err = drivers.NewSparkPost(cfg)
		if err != nil {
			return err
		}
	case "sendgrid":
		driver, err = drivers.NewSendGrid(cfg)
		if err != nil {
			return err
		}
	case "postal":
		driver, err = drivers.NewPostal(cfg)
		if err != nil {
			return err
		}
	case "postmark":
		driver, err = drivers.NewPostmark(cfg)
		if err != nil {
			return err
		}
	default:
		// Should not be reachable
		return err
		//log.Fatalf("unknown transport %s", transport)
	}

	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	tx := &mail.Transmission{
		Recipients: []string{msg.To},
		Subject:    msg.Subject,
		HTML:       formattedMessage,
		PlainText:  plainMessage,
	}

	// add attachments
	err = m.addAPIAttachments(msg, tx)
	if err != nil {
		return err
	}

	_, err = driver.Send(tx)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) addAPIAttachments(msg Message, tx *mail.Transmission) error {
	if len(msg.Attachments) > 0 {
		var attachments []mail.Attachment
		for _, attachment := range msg.Attachments {
			var attach mail.Attachment
			content, err := os.ReadFile(attachment)
			if err != nil {
				return err
			}

			fileName := filepath.Base(attachment)
			attach.Bytes = content
			attach.Filename = fileName
			attachments = append(attachments, attach)
		}
		tx.Attachments = attachments
	}
	return nil
}

// SendSMTPMessage builds and send an email message using SMTP. This is
// called by ListenForMail, and can also be called directly when necessary
func (m *Mail) SendSMTPMessage(msg Message) error {
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	server := smtpmail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := smtpmail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)
	email.SetBody(smtpmail.TextHTML, formattedMessage)
	email.AddAlternative(smtpmail.TextPlain, plainMessage)

	if len(msg.Attachments) > 0 {
		for _, attachment := range msg.Attachments {
			email.AddAttachment(attachment)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) getEncryption(e string) smtpmail.Encryption {
	switch e {
	case "tls":
		return smtpmail.EncryptionSTARTTLS
	case "ssl":
		return smtpmail.EncryptionSSL
	case "none":
		return smtpmail.EncryptionNone
	default:
		return smtpmail.EncryptionSTARTTLS
	}
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("%s/%s.html.gohtml", m.Templates, msg.Template)

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.Data); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}
func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("%s/%s.plain.gohtml", m.Templates, msg.Template)

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.Data); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}
	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}
