package rule

import (
	"go/token"
	"golang.org/x/tools/go/analysis"
)

type Rule interface {
	Name() string
	Check(msg string, pos token.Pos) []analysis.Diagnostic
}

type BaseRule struct {
	name    string
	enabled bool
}

func NewBaseRule(name string, enabled bool) BaseRule {
	return BaseRule{
		name:    name,
		enabled: enabled,
	}
}

func (b *BaseRule) Name() string {
	return b.name
}

func (b *BaseRule) IsEnabled() bool {
	return b.enabled
}
