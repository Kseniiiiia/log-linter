package rules

import (
	"go/token"
	"log"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"log-linter/config"
	"log-linter/rule"
)

const LowercaseName = "lowercase"

type LowercaseRule struct {
	rule.BaseRule
}

func NewLowercaseRule(cfg config.RuleConfig) (rule.Rule, error) {
	return &LowercaseRule{
		BaseRule: rule.NewBaseRule(LowercaseName, cfg.Enabled),
	}, nil
}

func (r *LowercaseRule) Check(msg string, pos token.Pos) []analysis.Diagnostic {
	if !r.IsEnabled() {
		return nil
	}

	rs := []rune(msg)
	i := 0
	for i < len(rs) && unicode.IsSpace(rs[i]) {
		i++
	}
	if i >= len(rs) {
		return nil
	}
	if unicode.IsLetter(rs[i]) && unicode.IsUpper(rs[i]) {
		return []analysis.Diagnostic{{
			Pos:     pos + token.Pos(i),
			Message: "log message must start with lowercase letter",
		}}
	}
	return nil
}

func init() {
	if err := rule.Global.Register(LowercaseName, NewLowercaseRule); err != nil {
		log.Printf("failed to register rule %q: %v", LowercaseName, err)
	}
}
