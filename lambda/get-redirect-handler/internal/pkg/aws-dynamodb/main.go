package awsdynamodb

import (
	"context"
	customerrors "get-redirect-handler/internal/custom-errors"
	"get-redirect-handler/internal/interfaces"
	"get-redirect-handler/internal/models"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamodbService struct {
	dynamodbClient    interfaces.DynamodbClientInterface
	dynamodbTableName string
}

func NewDynamodbService(dynamodbClient interfaces.DynamodbClientInterface, dynamodbTableName string) *DynamodbService {
	return &DynamodbService{
		dynamodbClient:    dynamodbClient,
		dynamodbTableName: dynamodbTableName,
	}
}

func (s *DynamodbService) GetOriginalURLWithContext(ctx context.Context, shortenedID string) (string, *customerrors.CustomError) {
	shortenedURL := models.ShortenedURL{}

	output, err := s.dynamodbClient.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.dynamodbTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"shortened_id": {
				S: aws.String(shortenedID),
			},
		},
	})

	if err != nil {
		log.Println("GetItemWithContext err:", err)
		switch err.(type) {
		case *dynamodb.ResourceNotFoundException: // table or index not found
			return "", ErrResourceNotFound
		default:
			return "", ErrInternalServer
		}
	}

	// no item exist with such partition key
	if output.Item == nil {
		return "", ErrItemNotExist
	}

	err = dynamodbattribute.UnmarshalMap(output.Item, &shortenedURL)

	if err != nil {
		log.Println("unmarshal error:", err.Error())
		return "", ErrInternalServer
	}

	return shortenedURL.OriginalURL, nil
}
