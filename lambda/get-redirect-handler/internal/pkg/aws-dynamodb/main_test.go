package awsdynamodb

import (
	"context"
	"get-redirect-handler/internal/mock"
	"get-redirect-handler/internal/models"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/golang/mock/gomock"
)

func TestDynamodbService_GetOriginalURLWithContext(t *testing.T) {
	t.Run(`should return ("", ErrResourceNotFound) if receive ResourceNotFoundException error from GetItemWithContext`, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockedDynamodbClient := mock.NewMockDynamodbClientInterface(ctrl)

		dynamodbService := NewDynamodbService(mockedDynamodbClient, "mocked_table_name")

		mockedDynamodbClient.
			EXPECT().
			GetItemWithContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, &dynamodb.ResourceNotFoundException{}).
			Times(1)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		output, err := dynamodbService.GetOriginalURLWithContext(ctx, "test_id")

		if output != "" {
			t.Errorf("should be empty string but got %+v\n", output)
		}

		if err != ErrResourceNotFound {
			t.Errorf("err should be ErrResourceNotFound but got %+v\n", err)
		}
	})

	t.Run(`should return ("", ErrInternalServer) if receive all other error from GetItemWithContext`, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockedDynamodbClient := mock.NewMockDynamodbClientInterface(ctrl)

		dynamodbService := NewDynamodbService(mockedDynamodbClient, "mocked_table_name")

		mockedDynamodbClient.
			EXPECT().
			GetItemWithContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, &dynamodb.LimitExceededException{}).
			Times(1)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		output, err := dynamodbService.GetOriginalURLWithContext(ctx, "test_id")

		if output != "" {
			t.Errorf("should be empty string but got %+v\n", output)
		}

		if err != ErrInternalServer {
			t.Errorf("err should be ErrInternalServer but got %+v\n", err)
		}
	})

	t.Run(`should return ("", ErrItemNotExist) if receive no error and empty Item from GetItemWithContext`, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockedDynamodbClient := mock.NewMockDynamodbClientInterface(ctrl)

		dynamodbService := NewDynamodbService(mockedDynamodbClient, "mocked_table_name")

		mockedDynamodbClient.
			EXPECT().
			GetItemWithContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&dynamodb.GetItemOutput{}, nil).
			Times(1)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		output, err := dynamodbService.GetOriginalURLWithContext(ctx, "test_id")

		if output != "" {
			t.Errorf("should be empty string but got %+v\n", output)
		}

		if err != ErrItemNotExist {
			t.Errorf("err should be ErrItemNotExist but got %+v\n", err)
		}
	})

	t.Run(`should return ("mocked_url", nil) if receive no error and non-empty Item from GetItemWithContext`, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockedDynamodbClient := mock.NewMockDynamodbClientInterface(ctrl)

		dynamodbService := NewDynamodbService(mockedDynamodbClient, "mocked_table_name")

		mockedShortenedURL := models.ShortenedURL{
			ShortenedID: "mocked_id",
			OriginalURL: "mocked_url",
		}
		attributeValue, _ := dynamodbattribute.MarshalMap(mockedShortenedURL)

		mockedDynamodbClient.
			EXPECT().
			GetItemWithContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&dynamodb.GetItemOutput{
				Item: attributeValue,
			}, nil).
			Times(1)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		output, err := dynamodbService.GetOriginalURLWithContext(ctx, "mocked_id")

		if output != "mocked_url" {
			t.Errorf("should be mocked_url but got %+v\n", output)
		}

		if err != err {
			t.Errorf("err should be nil but got %+v\n", err)
		}
	})

}
