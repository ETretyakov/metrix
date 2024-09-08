package analysers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const (
	specPackage      = "main"
	specImport       = "\"os\""
	specDeclaration  = "main"
	specFunction     = "Exit"
	specFunctionPrfx = "os"
)

var ExitCheckAnalyser = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for calling os.Exit()",
	Run:  run,
}

func run(pass *analysis.Pass) (res any, err error) {
	for _, file := range pass.Files {
		haveSpecImport := false
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.File:
				if x.Name.Name != specPackage {
					return false
				}
			case *ast.ImportSpec:
				if x.Path.Value == specImport {
					haveSpecImport = true
				}
			case *ast.FuncDecl:
				if !haveSpecImport {
					return false
				}
				if x.Name.Name != specDeclaration {
					return false
				}
				return true
			case *ast.CallExpr:
				if f, ok := x.Fun.(*ast.SelectorExpr); ok {
					if pkg, ok := f.X.(*ast.Ident); ok {
						if pkg.Name == specFunctionPrfx && f.Sel.Name == specFunction {
							pass.Reportf(pkg.NamePos, "calling os.Exit in function main")
							return false
						}
					}
				}
			}
			return true
		})
	}
	return
}
