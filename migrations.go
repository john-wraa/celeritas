package celeritas

import (
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (c *Celeritas) MigrateUp(dsn string) error {
	rootPath := filepath.ToSlash(c.RootPath)
	m, err := migrate.New("file://"+rootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer func(m *migrate.Migrate) {
		_, _ = m.Close()
	}(m)

	if err := m.Up(); err != nil {
		return err
	}
	return nil
}

func (c *Celeritas) MigrateDownAll(dsn string) error {
	rootPath := filepath.ToSlash(c.RootPath)
	m, err := migrate.New("file://"+rootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer func(m *migrate.Migrate) {
		_, _ = m.Close()
	}(m)

	if err := m.Down(); err != nil {
		return err
	}
	return nil
}

func (c *Celeritas) Steps(n int, dsn string) error {
	rootPath := filepath.ToSlash(c.RootPath)
	m, err := migrate.New("file://"+rootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer func(m *migrate.Migrate) {
		_, _ = m.Close()
	}(m)

	if err := m.Steps(n); err != nil {
		return err
	}
	return nil
}

func (c *Celeritas) MigrateForce(dsn string) error {
	rootPath := filepath.ToSlash(c.RootPath)
	m, err := migrate.New("file://"+rootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer func(m *migrate.Migrate) {
		_, _ = m.Close()
	}(m)

	if err := m.Force(-1); err != nil {
		return err
	}
	return nil
}
