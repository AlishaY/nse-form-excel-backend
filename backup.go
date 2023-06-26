package main

import (
	"io/ioutil"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	// "path/filepath"
	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	// "google.golang.org/api/googleapi"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Jobs struct {
	JobDate         string `json:"jobDate"`
	ReferenceNo     string `json:"referenceNo"`
	CustomerName    string `json:"customerName"`
	DeliveryLocName string `json:"deliveryLocName"`
	TruckNo         string `json:"truckNo"`
}

var dbs *sql.DB

// var (
// 	oauthConfig *oauth2.Config
// )

func m() {
	// Load the OAuth2 credentials file
	credentialsFile := "client_secret_465912978126-5ckh3s4qs84rhml2j0vqlqknjpf4spj7.apps.googleusercontent.com.json"

	// Set up OAuth2 configuration
	// oauthConfig = getOAuthConfig(credentialsFile)

	// Create a new Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(cors.Default())

	// Define the API endpoints
	router.GET("/joborder/:jobNo", getJob)
	// router.POST("/joborder/storeIssues", writeData)
	// Define the API endpoint to write data to Google Sheets
	router.POST("/joborder/storeIssues", func(c *gin.Context) {
		writeData(c, credentialsFile)
	})

	// Start the server
	router.Run(":3000")
}

func writeDatas(c *gin.Context, credentialsFile string) {
	// Read the OAuth2 credentials file
	credentials, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}

	// Parse the credentials file
	config, err := google.ConfigFromJSON(credentials, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse credentials: %v", err)
	}

	// Parse the request body to get the data
	var data struct {
		JobDate            string `json:"jobDate"`
		ReferenceNo        string `json:"referenceNo"`
		CustomerName       string `json:"customerName"`
		DeliveryLocName    string `json:"deliveryLocName"`
		Remarks            string `json:"remarks"`
		SelectedHappened   string `json:"selectedHappened"`
		SelectedIssue      string `json:"selectedIssue"`
		SelectedSettlement string `json:"selectedSettlement"`
		SettlementDate     string `json:"settlementDate"`
		TruckNo            string `json:"truckNo"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Create a new HTTP client
	client := getClient(config)

	// Create a new Sheets service
	service, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}

	// Write data to a Google Sheet
	spreadsheetID := "your-spreadsheet-id"
	rangeValue := "Sheet1!A1:B1"
	values := [][]interface{}{
		{data.JobDate, data.ReferenceNo, data.CustomerName, data.DeliveryLocName, data.Remarks, data.SelectedHappened, data.SelectedIssue, data.SelectedSettlement, data.SettlementDate, data.TruckNo},
	}
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	_, err = service.Spreadsheets.Values.Update(spreadsheetID, rangeValue, valueRange).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Unable to write data to Google Sheet: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data written successfully to Google Sheet",
	})
}

// getClient returns an HTTP client with the specified OAuth2 config
func getClients(config *oauth2.Config) *http.Client {
	token := getTokenFromWeb(config)
	return config.Client(context.Background(), token)
}

// getTokenFromWeb uses the provided OAuth2 configuration to obtain an access token.
func getTokenFromWebs(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Open the following URL in your browser and authorize the application:\n%s\n", authURL)

	fmt.Print("Enter the authorization code: ")
	var authCode string
	_, err := fmt.Scan(&authCode)
	if err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	return token
}

func getOAuthConfigs(credentialsFile string) *oauth2.Config {
	credentials, err := os.ReadFile(credentialsFile)
	fmt.Println("credentials")
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}
	fmt.Println("credentials 2")
	oauthCreds, err := google.ConfigFromJSON(credentials, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse credentials: %v", err)
	}
	fmt.Println("credentials 3 last")
	return oauthCreds
}

func createServicess(ctx context.Context, token *oauth2.Token) (*sheets.Service, error) {
	client := oauthConfig.Client(ctx, token)
	fmt.Println("token")

	// Create a token source that handles token refreshing
	tokenSource := oauthConfig.TokenSource(ctx, token)
	client = oauth2.NewClient(ctx, tokenSource)
	fmt.Println("token 2")

	service, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	fmt.Println("token 3")
	if err != nil {
		log.Fatalf("Unable to create Google Sheets service: %v", err)
		return nil, err
	}
	fmt.Println("token 4 last")
	return service, nil
}

// func writeData(c *gin.Context) {
// 	ctx := context.Background()
// 	fmt.Println("write")

// 	// Parse the request body to get the data
// 	var data struct {
// 		JobDate            string `json:"jobDate"`
// 		ReferenceNo        string `json:"referenceNo"`
// 		CustomerName       string `json:"customerName"`
// 		DeliveryLocName    string `json:"deliveryLocName"`
// 		Remarks            string `json:"remarks"`
// 		SelectedHappened   string `json:"selectedHappened"`
// 		SelectedIssue      string `json:"selectedIssue"`
// 		SelectedSettlement string `json:"selectedSettlement"`
// 		SettlementDate     string `json:"settlementDate"`
// 		TruckNo            string `json:"truckNo"`
// 	}
// 	if err := c.ShouldBindJSON(&data); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid request body",
// 		})
// 		return
// 	}
// 	fmt.Println("write 2")
// 	// Get the access token from the request headers or session
// 	accessToken := c.Request.Header.Get("Authorization")
// 	token := &oauth2.Token{
// 		AccessToken: accessToken,
// 	}
// 	fmt.Println("write 3")
// 	fmt.Printf("access-token: %s", accessToken)
// 	service, err := createService(ctx, token)
// 	fmt.Println(service)
// 	if err != nil {
// 		log.Fatalf("Unable to create Google Sheets service: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Internal Server Error",
// 		})
// 		return
// 	}
// 	fmt.Println("write 4")
// 	spreadsheetID := "1IKRI_CasrSyOCPbyhsDKQmE0tJGTlNbP0tVRF9FHPEQ"
// 	sheetName := "sample1"
// 	fmt.Println("write 5")
// 	values := [][]interface{}{
// 		{"jobDate", "referenceNo", "customerName", "deliveryLocName", "Remarks", "selectedHappened", "selectedIssue", "selectedSettlement", "settlementDate", "truckNo"},
// 		{data.JobDate, data.ReferenceNo, data.CustomerName, data.DeliveryLocName, data.Remarks, data.SelectedHappened, data.SelectedIssue, data.SelectedSettlement, data.SettlementDate, data.TruckNo},
// 	}
// 	fmt.Println("write 6")
// 	rangeValue := sheetName + "!A1:D2"
// 	valueRange := &sheets.ValueRange{
// 		Values: values,
// 	}
// 	fmt.Println("write 7")
// 	_, err = service.Spreadsheets.Values.Update(spreadsheetID, rangeValue, valueRange).ValueInputOption("USER_ENTERED").Do()
// 	if err != nil {
// 		log.Fatalf("Unable to write data to Google Sheets: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Internal Server Error",
// 		})
// 		return
// 	}
// 	fmt.Println("write 8 last")
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Data written successfully to Google Sheets",
// 	})
// }

// func createService(credentialsFile string) (*sheets.Service, error) {
// 	ctx := context.Background()

// 	service, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile), option.WithScopes(sheets.SpreadsheetsScope))
// 	if err != nil {
// 		log.Fatalf("Unable to create Google Sheets service: %v", err)
// 		return nil, err
// 	}

// 	return service, nil
// }

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

func getJobs(c *gin.Context) {
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
