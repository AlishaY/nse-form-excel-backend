package main

import (
	"io/ioutil"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"net/http"
	"strings"
	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)
//json uri redirect: http://localhost:8080/callback
type Job struct {
	JobNo         	string `json:"jobNo"`
	JobDate         string `json:"jobDate"`
	ReferenceNo     string `json:"referenceNo"`
	CustomerName    string `json:"customerName"`
	DeliveryLocName string `json:"deliveryLocName"`
	DeliveryPointName string `json:"deliveryPointName"`
	TruckNo         string `json:"truckNo"`
}
//103.230.124.241
var db *sql.DB

func main() {
	var err error
	db, err = sql.Open(
		"sqlserver",
		// "sqlserver://coadmin:alisha@1234@localhost?database=TODO&connection+timeout=30",
		"sqlserver://coadmin:tms@1234@103.230.124.241:1433?database=CoTMS&connection+timeout=30",
	)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	fmt.Println("Database connection successful")

	router := gin.Default()

	// Enable CORS
	router.Use(cors.Default())

	// Define the API endpoints
	router.GET("/joborder/:jobNo", getJob)
	// router.POST("/joborder/storeIssues", writeIssue)
	router.POST("/joborder/write", writeDataHandler)

	// Start the server
	router.Run(":3000")
}

func writeDataHandler(c *gin.Context) {
	// Path to your service account JSON key file
	keyFile := "service_account.json"

	// Create a new Google Sheets service client with the service account credentials
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(keyFile))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create Sheets service client",
		})
		return
	}

	// Parse the request body to get the data
	var data struct {
		Timestamp  			string `json:"timestamp"`
		Email				string `json:"email"`
		SelectedIssue       string `json:"selectedIssue"`
		SelectedHappened    string `json:"selectedHappened"`
		JobDate				string `json:"jobDate"`
		JobNo				string `json:"jobNo"`
		ReferenceNo			string `json:"referenceNo"`
		CustomerName		string `json:"customerName"`
		DeliveryPointName	string `json:"DeliveryPointName"`
		DeliveryLocName		string `json:"deliveryLocName"`
		Remarks             string `json:"remarks"`
		SelectedSettlement  string `json:"selectedSettlement"`
		SettlementDate  	string `json:"settlementDate"`
		TruckNo  			string `json:"truckNo"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	now := time.Now();
	currDate := now.Format("2006-01-02");
	currTime := now.Format("15:04:05");
	data.Timestamp = currDate + " " + currTime;
	fmt.Println("the timestamp", data.Timestamp)
	fmt.Println("the date", currDate)
	fmt.Println("the now", now)

	// Split the SettlementDate string at the "T" delimiter
	dateTimeParts := strings.Split(data.SettlementDate, "T")

	// Extract the date portion (index 0) from the resulting slice
	date := dateTimeParts[0]

	values := [][]interface{}{
		{data.Timestamp, data.Email, data.SelectedIssue, data.SelectedHappened, data.JobNo, data.JobDate, data.ReferenceNo, data.CustomerName, data.DeliveryPointName, data.DeliveryLocName, data.TruckNo, data.SelectedSettlement, date, data.Remarks, },
	}

	spreadsheetID := "1IKRI_CasrSyOCPbyhsDKQmE0tJGTlNbP0tVRF9FHPEQ"
	rangeValue := "Sheet1!A1:J2"
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Write the values to the spreadsheet
	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, rangeValue, valueRange).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to write data to spreadsheet" + err.Error(),
		})
		return
	}

	// Data successfully written
	c.JSON(http.StatusOK, gin.H{
		"message": "Data written to spreadsheet",
	})
}

func writeIssue(c *gin.Context) {
	ctx := context.Background()

	cl, err := getClient(ctx)
        if err != nil {
                log.Fatal(err)
        }

        sheetsService, err := sheets.New(cl)
        if err != nil {
                log.Fatal(err)
        }
	// The ID of the spreadsheet to update.
	spreadsheetID := "1IKRI_CasrSyOCPbyhsDKQmE0tJGTlNbP0tVRF9FHPEQ" // Replace with your spreadsheet ID

	// The A1 notation of the range to update.
	rangeValue := "Sheet1!A1:D2" // Replace with the desired range

	// The input data.
	rb := &sheets.ValueRange{
		Values: [][]interface{}{
			{"Value A1", "Value B1", "Value C1", "Value D1"},
			{"Value A2", "Value B2", "Value C2", "Value D2"},
		},
	}

	// Append the values to the spreadsheet
	resp, err := sheetsService.Spreadsheets.Values.Append(spreadsheetID, rangeValue, rb).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated %d cells.\n", resp.Updates.UpdatedCells)
}

func getClient(ctx context.Context) (*http.Client, error) {
	fmt.Println("getclient")
	tokenFile := "token.json"
	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %v", err)
	}
	fmt.Println("getclient 1")
	config, err := google.ConfigFromJSON(token, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}
	fmt.Println("getclient 2")
	client := config.Client(ctx, &oauth2.Token{
		AccessToken: string(token),
	})
	fmt.Println("getclient 3")
	return client, nil
}

func getJob(c *gin.Context) {
	fmt.Println("first line in getJob")
	jobNo := c.Param("jobNo")
	rows, err := db.Query("SELECT JobDate, JobNo, ReferenceNo, CustomerName, DeliveryLocName, DeliveryPointName, TruckNo FROM JobOrder WHERE JobNo=@JobNo", sql.Named("JobNo", jobNo))
	if err != nil {
		fmt.Println("this is error")
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Jobs"})
		return
	}
	defer rows.Close()
	fmt.Println("after defer rows close")

	jobs := []Job{}
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.JobDate, &job.JobNo, &job.ReferenceNo, &job.CustomerName, &job.DeliveryLocName, &job.DeliveryPointName, &job.TruckNo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		// job.JobDate = strings.TrimSpace(job.JobDate[:10]) // Slice to get the first 10 characters (date portion)
		parsedTime, err := time.Parse(time.RFC3339, job.JobDate)
		if err != nil {
			fmt.Println("Failed to parse date:", err)
			return
		}
		fmt.Println("parsedTime HERE", parsedTime)
		job.JobDate = parsedTime.Format("2006-01-02")
		fmt.Println("job dATE HERE", job.JobDate)

		jobs = append(jobs, job)
	}
	fmt.Println("end")
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, jobs)
	// c.JSON(http.StatusOK, jobs)
}

// func writeData(c *gin.Context) {
// 	apiKey := "AIzaSyCJNfVJMqWTCLAaWWcElM4t_5cxwj1B8_g"

// 	service, err := createService(apiKey)
// 	if err != nil {
// 		log.Fatalf("Unable to create Google Sheets service: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Internal Server Error",
// 		})
// 		return
// 	}

// 	spreadsheetID := "YOUR_SPREADSHEET_ID"
// 	sheetName := "YOUR_SHEET_NAME"

// 	// Parse the request body to get the data
// 	var data struct {
// 		JobDate				string `json:"jobDate"`
// 		ReferenceNo			string `json:"referenceNo"`
// 		CustomerName		string `json:"customerName"`
// 		DeliveryLocName		string `json:"deliveryLocName"`
// 		Remarks             string `json:"remarks"`
// 		SelectedHappened    string `json:"selectedHappened"`
// 		SelectedIssue       string `json:"selectedIssue"`
// 		SelectedSettlement  string `json:"selectedSettlement"`
// 		SettlementDate  	string `json:"settlementDate"`
// 		TruckNo  			string `json:"truckNo"`
// 	}
// 	if err := c.ShouldBindJSON(&data); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid request body",
// 		})
// 		return
// 	}

// 	values := [][]interface{}{
// 		{"jobDate", "referenceNo", "customerName", "deliveryLocName", "Remarks", "selectedHappened", "selectedIssue", "selectedSettlement", "settlementDate", "truckNo"},
// 		{data.JobDate, data.ReferenceNo, data.CustomerName, data.DeliveryLocName, data.Remarks, data.SelectedHappened, data.SelectedIssue, data.SelectedSettlement, data.SettlementDate, data.TruckNo},
// 	}

// 	rangeValue := sheetName + "!A1:D2"
// 	valueRange := &sheets.ValueRange{
// 		Values: values,
// 	}

// 	_, err = service.Spreadsheets.Values.Update(spreadsheetID, rangeValue, valueRange).ValueInputOption("USER_ENTERED").Do()
// 	if err != nil {
// 		if gErr, ok := err.(*googleapi.Error); ok {
// 			log.Fatalf("Unable to write data to Google Sheets: %v", gErr.Message)
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"error": "Internal Server Error",
// 			})
// 		} else {
// 			log.Fatalf("Unable to write data to Google Sheets: %v", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"error": "Internal Server Error",
// 			})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Data written successfully to Google Sheets",
// 	})
// }


