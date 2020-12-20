package awsdynamodb

import (
	"context"
	"post-register-handler/internal/mock"
	"post-register-handler/internal/models"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/golang/mock/gomock"
)

func TestDynamodbService_PutShortenedURLWithContext(t *testing.T) {
	t.Run(`should return (true, nil) if receive ResourceNotFoundException error from PutItemWithContext`, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockedDynamodbClient := mock.NewMockDynamodbClientInterface(ctrl)

		dynamodbService := NewDynamodbService(mockedDynamodbClient, "mocked_table_name")

		mockedDynamodbClient.
			EXPECT().
			PutItemWithContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, &dynamodb.ConditionalCheckFailedException{}).
			Times(1)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		needNewID, err := dynamodbService.PutShortenedURLWithContext(ctx, models.ShortenedURL{})

		if needNewID != true {
			t.Errorf("should be true but got %+v\n", needNewID)
		}

		if err != nil {
			t.Errorf("err should be nil but got %+v\n", err)
		}
	})

	t.Run(`should return (false, ErrInternalServer) if receive other error from PutItemWithContext`, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockedDynamodbClient := mock.NewMockDynamodbClientInterface(ctrl)

		dynamodbService := NewDynamodbService(mockedDynamodbClient, "mocked_table_name")

		mockedDynamodbClient.
			EXPECT().
			PutItemWithContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, &dynamodb.ResourceNotFoundException{}).
			Times(1)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		needNewID, err := dynamodbService.PutShortenedURLWithContext(ctx, models.ShortenedURL{})

		if needNewID != false {
			t.Errorf("should be false but got %+v\n", needNewID)
		}

		if err != ErrInternalServer {
			t.Errorf("err should be ErrInternalServer but got %+v\n", err)
		}
	})

	t.Run(`should return (false, nil) if PutItemWithContext successfully`, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockedDynamodbClient := mock.NewMockDynamodbClientInterface(ctrl)

		dynamodbService := NewDynamodbService(mockedDynamodbClient, "mocked_table_name")

		mockedDynamodbClient.
			EXPECT().
			PutItemWithContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, nil).
			Times(1)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		needNewID, err := dynamodbService.PutShortenedURLWithContext(ctx, models.ShortenedURL{})

		if needNewID != false {
			t.Errorf("should be false but got %+v\n", needNewID)
		}

		if err != nil {
			t.Errorf("err should be nil but got %+v\n", err)
		}
	})
}
