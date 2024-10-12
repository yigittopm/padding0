/*
Copyright © 2024 Mert Yiğittop <yigittopm@hotmail.com>
*/
package analysis

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
)

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

	getFiles(dir)
	return nil
}

func getFiles(dir string) {
	var wg sync.WaitGroup
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Dir not found: %v", err)
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			wg.Add(1)
			go func() {
				defer wg.Done()
				readFile(path)
			}()
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the path: %v", err)
	}

	wg.Wait()
}

func readFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		if err.Error() != "EOF" {
			log.Fatalf("Error opening file: %v", err)
		}
	}

	modifiedContent := ExtractStruct(string(file))
	writeFile(path, modifiedContent)
}

func writeFile(path string, content string) {
	content = strings.TrimSpace(content)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}
}

func ExtractStruct(content string) string {
	re := regexp.MustCompile(`type\s+(\w+)\s+struct\s*{([^}]*)}`)
	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {

			original := match[0]
			modified := ProcessStruct(match)
			content = strings.Replace(content, original, modified, 1)
		}
	}
	return content
}

type Field struct {
	Name    string
	Type    string
	Comment string
	Size    int
}

func ProcessStruct(match []string) string {
	structName := match[1]
	structBody := match[2]
	lines := strings.Split(structBody, "\n")

	var fields []Field
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := regexp.MustCompile(`(\w+)\s+(\w+)\s*(//.*)?`).FindStringSubmatch(line)
		if len(parts) > 2 {
			fieldName := parts[1]
			fieldType := parts[2]
			fieldComment := ""
			if len(parts) > 3 {
				fieldComment = parts[3]
			}

			fields = append(fields, Field{
				Name:    fieldName,
				Type:    fieldType,
				Comment: fieldComment,
				Size:    typeSizes[fieldType],
			})
		}

	}

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Size > fields[j].Size
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("type %s struct {\n", structName))
	for _, field := range fields {
		sb.WriteString(fmt.Sprintf("\t\t%s %s %s\n", field.Name, field.Type, field.Comment))
	}
	sb.WriteString("\t}")
	return sb.String()
}
