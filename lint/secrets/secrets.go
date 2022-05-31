package main

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type analyzerPlugin struct{}

// This must be implemented
func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		SecretsAnalyzer,
	}
}

// This must be defined and named 'AnalyzerPlugin'
var AnalyzerPlugin analyzerPlugin

var _regexMatch = regexp.MustCompile(`(?i)(passwd|pass|password|pwd|secret|private_key|token|bearer|credential|license)`)

var SecretsAnalyzer = &analysis.Analyzer{
	Name:     "secrets",
	Doc:      "Checks that no secret is hardcoded in code.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	// pass.ResultOf[inspect.Analyzer] will be set if we've added inspect.Analyzer to Requires.
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{ // filter needed nodes: visit only them
		(*ast.KeyValueExpr)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		keyValueExpr := node.(*ast.KeyValueExpr)

		key, ok := keyValueExpr.Key.(*ast.Ident)
		if !ok {
			return
		}

		// match key for specified words
		if !_regexMatch.MatchString(key.Name) {
			return
		}

		value, ok := keyValueExpr.Value.(*ast.CallExpr)
		if !ok {
			return
		}

		// check if value is a GetEnv function call
		valueFun, ok := value.Fun.(*ast.SelectorExpr)
		if !ok {
			// if GetEnv is not called with util.GetEnv but with GetEnv
			valueFunIdent, ok := value.Fun.(*ast.Ident)
			if !ok {
				return
			}

			if !strings.HasPrefix(valueFunIdent.Name, "GetEnv") {
				return
			}
		} else {
			if !strings.HasPrefix(valueFun.Sel.Name, "GetEnv") {
				return
			}
		}

		// assert functtion parameters length
		if len(value.Args) < 2 {
			return
		}

		// get 2nd param of GetEnv function -> defaultVaue
		defaultValue, ok := value.Args[1].(*ast.BasicLit)
		if !ok {
			return
		}

		if defaultValue.Kind != token.STRING {
			return
		}

		// value is '""' therefore check for <= 2
		if len(defaultValue.Value) <= 2 {
			return
		}

		envKey := ""
		envKeyBasicLit, ok := value.Args[0].(*ast.BasicLit)
		if ok {
			envKey = envKeyBasicLit.Value
		}

		pass.Reportf(node.Pos(), "default value for secret '%s' (%s) should be empty", key.Name, envKey)
	})

	return nil, nil
}
