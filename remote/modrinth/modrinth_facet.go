package modrinth

import (
	"strings"
)

type searchResultIndexing uint8

const (
	indexRelevance searchResultIndexing = iota
	indexDownloads
	indexFollows
	indexNewest
	indexUpdated
)

type facetItemOperation uint8

const (
	operationEq facetItemOperation = iota
	operationNeq
	operationLeq
	operationGeq
	operationLt
	operationGt
)

var facetItemOperationStrings = map[facetItemOperation]string{
	operationEq:  ":",
	operationNeq: "!=",
	operationLeq: "<=",
	operationGeq: ">=",
	operationLt:  "<",
	operationGt:  ">",
}

// facet is the data structure to construct an advanced Modrinth search. It
// does not contain all the search options, only the ones that are expected in this
// program.
//
// From Modrinth docs:
//
// In order to then use these facets, you need a value to filter by, as well as
// an operation to perform on this value. The most common operation is ':'
// (same as =), though you can also use !=, >=, >, <=, and <. Join together the
// type, operation, and value, and you’ve got your string.
//
// {type} {operation} {value}
//
// API Docs: https://docs.modrinth.com/api/operations/searchprojects/
type facetItem struct {
	Type      string
	Operation facetItemOperation
	Value     string
}

func (f *facetItem) String() string {
	return `"` + f.Type + facetItemOperationStrings[f.Operation] + f.Value + `"`
}

// facetItems is an array of facetItem. It represents an expression joined by OR statements.
// a complete facet is an array of facetItems, with each array joined by AND statements.
type facetItems []facetItem

// A facet is an array of facetItems, with each array joined by AND statements.
type facet struct {
	Expressions []facetItems
}

// There are no facet data structures, rather, a function is used to directly
// create a facet string that can be used in the URL.
func createFacet(expressions ...facetItems) string {
	var sb strings.Builder
	sb.WriteRune('[')
	for i, expression := range expressions {
		if i > 0 {
			sb.WriteRune(',')
		}
		sb.WriteRune('[')
		for j, item := range expression {
			if j > 0 {
				sb.WriteRune(',')
			}
			sb.WriteString(item.String())
		}
		sb.WriteRune(']')
	}
	sb.WriteRune(']')
	return sb.String()
}

// StringifyFacets builds multiple facet structs into a string that can be embedded
// into the url's facets props. It joins all facets with an AND operator.
//
// The facet prop uses a json-like format:
//
// {type} {operation} {value}
//
// categories = adventure
// versions != 1.20.1
// downloads <= 100
//
// You then join these strings together in arrays to signal AND and OR operators.
//
// OR
// All elements in a single array are considered to be joined by OR statements.
// For example, the search [["versions:1.16.5", "versions:1.17.1"]] translates
// to Projects that support 1.16.5 OR 1.17.1.
//
// AND
// Separate arrays are considered to be joined by AND statements. For example,
// the search [["versions:1.16.5"], ["project_type:modpack"]] translates to
// Projects that support 1.16.5 AND are modpacks.
func StringifyFacets(facets ...*facet) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, facet := range facets {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString("[")
		for j, items := range facet.Expressions {
			if j > 0 {
				sb.WriteString(",")
			}
			for k, item := range items {
				if k > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(`"`)
				sb.WriteString(item.Type)
				sb.WriteString(facetItemOperationStrings[item.Operation])
				sb.WriteString(item.Value)
				sb.WriteString(`"`)
			}
		}
		sb.WriteString("]")
	}
	sb.WriteString("]")
	return sb.String()
}

var facetAllLoaders = facetItems{
	{
		Type:      "categories",
		Operation: operationEq,
		Value:     "forge",
	},
	{
		Type:      "categories",
		Operation: operationEq,
		Value:     "fabric",
	},
	{
		Type:      "categories",
		Operation: operationEq,
		Value:     "quilt",
	},
	{
		Type:      "categories",
		Operation: operationEq,
		Value:     "liteloader",
	},
	{
		Type:      "categories",
		Operation: operationEq,
		Value:     "modloader",
	},
	{
		Type:      "categories",
		Operation: operationEq,
		Value:     "rift",
	},
	{
		Type:      "categories",
		Operation: operationEq,
		Value:     "neoforge",
	},
}

var facetForge = facetItems{
	{
		Type:      "categories",
		Operation: operationEq,
		Value:     "forge",
	},
}

var facetFabric = facetItems{
	{
		Type:      "categories",
		Operation: operationEq,
		Value:     "fabric",
	},
}

var facetServerSupported = facetItems{
	{
		Type:      "server_side",
		Operation: operationEq,
		Value:     "required",
	},
	{
		Type:      "server_side",
		Operation: operationEq,
		Value:     "optional",
	},
}

var facetClientSupported = facetItems{
	{
		Type:      "client_side",
		Operation: operationEq,
		Value:     "required",
	},
	{
		Type:      "client_side",
		Operation: operationEq,
		Value:     "optional",
	},
}

var facetBothRequired = []facetItems{
	{
		{
			Type:      "server_side",
			Operation: operationEq,
			Value:     "required",
		},
	},
	{
		{
			Type:      "client_side",
			Operation: operationEq,
			Value:     "required",
		},
	},
}

var facetBothSupported = []facetItems{
	{
		{
			Type:      "server_side",
			Operation: operationEq,
			Value:     "required",
		},
		{
			Type:      "server_side",
			Operation: operationEq,
			Value:     "optional",
		},
	},
	{
		{
			Type:      "client_side",
			Operation: operationEq,
			Value:     "required",
		},
		{
			Type:      "client_side",
			Operation: operationEq,
			Value:     "optional",
		},
	},
}
