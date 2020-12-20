package main

import (
	"context"
	awsdynamodb "get-redirect-handler/internal/pkg/aws-dynamodb"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	lambda.Start(lambdaEventHandler)
}

func lambdaEventHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("context:", ctx)
	log.Printf("%+v/n", ctx)
	log.Println("request:", request)
	log.Printf("%+v/n", request)

	shortenedID, ok := request.PathParameters["shortened_id"]
	log.Println("shortenedID:", shortenedID)

	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	dynamoSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dynamoClient := dynamodb.New(dynamoSession)

	dynamodbService := awsdynamodb.NewDynamodbService(dynamoClient, os.Getenv("DYNAMODB_TABLE_NAME"))

	originalURL, err := dynamodbService.GetOriginalURLWithContext(ctx, shortenedID)

	if err != nil {
		log.Println("GetOriginalURLWithContext err:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: err.HTTPStatusCode,
			Headers: map[string]string{
				"Cache-Control": "max-age=10",
			},
			Body: err.Message,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMovedPermanently,
		Headers: map[string]string{
			"Cache-Control": "max-age=8640000",
			"Location":      originalURL,
		},
	}, nil
}
