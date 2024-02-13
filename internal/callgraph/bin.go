package callgraph

import (
	"bytes"
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"fmt"
	"os"
	"slices"

	"github.com/dominikbraun/graph"
	"github.com/samber/lo"
)

// rodata returns the read-only data section of the provided binary.
func rodata(file string) ([]byte, error) {
	f := lo.Must(os.Open(file))

	defer f.Close()

	// Look at first few bytes to determine the file format.
	prefix := make([]byte, 4)
	lo.Must(f.ReadAt(prefix, 0))

	// Handle the file formats we support.
	switch {
	case bytes.HasPrefix(prefix, []byte("\x7FELF")): // Linux
		f := lo.Must(elf.NewFile(f))
		return f.Section(".rodata").Data()
	case bytes.HasPrefix(prefix, []byte("MZ")): // Windows
		f := lo.Must(pe.NewFile(f))
		return f.Section(".rdata").Data()
	// case bytes.HasPrefix(prefix, []byte("\xFE\xED\xFA")): // MacOS
	// 	f := lo.Must(macho.NewFile(f))
	// 	return f.Section("__rodata").Data()
	case bytes.HasPrefix(prefix[1:], []byte("\xFA\xED\xFE")): // MacOS
		f := lo.Must(macho.NewFile(f))
		return f.Section("__rodata").Data()
	default:
		return nil, fmt.Errorf("unknown format")
	}
}

// ReadComponentGraph reads component graph information from the specified
// binary. It returns a slice of components and a component graph whose nodes
// are indices into that slice.
func ReadComponentGraph(file string) (graph.Graph[string, string], error) {
	data := lo.Must(rodata(file))

	es := ParseEdges(data)
	g := graph.New(graph.StringHash, graph.Directed())

	// NOTE: initially, all node numbers are zero.
	nodeMap := map[string]int{}
	for _, e := range es {
		nodeMap[e[0]] = 0
		nodeMap[e[1]] = 0
	}
	// Assign node numbers.
	components := lo.Keys(nodeMap)
	slices.Sort(components)
	for _, c := range components {
		lo.Must0(g.AddVertex(c))
	}

	// Convert component edges into graph edges.
	for _, e := range es {
		src := e[0]
		dst := e[1]
		lo.Must0(g.AddEdge(src, dst))
	}

	return g, nil
}
