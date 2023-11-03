package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type GroupRequest struct {
	UserID  string `json:"userId"`
	GroupID string `json:"groupId"`
}

func joinGroup(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var requestData GroupRequest
	err := json.Unmarshal([]byte(request.Body), &requestData)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid input"}
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Users"),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(requestData.UserID),
			},
		},
		UpdateExpression:          aws.String("ADD Groups :group"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":group": {SS: []*string{aws.String(requestData.GroupID)}}},
	}

	_, err = db.UpdateItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error joining group"}
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Successfully joined group"}
}

func leaveGroup(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var requestData GroupRequest
	err := json.Unmarshal([]byte(request.Body), &requestData)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid input"}
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Users"),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(requestData.UserID),
			},
		},
		UpdateExpression:          aws.String("REMOVE Groups :group"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":group": {SS: []*string{aws.String(requestData.GroupID)}}},
	}

	_, err = db.UpdateItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error leaving group"}
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Successfully left group"}
}
