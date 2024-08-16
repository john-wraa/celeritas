package main

import (
	"fmt"
	"strings"
	"time"
)

func doSessionTable() error {
	dbType := cel.DB.DatabaseType
	//goland:noinspection SpellCheckingInspection
	switch strings.ToLower(dbType) {
	case "maria":
	case "mariadb":
	case "mysqldb":
		dbType = "mysql"
	case "pg":
	case "postgressql":
		dbType = "postgres"
	}
	fileName := fmt.Sprintf("%d_create_sessions.table", time.Now().UTC().UnixMicro())
	upFile := fmt.Sprintf("%s/migrations/%s.%s.up.sql", cel.RootPath, fileName, dbType)
	downFile := fmt.Sprintf("%s/migrations/%s.%s.down.sql", cel.RootPath, fileName, dbType)

	err := copyFileFromTemplate(fmt.Sprintf("templates/migrations/%s_session.sql", dbType), upFile)
	if err != nil {
		exitGracefully(err)
	}
	err = copyDataToFile([]byte("drop table sessions"), downFile)
	if err != nil {
		exitGracefully(err)
	}

	err = doMigrate("up", "")

	return nil
}
