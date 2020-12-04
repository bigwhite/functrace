package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

var (
	wrote bool
)

func init() {
	flag.BoolVar(&wrote, "w", false, "write result to (source) file instead of stdout")
}

func usage() {
	fmt.Println("gen [-w] xxx.go")
	flag.PrintDefaults()
}

func main() {
	fmt.Println(os.Args)
	flag.Usage = usage
	flag.Parse()

	if len(os.Args) < 2 {
		usage()
		return
	}

	var file string
	if len(os.Args) == 3 {
		file = os.Args[2]
	}

	if len(os.Args) == 2 {
		file = os.Args[1]
	}

	if !strings.Contains(file, ".go") {
		usage()
		return
	}

	newSrc, err := rewrite(file, nil)
	if err != nil {
		panic(err)
	}

	if !wrote {
		fmt.Println(string(newSrc))
		return
	}

	// write to the source file
	f, err := os.OpenFile(file, os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("open %s error: %v\n", file, err)
		return
	}
	defer f.Close()

	f.Seek(0, os.SEEK_SET)
	f.Truncate(0)
	f.Write(newSrc)
	fmt.Printf("add trace for %s ok\n", file)
}

func hasFuncDecl(f *ast.File) bool {
	if len(f.Decls) == 0 {
		return false
	}

	for _, decl := range f.Decls {
		_, ok := decl.(*ast.FuncDecl)
		if ok {
			return true
		}
	}

	return false
}

func rewrite(filename string, oldSource []byte) ([]byte, error) {
	fset := token.NewFileSet()
	oldAST, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", filename, err)
	}
	//fmt.Printf("%#v\n", *oldAST)

	if !hasFuncDecl(oldAST) {
		return nil, nil
	}

	// add import declaration
	astutil.AddImport(fset, oldAST, "github.com/bigwhite/functrace")
	//fmt.Printf("added=%#v\n", added)

	// inject code into each function declaration
	addDeferTraceIntoFuncDecls(oldAST)

	buf := &bytes.Buffer{}
	err = format.Node(buf, fset, oldAST)
	if err != nil {
		return nil, fmt.Errorf("error formatting new code: %w", err)
	}
	return buf.Bytes(), nil
}

func addDeferTraceIntoFuncDecls(f *ast.File) {
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if ok {
			// inject code to fd
			addDeferStmt(fd)
		}
	}
}

func addDeferStmt(fd *ast.FuncDecl) (added bool) {
	stmts := fd.Body.List

	// check whether "defer functrace.Trace()()" has already exists
	for _, stmt := range stmts {
		ds, ok := stmt.(*ast.DeferStmt)
		if !ok {
			continue
		}
		// it is a defer stmt
		ce, ok := ds.Call.Fun.(*ast.CallExpr)
		if !ok {
			continue
		}

		se, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		x, ok := se.X.(*ast.Ident)
		if !ok {
			continue
		}
		if (x.Name == "functrace") && (se.Sel.Name == "Trace") {
			// already exist , return
			return false
		}
	}

	// not found "defer functrace.Trace()()"
	// add one
	ds := &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "functrace",
					},
					Sel: &ast.Ident{
						Name: "Trace",
					},
				},
			},
		},
	}

	newList := make([]ast.Stmt, len(stmts)+1)
	copy(newList[1:], stmts)
	newList[0] = ds
	fd.Body.List = newList
	return true
}
