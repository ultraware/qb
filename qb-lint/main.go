package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path"
	"strings"

	"golang.org/x/tools/go/loader"
)

const aliasLimit = 6

var fset *token.FileSet
var info = &types.Info{
	Types: make(map[ast.Expr]types.TypeAndValue),
}

func fatal(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(2)
}

func main() {
	fixArgs()
	okay := true

	var conf loader.Config
	conf.ImportWithTests(os.Args[1])

	prog, err := conf.Load()
	if err != nil {
		fatal(err)
	}

	pkg := prog.Package(os.Args[1])

	fset = prog.Fset
	info = &pkg.Info

	for _, v := range pkg.Files {
		parseFile(v)
	}

	if !okay {
		os.Exit(1)
	}
}

func fixArgs() {
	var target string
	if len(os.Args) == 1 || os.Args[1] == `` {
		target = `.`
	} else {
		target = os.Args[1]
	}

	if target[0] == '.' {
		gopath := ``
		cwd, _ := os.Getwd()

		target = path.Clean(path.Join(cwd, target))

		for _, v := range strings.Split(os.Getenv(`GOPATH`), `:`) {
			path := path.Join(v, `src`)
			if strings.HasPrefix(cwd, path) {
				gopath = path
			}
		}

		target = strings.TrimPrefix(strings.TrimPrefix(target, gopath), `/`)
	}
	os.Args = append(os.Args[:1], target)
}

type field struct {
	table  ast.Ident
	column ast.Ident
}

func (f field) String() string {
	return f.table.Name + `.` + f.column.Name
}

func (f field) PositionString() string {
	pos := fset.Position(f.table.Pos())
	return pos.String()
}

var usedFields map[field]bool

func parseFile(file *ast.File) bool {
	okay := true

	for _, v := range file.Decls {
		f, ok := v.(*ast.FuncDecl)
		if !ok {
			continue
		}

		usedFields = make(map[field]bool)
		translate = make(map[interface{}]*ast.Ident)

		if f.Name.Obj == nil {
			continue
		}
		initialFunc = f.Name.Obj.Decl
		currentFunc = initialFunc

		ast.Inspect(f, findTables)
		ast.Inspect(f, findSelected)
		ast.Inspect(f, findUsed)

		for k, v := range usedFields {
			if v {
				continue
			}
			fmt.Println(k.PositionString()+`:Field`, k, `retrieved but never used`)
			okay = false
		}
	}

	return okay
}

var initialFunc interface{}
var currentFunc interface{}

var translate map[interface{}]*ast.Ident

func tryFunc(n ast.Node, f func(ast.Node) bool) {
	if v, ok := n.(*ast.CallExpr); ok {
		i, ok := v.Fun.(*ast.Ident)
		if !ok || i.Obj == nil || i.Obj.Decl == nil {
			return
		}
		if i.Obj.Decl == initialFunc {
			return
		}
		if fmt.Sprint(i.Obj.Kind) == `func` {
			param := 0
			for _, p := range i.Obj.Decl.(*ast.FuncDecl).Type.Params.List {
				for _, i := range p.Names {
					if _, ok := v.Args[param].(*ast.Ident); !ok {
						continue
					}
					translate[i.Obj.Decl] = v.Args[param].(*ast.Ident)
					param++
				}
			}

			oldFunc := currentFunc
			currentFunc = i.Obj.Decl
			ast.Inspect(i.Obj.Decl.(*ast.FuncDecl), f)

			currentFunc = oldFunc
		}
	}
}

func findTables(n ast.Node) bool {
	assign, ok := n.(*ast.AssignStmt)
	if !ok {
		return true
	}
	if len(assign.Lhs) != len(assign.Rhs) {
		return false
	}
	for k := range assign.Lhs {
		if t, ok := tryQbStruct(assign.Rhs[k]); !ok || !checkQbStruct(t) {
			continue
		}
		rhs := assign.Lhs[k].(*ast.Ident)
		if len(rhs.Name) > aliasLimit {
			fmt.Printf("%v:Found long alias '%s' (length %d). An alias should be short (<= %d)\n",
				fset.Position(rhs.Pos()), rhs.Name, len(rhs.Name), aliasLimit)
		}
	}
	return false
}

func findSelected(n ast.Node) bool {
	if expr, ok := n.(ast.Expr); ok {
		return traverseExpr(nil, expr)
	}
	return true
}

var skip int

func findUsed(n ast.Node) bool {
	if len(usedFields) > 0 {
		tryFunc(n, findUsed)
	}

	if skip > 0 {
		skip--
		return true
	}
	if _, ok := n.(*ast.AssignStmt); ok {
		skip = 1
		return true
	}

	return tryField(n)
}

func traverseExpr(parent ast.Expr, e ast.Expr) bool {
	if expr, ok := e.(*ast.SelectorExpr); ok {
		switch fmt.Sprint(expr.Sel) {
		case `SubQuery`:
			break
		case `Returning`:
			checkFields(parent, 1)
		case `Select`:
			checkFields(parent, 0)
		default:
			return traverseExpr(expr, expr.X)
		}
		return false
	}

	call, ok := e.(*ast.CallExpr)
	if !ok {
		return true
	}

	return traverseExpr(e, call.Fun)
}

func checkFields(e ast.Expr, skip int) {
	call := e.(*ast.CallExpr)
	args := call.Args[skip:]

argsLoop:
	for k := range args {
		v, ok := args[k].(*ast.SelectorExpr)
		if !ok {
			continue
		}
		for _, val := range args[:k] {
			if prev, ok := val.(*ast.SelectorExpr); ok && compareFields(v, prev) {
				fmt.Printf("%v:Field %s.%s selected multiple times.\n",
					fset.Position(v.Pos()), v.X.(*ast.Ident).Name, v.Sel.Name)
				continue argsLoop
			}
		}
		field := getField(v)
		if field == nil {
			continue
		}

		usedFields[*field] = false
	}
}

func getField(e ast.Expr) *field {
	if expr, ok := e.(*ast.SelectorExpr); ok {
		if t, ok := tryQbStruct(expr.X); !ok || !checkQbStruct(t) {
			// False positive. The found type is not a type from qb-generator
			return nil
		}

		return &field{
			table:  *expr.X.(*ast.Ident),
			column: *expr.Sel,
		}
	}
	return nil
}

func tryQbStruct(e ast.Expr) (*types.Struct, bool) {
	t := info.TypeOf(e)
	ptr, ok := t.(*types.Pointer)
	if !ok {
		return nil, false
	}
	t = ptr.Elem().Underlying()
	s, ok := t.(*types.Struct)
	return s, ok
}

func checkQbStruct(t *types.Struct) bool {
	if t.Field(0).Name() != `Data` {
		return false
	}
	if t.Field(t.NumFields()-1).Name() != `table` {
		return false
	}
	return true
}

func tryField(n ast.Node) bool {
	expr, ok := n.(*ast.SelectorExpr)
	if !ok {
		return true
	}

	sel, ok := expr.X.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != `Data` {
		return true
	}

	x := sel.X.(*ast.Ident)
	if _, ok := x.Obj.Decl.(*ast.Field); ok {
		new, ok := translate[x.Obj.Decl]
		if !ok {
			return false
		}
		x = new
	}

	if !retrieved(expr, x) && len(usedFields) > 0 {
		if _, ok := translate[sel.X.(*ast.Ident).Obj.Decl]; currentFunc == initialFunc || ok {
			fmt.Printf("%v:Field %s used, but it was never selected\n", fset.Position(sel.X.Pos()),
				fmt.Sprint(sel.X, `.Data.`, expr.Sel))
		}
	}
	return false
}

func retrieved(expr *ast.SelectorExpr, x *ast.Ident) bool {
	var target *field
	for k := range usedFields {
		if x.Obj != k.table.Obj || // table object doesn't match
			expr.Sel.Name != k.column.Name || // field name doesn't match
			expr.X.Pos() < k.table.Pos() { // The field was used before it was selected
			continue
		}

		if target == nil || target.table.Pos() < k.table.Pos() {
			t := k
			target = &t
		}
	}
	if target != nil {
		usedFields[*target] = true
		return true
	}
	return false
}

func compareFields(f1 *ast.SelectorExpr, f2 *ast.SelectorExpr) bool {
	x1, ok := f1.X.(*ast.Ident)
	if !ok {
		return false
	}
	x2, ok := f2.X.(*ast.Ident)
	if !ok {
		return false
	}

	return x1.Obj == x2.Obj && f1.Sel.Name == f2.Sel.Name
}
