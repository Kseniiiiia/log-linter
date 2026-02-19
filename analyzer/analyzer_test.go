package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"log-linter/analyzer"
	"log-linter/config"
	"log-linter/rule"
	_ "log-linter/rule/rules"
)

func TestLowercaseRule(t *testing.T) {
	cfg := config.DefaultConfig()
	for name := range cfg.Rules {
		if name != "lowercase" {
			cfg.Rules[name] = config.RuleConfig{
				Enabled: false,
				Options: map[string]interface{}{},
			}
		}
	}

	rules := createRulesFromConfig(t, cfg)
	analyzer := analyzer.NewAnalyzer(rules)

	analysistest.Run(t, analysistest.TestData(), analyzer, "lowercase")
}

func TestEnglishRule(t *testing.T) {
	cfg := config.DefaultConfig()

	for name := range cfg.Rules {
		if name != "english" {
			cfg.Rules[name] = config.RuleConfig{
				Enabled: false,
				Options: map[string]interface{}{},
			}
		}
	}

	rules := createRulesFromConfig(t, cfg)
	analyzer := analyzer.NewAnalyzer(rules)

	analysistest.Run(t, analysistest.TestData(), analyzer, "english")
}

func TestSymbolsRule(t *testing.T) {
	cfg := config.DefaultConfig()

	for name := range cfg.Rules {
		if name != "symbols" {
			cfg.Rules[name] = config.RuleConfig{
				Enabled: false,
				Options: map[string]interface{}{},
			}
		}
	}

	rules := createRulesFromConfig(t, cfg)
	analyzer := analyzer.NewAnalyzer(rules)

	analysistest.Run(t, analysistest.TestData(), analyzer, "symbols")
}

func TestSensitiveRule(t *testing.T) {
	cfg := config.DefaultConfig()
	for name := range cfg.Rules {
		if name != "sensitive" {
			cfg.Rules[name] = config.RuleConfig{
				Enabled: false,
				Options: map[string]interface{}{},
			}
		}
	}

	rules := createRulesFromConfig(t, cfg)
	analyzer := analyzer.NewAnalyzer(rules)

	analysistest.Run(t, analysistest.TestData(), analyzer, "sensitive")
}

func TestAllRules(t *testing.T) {
	cfg := config.DefaultConfig()
	rules := createRulesFromConfig(t, cfg)
	analyzer := analyzer.NewAnalyzer(rules)

	analysistest.Run(t, analysistest.TestData(), analyzer, "multi")
}

func createRulesFromConfig(t *testing.T, cfg *config.Config) []rule.Rule {
	t.Helper()

	var rules []rule.Rule
	for name, ruleCfg := range cfg.Rules {
		r, err := rule.Global.Create(name, ruleCfg)
		if err != nil {
			t.Fatalf("cannot create rule %q: %v", name, err)
		}
		rules = append(rules, r)
	}

	return rules
}
