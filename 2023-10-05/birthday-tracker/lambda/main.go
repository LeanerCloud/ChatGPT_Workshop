package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var db = dynamodb.New(session.New(), aws.NewConfig())

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Routing based on the method and path.
	switch request.HTTPMethod {
	case "POST":
		if request.Path == "/birthday" {
			return addOrUpdateBirthday(request), nil
		} else if request.Path == "/group/join" {
			return joinGroup(request), nil
		} else if request.Path == "/group/leave" {
			return leaveGroup(request), nil
		}
	case "GET":
		if request.Path == "/birthdays" {
			return listBirthdays(request), nil
		}
	}

	return events.APIGatewayProxyResponse{StatusCode: 404}, nil
}

func main() {
	lambda.Start(handler)
}
