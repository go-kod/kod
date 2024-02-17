package internal

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type makeInterfaceFile struct {
	DirPath           string
	PkgName           string
	Structs           []string
	Struct2Interfaces map[string]string
	TypeDoc           map[string]string
	AllMethods        map[string][]string
	AllImports        []string
}

type Method struct {
	Code string
	Docs []string
}

func (m *Method) Lines() []string {
	var lines []string
	lines = append(lines, m.Docs...)
	lines = append(lines, m.Code)
	return lines
}

func getReceiverTypeName(src []byte, fl interface{}) (string, *ast.FuncDecl) {
	fd, ok := fl.(*ast.FuncDecl)
	if !ok {
		return "", nil
	}
	t, err := getReceiverType(fd)
	if err != nil {
		return "", nil
	}
	st := string(src[t.Pos()-1 : t.End()-1])
	if len(st) > 0 && st[0] == '*' {
		st = st[1:]
	}
	return st, fd
}

func getReceiverType(fd *ast.FuncDecl) (ast.Expr, error) {
	if fd.Recv == nil {
		return nil, fmt.Errorf("fd is not a method, it is a function")
	}
	return fd.Recv.List[0].Type, nil
}

func formatFieldList(src []byte, fl *ast.FieldList) []string {
	if fl == nil {
		return nil
	}
	var parts []string
	for _, l := range fl.List {
		names := make([]string, len(l.Names))
		for i, n := range l.Names {
			names[i] = n.Name
		}
		t := string(src[l.Type.Pos()-1 : l.Type.End()-1])

		if len(names) > 0 {
			typeSharingArgs := strings.Join(names, ", ")
			parts = append(parts, fmt.Sprintf("%s %s", typeSharingArgs, t))
		} else {
			parts = append(parts, t)
		}
	}
	return parts
}

type structInfo struct {
	pkgName           string
	structs           []string
	struct2Interfaces map[string]string
	methods           map[string][]Method
	imports           []string
	typeDoc           map[string]string
}

func parseStruct(src []byte) (structInfo structInfo, err error) {
	fset := token.NewFileSet()
	a, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return
	}

	structInfo.pkgName = a.Name.Name

	for _, i := range a.Imports {
		if i.Name != nil {
			structInfo.imports = append(structInfo.imports, fmt.Sprintf("%s %s", i.Name.String(), i.Path.Value))
		} else {
			structInfo.imports = append(structInfo.imports, i.Path.Value)
		}
	}

	structInfo.methods = make(map[string][]Method)
	structInfo.struct2Interfaces = make(map[string]string, 0)
	for _, d := range a.Decls {

		// 获取kod.Implements[]的interface名
		if gd, ok := d.(*ast.GenDecl); ok {
			if len(gd.Specs) > 0 {
				if ts, ok := gd.Specs[0].(*ast.TypeSpec); ok {
					if st, ok := ts.Type.(*ast.StructType); ok {
						for _, f := range st.Fields.List {
							if ie, ok := f.Type.(*ast.IndexExpr); ok {
								if se, ok := ie.X.(*ast.SelectorExpr); ok {
									if se.X.(*ast.Ident).Name == "kod" && se.Sel.Name == "Implements" {
										structName := gd.Specs[0].(*ast.TypeSpec).Name.Name
										if intfName, ok := ie.Index.(*ast.Ident); ok {
											if _, ok := structInfo.struct2Interfaces[structName]; !ok {
												structInfo.structs = append(structInfo.structs, structName)
											}
											structInfo.struct2Interfaces[structName] = intfName.Name
										}
									}
								}
							}
						}
					}
				}
			}
		}

		if structName, fd := getReceiverTypeName(src, d); structName != "" {
			// Main的启动方法，不需要生成interface
			if fd.Name.Name == "Run" {
				continue
			}

			// 私有方法
			if !fd.Name.IsExported() {
				continue
			}
			// 初始化方法
			if fd.Name.Name == "Init" {
				continue
			}
			// 退出方法
			if fd.Name.Name == "Stop" {
				continue
			}
			// Hooks
			if fd.Name.Name == "Hooks" || fd.Name.Name == "Interceptors" {
				continue
			}

			params := formatFieldList(src, fd.Type.Params)
			ret := formatFieldList(src, fd.Type.Results)
			method := fmt.Sprintf("%s(%s) (%s)", fd.Name.String(), strings.Join(params, ", "), strings.Join(ret, ", "))
			var docs []string
			if fd.Doc != nil {
				for _, d := range fd.Doc.List {
					docs = append(docs, d.Text)
				}
			}

			structInfo.methods[structName] = append(structInfo.methods[structName], Method{
				Code: method,
				Docs: docs,
			})
		}
	}

	structInfo.typeDoc = make(map[string]string)
	for _, t := range doc.New(&ast.Package{Files: map[string]*ast.File{"": a}}, "", doc.AllDecls).Types {
		structInfo.typeDoc[t.Name] = strings.TrimSuffix(t.Doc, "\n")
	}

	return
}

func makeInterfaceHead(pkgName string, imports []string) []string {
	output := []string{
		"// Code generated by kod struct2interface; DO NOT EDIT.",
		"",
		"package " + pkgName,
		"import (",
	}
	output = append(output, imports...)
	output = append(output,
		")",
		"",
	)
	return output
}

func makeInterfaceBody(output []string, ifaceComment map[string]string, structName, intfName string, methods []string) []string {

	comment := strings.TrimSuffix(strings.Replace(ifaceComment[structName], "\n", "\n//\t", -1), "\n//\t")
	if len(strings.TrimSpace(comment)) > 0 {
		output = append(output, fmt.Sprintf("// %s", comment))
	}

	output = append(output, fmt.Sprintf("type %s interface {", intfName))
	output = append(output, methods...)
	output = append(output, "}")
	output = append(output, "\n")
	return output
}

func createFile(cmd *cobra.Command, objs map[string]*makeInterfaceFile) error {
	for _, obj := range objs {
		if obj == nil {
			continue
		}

		if len(obj.Struct2Interfaces) == 0 {
			continue
		}

		var (
			pkgName          = obj.PkgName
			structAllImports = obj.AllImports
		)

		structAllImports = lo.Uniq(structAllImports)

		output := makeInterfaceHead(pkgName, structAllImports)

		for _, structName := range obj.Structs {
			output = makeInterfaceBody(output, obj.TypeDoc, structName, obj.Struct2Interfaces[structName], obj.AllMethods[structName])
		}

		code := strings.Join(output, "\n")
		result, err := ImportsCode(code)
		if err != nil {
			fmt.Printf("[struct2interface] %s \n", "formatCode error")
			return err
		}
		var fileName = filepath.Join(obj.DirPath, "kod_gen_interface.go")
		if err = os.WriteFile(fileName, result, 0644); err != nil {
			return fmt.Errorf("write file error: %s", err.Error())
		}

		if commandExists("mockgen") {
			cmd := exec.Command("mockgen", "-source", fileName, "-destination", filepath.Join(obj.DirPath, "kod_gen_mock.go"), "-package", pkgName)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			if err = cmd.Run(); err != nil {
				return fmt.Errorf("mockgen error: %s", err.Error())
			}
		}

		// fmt.Printf("[struct2interface] %s %s %s \n", "parsing", time.Since(startTime).String(), fileName)
	}

	return nil
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)

	return err == nil
}

func makeFile(file string) (*makeInterfaceFile, error) {
	var (
		allMethods = make(map[string][]string)
		allImports = make([]string, 0)
		iset       = make(map[string]struct{})
		typeDoc    = make(map[string]string)
	)

	src := lo.Must(os.ReadFile(file))

	structInfo, err := parseStruct(src)
	if err != nil {
		return nil, fmt.Errorf("parseStruct error: %s", err.Error())
	}

	for _, i := range structInfo.imports {
		if _, ok := iset[i]; !ok {
			allImports = append(allImports, i)
			iset[i] = struct{}{}
		}
	}

	for structName, methods := range structInfo.methods {
		typeDoc[structName] = fmt.Sprintf("%s is a component that implements %s.\n%s",
			structName, structInfo.struct2Interfaces[structName], structInfo.typeDoc[structName])
		for _, m := range methods {
			allMethods[structName] = append(allMethods[structName], m.Lines()...)
		}
	}

	return &makeInterfaceFile{
		DirPath:           filepath.Dir(file),
		PkgName:           structInfo.pkgName,
		Structs:           structInfo.structs,
		Struct2Interfaces: structInfo.struct2Interfaces,
		TypeDoc:           typeDoc,
		AllMethods:        allMethods,
		AllImports:        allImports,
	}, nil
}

func Struct2Interface(cmd *cobra.Command, dir string) error {
	var mapDirPath = make(map[string]*makeInterfaceFile)
	lo.Must0(filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}

		if strings.HasPrefix(filepath.Base(path), "kod_gen") {
			return nil
		}
		if !strings.HasSuffix(filepath.Base(path), ".go") {
			return nil
		}

		result, err := makeFile(path)
		if err != nil {
			return fmt.Errorf("makeFile error: %s", err.Error())
		}

		if result == nil {
			return nil
		}

		if obj, ok := mapDirPath[filepath.Dir(path)+result.PkgName]; ok {

			obj.AllImports = append(obj.AllImports, result.AllImports...)
			obj.Structs = append(obj.Structs, result.Structs...)
			for k, v := range result.Struct2Interfaces {
				obj.Struct2Interfaces[k] = v
			}

			for k, v := range result.TypeDoc {
				obj.TypeDoc[k] = v
			}
			for k, v := range result.AllMethods {
				if vv, ok := obj.AllMethods[k]; ok {
					obj.AllMethods[k] = append(vv, v...)
				} else {
					obj.AllMethods[k] = v
				}
			}
		} else {
			mapDirPath[filepath.Dir(path)+result.PkgName] = result
		}

		return nil
	}))

	return createFile(cmd, mapDirPath)
}

var struct2interface = &cobra.Command{
	Use:   "struct2interface",
	Short: "generate interface from struct for your kod application.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			args = []string{"."}
		}

		startTime := time.Now()

		for _, arg := range args {
			lo.Must0(Struct2Interface(cmd, arg))
		}

		fmt.Printf("[struct2interface] %s \n", time.Since(startTime).String())
	},
}

func init() {
	rootCmd.AddCommand(struct2interface)
}
