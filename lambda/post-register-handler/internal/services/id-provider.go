package services

import (
	"math/rand"
	"time"
)

type IDProvider struct {
	maximumDigits int
	minimalDigits int
	validElements []string
}

func NewIDProvider(minimalDigits int, maximumDigits int, validElements []string) *IDProvider {
	return &IDProvider{
		maximumDigits: maximumDigits,
		minimalDigits: minimalDigits,
		validElements: validElements,
	}
}

func (p *IDProvider) GenerateID() string {
	rand.Seed(time.Now().UnixNano())
	result := ""

	idLength := rand.Intn(p.maximumDigits-p.minimalDigits+1) + p.minimalDigits

	validElementsLength := len(p.validElements)
	for i := 0; i < idLength; i++ {
		result += p.validElements[rand.Intn(validElementsLength)]
	}
	return result
}
