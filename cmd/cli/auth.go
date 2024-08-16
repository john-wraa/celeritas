package main

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"time"
)

func doAuth() error {
	//create migrations
	dbType := cel.DB.DatabaseType
	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UTC().UnixMicro())
	upFile := cel.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := cel.RootPath + "/migrations/" + fileName + ".down.sql"
	log.Println(dbType, upFile, downFile)
	err := copyFileFromTemplate(fmt.Sprintf("templates/migrations/auth_tables.%s.sql", dbType), upFile)
	if err != nil {
		exitGracefully(err)
	}
	err = copyDataToFile([]byte("drop table if exists users cascade; drop table if exists tokens cascade; drop table if exists remember_tokens;"), downFile)
	if err != nil {
		exitGracefully(err)
	}

	// run migrations
	err = doMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	// TODO: Somewhere down the road, get rid of the `/myapp` hardcoding

	// copy models
	err = copyFileFromTemplate("templates/data/user.go.txt", cel.RootPath+"/myapp/data/user.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/data/token.go.txt", cel.RootPath+"/myapp/data/token.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/data/remember_token.go.txt", cel.RootPath+"/myapp/data/remember_token.go")
	if err != nil {
		exitGracefully(err)
	}

	// copy middleware
	err = copyFileFromTemplate("templates/middleware/auth.go.txt", cel.RootPath+"/myapp/middleware/auth.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/middleware/auth-token.go.txt", cel.RootPath+"/myapp/middleware/auth-token.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/middleware/remember.go.txt", cel.RootPath+"/myapp/middleware/remember.go")
	if err != nil {
		exitGracefully(err)
	}

	// copy handlers
	err = copyFileFromTemplate("templates/handlers/auth-handlers.go.txt", cel.RootPath+"/myapp/handlers/auth-handlers.go")
	if err != nil {
		exitGracefully(err)
	}

	// copy emails
	err = copyFileFromTemplate("templates/mailer/password-reset.html.gohtml", cel.RootPath+"/myapp/mail/password-reset.html.gohtml")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/mailer/password-reset.plain.gohtml", cel.RootPath+"/myapp/mail/password-reset.plain.gohtml")
	if err != nil {
		exitGracefully(err)
	}

	// copy views
	err = copyFileFromTemplate("templates/views/login.jet", cel.RootPath+"/myapp/views/login.jet")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/views/forgot.jet", cel.RootPath+"/myapp/views/forgot.jet")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/views/reset-password.jet", cel.RootPath+"/myapp/views/reset-password.jet")
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow("  - users, tokens and remember_tokens migrations created and executed successfully")
	color.Yellow("  - users and token models created")
	color.Yellow("  - auth middleware created")
	color.Yellow("")
	color.Yellow("Don't forget to add user and token models in data/models.com, and to add appropriate middleware to your routes!")

	return nil
}
