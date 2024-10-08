package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

func doMake(arg2, arg3, arg4 string) error {

	switch arg2 {
	case "key":
		rnd := cel.RandomString(32)
		color.Yellow("32 character encryption key: %s", rnd)

	case "migration":
		checkForDB()
		//dbType := cel.DB.DataType
		if arg3 == "" {
			exitGracefully(errors.New("you must give the migration a name"))
		}

		// default to migration type of fizz
		migrationType := "fizz"
		var up, down string

		// are we doing fizz or sql
		if arg4 == "fizz" || arg4 == "" {
			upBytes, _ := templateFS.ReadFile("templates/migrations/migration_up.fizz")
			downBytes, _ := templateFS.ReadFile("templates/migrations/migration_down.fizz")
			up = string(upBytes)
			down = string(downBytes)
		} else {
			migrationType = "sql"
		}

		// create the migrations for either fizz or sql
		err := cel.CreatePopMigration([]byte(up), []byte(down), arg3, migrationType)
		if err != nil {
			exitGracefully(err)
		}

	case "migration_old":
		dbType := cel.DB.DataType
		if arg3 == "" {
			exitGracefully(errors.New("you must give the migration a name"))
		}

		fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), arg3)

		upFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
		downFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

		err := copyFilefromTemplate("templates/migrations/migration."+dbType+".up.sql", upFile)
		if err != nil {
			exitGracefully(err)
		}

		err = copyFilefromTemplate("templates/migrations/migration."+dbType+".down.sql", downFile)
		if err != nil {
			exitGracefully(err)
		}

	case "auth":
		err := doAuth()
		if err != nil {
			exitGracefully(err)
		}

	case "handler":
		if arg3 == "" {
			exitGracefully(errors.New("you must give the handler a name"))
		}

		fileName := cel.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exists!"))
		}

		data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		handler := string(data)
		//goland:noinspection SpellCheckingInspection
		handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(arg3))

		err = os.WriteFile(fileName, []byte(handler), 0644)
		if err != nil {
			exitGracefully(err)
		}

	case "model":
		if arg3 == "" {
			exitGracefully(errors.New("you must give the model a name"))
		}

		data, err := templateFS.ReadFile("templates/data/model.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		model := string(data)

		plural := pluralize.NewClient()

		var modelName = arg3
		var tableName = arg3

		if plural.IsPlural(arg3) {
			modelName = plural.Singular(arg3)
			tableName = strings.ToLower(tableName)
		} else {
			tableName = strings.ToLower(plural.Plural(arg3))
		}

		fileName := cel.RootPath + "/data/" + strings.ToLower(modelName) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exists!"))
		}

		//goland:noinspection SpellCheckingInspection
		model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
		//goland:noinspection SpellCheckingInspection
		model = strings.ReplaceAll(model, "$TABLENAME$", tableName)

		err = copyDataToFile([]byte(model), fileName)
		if err != nil {
			exitGracefully(err)
		}

	case "mail":
		if arg3 == "" {
			exitGracefully(errors.New("you must give the mail template a name"))
		}
		htmlMail := cel.RootPath + "/mail/" + strings.ToLower(arg3) + ".html.tmpl"
		plainMail := cel.RootPath + "/mail/" + strings.ToLower(arg3) + ".plain.tmpl"

		err := copyFilefromTemplate("templates/mailer/mail.html.tmpl", htmlMail)
		if err != nil {
			exitGracefully(err)
		}

		err = copyFilefromTemplate("templates/mailer/mail.plain.tmpl", plainMail)
		if err != nil {
			exitGracefully(err)
		}

	case "session":
		err := doSessionTable()
		if err != nil {
			exitGracefully(err)
		}
	}

	return nil
}
