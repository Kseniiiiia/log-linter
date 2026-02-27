package rules

import (
	"go/token"
	"log"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"log-linter/config"
	"log-linter/rule"
)

const SymbolsName = "symbols"

type SymbolsRule struct {
	rule.BaseRule
	allowed map[rune]bool
}

func NewSymbolsRule(cfg config.RuleConfig) (rule.Rule, error) {
	allowed := make(map[rune]bool)

	if cfg.Options != nil {
		if val, exists := cfg.Options["allowed"]; exists {
			if allowedList, ok := val.([]interface{}); ok {
				for _, v := range allowedList {
					if s, ok := v.(string); ok && len(s) == 1 {
						for _, ch := range s {
							allowed[ch] = true
						}
					}
				}
			}

			if allowedList, ok := val.([]string); ok {
				for _, s := range allowedList {
					if len(s) == 1 {
						for _, ch := range s {
							allowed[ch] = true
						}
					}
				}
			}
		}
	}

	return &SymbolsRule{
		BaseRule: rule.NewBaseRule(SymbolsName, cfg.Enabled),
		allowed:  allowed,
	}, nil
}

func (r *SymbolsRule) Check(msg string, pos token.Pos) []analysis.Diagnostic {
	if !r.IsEnabled() {
		return nil
	}

	for i, ch := range msg {
		if (unicode.IsSymbol(ch) || unicode.IsPunct(ch)) && !r.allowed[ch] {
			return []analysis.Diagnostic{{
				Pos:     pos + token.Pos(i),
				Message: "log message must not contain special characters or emojis",
			}}
		}
	}

	return nil
}

func init() {
	if err := rule.Global.Register(SymbolsName, NewSymbolsRule); err != nil {
		log.Printf("failed to register rule %q: %v", SymbolsName, err)
	}
}
