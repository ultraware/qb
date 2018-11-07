package tests

import (
	"fmt"
	"os"
)

func getHost(prefix string) string {
	if e := os.Getenv(prefix + `_TESTHOST`); e != `` {
		return e
	}
	if e := os.Getenv(`TESTHOST`); e != `` {
		return e
	}
	return `127.0.0.1`
}

func getPostgresDBString() string {
	return fmt.Sprintf("host=%s dbname=qb_test user=qb_test sslmode=disable", getHost(`POSTGRES`))
}

func getMysqlDBString() string {
	return fmt.Sprintf("qb_test:qb_test@tcp(%s:3306)/qb_test?multiStatements=true&parseTime=true", getHost(`MYSQL`))
}

func getMssqlDBString() string {
	return fmt.Sprintf("sqlserver://qb_test:qb_test@%s:1433?database=qb_test", getHost(`MSSQL`))
}

const createSQL = `
CREATE TABLE one (
	ID int PRIMARY KEY,
	Name varchar(50) NOT NULL,
	CreatedAt timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE two (
	OneID int,
	Number int,
	Comment varchar(100) NOT NULL,
	ModifiedAt timestamp,
	PRIMARY KEY (OneID, Number)
);`
