package internal

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-kod/kod/internal/callgraph"
	"github.com/samber/lo"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

// TODO(rgrandl): Modify the generator code to use only the types package. Right
// now we are doing code generation relying both on the types and ast packages,
// which can be confusing and also we might do unnecessary work.

const (
	generatedCodeFile = "kod_gen.go"

	Usage = `Generate code for a Kod application.

Usage:
  kod generate [packages]

Description:
  "kod generate" generates code for the Kod applications in the
  provided packages. For example, "kod generate . ./foo" will generate code
  for the Kod applications in the current directory and in the ./foo
  directory. For every package, the generated code is placed in a kod_gen.go
  file in the package's directory. For example, "kod generate . ./foo" will
  create ./kod_gen.go and ./foo/kod_gen.go.

  You specify packages for "kod generate" in the same way you specify
  packages for go build, go test, go vet, etc. See "go help packages" for more
  information.

  Rather than invoking "kod generate" directly, you can place a line of the
  following form in one of the .go files in the package:

      //go:generate kod generate

  and then use the normal "go generate" command.

Examples:
  # Generate code for the package in the current directory.
  kod generate

  # Same as "kod generate".
  kod generate .

  # Generate code for the package in the ./foo directory.
  kod generate ./foo

  # Generate code for all packages in all subdirectories of current directory.
  kod generate ./...`
)

// Options controls the operation of Generate.
type Options struct {
	// If non-nil, use the specified function to report warnings.
	Warn func(error)
}

// Generate generates Kod code for the specified packages.
// The list of supplied packages are treated similarly to the arguments
// passed to "go build" (see "go help packages" for details).
func Generate(dir string, pkgs []string, opt Options) error {
	if opt.Warn == nil {
		opt.Warn = func(err error) { fmt.Fprintln(os.Stderr, err) }
	}
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedSyntax | packages.NeedImports | packages.NeedTypes | packages.NeedTypesInfo,
		Dir:        dir,
		Fset:       fset,
		ParseFile:  parseNonKodGenFile,
		BuildFlags: []string{"--tags=ignoreKodGen"},
	}
	pkgList, err := packages.Load(cfg, pkgs...)
	if err != nil {
		return fmt.Errorf("packages.Load: %w", err)
	}

	var automarshals typeutil.Map
	var errs []error
	for _, pkg := range pkgList {
		g, err := newGenerator(opt, pkg, fset, &automarshals)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if err := g.generate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// parseNonKodGenFile parses a Go file, except for kod_gen.go files whose
// contents are ignored since those contents may reference types that no longer
// exist.
func parseNonKodGenFile(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
	if filepath.Base(filename) == generatedCodeFile {
		return parser.ParseFile(fset, filename, src, parser.PackageClauseOnly)
	}
	return parser.ParseFile(fset, filename, src, parser.ParseComments|parser.DeclarationErrors)
}

type generator struct {
	opt        Options
	pkg        *packages.Package
	tset       *typeSet
	fileset    *token.FileSet
	components []*component
}

// errorf is like fmt.Errorf but prefixes the error with the provided position.
func errorf(fset *token.FileSet, pos token.Pos, format string, args ...interface{}) error {
	// Rewrite the position's filename relative to the current directory. This
	// replaces long filenames like "/home/foo/go-kod/kod/kod.go"
	// with much shorter filenames like "./kod.go".
	position := fset.Position(pos)
	if cwd, err := filepath.Abs("."); err == nil {
		if filename, err := filepath.Rel(cwd, position.Filename); err == nil {
			position.Filename = filename
		}
	}

	prefix := position.String()
	// if colors.Enabled() {
	// 	// Color the filename red when colors are enabled.
	// 	prefix = fmt.Sprintf("%s%v%s", colors.Color256(160), position, colors.Reset)
	// }
	return fmt.Errorf("%s: %w", prefix, fmt.Errorf(format, args...))
}

func newGenerator(opt Options, pkg *packages.Package, fset *token.FileSet, automarshals *typeutil.Map) (*generator, error) {
	// Abort if there were any errors loading the package.
	var errs []error
	for _, err := range pkg.Errors {
		errs = append(errs, err)
	}
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	// Search every file in the package for types that embed the
	// kod.AutoMarshal struct.
	tset := newTypeSet(pkg, automarshals, &typeutil.Map{})

	// Find and process all components.
	components := map[string]*component{}
	for _, file := range pkg.Syntax {
		filename := fset.Position(file.Package).Filename
		if filepath.Base(filename) == generatedCodeFile {
			// Ignore kod_gen.go files.
			continue
		}

		fileComponents, err := findComponents(opt, pkg, file, tset)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		for _, c := range fileComponents {
			// Check for component duplicates, two components that embed the
			// same kod.Implements[T].
			//
			// TODO(mwhittaker): This code relies on the fact that a component
			// interface and component implementation have to be in the same
			// package. If we lift this requirement, then this code will break.
			if existing, ok := components[c.fullIntfName()]; ok {
				errs = append(errs, errorf(pkg.Fset, c.impl.Obj().Pos(),
					"Duplicate implementation for component %s, other declaration: %v",
					c.fullIntfName(), fset.Position(existing.impl.Obj().Pos())))
				continue
			}
			components[c.fullIntfName()] = c
		}
	}
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &generator{
		opt:        opt,
		pkg:        pkg,
		tset:       tset,
		fileset:    fset,
		components: lo.Values(components),
	}, nil
}

// findComponents returns the components in the provided file. For example,
// findComponents will find and return the following component.
//
//	type something struct {
//	    kod.Implements[SomeComponentType]
//	    ...
//	}
func findComponents(opt Options, pkg *packages.Package, f *ast.File, tset *typeSet) ([]*component, error) {
	var components []*component
	var errs []error
	for _, d := range f.Decls {
		gendecl, ok := d.(*ast.GenDecl)
		if !ok || gendecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range gendecl.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			component, err := extractComponent(opt, pkg, f, tset, ts)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if component != nil {
				components = append(components, component)
			}
		}
	}
	return components, errors.Join(errs...)
}

// extractComponent attempts to extract a component from the provided TypeSpec.
// It returns a nil component if the TypeSpec doesn't define a component.
func extractComponent(opt Options, pkg *packages.Package, file *ast.File, tset *typeSet, spec *ast.TypeSpec) (*component, error) {
	// Check that the type spec is of the form `type t struct {...}`.
	s, ok := spec.Type.(*ast.StructType)
	if !ok {
		// This type declaration does not involve a struct. For example, it
		// might look like `type t int`. These non-struct type declarations
		// cannot be components.
		return nil, nil
	}
	def, ok := pkg.TypesInfo.Defs[spec.Name]
	if !ok {
		panic(errorf(pkg.Fset, spec.Pos(), "name %v not found", spec.Name))
	}
	impl, ok := def.Type().(*types.Named)
	if !ok {
		// For type aliases like `type t = struct{}`, t has type *types.Struct
		// and not type *types.Named. We ignore these.
		return nil, nil
	}

	// Find any kod.Implements[T] or kod.WithRouter[T] embedded fields.
	var intf *types.Named // The component interface type
	var isMain bool       // Is intf kod.Main?

	// var typ componentType   // ComponentType
	var refs []*types.Named // T for which kod.Ref[T] exists in struct
	for _, f := range s.Fields.List {
		typeAndValue, ok := pkg.TypesInfo.Types[f.Type]
		if !ok {
			panic(errorf(pkg.Fset, f.Pos(), "type %v not found", f.Type))
		}
		t := typeAndValue.Type

		if isKodRef(t) {
			// The field f has type kod.Ref[T].
			arg := t.(*types.Named).TypeArgs().At(0)
			if isKodMain(arg) {
				return nil, errorf(pkg.Fset, f.Pos(),
					"components cannot contain a reference to kod.Main")
			}
			named, ok := arg.(*types.Named)
			if !ok {
				return nil, errorf(pkg.Fset, f.Pos(),
					"kod.Ref argument %s is not a named type.",
					formatType(pkg, arg))
			}
			refs = append(refs, named)
		}

		if len(f.Names) != 0 {
			// Ignore unembedded fields.
			//
			// TODO(mwhittaker): Warn the user about unembedded
			// kod.Implements, kod.WithConfig, or kod.WithRouter?
			continue
		}

		switch {
		// The field f is an embedded kod.Implements[T].
		case isKodImplements(t):
			// Check that T is a named interface type inside the package.
			arg := t.(*types.Named).TypeArgs().At(0)
			named, ok := arg.(*types.Named)
			if !ok {
				return nil, errorf(pkg.Fset, f.Pos(),
					"kod.Implements argument %s is not a named type.",
					formatType(pkg, arg))
			}
			// typ = getComponentType(arg)
			isMain = isKodMain(arg)
			if !isMain && named.Obj().Pkg() != pkg.Types {
				return nil, errorf(pkg.Fset, f.Pos(),
					"kod.Implements argument %s is a type outside the current package. A component interface and implementation must be in the same package. If you can't move them into the same package, you can add `type %s %v` to the implementation's package and embed `kod.Implements[%s]` instead of `kod.Implements[%s]`.",
					formatType(pkg, named), named.Obj().Name(), formatType(pkg, named), named.Obj().Name(), formatType(pkg, named))
			}
			if _, ok := named.Underlying().(*types.Interface); !ok {
				return nil, errorf(pkg.Fset, f.Pos(),
					"kod.Implements argument %s is not an interface.",
					formatType(pkg, named))
			}
			intf = named
		}
	}

	if intf == nil {
		// TODO(mwhittaker): Warn the user if they embed kod.WithRouter or
		// kod.WithConfig but don't embed kod.Implements.
		return nil, nil
	}

	// Check that that the component implementation implements the component
	// interface.
	if !types.Implements(types.NewPointer(impl), intf.Underlying().(*types.Interface)) {
		return nil, errorf(pkg.Fset, spec.Pos(),
			"type %s embeds kod.Implements[%s] but does not implement interface %s.",
			formatType(pkg, impl), formatType(pkg, intf), formatType(pkg, intf))
	}

	// Disallow generic component implementations.
	if spec.TypeParams != nil && spec.TypeParams.NumFields() != 0 {
		return nil, errorf(pkg.Fset, spec.Pos(),
			"component implementation %s is generic. Component implements cannot be generic.",
			formatType(pkg, impl))
	}

	// Validate the component's methods.
	if err := validateMethods(pkg, tset, intf); err != nil {
		return nil, err
	}

	// Warn the user if the component has a mistyped Init method. Init methods
	// are supposed to have type "func(context.Context) error", but it's easy
	// to forget to add a context.Context argument or error return. Without
	// this warning, the component's Init method will be silently ignored. This
	// can be very frustrating to debug.
	if err := checkMistypedInit(pkg, tset, impl); err != nil {
		opt.Warn(err)
	}

	comp := &component{
		intf:   intf,
		impl:   impl,
		isMain: isMain,
		refs:   refs,
	}

	return comp, nil
}

// component represents a Kod component.
//
// A component is divided into an interface and implementation. For example, in
// the following code, Adder is the component interface, and adder is the
// component implementation. router is the router type.
//
//	type Adder interface{}
//	type adder struct {
//	    kod.Implements[Adder]
//	    kod.WithRouter[router]
//	}
//	type router struct{}
type component struct {
	intf   *types.Named   // component interface
	impl   *types.Named   // component implementation
	isMain bool           // intf is kod.Main
	refs   []*types.Named // List of T where a kod.Ref[T] field is in impl struct
}

func fullName(t *types.Named) string {
	return path.Join(t.Obj().Pkg().Path(), t.Obj().Name())
}

// intfName returns the component interface name.
func (c *component) intfName() string {
	return c.intf.Obj().Name()
}

// implName returns the component implementation name.
func (c *component) implName() string {
	return c.impl.Obj().Name()
}

// fullIntfName returns the full package-prefixed component interface name.
func (c *component) fullIntfName() string {
	return fullName(c.intf)
}

// methods returns the component interface's methods.
func (c *component) methods() []*types.Func {
	underlying := c.intf.Underlying().(*types.Interface)
	methods := make([]*types.Func, underlying.NumMethods())
	for i := 0; i < underlying.NumMethods(); i++ {
		methods[i] = underlying.Method(i)
	}

	// Sort the component's methods deterministically. This allows a developer
	// to re-order the interface methods without the generated code changing.
	sort.Slice(methods, func(i, j int) bool {
		return methods[i].Name() < methods[j].Name()
	})
	return methods
}

// validateMethods validates that the provided component's methods are all
// valid component methods.
func validateMethods(pkg *packages.Package, tset *typeSet, intf *types.Named) error {
	var errs []error
	underlying := intf.Underlying().(*types.Interface)
	for i := 0; i < underlying.NumMethods(); i++ {
		m := underlying.Method(i)
		t, ok := m.Type().(*types.Signature)
		if !ok {
			panic(errorf(pkg.Fset, m.Pos(), "method %s doesn't have a signature", m.Name()))
		}

		// Disallow unexported methods.
		if !m.Exported() {
			errs = append(errs, errorf(pkg.Fset, m.Pos(),
				"Method `%s%s %s` of Kod component %q is unexported. Every method in a component interface must be exported.",
				m.Name(), formatType(pkg, t.Params()), formatType(pkg, t.Results()), intf.Obj().Name()))
			continue
		}

		// // bad is a helper function for producing helpful error messages.
		// bad := func(bad, format string, arg ...any) error {
		// 	err := fmt.Errorf(format, arg...)
		// 	return errorf(
		// 		pkg.Fset, m.Pos(),
		// 		"Method `%s%s %s` of Kod component %q has incorrect %s types. %w",
		// 		m.Name(), formatType(pkg, t.Params()), formatType(pkg, t.Results()),
		// 		intf.Obj().Name(), bad, err)
		// }

		// do no process kod.Main
		if isKodMain(intf) {
			continue
		}

		// // do no process kod.Controller or kod.Component
		// if typ == ComponentTypeController || typ == ComponentTypeComponent {
		// 	continue
		// }

		// // Method must have one parameter at least.
		// if t.Params().Len() == 0 {
		// 	errs = append(errs, bad("argument", "The method must have one parameter at least."))
		// 	continue
		// }

		// // First argument must be context.Context.
		// if !isContext(t.Params().At(0).Type()) {
		// 	errs = append(errs, bad("argument", "The first argument must have type context.Context."))
		// }

		// // do no process kod.Controller or kod.Component
		// if typ == ComponentTypeRepository {
		// 	continue
		// }

		// // Method must have two parameters at most.
		// if t.Params().Len() > 2 {
		// 	errs = append(errs, bad("argument", "The method must have two parameters at most."))
		// 	continue
		// }

		// // Second argument must be a pointer to a struct
		// if t.Params().Len() == 2 && !isPointerToStruct(t.Params().At(1).Type()) {
		// 	errs = append(errs, bad("argument", "The second argument must be a pointer to a struct."))
		// }

		// // Result must have two values at most.
		// if t.Results().Len() > 2 {
		// 	errs = append(errs, bad("return", "The method must return two values at most."))
		// 	continue
		// }

		// // First result must be a pointer to a struct.
		// if t.Results().Len() == 2 && !isPointerToStruct(t.Results().At(0).Type()) {
		// 	errs = append(errs, bad("return", "The first return must be a pointer to a struct."))
		// }

		// // Last result must be error.
		// if t.Results().Len() == 1 && t.Results().At(t.Results().Len()-1).Type().String() != "error" {
		// 	// TODO(mwhittaker): If the function doesn't return anything, don't
		// 	// print t.Results.
		// 	errs = append(errs, bad("return", "The last return must have type error."))
		// }

		// if t.Results().Len() == 0 {
		// 	errs = append(errs, bad("return", "The last return must have type error."))
		// }
	}
	return errors.Join(errs...)
}

// checkMistypedInit returns an error if the provided component implementation
// has an Init method that does not have type "func(context.Context) error".
func checkMistypedInit(pkg *packages.Package, tset *typeSet, impl *types.Named) error {
	for i := 0; i < impl.NumMethods(); i++ {
		m := impl.Method(i)
		if m.Name() != "Init" {
			continue
		}

		// TODO(mwhittaker): Highlight the warning yellow instead of red.
		sig := m.Type().(*types.Signature)
		err := errorf(pkg.Fset, m.Pos(),
			`WARNING: Component %v's Init method has type "%v", not type "func(context.Context) error". It will be ignored.`,
			impl.Obj().Name(), sig)

		// Check Init's parameters.
		if sig.Params().Len() != 1 || !isContext(sig.Params().At(0).Type()) {
			return err
		}

		// Check Init's returns.
		if sig.Results().Len() != 1 || sig.Results().At(0).Type().String() != "error" {
			return err
		}
		return nil
	}
	return nil
}

type printFn func(format string, args ...interface{})

// TODO(mwhittaker): Have generate return an error.
func (g *generator) generate() error {
	if len(g.components) == 0 {
		// There's nothing to generate.
		return nil
	}

	// Process components in deterministic order.
	sort.Slice(g.components, func(i, j int) bool {
		return g.components[i].intfName() < g.components[j].intfName()
	})

	// Generate the file body.
	var body bytes.Buffer
	{
		fn := func(format string, args ...interface{}) {
			fmt.Fprintln(&body, fmt.Sprintf(format, args...))
		}
		// g.generateVersionCheck(fn)
		g.generateRegisteredComponents(fn)
		g.generateInstanceChecks(fn)
		g.generateLocalStubs(fn)

	}

	// Generate the file header. This should be done at the very end to ensure
	// that all types added to the body have been imported.
	var header bytes.Buffer
	{
		fn := func(format string, args ...interface{}) {
			fmt.Fprintln(&header, fmt.Sprintf(format, args...))
		}
		g.generateImports(fn)
	}

	// Create a generated file.
	filename := filepath.Join(g.pkgDir(), generatedCodeFile)
	dst := NewWriter(filename)
	defer dst.Cleanup()

	fmtAndWrite := func(buf bytes.Buffer) error {
		// Format the code.
		b := buf.Bytes()
		formatted, err := format.Source(b)
		if err != nil {
			fmt.Println(string(b))
			return fmt.Errorf("format.Source: %w", err)
		}
		b = formatted

		// Write to dst.
		_, err = io.Copy(dst, bytes.NewReader(b))
		return err
	}

	if err := fmtAndWrite(header); err != nil {
		return err
	}
	if err := fmtAndWrite(body); err != nil {
		return err
	}
	return dst.Close()
}

// pkgDir returns the directory of the package.
func (g *generator) pkgDir() string {
	if len(g.pkg.Syntax) == 0 {
		panic(fmt.Errorf("package %v has no source files", g.pkg))
	}
	f := g.pkg.Syntax[0]
	fname := g.fileset.Position(f.Package).Filename
	return filepath.Dir(fname)
}

// componentRef returns the string to use to refer to the interface
// implemented by a component in generated code.
func (g *generator) componentRef(comp *component) string {
	if comp.isMain {
		return g.kod().qualify("Main")
	}
	return comp.intfName() // We already checked that interface is in the same package.
}

// generateImports generates code to import all the dependencies.
func (g *generator) generateImports(p printFn) {
	p(`// Code generated by "kod generate". DO NOT EDIT.`)
	p("//go:build !ignoreKodGen")
	p("")
	p("package %s", g.pkg.Name)
	p("")
	p(`import (`)
	for _, imp := range g.tset.imports() {
		if imp.alias == "" {
			p(`	%s`, strconv.Quote(imp.path))
		} else {
			p(`	%s %s`, imp.alias, strconv.Quote(imp.path))
		}
	}
	p(`)`)
}

// generateInstanceChecks generates code that checks that every component
// implementation type implements kod.InstanceOf[T] for the appropriate T.
func (g *generator) generateInstanceChecks(p printFn) {
	// If someone deletes a kod.Implements annotation and forgets to re-run
	// `kod generate`, these checks will fail to build. Similarly, if a user
	// changes the interface in a kod.Implements and forgets to re-run
	// `kod generate`, these checks will fail to build.
	p(``)
	p(`// kod.InstanceOf checks.`)
	for _, c := range g.components {
		// e.g., var _ kod.InstanceOf[Odd] = &odd{}
		p(`var _ %s[%s] = (*%s)(nil)`, g.kod().qualify("InstanceOf"), g.tset.genTypeString(c.intf), g.tset.genTypeString(c.impl))
	}
}

// generateRegisteredComponents generates code that registers the components with Kod.
func (g *generator) generateRegisteredComponents(p printFn) {
	if len(g.components) == 0 {
		return
	}

	p(``)
	p(`func init() {`)
	for _, comp := range g.components {
		name := comp.intfName()
		myName := comp.fullIntfName()
		var b strings.Builder
		var inits strings.Builder

		inits.Reset()

		b.Reset()

		g.interceptor()
		localStubFn := fmt.Sprintf(`func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			var interceptors []kod.Interceptor
			if h, ok := info.Impl.(interface{ Interceptors() []kod.Interceptor }); ok {
				interceptors = h.Interceptors()
			}

			%s
			return %s_local_stub{
				impl: info.Impl.(%s),
				interceptor: interceptor.Chain(interceptors),
				name: info.Name,
				caller: info.Caller%s,
			} }`,
			inits.String(), notExported(name), g.componentRef(comp), b.String())
		refNames := make([]string, 0, len(comp.refs))
		for _, ref := range comp.refs {
			refNames = append(refNames, callgraph.MakeEdgeString(comp.fullIntfName(), fullName(ref)))
		}

		reflect := g.tset.importPackage("reflect", "reflect")
		p(`	%s(%s{`, g.codegen().qualify("Register"), g.codegen().qualify("Registration"))
		p(`		Name: %q,`, myName)
		// To get a reflect.Type for an interface, we have to first get a type
		// of its pointer and then resolve the underlying type. See:
		//   https://pkg.go.dev/reflect#example-TypeOf
		p(`		Iface: %s((*%s)(nil)).Elem(),`, reflect.qualify("TypeOf"), g.componentRef(comp))
		p(`		Impl: %s(%s{}),`, reflect.qualify("TypeOf"), comp.implName())
		p("		Refs: `%s`,", strings.Join(refNames, ",\n"))
		p(`		LocalStubFn: %s,`, localStubFn)
		p(`	})`)
	}
	p(`}`)
}

// kod imports and returns the kod package.
func (g *generator) kod() importPkg {
	return g.tset.importPackage(kodPackagePath, "kod")
}

// kod imports and returns the kod package.
func (g *generator) interceptor() importPkg {
	return g.tset.importPackage("github.com/go-kod/kod/core/interceptor", "interceptor")
}

// codegen imports and returns the codegen package.
func (g *generator) codegen() importPkg {
	path := kodPackagePath
	return g.tset.importPackage(path, "kod")
}

// ginContext imports and returns the kerror package.
func (g *generator) ginContext() importPkg {
	return g.tset.importPackage("github.com/gin-gonic/gin", "gin")
}

// formatType pretty prints the provided type, encountered in the provided
// currentPackage.
func formatType(currentPackage *packages.Package, t types.Type) string {
	qualifier := func(pkg *types.Package) string {
		if pkg == currentPackage.Types {
			return ""
		}
		return pkg.Name()
	}
	return types.TypeString(t, qualifier)
}

// generateLocalStubs generates code that creates stubs for the local components.
func (g *generator) generateLocalStubs(p printFn) {
	p(``)
	p(``)
	p(`// Local stub implementations.`)
	g.tset.importPackage("context", "context")

	var b strings.Builder
	for _, comp := range g.components {

		stub := notExported(comp.intfName()) + "_local_stub"
		p(``)
		p(`type %s struct{`, stub)
		p(`	impl %s`, g.componentRef(comp))
		p(`	name   string`)
		p(` caller string`)
		p(`	interceptor kod.Interceptor`)
		p(`}`)

		p(``)
		p(`// Check that %s implements the %s interface.`, stub, g.tset.genTypeString(comp.intf))
		p(`var _ %s = (*%s)(nil)`, g.tset.genTypeString(comp.intf), stub)
		p(``)

		for _, m := range comp.methods() {

			mt := m.Type().(*types.Signature)
			p(``)
			p(`func (s %s) %s(%s) (%s) {`, stub, m.Name(), g.args(mt), g.returns(mt))

			if haveGinContext(mt) {
				p(`var err error`)
				g.ginContext()
				p(`	ctx := a0.Request.Context()`)
			}

			if isHttpHandler(mt) {
				p(` ctx := a1.Context()`)
			}

			p(`info := kod.CallInfo{
					Component:  s.name,
					FullMethod: "%s.%s",
					Caller:     s.caller,
				}`, comp.fullIntfName(), m.Name())

			p(`
				if s.interceptor == nil {
					%s s.impl.%s(%s)
					return
				}
			`, g.returnsList(mt), m.Name(), g.argList(comp, mt))

			p(`call := func(ctx context.Context, info kod.CallInfo, req, res []any) (err error) {`)

			if haveGinContext(mt) {
				p(` a0.Request = a0.Request.WithContext(ctx)`)
			}
			if isHttpHandler(mt) {
				p(` a1 = a1.WithContext(ctx)`)
			}

			p(`	%s s.impl.%s(%s)
					%s
					return
				}`, g.returnsList(mt), m.Name(), g.argList(comp, mt), g.setReturnsList(mt))

			p(`			
				%s = s.interceptor(%s, info, []any{%s}, []any{%s}, call)`,
				lo.If(haveError(mt) || haveGinContext(mt), "err").Else("_"),
				lo.If(haveContext(mt) || haveGinContext(mt) || isHttpHandler(mt), "ctx").Else("context.Background()"),
				g.argsReflectList(comp, mt), g.returnsReflectList(mt))

			// Call the local method.
			b.Reset()

			if haveGinContext(mt) {
				g.ginContext()
				p(`if err != nil {`)
				p(`	a0.Error(err)`)
				p(`}`)
			}
			if mt.Results().Len() > 0 {
				p(`	return`)
			}
			p(`}`)
		}
	}
}

func (g *generator) setReturnsList(sig *types.Signature) string {
	var returns strings.Builder
	for i := 0; i < sig.Results().Len(); i++ {
		rt := sig.Results().At(i).Type()
		if g.tset.genTypeString(rt) == "error" {
			// fmt.Fprintf(&returns, "res[%d] = err\n", i)
			continue
		}
		fmt.Fprintf(&returns, "res[%d] = r%d\n", i, i)
	}

	return returns.String()
}

func (g *generator) argList(comp *component, sig *types.Signature) string {
	if sig.Params().Len() == 0 {
		return ""
	}

	var b strings.Builder
	if haveContext(sig) {
		fmt.Fprintf(&b, "ctx")
	} else {
		fmt.Fprintf(&b, "a0")
	}

	for i := 1; i < sig.Params().Len(); i++ {
		if sig.Variadic() && i == sig.Params().Len()-1 {
			fmt.Fprintf(&b, ", a%d...", i)
		} else {
			fmt.Fprintf(&b, ", a%d", i)
		}
	}
	return b.String()
}

func (g *generator) argsReflectList(comp *component, sig *types.Signature) string {
	if sig.Params().Len() == 0 {
		return ""
	}

	var b strings.Builder

	for i := 0; i < sig.Params().Len(); i++ {
		// filter context.Context
		if g.tset.genTypeString(sig.Params().At(i).Type()) == "context.Context" {
			continue
		}

		if sig.Variadic() {
			fmt.Fprintf(&b, "a%d...", i)
		} else {
			fmt.Fprintf(&b, "a%d", i)
		}

		if i != sig.Params().Len()-1 {
			fmt.Fprintf(&b, ",")
		}
	}
	return b.String()
}

// args returns a textual representation of the arguments of the provided
// signature. The first argument must be a context.Context. The returned code
// names the first argument ctx and all subsequent arguments a0, a1, and so on.
func (g *generator) args(sig *types.Signature) string {
	if sig.Params().Len() == 0 {
		return ""
	}

	var args strings.Builder
	for i := 0; i < sig.Params().Len(); i++ { // Skip initial context.Context
		at := sig.Params().At(i).Type()
		if !sig.Variadic() || i != sig.Params().Len()-1 {
			if g.tset.genTypeString(at) == "context.Context" {
				fmt.Fprintf(&args, ",ctx context.Context")
				continue
			}

			fmt.Fprintf(&args, ", a%d %s", i, g.tset.genTypeString(at))
			continue
		}
		// For variadic functions, the final argument is guaranteed to be a
		// slice. Instead of passing an argument of type []t, we pass ...t.
		subtype := at.(*types.Slice).Elem()
		if g.tset.genTypeString(subtype) == "context.Context" {
			fmt.Fprintf(&args, ",ctx ...context.Context")
			continue
		}

		fmt.Fprintf(&args, ", a%d ...%s", i, g.tset.genTypeString(subtype))
	}
	return args.String()[1:]
}

// returns returns a textual representation of the returns of the provided
// signature. The last return must be an error. The returned code names the
// returns r0, r1, and so on. The returned error is called err.
func (g *generator) returns(sig *types.Signature) string {
	var returns strings.Builder
	for i := 0; i < sig.Results().Len(); i++ {
		rt := sig.Results().At(i).Type()
		if g.tset.genTypeString(rt) == "error" {
			fmt.Fprintf(&returns, "err error")
			continue
		}
		fmt.Fprintf(&returns, "r%d %s, ", i, g.tset.genTypeString(rt))
	}
	return returns.String()
}

func (g *generator) returnsList(sig *types.Signature) string {
	var returns strings.Builder
	for i := 0; i < sig.Results().Len(); i++ {
		rt := sig.Results().At(i).Type()
		if g.tset.genTypeString(rt) == "error" {
			fmt.Fprintf(&returns, "err")
			continue
		}
		fmt.Fprintf(&returns, "r%d", i)
		if i != sig.Results().Len()-1 {
			fmt.Fprintf(&returns, ", ")
		}
	}
	if returns.Len() == 0 {
		return ""
	}
	return returns.String() + "="
}

func (g *generator) returnsReflectList(sig *types.Signature) string {
	var returns strings.Builder
	for i := 0; i < sig.Results().Len(); i++ {
		rt := sig.Results().At(i).Type()
		if g.tset.genTypeString(rt) == "error" {
			continue
		}
		fmt.Fprintf(&returns, "r%d", i)
		if i != sig.Results().Len()-1 {
			fmt.Fprintf(&returns, ", ")
		}
	}

	return returns.String()
}

// notExported sets the first character in the string to lowercase.
func notExported(name string) string {
	if len(name) == 0 {
		return name
	}
	a := []rune(name)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

// func isPointerToStruct(t types.Type) bool {
// 	ptr, ok := t.(*types.Pointer)
// 	if !ok {
// 		return false
// 	}
// 	_, ok = ptr.Elem().Underlying().(*types.Struct)
// 	return ok
// }

// func needInterceptor(m *types.Signature) bool {
// 	return haveContext(m) || haveGinContext(m)
// }

func haveContext(mt *types.Signature) bool {
	for i := 0; i < mt.Params().Len(); i++ {
		if isContext(mt.Params().At(i).Type()) {
			return true
		}
	}

	return false
}

func haveGinContext(mt *types.Signature) bool {
	if mt.Params().Len() == 1 {
		if isGinContext(mt.Params().At(0).Type()) {
			return true
		}
	}

	return false
}

func isHttpHandler(mt *types.Signature) bool {
	if mt.Params().Len() == 2 {
		if isHttpResponseWriter(mt.Params().At(0).Type()) && isHttpRequest(mt.Params().At(1).Type()) {
			return true
		}
	}

	return false
}

func haveError(mt *types.Signature) bool {
	for i := 0; i < mt.Results().Len(); i++ {
		if isError(mt.Results().At(i).Type()) {
			return true
		}
	}
	return false
}
