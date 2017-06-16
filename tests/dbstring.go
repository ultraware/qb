package tests

import (
	"fmt"
	"os"
)

func getHost() string {
	if os.Getenv(`TESTHOST`) != `` {
		return os.Getenv(`TESTHOST`)
	}
	return `127.0.0.1`
}

func getPostgresDBString() string {
	return fmt.Sprintf("host=%s dbname=test user=test sslmode=disable", getHost())
}

func getMysqlDBString() string {
	return fmt.Sprintf("test:test@tcp(%s:3306)/test?multiStatements=true&parseTime=true", getHost())
}

const createSQL = `
DROP TABLE IF EXISTS one, two;

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
