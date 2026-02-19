package analyzer

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"strings"

	"log-linter/rule"
)

func NewAnalyzer(rules []rule.Rule) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "loglint",
		Doc:  "Check log-lint rules",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			runAnalyzer(pass, rules)
			return nil, nil
		},
	}
}

func runAnalyzer(pass *analysis.Pass, rules []rule.Rule) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			if !isLogFunction(pass, callExpr) {
				return true
			}

			msg, pos, _ := extractLogMessage(callExpr)
			if msg == "" {
				return true
			}

			for _, r := range rules {
				diags := r.Check(msg, pos)
				for _, diag := range diags {
					pass.Report(diag)
				}
			}
			return true
		})
	}
}

func isLogFunction(pass *analysis.Pass, call *ast.CallExpr) bool {
	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		var pkgName string
		if pkgIdent, ok := fun.X.(*ast.Ident); ok {
			pkgName = pkgIdent.Name
		}

		switch pkgName {
		case "log", "slog":
			return true
		case "zap":
			if obj := pass.TypesInfo.ObjectOf(fun.Sel); obj != nil {
				if strings.Contains(obj.Pkg().Path(), "zap") {
					return true
				}
			}
			return true
		}

		switch fun.Sel.Name {
		case "Info", "Infof", "Infoln",
			"Error", "Errorf", "Errorln",
			"Debug", "Debugf", "Debugln",
			"Warn", "Warnf", "Warnln",
			"Print", "Printf", "Println",
			"Fatal", "Fatalf", "Fatalln",
			"Panic", "Panicf", "Panicln",
			"DPanic", "DPanicf", "DPanicln":
			if obj := pass.TypesInfo.ObjectOf(fun.Sel); obj != nil {
				pkgPath := obj.Pkg().Path()
				if strings.Contains(pkgPath, "log") ||
					strings.Contains(pkgPath, "slog") ||
					strings.Contains(pkgPath, "zap") {
					return true
				}
			}
		}
	}
	return false
}

func extractLogMessage(call *ast.CallExpr) (string, token.Pos, ast.Node) {
	if len(call.Args) == 0 {
		return "", token.NoPos, nil
	}

	firstArg := call.Args[0]

	if lit, ok := firstArg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		msg := strings.Trim(lit.Value, "\"`")
		return msg, lit.Pos(), lit
	}

	if binary, ok := firstArg.(*ast.BinaryExpr); ok && binary.Op == token.ADD {
		if lit, ok := binary.X.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			msg := strings.Trim(lit.Value, "\"`")
			return msg, lit.Pos(), lit
		}
	}

	return "", token.NoPos, nil
}
