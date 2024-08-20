package main

import "github.com/gobuffalo/pop"

func doMigrate(arg2, arg3 string) error {
	//dsn := getDSN()
	checkForDB()
	tx, err := cel.PopConnect()
	if err != nil {
		exitGracefully(err)
	}
	defer func(tx *pop.Connection) {
		_ = tx.Close()
	}(tx)

	// run the migration command
	switch arg2 {
	case "up":
		err := cel.RunPopMigrations(tx)
		if err != nil {
			return err
		}

	case "down":
		if arg3 == "all" {
			err := cel.PopMigrateDown(tx, -1)
			if err != nil {
				return err
			}
		} else {
			err := cel.PopMigrateDown(tx, 1)
			if err != nil {
				return err
			}
		}
	case "reset":
		err := cel.PopMigrateReset(tx)
		if err != nil {
			return err
		}
	default:
		showHelp()
	}

	return nil
}
