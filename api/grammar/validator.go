// api/generator/validator.go

package grammar

import (
	"fmt"
	"strings"
)

// validateGrammar checks for undefined non-terminals in the grammar rules
func (rtg *RandomTextGenerator) validateGrammar() error {
	// Check for undefined non-terminals
	for _, productions := range rtg.GrammarRules {
		for _, prod := range productions {
			symbols := strings.Fields(prod)
			for _, sym := range symbols {
				if strings.HasPrefix(sym, "<") && strings.HasSuffix(sym, ">") {
					nonTerm := strings.Trim(sym, "<>")
					if _, exists := rtg.GrammarRules[nonTerm]; !exists {
						return fmt.Errorf("undefined non-terminal: %s", nonTerm)
					}
				}
			}
		}
	}
	return nil
}
