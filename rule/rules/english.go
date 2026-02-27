package rules

import (
	"go/token"
	"log"
	"log-linter/rule"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"log-linter/config"
)

const EnglishName = "english"

type EnglishRule struct {
	rule.BaseRule
}

func NewEnglishRule(cfg config.RuleConfig) (rule.Rule, error) {
	return &EnglishRule{
		BaseRule: rule.NewBaseRule(EnglishName, cfg.Enabled),
	}, nil
}

func (r *EnglishRule) Check(msg string, pos token.Pos) []analysis.Diagnostic {
	if !r.IsEnabled() {
		return nil
	}

	for i, ch := range msg {
		if unicode.IsLetter(ch) && ch > unicode.MaxASCII {
			return []analysis.Diagnostic{{
				Pos:     pos + token.Pos(i),
				Message: "log message must contain only English letters",
			}}
		}
	}
	return nil
}

func init() {
	if err := rule.Global.Register(EnglishName, NewEnglishRule); err != nil {
		log.Printf("failed to register rule %q: %v", EnglishName, err)
	}
}
