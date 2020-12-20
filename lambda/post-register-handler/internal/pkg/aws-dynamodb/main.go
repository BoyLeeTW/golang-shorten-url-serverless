package awsdynamodb

import (
	"context"
	"log"
	customerrors "post-register-handler/internal/custom-errors"
	"post-register-handler/internal/interfaces"
	"post-register-handler/internal/models"

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

func (s *DynamodbService) PutShortenedURLWithContext(ctx context.Context, shortenedURL models.ShortenedURL) (bool, *customerrors.CustomError) {
	attributeValue, err := dynamodbattribute.MarshalMap(shortenedURL)
	if err != nil {
		log.Println("marshal error:", err)
		return false, ErrInternalServer
	}

	_, err = s.dynamodbClient.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		// check only if id doesn't exist
		ConditionExpression: aws.String("attribute_not_exists(shortened_id)"),
		Item:                attributeValue,
		TableName:           aws.String(s.dynamodbTableName),
	})

	if err != nil {
		switch err.(type) {
		case *dynamodb.ConditionalCheckFailedException:
			log.Println("dynamodb.ConditionalCheckFailedException")
			// regenerate id again
			return true, nil
		default:
			log.Println("put item error:", err)
			return false, ErrInternalServer
		}
	}

	return false, nil
}
