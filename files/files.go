package files

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type Work struct {
	file    string
	pattern *regexp.Regexp
	result  chan Result
}

type Result struct {
	file       string
	lineNumber int
	text       string
}

func worker(jobs chan Work) {
	for work := range jobs {
		f, err := os.Open(work.file)
		if err != nil {
			fmt.Println(err)
			continue
		}
		scn := bufio.NewScanner(f)
		lineNumber := 1
		for scn.Scan() {
			result := work.pattern.Find(scn.Bytes())
			if len(result) > 0 {
				work.result <- Result{
					file:       work.file,
					lineNumber: lineNumber,
					text:       string(result),
				}
			}
			lineNumber++
		}
		if err := f.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
		close(work.result)
	}
}

func ReadDir(filePath string, pattern string) {
	jobs := make(chan Work)
	wg := sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(jobs)
		}()
	}

	rex, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}

	allResults := make(chan chan Result)
	go func() {
		defer close(allResults)
		err := filepath.Walk(filePath, func(path string, d fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				if strings.Contains(path, "node_modules") || strings.Contains(path, "vendor") {
					return nil
				}

				ch := make(chan Result)
				jobs <- Work{file: path, pattern: rex, result: ch}
				allResults <- ch
			}
			return nil
		})

		if err != nil {
			panic("Failed search file")
		}
	}()

	for resultCh := range allResults {
		for result := range resultCh {
			fmt.Printf("%s #%d: %s\n", result.file, result.lineNumber, result.text)
		}
	}

	close(jobs)
	wg.Wait()
}
