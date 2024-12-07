package analysis

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type FieldWithComments struct {
	Names   []*ast.Ident
	Field   *ast.Field
	Doc     *ast.CommentGroup
	Comment *ast.CommentGroup
}

func Start() error {
	dirFlag := flag.String("d", "", "Directory to search for .go files")
	flag.Parse()

	var dir string
	if *dirFlag == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			log.Fatalf("Dir not found: %v", err)
		}
	} else {
		dir = *dirFlag
	}

	if err := getFiles(dir); err != nil {
		return err
	}

	return nil
}

func getFiles(dir string) error {
	var wg sync.WaitGroup
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") && !strings.HasSuffix(info.Name(), "_test.go") {
			wg.Add(1)
			go func() {
				defer wg.Done()
				readFile(path)
			}()
		}

		return nil
	})

	if err != nil {
		return err
	}

	wg.Wait()
	return nil
}

func readFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		if err.Error() != "EOF" {
			log.Fatalf("Error opening file: %v", err)
		}
		return nil
	}

	modifiedContent, err := TokenizeStructFields(string(file))
	if err != nil {
		return err
	}

	err = writeFile(path, modifiedContent)
	if err != nil {
		return err
	}

	return nil
}

func writeFile(path string, content string) error {
	content = strings.TrimSpace(content)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}

	err = formatTheWrittenFile(path)
	if err != nil {
		return err
	}

	return nil
}

func formatTheWrittenFile(path string) error {
	cmd := exec.Command("gofmt", "-w", path)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

var structDefinitions = make(map[string]*ast.StructType)

func TokenizeStructFields(content string) (string, error) {
	fileSet := token.NewFileSet()

	node, err := parser.ParseFile(fileSet, "", content, parser.AllErrors)
	if err != nil {
		return content, err
	}

	ast.Inspect(node, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if structType, ok := ts.Type.(*ast.StructType); ok {
			structDefinitions[ts.Name.Name] = structType
		} else {
			typeSizes[ts.Name.Name] = typeSizes[ts.Type.(*ast.Ident).Name]
		}

		return true
	})

	ast.Inspect(node, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		sortedField := sortFieldsBySize(structType)

		for _, field := range sortedField {
			fmt.Println(field.Names, field.Comment.Text())
		}

		return false
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
		sortedFields[i] = sortedField.Field

		sortedFields[i].Doc = sortedField.Doc
		sortedFields[i].Comment = sortedField.Comment
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
