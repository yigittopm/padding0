package analysis

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"sort"
	"strings"
)

type FieldWithComments struct {
	Names   []*ast.Ident
	Field   *ast.Field
	Tag     *ast.BasicLit
	Doc     *ast.CommentGroup
	Comment *ast.CommentGroup
}

var structDefinitions = make(map[string]*ast.StructType)

func TokenizeStructFields(content string) (string, error) {
	fileSet := token.NewFileSet()

	node, err := parser.ParseFile(fileSet, "", content, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return content, err
	}

	// Bug is here but I don't know how to fix it
	ast.Inspect(node, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if structType, ok := ts.Type.(*ast.StructType); ok {
			structDefinitions[ts.Name.Name] = structType
			_ = sortFieldsBySize(structType)
		} else if ident, ok := ts.Type.(*ast.Ident); ok {
			if size, found := typeSizes[ident.Name]; found {
				typeSizes[ts.Name.Name] = size
			}
		}

		return true
	})

	formattedContent := new(strings.Builder)
	err = printer.Fprint(formattedContent, fileSet, node)
	if err != nil {
		return content, err
	}

	return formattedContent.String(), nil
}

func sortFieldsBySize(structType *ast.StructType) []*FieldWithComments {
	fields := structType.Fields.List

	fieldsWithComments := make([]*FieldWithComments, len(fields))
	for i, field := range fields {
		fieldsWithComments[i] = &FieldWithComments{
			Names:   field.Names,
			Field:   field,
			Tag:     field.Tag,
			Doc:     field.Doc,
			Comment: field.Comment,
		}
	}

	sort.SliceStable(fieldsWithComments, func(i, j int) bool {
		sizeI := getSizeOfType(fieldsWithComments[i].Field)
		sizeJ := getSizeOfType(fieldsWithComments[j].Field)

		return sizeI > sizeJ
	})

	sortedFields := make([]*ast.Field, len(fieldsWithComments))
	for i, sortedField := range fieldsWithComments {
		sortedFields[i] = &ast.Field{
			Names:   sortedField.Names,
			Type:    sortedField.Field.Type,
			Tag:     sortedField.Field.Tag,
			Doc:     sortedField.Doc,
			Comment: sortedField.Comment,
		}
	}

	structType.Fields.List = sortedFields

	for _, field := range fieldsWithComments {
		switch typ := field.Field.Type.(type) {
		case *ast.Ident:
			if innerStruct, found := structDefinitions[typ.Name]; found {
				sortFieldsBySize(innerStruct)
			}
		case *ast.StructType:
			sortFieldsBySize(typ)
		}
	}

	return fieldsWithComments
}

func getSizeOfType(field *ast.Field) int {
	switch expr := field.Type.(type) {
	case *ast.Ident:
		if size, ok := typeSizes[expr.Name]; ok {
			return size
		} else if structType, ok := structDefinitions[expr.Name]; ok {
			return calculateStructSize(structType)
		} else {
			return typeSizes[expr.Name]
		}
	case *ast.StructType:
		return calculateStructSize(expr)
	}
	return 0
}

func calculateStructSize(structType *ast.StructType) int {
	totalSize := 0
	for _, field := range structType.Fields.List {
		totalSize += getSizeOfType(field)
	}
	return totalSize
}
