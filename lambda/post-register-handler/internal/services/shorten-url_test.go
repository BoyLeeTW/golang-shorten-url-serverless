package services

import (
	"context"
	"post-register-handler/internal/mock"
	"post-register-handler/internal/models"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestShortenURLService_RegisterURLWithContext(t *testing.T) {
	t.Run("test if id from idProvider is passed to dynamodbService correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockedDynamodbService := mock.NewMockDynamoServiceInterface(ctrl)
		mockedIDProvider := mock.NewMockIDProviderInterface(ctrl)

		shortenURLService := NewShortenURLService(
			mockedDynamodbService,
			mockedIDProvider,
		)

		mockedIDProvider.
			EXPECT().
			GenerateID().
			Return("mocked_id").
			Times(1)

		mockedDynamodbService.
			EXPECT().
			PutShortenedURLWithContext(gomock.Any(), models.ShortenedURL{
				OriginalURL: "register_url",
				ShortenedID: "mocked_id",
			}).
			Times(1).Return(false, nil)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		shortenURLService.RegisterURLWithContext(ctx, "register_url")
	})
}
