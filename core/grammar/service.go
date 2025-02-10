// components/grammar/service.go
package grammar

import "fmt"

type Service struct {
	generator *RandomTextGenerator
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Generate(grammarContent string) (string, error) {
	generator, err := NewRandomTextGenerator(grammarContent)
	if err != nil {
		return "", fmt.Errorf("failed to create generator: %w", err)
	}

	text := generator.Run()
	if text == "" {
		return "", fmt.Errorf("generated text is empty")
	}

	return text, nil
}
