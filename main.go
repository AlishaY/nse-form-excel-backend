package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
	"github.com/gin-gonic/gin"
)

type Job struct {
	// ID       	int  `json:"id"`
	// Item       	string  `json:"item"`
	JobDate       		string  `json:"jobDate"`
	ReferenceNo 		string 	`json:"referenceNo"`
	CustomerName		string	`json:"customerName"`
	DeliveryLocName		string	`json:"deliveryLocName"`
	TruckNo				string	`json:"truckNo"`
}

var db *sql.DB

func main() {
	// Database connection configuration
	// server := "WIN-OMPBT533FJ7"
	// port := 1433
	// user := "coadmin"
	// password := "tms@1234"
	// database := "CoTMS"

	var err error
	// db, err = sql.Open(
	// 	"sqlserver",
	// 	"sqlserver://coadmin:tms@1234@103.230.124.241:1433?database=CoTMS&connection+timeout=30",
	// )
	db, err = sql.Open(
		"sqlserver",
		"sqlserver://coadmin:alisha@1234@localhost?database=TODO&connection+timeout=30",
	)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Create a new Gin router
	router := gin.Default()

	// Define your API endpoints
	router.GET("/joborder", getJob)

	// Start the server
	router.Run(":3000")
}

func getJob(c *gin.Context) {
	fmt.Println("first line in getJob")
	rows, err := db.Query("SELECT JobDate, ReferenceNo, CustomerName, DeliveryLocName, TruckNo FROM JobOrder WHERE JobNo='TDN0098789'")
	if err != nil {
		fmt.Println("this is error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Jobs"})
		return
	}
	defer rows.Close()
	fmt.Println("after defer rows close")

	jobs := []Job{}
	for rows.Next() {
		fmt.Println("scan")
		var job Job
		err := rows.Scan(&job.JobDate, &job.ReferenceNo, &job.CustomerName, &job.DeliveryLocName, &job.TruckNo)
		fmt.Println("scan 2")
		if err != nil {
			fmt.Println("scan 3")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		fmt.Println("scan 4")
		jobs = append(jobs, job)
	}
	fmt.Println("end")
	c.JSON(http.StatusOK, jobs)
}
