package rules

import (
	"fmt"
	"go/token"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"log-linter/config"
	"log-linter/rule"
)

const SensitiveName = "sensitive"

type SensitiveRule struct {
	rule.BaseRule
	patterns []*regexp.Regexp
}

func NewSensitiveRule(cfg config.RuleConfig) (rule.Rule, error) {
	var patterns []*regexp.Regexp

	if cfg.Options != nil {
		if patternsList, ok := cfg.Options["patterns"].([]string); ok {
			for _, patternStr := range patternsList {
				re, err := regexp.Compile(patternStr)
				if err != nil {
					continue
				}
				patterns = append(patterns, re)
			}
			return &SensitiveRule{
				BaseRule: rule.NewBaseRule(SensitiveName, cfg.Enabled),
				patterns: patterns,
			}, nil
		}

		if patternsList, ok := cfg.Options["patterns"].([]interface{}); ok {
			for _, p := range patternsList {
				if patternStr, ok := p.(string); ok {
					re, err := regexp.Compile(patternStr)
					if err != nil {
						continue
					}
					patterns = append(patterns, re)
				}
			}
		}
	}

	return &SensitiveRule{
		BaseRule: rule.NewBaseRule(SensitiveName, cfg.Enabled),
		patterns: patterns,
	}, nil
}

func (r *SensitiveRule) Check(msg string, pos token.Pos) []analysis.Diagnostic {
	if !r.IsEnabled() || len(r.patterns) == 0 {
		return nil
	}

	var diagnostics []analysis.Diagnostic

	for _, re := range r.patterns {
		matches := re.FindAllStringIndex(msg, -1)

		for _, match := range matches {
			start := match[0]

			diagnostics = append(diagnostics, analysis.Diagnostic{
				Pos:     pos + token.Pos(start),
				Message: fmt.Sprintf("log message contains sensitive data"),
			})
		}
	}

	return diagnostics
}

func init() {
	rule.Global.Register(SensitiveName, NewSensitiveRule)
}
