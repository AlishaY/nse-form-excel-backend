func writeDatas(c *gin.Context) {
		apiKey := "AIzaSyCJNfVJMqWTCLAaWWcElM4t_5cxwj1B8_g"
	
		service, err := createService(apiKey)
		if err != nil {
			log.Fatalf("Unable to create Google Sheets service: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}
	
		spreadsheetID := "YOUR_SPREADSHEET_ID"
		sheetName := "YOUR_SHEET_NAME"
	
		// Parse the request body to get the data
		var data struct {
			JobDate				string `json:"jobDate"`
			ReferenceNo			string `json:"referenceNo"`
			CustomerName		string `json:"customerName"`
			DeliveryLocName		string `json:"deliveryLocName"`
			Remarks             string `json:"remarks"`
			SelectedHappened    string `json:"selectedHappened"`
			SelectedIssue       string `json:"selectedIssue"`
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
	
		values := [][]interface{}{
			{"jobDate", "referenceNo", "customerName", "deliveryLocName", "Remarks", "selectedHappened", "selectedIssue", "selectedSettlement", "settlementDate", "truckNo"},
			{data.JobDate, data.ReferenceNo, data.CustomerName, data.DeliveryLocName, data.Remarks, data.SelectedHappened, data.SelectedIssue, data.SelectedSettlement, data.SettlementDate, data.TruckNo},
		}
	
		rangeValue := sheetName + "!A1:D2"
		valueRange := &sheets.ValueRange{
			Values: values,
		}
	
		_, err = service.Spreadsheets.Values.Update(spreadsheetID, rangeValue, valueRange).ValueInputOption("USER_ENTERED").Do()
		if err != nil {
			if gErr, ok := err.(*googleapi.Error); ok {
				log.Fatalf("Unable to write data to Google Sheets: %v", gErr.Message)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
			} else {
				log.Fatalf("Unable to write data to Google Sheets: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
			}
			return
		}
	
		c.JSON(http.StatusOK, gin.H{
			"message": "Data written successfully to Google Sheets",
		})
	}