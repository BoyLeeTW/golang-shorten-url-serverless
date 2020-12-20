package interfaces

import (
	"context"
	customerrors "post-register-handler/internal/custom-errors"
	"post-register-handler/internal/models"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoServiceInterface interface {
	PutShortenedURLWithContext(ctx context.Context, shortenedURL models.ShortenedURL) (bool, *customerrors.CustomError)
}

type DynamodbClientInterface interface {
	PutItemWithContext(ctx context.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error)
}
