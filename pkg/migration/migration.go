package migration

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

type Migration struct {
	mgrt  *migrate.Migrate
	sqlDB *sql.DB
}

func NewMigration(sqlDB *sql.DB, dbName, sourceFile string) (*Migration, error) {
	dbDriver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}
	migrate, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", sourceFile), dbName, dbDriver)
	if err != nil {
		return nil, fmt.Errorf("migration failed: %v", err)
	}
	return &Migration{
		mgrt:  migrate,
		sqlDB: sqlDB,
	}, nil
}

func (m *Migration) Up() error {
	logrus.Println("Start schema migration . . .")
	_, isDirty, err := m.mgrt.Version()
	if err != nil && err != migrate.ErrNilVersion {
		logrus.Println("Database schema migration failed")
		return err
	}
	if isDirty {
		logrus.Println("Database schema migration failed")
		logrus.Println("In order to protect your data integrity please fix you database schema using this tools.")
		logrus.Println("https://github.com/golang-migrate/migrate/releases")
	}
	if err = m.mgrt.Up(); err == migrate.ErrNoChange {
		logrus.Println("No changes, database schema is up to date")
		return nil
	}
	if err != nil {
		logrus.Println("Database schema migration failed")
		return err
	}
	return nil
}
