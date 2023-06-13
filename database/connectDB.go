package database

import (
	// "database/sql"
	"fmt"
	"log"

	// _ "github.com/denisenkom/go-mssqldb"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// var db *gorm.DB
// var server = "ALISHAYAACOB"
// var port = 1433
// var user = "coadmin"
// var password = "alisha@1234"
// var database = "CoKPI"

var db *gorm.DB
var server = "ALISHAYAACOB"
var port = 1433
var user = "coadmin"
var password = "alisha@1234"
var database = "CoKPI"

// ConnectDB opens a connection to the database and returns the connection object
func ConnectDB() (*gorm.DB, error) {
	// Build connection string
	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", user, password, server, port, database)

	var err error
	// Create connection pool
	db, err = gorm.Open(sqlserver.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err.Error())
		return nil, err
	}

	fmt.Println("Connected to the database")
	return db, nil
}