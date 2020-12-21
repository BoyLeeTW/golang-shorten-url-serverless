package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"post-register-handler/internal/models"
	awsdynamodb "post-register-handler/internal/pkg/aws-dynamodb"
	"post-register-handler/internal/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	lambda.Start(lambdaEventHandler)
}

func lambdaEventHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("request:", request)
	log.Printf("%+v/n", request)

	body := models.PostRequestBody{}
	err := json.Unmarshal([]byte(request.Body), &body)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
	}
	log.Println("body:", body)

	// validate register url
	if _, err := url.ParseRequestURI(body.RegisterURL); err != nil {
		return events.APIGatewayProxyResponse{Body: "invalid register_url", StatusCode: http.StatusBadRequest}, nil
	}

	dynamoSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dynamoClient := dynamodb.New(dynamoSession)

	dynamodbService := awsdynamodb.NewDynamodbService(dynamoClient, os.Getenv("DYNAMODB_TABLE_NAME"))

	idProvider := services.NewIDProvider(2, 10, []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"})

	shortenURLService := services.NewShortenURLService(dynamodbService, idProvider)

	shortenedID, customError := shortenURLService.RegisterURLWithContext(ctx, body.RegisterURL)
	if customError != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: customError.HTTPStatusCode,
			Body:       customError.Message,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"shortened_id":"` + shortenedID + `"}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
