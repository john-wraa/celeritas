package mailer

import (
	"errors"
	"testing"
)

var testMsg = Message{
	From:        "me@here.com",
	FromName:    "Joe",
	To:          "you@there.com",
	Subject:     "Test",
	Template:    "test",
	Attachments: []string{"./testdata/mail/test.html.gohtml"},
}

func TestMail_SendSMTPMessage(t *testing.T) {
	err := mailer.SendSMTPMessage(testMsg)
	if err != nil {
		t.Error(err)
	}
}

func TestMail_SendUsingChan(t *testing.T) {
	mailer.Jobs <- testMsg
	res := <-mailer.Results
	if res.Error != nil {
		t.Error(errors.New("failed to send using chan"))
	}
	originalTestMsg := testMsg
	testMsg.To = "not_an_email_address"
	mailer.Jobs <- testMsg
	res = <-mailer.Results
	if res.Error == nil {
		t.Error(errors.New("no error received with invalid to address"))
	}
	testMsg = originalTestMsg
}

func TestMail_SendUsingAPI(t *testing.T) {
	msg := Message{
		To:          "you@there.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.gohtml"},
	}

	mailer.API = "unknown"
	mailer.APIKey = "abs123"
	mailer.APIUrl = "https://www.fake.fake"

	err := mailer.SendUsingAPI(msg, "unknown")
	if err == nil {
		t.Error(errors.New("no error received with invalid API"))
	}
	mailer.API = ""
	mailer.APIKey = ""
	mailer.APIUrl = ""
}

func TestMail_buildHTMLMessage(t *testing.T) {
	_, err := mailer.buildHTMLMessage(testMsg)
	if err != nil {
		t.Error(err)
	}
}

func TestMail_buildPlainTextMessage(t *testing.T) {
	_, err := mailer.buildPlainTextMessage(testMsg)
	if err != nil {
		t.Error(err)
	}
}

func TestMail_send(t *testing.T) {
	err := mailer.Send(testMsg)
	if err != nil {
		t.Error(err)
	}

	mailer.API = "unknown"
	mailer.APIKey = "abs123"
	mailer.APIUrl = "https://www.fake.fake"

	err = mailer.Send(testMsg)
	if err == nil {
		t.Error(err)
	}

	mailer.API = ""
	mailer.APIKey = ""
	mailer.APIUrl = ""
}

func TestMail_ChooseAPI(t *testing.T) {
	mailer.API = "unknown"
	err := mailer.ChooseAPI(testMsg)
	if err == nil {
		t.Error(err)
	}

}
