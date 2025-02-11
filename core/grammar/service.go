// components/grammar/service.go
package grammar

import (
	"fmt"
	"sync"
)

type Service struct {
	generator *RandomTextGenerator
}

func NewGrammarGenService() *Service {
	return &Service{}
}

func (s *Service) ExecuteGrammarGen(grammarContent string) (string, error) {
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

// GenerateMultiple generates n texts concurrently
func (s *Service) GenerateMultiple(grammarContent string, count int) ([]string, error) {
    generator, err := NewRandomTextGenerator(grammarContent)
    if err != nil {
        return nil, fmt.Errorf("failed to create generator: %w", err)
    }

    messages := make([]string, count)
    var wg sync.WaitGroup
    errChan := make(chan error, count)

    // Generate texts concurrently
    for i := 0; i < count; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()
            text := generator.Run()
            if text == "" {
                errChan <- fmt.Errorf("generated text is empty at index %d", index)
                return
            }
            messages[index] = text
        }(i)
    }

    // Wait for all generations to complete
    wg.Wait()
    close(errChan)

    // Check for any errors
    if len(errChan) > 0 {
        return nil, <-errChan // Return the first error encountered
    }

    return messages, nil
}
