package plugin

import (
	"fmt"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"log-linter/analyzer"
	"log-linter/config"
	"log-linter/rule"
	_ "log-linter/rule/rules"
)

func init() {
	register.Plugin("log-linter", New)
}

type Plugin struct {
	cfg *config.Config
}

func New(settings any) (register.LinterPlugin, error) {
	p := &Plugin{
		cfg: config.DefaultConfig(),
	}

	if settings != nil {
		settingsMap, ok := settings.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("settings must be a map, got %T", settings)
		}

		if err := p.loadConfig(settingsMap); err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}

	return p, nil
}

func (p *Plugin) loadConfig(settings map[string]interface{}) error {
	rulesMap, ok := settings["rules"].(map[string]interface{})
	if !ok {
		return nil
	}

	for name, ruleSettings := range rulesMap {
		ruleCfg, ok := ruleSettings.(map[string]interface{})
		if !ok {
			continue
		}

		cfg := config.RuleConfig{
			Enabled: true,
			Options: make(map[string]interface{}),
		}

		if enabled, ok := ruleCfg["enabled"].(bool); ok {
			cfg.Enabled = enabled
		}

		if options, ok := ruleCfg["options"].(map[string]interface{}); ok {
			cfg.Options = options
		}

		p.cfg.Rules[name] = cfg
	}

	return nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	var rules []rule.Rule

	for name, ruleCfg := range p.cfg.Rules {
		if !ruleCfg.Enabled {
			continue
		}

		r, err := rule.Global.Create(name, ruleCfg)
		if err != nil {
			continue
		}
		rules = append(rules, r)
	}

	if len(rules) == 0 {
		return []*analysis.Analyzer{}, nil
	}

	analyzer := analyzer.NewAnalyzer(rules)
	return []*analysis.Analyzer{analyzer}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
