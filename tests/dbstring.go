package tests

import (
	"fmt"
	"os"
)

func getDBString() string {
	host := os.Getenv(`TESTHOST`)
	if host == `` {
		host = `127.0.0.1`
	}
	return fmt.Sprintf("host=%s dbname=test user=test sslmode=disable", host)
}

const createSQL = `
DROP TABLE IF EXISTS one;
DROP TABLE IF EXISTS two;

CREATE TABLE one (
	ID int PRIMARY KEY,
	Name varchar(50) NOT NULL,
	CreatedAt timestamp NOT NULL DEFAULT now()
);

CREATE TABLE two (
	OneID int,
	Number int,
	Comment varchar(100) NOT NULL,
	ModifiedAt timestamp,
	PRIMARY KEY (OneID, Number)
);`
