// package mdext is a extension for the goldmark(http://github.com/yuin/goldmark).
// The tom extension parses TOML metadata blocks and store metadata to a
// parser.Context.
// It is copied slavishly from http://github.com/yuin/goldmark-meta
package tomltc

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	//"gopkg.in/yaml.v2"
	"github.com/BurntSushi/toml"
)

// MapItem is an item in a MapSlice.
// From https://github.com/go-yaml/yaml/blob/v2/yaml.go
type MapItem struct {
	Key, Value interface{}
}

// MapSlice encodes and decodes as a TOML map.
// The order of keys is preserved when encoding and decoding.
// From https://github.com/go-yaml/yaml/blob/v2/yaml.go
type MapSlice []MapItem

type tomldata struct {
	Map   map[string]interface{}
	Items MapSlice
	Error error
	Node  gast.Node
}

var contextKey = parser.NewContextKey()

// Get returns a TOML metadata.
func Get(pc parser.Context) map[string]interface{} {
	v := pc.Get(contextKey)
	if v == nil {
		return nil
	}
	// TC:
	fmt.Printf("Get(): v.(*tomldata) is %#v\n", v.(*tomldata))
	d := v.(*tomldata)
	return d.Map
}

// TryGet tries to get a TOML metadata.
// If there are TOML parsing errors, then nil and error are returned
func TryGet(pc parser.Context) (map[string]interface{}, error) {
	fmt.Printf("TryGet contextKey: [%v]\n", contextKey)
	dtmp := pc.Get(contextKey)
	if dtmp == nil {
		//fmt.Printf("TryGet(): %#v\n", dtmp);
		return nil, nil
	}
	d := dtmp.(*tomldata)
	if d.Error != nil {
		return nil, d.Error
	}
	//fmt.Printf("TryGet(): %#v\n", d.Map);
	return d.Map, nil
}

// GetItems returns a TOML metadata.
// GetItems preserves defined key order.
func GetItems(pc parser.Context) MapSlice {
	v := pc.Get(contextKey)
	if v == nil {
		return nil
	}
	d := v.(*tomldata)
	// TC:
	//fmt.Printf("GetItems(): %#v\n", d.Items);
	return d.Items
}

// TryGetItems returns a TOML metadata.
// TryGetItems preserves defined key order.
// If there are TOML parsing errors, then nil and error are returned.
func TryGetItems(pc parser.Context) (MapSlice, error) {
	fmt.Printf("TryGetItems()")
	dtmp := pc.Get(contextKey)
	if dtmp == nil {
		return nil, nil
	}
	d := dtmp.(*tomldata)
	if d.Error != nil {
		return nil, d.Error
	}
	fmt.Printf("\t%#v\n", d.Items)
	return d.Items, nil
}

type metaParser struct {
}

var defaultMetaParser = &metaParser{}

// NewParser returns a BlockParser that can parse TOML metadata blocks.
func NewParser() parser.BlockParser {
	return defaultMetaParser
}

func isSeparator(line []byte) bool {
	line = util.TrimRightSpace(util.TrimLeftSpace(line))
	for i := 0; i < len(line); i++ {
		if line[i] != '-' {
			return false
		}
	}
	return true
}

func (b *metaParser) Trigger() []byte {
	return []byte{'-'}
}

func (b *metaParser) Open(parent gast.Node, reader text.Reader, pc parser.Context) (gast.Node, parser.State) {
	linenum, _ := reader.Position()
	if linenum != 0 {
		return nil, parser.NoChildren
	}
	line, _ := reader.PeekLine()
	if isSeparator(line) {
		return gast.NewTextBlock(), parser.NoChildren
	}
	return nil, parser.NoChildren
}

func (b *metaParser) Continue(node gast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()
	if isSeparator(line) && !util.IsBlank(line) {
		// TC: I think this is where at least 1 line
		// of the metadata was found
		reader.Advance(segment.Len())
		return parser.Close
	}
	node.Lines().Append(segment)
	return parser.Continue | parser.NoChildren
}

func (b *metaParser) Close(node gast.Node, reader text.Reader, pc parser.Context) {
	lines := node.Lines()
	var buf bytes.Buffer
	for i := 0; i < lines.Len(); i++ {
		segment := lines.At(i)
		buf.Write(segment.Value(reader.Source()))
	}
	d := &tomldata{}
	d.Node = node
	meta := map[string]interface{}{}
	// meta := make(map[string]interface{})
	fmt.Printf("Close(): about to unmarshal\n\t%v\n", string(buf.Bytes()))
	if err := toml.Unmarshal(buf.Bytes(), &meta); err != nil {
		fmt.Printf("\t error unmarshalling\n%v\n", string(buf.Bytes()))
		d.Error = err
	} else {
		d.Map = meta
	}

	metaMapSlice := MapSlice{}
	if err := toml.Unmarshal(buf.Bytes(), &metaMapSlice); err != nil {
		d.Error = err
	} else {
		d.Items = metaMapSlice
	}

	pc.Set(contextKey, d)

	if d.Error == nil {
		node.Parent().RemoveChild(node.Parent(), node)
	}
}

func (b *metaParser) CanInterruptParagraph() bool {
	return false
}

func (b *metaParser) CanAcceptIndentedLine() bool {
	return false
}

type astTransformer struct {
}

var defaultASTTransformer = &astTransformer{}

func (a *astTransformer) Transform(node *gast.Document, reader text.Reader, pc parser.Context) {
	fmt.Println("Transform() -- START")
	fmt.Println("\tpc.Get()")
	dtmp := pc.Get(contextKey)
	if dtmp == nil {
		return
	}
	d := dtmp.(*tomldata)
	// xxx
	if d.Error != nil {
		fmt.Println("\tError in tmp")
		msg := gast.NewString([]byte(fmt.Sprintf("<!-- %s -->", d.Error)))
		msg.SetCode(true)
		d.Node.AppendChild(d.Node, msg)
		return
	}

	fmt.Println("\tpc.GetItems()")
	meta := GetItems(pc)
	if meta == nil {
		return
	}
	table := east.NewTable()
	fmt.Println("\ttable: %#v", table)
	alignments := []east.Alignment{}
	for range meta {
		alignments = append(alignments, east.AlignNone)
	}
	row := east.NewTableRow(alignments)
	for _, item := range meta {
		cell := east.NewTableCell()
		cell.AppendChild(cell, gast.NewString([]byte(fmt.Sprintf("%v", item.Key))))
		row.AppendChild(row, cell)
	}
	table.AppendChild(table, east.NewTableHeader(row))

	row = east.NewTableRow(alignments)
	for _, item := range meta {
		cell := east.NewTableCell()
		cell.AppendChild(cell, gast.NewString([]byte(fmt.Sprintf("%v", item.Value))))
		row.AppendChild(row, cell)
	}
	table.AppendChild(table, row)
	node.InsertBefore(node, node.FirstChild(), table)
	fmt.Println("\tcomplete table: %#v", table)
	fmt.Println("Transform() -- END")
}

// Option is a functional option type for this extension.
type Option func(*meta)

// WithTable is a functional option that renders a TOML metadata as a table.
func WithTable() Option {
	return func(m *meta) {
		m.Table = true
	}
}

type meta struct {
	Table bool
}

// Meta is a extension for the goldmark.
var Meta = &meta{}

// New returns a new Meta extension.
func New(opts ...Option) goldmark.Extender {
	e := &meta{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *meta) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(NewParser(), 0),
		),
	)
	if e.Table {
		m.Parser().AddOptions(
			parser.WithASTTransformers(
				util.Prioritized(defaultASTTransformer, 0),
			),
		)
	}
}
