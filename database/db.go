package database

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
)

var Data *sqlx.DB

func ConnectAndMigrate(host, port, databaseName, user, password string, sslMode string) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, databaseName, sslMode)
	DB, err := sqlx.Open("postgres", connStr)

	if err != nil {
		log.Printf("ConnectAndMigrate : Error in database connection.")
		return err
	}

	Data = DB
	return migrateUp(DB)
}

func migrateUp(data *sqlx.DB) error {
	//log.Println(data.Driver())
	driver, err := postgres.WithInstance(data.DB, &postgres.Config{})
	if err != nil {
		log.Printf("migrateUp : Error in retrieving database driver.")
		return err
	}
	m, instanceErr := migrate.NewWithDatabaseInstance(
		"file://database/migration",
		"postgres", driver)

	if instanceErr != nil {
		log.Printf("migrateUp : Error in creating migrate. %s", instanceErr)
		return instanceErr
	}

	if migrateErr := m.Up(); migrateErr != nil && migrateErr != migrate.ErrNoChange {
		log.Printf("migrateUp : Error in migrating up. %s", migrateErr)
		return migrateErr
	}
	return nil
}

func Tx(fn func(tx *sqlx.Tx) error) error {
	tx, err := Data.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				logrus.Errorf("failed to rollback tx: %s", rollBackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			logrus.Errorf("failed to commit tx: %s", commitErr)
		}
	}()
	err = fn(tx)
	return err
}
