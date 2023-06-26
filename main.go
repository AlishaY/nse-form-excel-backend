package main

import (
	"io/ioutil"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	// "errors"
	// "os"
	"strings"
	// "path/filepath"
	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	// "golang.org/x/net/context"
	// "google.golang.org/api/googleapi"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)
//json uri redirect: http://localhost:8080/callback
type Job struct {
	JobDate         string `json:"jobDate"`
	ReferenceNo     string `json:"referenceNo"`
	CustomerName    string `json:"customerName"`
	DeliveryLocName string `json:"deliveryLocName"`
	TruckNo         string `json:"truckNo"`
}

var db *sql.DB

func main() {
	router := gin.Default()

	// Enable CORS
	router.Use(cors.Default())

	// Define the API endpoints
	router.GET("/joborder/:jobNo", getJob)
	router.POST("/joborder/storeIssues", writeIssue)
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

	// Specify the spreadsheet ID and range
	spreadsheetID := "1IKRI_CasrSyOCPbyhsDKQmE0tJGTlNbP0tVRF9FHPEQ"
	rangeValue := "sample1!A1:B3"

	// Create the value range to write
	values := [][]interface{}{
		{"Value A1", "Value B1"},
		{"Value A2", "Value B2"},
		{"Value A2", "Value B2"},
	}
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Write the values to the spreadsheet
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, rangeValue, valueRange).ValueInputOption("USER_ENTERED").Do()
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

	// service, err := getServiceClient(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

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

func getServiceClient(ctx context.Context) (*sheets.Service, error) {
	credentialsFile := "credentials.json" // Update the filename to match your credentials file

	// Read the credentials file
	credentials, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %v", err)
	}

	// Create a Config object from the credentials file
	config, err := google.ConfigFromJSON(credentials, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %v", err)
	}

	token := &oauth2.Token{
		// If you have an access token saved in the token.json file, you can assign it here
		AccessToken: "ya29.a0AWY7CklmInTipDPAw3f8S1wzhBHjd0CkTzNH5rVp7Fc6d6CEWMdGbdrsJjRoB29c1lTt_BmpGF8UyiHSPAlGWzAlo-Ndfc_6utC4BRjGsorW8HYznAD_b16jPNetKtS3ijyR39H_xanEeUFqjjqq8ZvIL426aCgYKAd8SARISFQG1tDrpPppg6SjW_qI8xH9-5M7HiQ0163",
	}

	// Obtain a token source from the Config
	tokenSource := config.TokenSource(ctx, token)

	// Create a new HTTP client using the token source
	client := oauth2.NewClient(ctx, tokenSource)

	// Create a new Google Sheets service client
	service, err := sheets.New(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create Sheets client: %v", err)
	}

	return service, nil
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
	rows, err := db.Query("SELECT JobDate, ReferenceNo, CustomerName, DeliveryLocName, TruckNo FROM JobOrder WHERE JobNo=@JobNo", sql.Named("JobNo", jobNo))
	if err != nil {
		fmt.Println("this is error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Jobs"})
		return
	}
	defer rows.Close()
	fmt.Println("after defer rows close")

	jobs := []Job{}
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.JobDate, &job.ReferenceNo, &job.CustomerName, &job.DeliveryLocName, &job.TruckNo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		job.JobDate = strings.TrimSpace(job.JobDate[:10]) // Slice to get the first 10 characters (date portion)
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


