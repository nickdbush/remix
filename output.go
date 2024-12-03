package remix

import (
	"fmt"
	"go/format"
	"go/types"
	"os"
)

type Output interface {
	Add(string, *Node)
	Convert(*Node, *Node)
	Finish()
}

type Source struct {
	pkg    Package
	file   string
	buffer string
}

func ToSource(pkg Package, file string) Source {
	return Source{
		pkg:    pkg,
		file:   file,
		buffer: fmt.Sprintf("package %s\n", pkg.Name()),
	}
}

func (s *Source) Add(name string, node *Node) {

}

func (s *Source) Convert(name string, from, to *Node, casts *Casts) {
	s.buffer += Conversion{
		From:      from,
		To:        to,
		Qualifier: types.RelativeTo(s.pkg.pkg),
		Casts:     casts,
	}.source(name)
}

func (s *Source) Finish() {
	bytes := s.finish()
	err := os.WriteFile(s.file, bytes, 0666)
	if err != nil {
		panic(err)
	}
}

func (s *Source) finish() []byte {
	bytes := []byte(s.buffer)
	formatted, formatErr := format.Source(bytes)
	if formatErr != nil {
		return bytes
	}
	return formatted
}

type Dot struct {
	nodes       map[int]dotNode
	outputs     map[string]int
	edges       map[dotEdge]struct{}
	conversions map[dotEdge]struct{}
}

type dotNode struct {
	label        string
	intermediate bool
}

type dotEdge struct {
	from, to int
}

func ToDot() *Dot {
	return &Dot{
		nodes:       make(map[int]dotNode),
		outputs:     make(map[string]int),
		edges:       make(map[dotEdge]struct{}),
		conversions: make(map[dotEdge]struct{}),
	}
}

func (w *Dot) Add(name string, node *Node) {
	w.insertNode(node)
	w.outputs[name] = node.id
}

func (w *Dot) Convert(from, to *Node) {
	w.insertNode(from)
	w.insertNode(to)
	edge := dotEdge{from: from.id, to: to.id}
	w.conversions[edge] = struct{}{}
}

func (w *Dot) Finish() {
	out := "digraph G {\n"
	for id, node := range w.nodes {
		if node.intermediate {
			out += fmt.Sprintf("\tn%d [label=%q shape=box];\n", id, node.label)
		} else {
			out += fmt.Sprintf("\tn%d [label=%q];\n", id, node.label)
		}
	}

	for name, id := range w.outputs {
		out += fmt.Sprintf("\to%d [label=%q shape=hexagon];\n", id, name)
		out += fmt.Sprintf("\tn%d -> o%d;\n", w.outputs[name], id)
	}

	for edge := range w.edges {
		out += fmt.Sprintf("\tn%d -> n%d;\n", edge.from, edge.to)
	}

	for edge := range w.conversions {
		out += fmt.Sprintf("\tn%d -> n%d [style=dashed];\n", edge.from, edge.to)
	}

	out += "}\n"

	fmt.Println(out)
}

func (w *Dot) insertNode(node *Node) {
	if op, ok := node.Operation.(loadOp); ok {
		w.nodes[node.id] = dotNode{label: op.debug(), intermediate: false}
	} else {
		w.nodes[node.id] = dotNode{label: node.debug(), intermediate: true}
	}

	for _, parent := range node.parents() {
		w.insertNode(parent)
		edge := dotEdge{from: parent.id, to: node.id}
		w.edges[edge] = struct{}{}
	}
}
