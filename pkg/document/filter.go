package document

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"

	"sigs.k8s.io/kustomize/v3/pkg/resid"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
)

// Visitor object for traversing through Abstract Syntax Tree
type Visitor struct {
	Bundle Bundle
	err    error
}

// Visit mathod is executed by ast.Walk on each node of Abstract Syntax Tree
func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	var be *ast.BinaryExpr
	switch nodeT := node.(type) {
	case *ast.BinaryExpr:
		be = nodeT
	case *ast.ParenExpr:
		return v
	default:
		return nil
	}
	switch be.Op {
	case token.EQL:
		v.Bundle, v.err = filterByFieldVal(be, v.Bundle, reflect.DeepEqual)
	case token.NEQ:
		funcNotEq := func(x, y interface{}) bool { return !reflect.DeepEqual(x, y) }
		v.Bundle, v.err = filterByFieldVal(be, v.Bundle, funcNotEq)
	case token.LOR:
		v.Bundle, v.err = filterByLogicalExpr(be, v.Bundle, union)
	case token.LAND:
		v.Bundle, v.err = filterByLogicalExpr(be, v.Bundle, intersection)
	}
	return nil
}

func intersection(left, right Bundle) (Bundle, error) {
	result := &BundleFactory{
		KustomizeBuildOptions: left.GetKustomizeBuildOptions(),
		FileSystem:            left.GetFileSystem(),
	}

	resourceHash := make(map[resid.ResId]bool)
	for _, resID := range left.GetKustomizeResourceMap().AllIds() {
		resourceHash[resID] = true
	}

	resourceMap := resmap.New()
	for _, res := range right.GetKustomizeResourceMap().Resources() {
		if _, found := resourceHash[res.CurId()]; found {
			if err := resourceMap.Append(res); err != nil {
				return nil, err
			}
		}
	}
	if err := result.SetKustomizeResourceMap(resourceMap); err != nil {
		return nil, err
	}
	return result, nil
}

func union(left, right Bundle) (Bundle, error) {
	result := &BundleFactory{
		KustomizeBuildOptions: left.GetKustomizeBuildOptions(),
		FileSystem:            left.GetFileSystem(),
		ResMap:                left.GetKustomizeResourceMap(),
	}
	for _, res := range right.GetKustomizeResourceMap().Resources() {
		// Ignore errors on duplicated resources
		_ = result.ResMap.Append(res) //nolint:errcheck
	}
	return result, nil
}

func filterByLogicalExpr(
	expr *ast.BinaryExpr,
	bundle Bundle,
	joinF func(Bundle, Bundle) (Bundle, error),
) (Bundle, error) {
	lv, rv := &Visitor{Bundle: bundle}, &Visitor{Bundle: bundle}
	ast.Walk(lv, expr.X)
	if lv.err != nil {
		return nil, lv.err
	}
	ast.Walk(rv, expr.Y)
	if rv.err != nil {
		return nil, rv.err
	}
	return joinF(lv.Bundle, rv.Bundle)
}

func makeJSONPath(expr ast.Expr) string {
	switch exprType := expr.(type) {
	case *ast.SelectorExpr:
		switch xType := exprType.X.(type) {
		case *ast.Ident:
			return xType.String() + "." + exprType.Sel.String()
		case *ast.SelectorExpr, *ast.IndexExpr:
			return makeJSONPath(xType) + "." + exprType.Sel.String()
		}
	case *ast.IndexExpr:
		idx, ok := exprType.Index.(*ast.BasicLit)
		if !ok {
			return ""
		}
		idxVal := "[" + strings.Trim(idx.Value, "\"") + "]"
		switch xType := exprType.X.(type) {
		case *ast.Ident:
			return xType.String() + idxVal
		case *ast.SelectorExpr, *ast.IndexExpr:
			return makeJSONPath(xType) + idxVal
		}
	}
	return ""
}

func filterByFieldVal(
	expr *ast.BinaryExpr,
	bundle Bundle,
	filterFunc func(a, b interface{}) bool,
) (Bundle, error) {
	var jPath, value string
	for _, operand := range [...]ast.Expr{expr.X, expr.Y} {
		switch opType := operand.(type) {
		case *ast.SelectorExpr, *ast.IndexExpr:
			jPath = makeJSONPath(opType)
		case *ast.Ident:
			jPath = opType.String()
		case *ast.BasicLit:
			value = strings.Trim(opType.Value, "\"")
		default:
			return nil, fmt.Errorf("Faled to convert operands %#v and %#v", expr.X, expr.Y)
		}
	}
	return bundle.GetByFieldValue(jPath, func(x interface{}) bool { return filterFunc(x, value) })
}

// EvaluateExpressionFilter walks through syntax tree and filters document
// bundle according to node expressions
func EvaluateExpressionFilter(expr string, bundle Bundle) (Bundle, error) {
	astExpr, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, err
	}
	v := &Visitor{Bundle: bundle}
	ast.Walk(v, astExpr)
	if v.err != nil {
		return nil, v.err
	}
	return v.Bundle, nil
}
