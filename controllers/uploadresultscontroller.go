package controllers

import (
	"awesomeProject/auth"
	"awesomeProject/database"
	"awesomeProject/models"
	"awesomeProject/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

func Upload(context *gin.Context) {

	userTestResultsPhoto, _ := context.FormFile("user_photo")
	partnerTestResultsPhoto, _ := context.FormFile("partner_photo")

	tokenString := context.GetHeader("Authorization")

	claims, err := auth.GetUserDetailsFromToken(tokenString)

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	email := claims.Email
	firstName := claims.FirstName
	lastName := claims.LastName
	phone := claims.Phone
	userId := claims.Id

	var userImageLink string
	if userTestResultsPhoto != nil {
		userImageLink, err = utils.SavePhoto(context, userTestResultsPhoto, phone)
		if err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading user test image"})
			return
		}
	}

	var partnerImageLink string
	if partnerTestResultsPhoto != nil {
		partnerImageLink, err = utils.SavePhoto(context, partnerTestResultsPhoto, phone)
		if err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading partner test image"})
			return
		}
	}

	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	credentialsFile := "./credentials.json"
	client, err := utils.GetClient(credentialsFile)
	if err != nil {
		log.Fatalf("Error getting Google Sheets client: %v", err)
	}

	nairobiLocation, err := time.LoadLocation("Africa/Nairobi")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}
	now := time.Now().In(nairobiLocation)
	formattedDateTime := now.Format("02/01/2006 15:04")

	results := models.Results{
		Results:        context.PostForm("results"),
		PartnerResults: context.PostForm("partner_results"),
		Image:          userImageLink,
		PartnerImage:   partnerImageLink,
		CareOption:     context.PostForm("care_option"),
		Date:           formattedDateTime,
		UserId:         userId,
	}

	sheetRange := "Sheet1!A1:J5"
	values := [][]interface{}{
		{firstName, lastName, phone, email, results.Results, results.PartnerResults, results.Image, results.PartnerImage, results.CareOption, formattedDateTime},
	}

	err = utils.WriteDataToSpreadsheet(client, spreadsheetID, sheetRange, values)
	if err != nil {
		log.Fatalf("Error writing data to spreadsheet: %v", err)
	}
	record := database.Instance.Create(&results)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": record.Error.Error()})
		context.Abort()
		return
	}

	fmt.Println("Data written to the spreadsheet successfully!")
	context.JSON(http.StatusOK, gin.H{"message": "Test submitted successfully"})
}
