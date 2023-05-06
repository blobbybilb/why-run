package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const defaultValue = "__default__"

func main() {
	helpFlag := flag.Bool("h", false, "display help")
	flag.Parse()

	if *helpFlag {
		help()
		return
	}

	config, err := parse("why.run")
	if err != nil {
		panic(err)
	}

	if len(flag.Args()) > 0 {
		run(config, flag.Args()[0])
	} else {
		run(config, defaultValue)
	}

	if err != nil {
		panic(err)
	}
}

func parse(path string) (map[string][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	result := map[string][]string{}

	currentKey := defaultValue
	currentValue := []string{}

	for scanner.Scan() {
		line := scanner.Text()

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		isKey := strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]")

		if isKey {
			result[currentKey] = currentValue
			currentValue = []string{}

			currentKey = line[1 : len(line)-1]
		} else {
			currentValue = append(currentValue, line)
		}
	}

	result[currentKey] = currentValue

	return result, nil
}

func run(config map[string][]string, task string) {
	tasks := config[task]
	wg := sync.WaitGroup{}

	for _, t := range tasks {
		if strings.HasPrefix(t, "- ") {
			t = t[2:]
			fmt.Println("Running (async):", t)
			wg.Add(1)
			go func() {
				exec.Command("bash", "-c", t).Run()
				wg.Done()
			}()
		} else {
			fmt.Println("Running:", t)
			exec.Command("bash", "-c", t).Run()

		}
	}

	wg.Wait()
}

func help() {
	fmt.Println("Soon...")
}
