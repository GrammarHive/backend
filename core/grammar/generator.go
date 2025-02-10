// core/generator/generator.go

package grammar

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// RandomTextGenerator represents a context-free grammar based text generator
type RandomTextGenerator struct {
	GrammarRules map[string][]string
	StartSymbol  string
}

// NewRandomTextGenerator creates and initializes a new RandomTextGenerator instance
func NewRandomTextGenerator(grammarFileContent string) (*RandomTextGenerator, error) {
	rtg := &RandomTextGenerator{
		GrammarRules: make(map[string][]string),
		StartSymbol:  "start", // looking at non-terminal without `<>`
	}
	lines := strings.Split(strings.TrimSpace(grammarFileContent), "\n")
	rtg.readGrammarRules(lines)

	// Validate grammar after reading rules
    if err := rtg.validateGrammar(); err != nil {
        return nil, err
    }

	return rtg, nil
}

// readGrammarRules parses the grammar rules from the input lines
func (rtg *RandomTextGenerator) readGrammarRules(lines []string) {
	var currentNonTerminal string
	var productions []string
	inRule := false

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		
		if line == "" {
			continue
		}

		switch {
		case line == "{":
			inRule = true
			productions = make([]string, 0)
		case line == "}":
			if currentNonTerminal != "" && len(productions) > 0 {
				rtg.GrammarRules[currentNonTerminal] = productions
			}
			inRule = false
			currentNonTerminal = ""
			productions = nil
		case inRule && strings.HasPrefix(line, "<") && strings.HasSuffix(line, ">"):
			currentNonTerminal = strings.Trim(line, "<>")
		case inRule && strings.HasSuffix(line, ";"):
			production := strings.TrimSuffix(line, ";")
			production = strings.TrimSpace(production)
			if production != "" {
				productions = append(productions, production)
			}
		}
	}
}

// expandSymbol recursively expands a grammar symbol
func (rtg *RandomTextGenerator) expandSymbol(symbol string, depth *int) string {
	if *depth > 800 { // configurable max depth
		return "Error: Maximum recursion depth exceeded"
	}

	if !strings.HasPrefix(symbol, "<") || !strings.HasSuffix(symbol, ">") {
		return symbol
	}

	nonTerminal := strings.Trim(symbol, "<>")
	productions, exists := rtg.GrammarRules[nonTerminal]
	
	if !exists {
		fmt.Printf("Warning: No production rules found for non-terminal: %s\n", nonTerminal)
		return symbol
	}

	rand.Seed(time.Now().UnixNano())
	production := productions[rand.Intn(len(productions))]

	symbols := strings.Fields(production)
	var result []string

	for _, sym := range symbols {
		*depth++
		result = append(result, rtg.expandSymbol(sym, depth))
	}

	return strings.Join(result, " ")
}

// Run generates random text by expanding the start symbol
func (rtg *RandomTextGenerator) Run() string {
	if len(rtg.GrammarRules) == 0 {
		return "Error: Grammar rules not properly initialized"
	}
	
	depthCount := 0
	result := rtg.expandSymbol("<" + rtg.StartSymbol + ">", &depthCount)
	return strings.TrimSpace(result)
}

