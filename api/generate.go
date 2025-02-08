package api

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// RandomTextGenerator represents a context-free grammar based text generator
// that produces random text based on predefined grammar rules
type RandomTextGenerator struct {
	GrammarRules map[string][]string // Maps non-terminals to their production rules
	StartSymbol  string              // The starting symbol for text generation
}

// NewRandomTextGenerator creates and initializes a new RandomTextGenerator instance
// using the provided grammar file content
//
// Parameters:
//   - grammarFileContent: string containing the grammar rules in the specified format
//
// Returns:
//   - *RandomTextGenerator: pointer to the initialized generator
func NewRandomTextGenerator(grammarFileContent string) *RandomTextGenerator {
	rtg := &RandomTextGenerator{
		GrammarRules: make(map[string][]string),
		StartSymbol:  "start",
	}
	lines := strings.Split(strings.TrimSpace(grammarFileContent), "\n")
	rtg.readGrammarRules(lines)
	return rtg
}

// readGrammarRules parses the grammar rules from the input lines and populates
// the GrammarRules map
//
// Parameters:
//   - lines: slice of strings containing the grammar rules
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

// expandSymbol recursively expands a grammar symbol (terminal or non-terminal)
// according to the grammar rules
//
// Parameters:
//   - symbol: string representing the grammar symbol to expand
//
// Returns:
//   - string: the expanded text
func (rtg *RandomTextGenerator) expandSymbol(symbol string) string {
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
		result = append(result, rtg.expandSymbol(sym))
	}
	
	return strings.Join(result, " ")
}

// Run generates random text by expanding the start symbol according to
// the grammar rules
//
// Returns:
//   - string: the generated text, or an error message if generation fails
func (rtg *RandomTextGenerator) Run() string {
	if len(rtg.GrammarRules) == 0 {
		return "Error: Grammar rules not properly initialized"
	}
	
	result := rtg.expandSymbol("<" + rtg.StartSymbol + ">")
	return strings.TrimSpace(result)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	githubRawURL := "https://raw.githubusercontent.com/HarryZ10/api.resumes.guide/main/static/resume.g"
	resp, err := http.Get(githubRawURL)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Failed to fetch file", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read file contents", http.StatusInternalServerError)
		return
	}

	rtg := NewRandomTextGenerator(string(body))
	generatedText := rtg.Run()

	if strings.HasPrefix(generatedText, "Error:") {
		http.Error(w, generatedText, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": %q, "status": "OK"}`, generatedText)
}
