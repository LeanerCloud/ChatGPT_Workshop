package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type BirthdayRequest struct {
	UserID   string `json:"userId"`
	Birthday string `json:"birthday"`
}

func addOrUpdateBirthday(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var requestData BirthdayRequest
	err := json.Unmarshal([]byte(request.Body), &requestData)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid input"}
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Users"),
		Item: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(requestData.UserID),
			},
			"Birthday": {
				S: aws.String(requestData.Birthday),
			},
		},
	}

	_, err = db.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error saving birthday"}
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Birthday saved successfully"}
}

func listBirthdays(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	userID := request.QueryStringParameters["userId"]

	input := &dynamodb.GetItemInput{
		TableName: aws.String("Users"),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(userID),
			},
		},
	}

	result, err := db.GetItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error fetching data"}
	}

	// Convert the birthdays to JSON and return
	birthdayData, err := json.Marshal(result.Item)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error parsing data"}
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(birthdayData)}
}
