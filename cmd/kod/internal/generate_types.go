package internal

import (
	"fmt"
	"go/types"
	"sort"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

const kodPackagePath = "github.com/go-kod/kod"

// typeSet holds type information needed by the code generator.
type typeSet struct {
	pkg            *packages.Package
	imported       []importPkg          // imported packages
	importedByPath map[string]importPkg // imported, indexed by path
	importedByName map[string]importPkg // imported, indexed by name

	// If checked[t] != nil, then checked[t] is the cached result of calling
	// check(pkg, t, string[]{}). Otherwise, if checked[t] == nil, then t has
	// not yet been checked for serializability. Read typeutil.Map's
	// documentation for why checked shouldn't be a map[types.Type]bool.
	// checked typeutil.Map

	// If measurable[t] != nil, then measurable[t] == isMeasurableType(t).
	//nolint
	measurable typeutil.Map
}

// importPkg is a package imported by the generated code.
type importPkg struct {
	path  string // e.g., "github.com/go-kod/kod"
	pkg   string // e.g., "kod", "context", "time"
	alias string // e.g., foo in `import foo "context"`
	local bool   // are we in this package?
}

// name returns the name by which the imported package should be referenced in
// the generated code. If the package is imported without an alias, like this:
//
//	import "context"
//
// then the name is the same as the package name (e.g., "context"). However, if
// a package is imported with an alias, then the name is the alias:
//
//	import thisIsAnAlias "context"
//
// If the package is local, an empty string is returned.
func (i importPkg) name() string {
	if i.local {
		return ""
	} else if i.alias != "" {
		return i.alias
	}
	return i.pkg
}

// qualify returns the provided member of the package, qualified with the
// package name. For example, the "Context" type inside the "context" package
// is qualified "context.Context". The "Now" function inside the "time" package
// is qualified "time.Now". Note that the package name is not prefixed when
// qualifying members of the local package.
func (i importPkg) qualify(member string) string {
	if i.local {
		return member
	}
	return fmt.Sprintf("%s.%s", i.name(), member)
}

// newTypeSet returns the container for types found in pkg.
func newTypeSet(pkg *packages.Package, automarshals, automarshalCandidates *typeutil.Map) *typeSet {
	return &typeSet{
		pkg:            pkg,
		imported:       []importPkg{},
		importedByPath: map[string]importPkg{},
		importedByName: map[string]importPkg{},
	}
}

// importPackage imports a package with the provided path and package name. The
// package is imported with an alias if there is a package name clash.
func (tset *typeSet) importPackage(path, pkg string) importPkg {
	newImportPkg := func(path, pkg, alias string, local bool) importPkg {
		i := importPkg{path: path, pkg: pkg, alias: alias, local: local}
		tset.imported = append(tset.imported, i)
		tset.importedByPath[i.path] = i
		tset.importedByName[i.name()] = i
		return i
	}

	if imp, ok := tset.importedByPath[path]; ok {
		// This package has already been imported.
		return imp
	}

	if _, ok := tset.importedByName[pkg]; !ok {
		// Import the package without an alias.
		return newImportPkg(path, pkg, "", path == tset.pkg.PkgPath)
	}

	// Find an unused alias.
	var alias string
	counter := 1
	for {
		alias = fmt.Sprintf("%s%d", pkg, counter)
		if _, ok := tset.importedByName[alias]; !ok {
			break
		}
		counter++
	}
	return newImportPkg(path, pkg, alias, path == tset.pkg.PkgPath)
}

// imports returns the list of packages to import in generated code.
func (tset *typeSet) imports() []importPkg {
	sort.Slice(tset.imported, func(i, j int) bool {
		return tset.imported[i].path < tset.imported[j].path
	})
	return tset.imported
}

// genTypeString returns the string representation of t as to be printed
// in the generated code, updating import definitions to account for the
// returned type string.
//
// Since this call has side-effects (i.e., updating import definitions), it
// should only be called when the returned type string is written into
// the generated file; otherwise, the generated code may end up with spurious
// imports.
func (tset *typeSet) genTypeString(t types.Type) string {
	// qualifier is passed to types.TypeString(Type, Qualifier) to determine
	// how packages are printed when pretty printing types. For this qualifier,
	// types in the root package are printed without their package name, while
	// types outside the root package are printed with their package name. For
	// example, if we're in root package foo, then the type foo.Bar is printed
	// as Bar, while the type io.Reader is printed as io.Reader. See [1] for
	// more information on qualifiers and pretty printing types.
	//
	// [1]: https://github.com/golang/example/tree/master/gotypes#formatting-support
	var qualifier = func(pkg *types.Package) string {
		if pkg == tset.pkg.Types {
			return ""
		}
		return tset.importPackage(pkg.Path(), pkg.Name()).name()
	}
	return types.TypeString(t, qualifier)
}

// isKodType returns true iff t is a named type from the kod package with
// the specified name and n type arguments.
func isKodType(t types.Type, name string, n int) bool {
	named, ok := t.(*types.Named)
	return ok &&
		named.Obj().Pkg() != nil &&
		named.Obj().Pkg().Path() == kodPackagePath &&
		named.Obj().Name() == name &&
		named.TypeArgs().Len() == n
}

func isKodImplements(t types.Type) bool {
	return isKodType(t, "Implements", 1)
}

func isKodRef(t types.Type) bool {
	return isKodType(t, "Ref", 1)
}

func isKodMain(t types.Type) bool {
	return isKodType(t, "Main", 0)
}

func isContext(t types.Type) bool {
	n, ok := t.(*types.Named)
	if !ok {
		return false
	}
	return n.Obj().Pkg().Path() == "context" && n.Obj().Name() == "Context"
}

func isGinContext(t types.Type) bool {
	p, ok := t.(*types.Pointer)
	if !ok {
		return false
	}
	n, ok := p.Elem().(*types.Named)
	if !ok {
		return false
	}
	return n.Obj().Pkg().Path() == "github.com/gin-gonic/gin" && n.Obj().Name() == "Context"
}

func isHttpResponseWriter(t types.Type) bool {
	n, ok := t.(*types.Named)
	if !ok {
		return false
	}
	return n.Obj().Pkg().Path() == "net/http" && n.Obj().Name() == "ResponseWriter"
}

func isHttpRequest(t types.Type) bool {
	p, ok := t.(*types.Pointer)
	if !ok {
		return false
	}
	n, ok := p.Elem().(*types.Named)
	if !ok {
		return false
	}
	return n.Obj().Pkg().Path() == "net/http" && n.Obj().Name() == "Request"
}

func isError(t types.Type) bool {
	n, ok := t.(*types.Named)
	if !ok {
		return false
	}
	return n.Obj().Name() == "error"
}
