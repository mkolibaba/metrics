// Package osexitusage реализует анализатор, запрещающий
// использовать прямой вызов os.Exit в функции main пакета main.
package osexitusage

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
)

const (
	mainPackageName  = "main"
	mainFunctionName = "main"
	osPackageName    = "os"
	exitFunctionName = "Exit"
	usageWarning     = "os.Exit is used in main"
)

var Analyzer = &analysis.Analyzer{
	Name: "osexitusage",
	Doc:  "checks whether os.Exit is used incorrectly",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name != mainPackageName {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				if x.Name.Name != mainFunctionName {
					return false
				}
			case *ast.CallExpr:
				if isOSExitInMain(x, pass.TypesInfo) {
					pass.Reportf(x.Pos(), usageWarning)
				}
			}
			return true
		})
	}
	return nil, nil
}

func isOSExitInMain(x *ast.CallExpr, info *types.Info) bool {
	switch f := x.Fun.(type) {
	case *ast.SelectorExpr:
		if ident, ok := f.X.(*ast.Ident); ok {
			if ident.Name == osPackageName && f.Sel.Name == exitFunctionName {
				return true
			}
		}
	case *ast.Ident:
		if f.Name == exitFunctionName {
			if obj, ok := info.Uses[f]; ok && obj.Pkg().Name() == osPackageName {
				return true
			}
		}
	}

	return false
}
