package database

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	logger "github.com/sirupsen/logrus"
	"os"
)

func ConnectPSQLFromEnv() (*sqlx.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	db := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	return ConnectPSQL(host, user, pass, db, port)
}

func ConnectPSQL(host, user, pass, db, port string) (*sqlx.DB, error) {
	// Now you can use the updated variables for your database connection
	connectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, db)
	logger.Debug("Connect to %s\n", connectString)

	dbh, err := sqlx.Open("postgres", connectString)
	if err != nil {
		return nil, err
	} else if dbh == nil {
		return nil, errors.New("postgres Open returned nil")
	}

	err = dbh.Ping()
	if err != nil {
		return nil, err
	}

	return dbh, nil
}
