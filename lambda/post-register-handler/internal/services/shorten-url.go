package services

import (
	"context"
	customerrors "post-register-handler/internal/custom-errors"
	"post-register-handler/internal/interfaces"
	"post-register-handler/internal/models"
)

type ShortenURLService struct {
	dynamodbService interfaces.DynamoServiceInterface
	idProvider      interfaces.IDProviderInterface
}

func NewShortenURLService(dynamodbService interfaces.DynamoServiceInterface, idProvider interfaces.IDProviderInterface) *ShortenURLService {
	return &ShortenURLService{
		dynamodbService: dynamodbService,
		idProvider:      idProvider,
	}
}

func (s *ShortenURLService) RegisterURLWithContext(ctx context.Context, registerURL string) (string, *customerrors.CustomError) {
	for {
		id := s.idProvider.GenerateID()

		shortenedURL := models.ShortenedURL{
			ShortenedID: id,
			OriginalURL: registerURL,
		}

		needNewID, err := s.dynamodbService.PutShortenedURLWithContext(ctx, shortenedURL)

		// error occur
		if !needNewID && err != nil {
			return "", err
		}

		// successfully register
		if !needNewID && err == nil {
			return id, nil
		}
	}
}
