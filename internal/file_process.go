package analysis

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
			log.Printf("Dir not found: %v", err)
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
	errCh := make(chan error, 10)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), "example.go") && !strings.HasSuffix(info.Name(), "_test.go") {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := readFile(path); err != nil {
					errCh <- err
				}
			}()
		}

		return nil
	})

	if err != nil {
		return err
	}

	wg.Wait()
	close(errCh)

	for e := range errCh {
		log.Printf("Error processing file: %v", e)
	}

	return nil
}

func readFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		if err.Error() != "EOF" {
			log.Printf("Error opening file: %v", err)
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
	fmt.Println(content)
	/*err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}

	err = formatTheWrittenFile(path)
	if err != nil {
		return err
	}*/

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
