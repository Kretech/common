package p

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"runtime"
	"strings"
)

// VarName 用来获取变量的名字
// VarName(a, b) => []string{"a", "b"}
func VarName(args ...interface{}) []string {
	return varNameDepth(1, args...)
}

func varNameDepth(skip int, args ...interface{}) (c []string) {
	pc, _, _, _ := runtime.Caller(skip)
	callFunc := runtime.FuncForPC(pc)
	ss := strings.Split(callFunc.Name(), "/")

	// 用户通过这个方法来获取变量名。
	// 可能有几种写法：p.F() alias.F() .F()，我们需要解析 import 来确定
	shouldCallName := ss[len(ss)-1]
	shouldCallPkg := callFunc.Name()[:strings.LastIndex(callFunc.Name(), `.`)]

	_, file, line, _ := runtime.Caller(skip + 1)

	// todo 一行多次调用时，还需根据 runtime 找到 column 一起定位
	cacheKey := fmt.Sprintf("%s:%d@%s", file, line, shouldCallName)
	return cacheGet(cacheKey, func() interface{} {

		r := []string{}
		found := false

		fset := token.NewFileSet()
		f, _ := parser.ParseFile(fset, file, nil, 0)

		// import alias
		aliasImport := make(map[string]string)
		for _, decl := range f.Decls {
			decl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range decl.Specs {
				is, ok := spec.(*ast.ImportSpec)
				if !ok {
					continue
				}

				if is.Name != nil && strings.Trim(is.Path.Value, `""`) == shouldCallPkg {
					aliasImport[is.Name.Name] = shouldCallPkg
					shouldCallName = is.Name.Name + "." + strings.Split(shouldCallName, ".")[1]

					shouldCallName = strings.TrimLeft(shouldCallName, `.`)
				}
			}
		}

		// q.Q(shouldCallName, shouldCallPkg, aliasImport)

		// q.Q(f)
		// q.Q(f.Decls[1].(*ast.FuncDecl).Body.List[1].(*ast.ExprStmt).X.(*ast.CallExpr).Args[0].(*ast.CallExpr).Args[0].(*ast.Ident).Obj)
		ast.Inspect(f, func(node ast.Node) bool {
			if found {
				return false
			}

			if node == nil {
				return false
			}

			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			// q.Q(call)
			isArgsCall := func(expr *ast.CallExpr, shouldCallName string) bool {
				if strings.Contains(shouldCallName, ".") {
					fn, ok := call.Fun.(*ast.SelectorExpr)
					if !ok {
						return false
					}

					currentName := fn.X.(*ast.Ident).Name + "." + fn.Sel.Name

					// q.Q(shouldCallName, currentName)
					return shouldCallName == currentName
				} else {
					fn, ok := call.Fun.(*ast.Ident)
					if !ok {
						return false
					}

					return fn.Name == shouldCallName
				}

				return false
			}

			if !isArgsCall(call, shouldCallName) {
				return true
			}

			if fset.Position(call.End()).Line != line {
				return true
			}

			for _, arg := range call.Args {
				r = append(r, arg.(*ast.Ident).Name)
			}

			found = true
			return false
		})

		return r
	}).([]string)
}

// Compact 将多个变量打包到一个字典里
// a,b:=1,2 Comapct(a, b) => {"a":1,"b":2}
// 参考自 http://php.net/manual/zh/function.compact.php
func Compact(args ...interface{}) (r map[string]interface{}) {
	return DepthCompact(1, args...)
}

func DepthCompact(depth int, args ...interface{}) (r map[string]interface{}) {
	ps := varNameDepth(depth+1, args...)

	r = make(map[string]interface{}, len(ps))
	for idx, param := range ps {
		r[param] = args[idx]
	}

	return
}

var m = newRWMap()

func cacheGet(key string, backup func() interface{}) interface{} {

	v := m.Get(key)

	if v == nil {
		v = backup()
		m.Set(key, v)
	}

	return v
}