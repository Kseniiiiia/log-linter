package main

import (
	"flag"
	"log"

	"golang.org/x/tools/go/analysis/singlechecker"

	"log-linter/analyzer"
	"log-linter/config"
	"log-linter/rule"
	_ "log-linter/rule/rules"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config file (optional)")
	flag.Parse()

	cfg := config.Load(configPath)

	var rules []rule.Rule
	for name, ruleCfg := range cfg.Rules {
		if !ruleCfg.Enabled {
			continue
		}
		r, err := rule.Global.Create(name, ruleCfg)
		if err != nil {
			log.Printf("cannot create rule %q %v", name, err)
			continue
		}
		rules = append(rules, r)
	}

	analyzer := analyzer.NewAnalyzer(rules)
	singlechecker.Main(analyzer)
}
