package models

type ShortenedURL struct {
	ShortenedID string `dynamodbav:"shortened_id"`
	OriginalURL string `dynamodbav:"original_url"`
}
