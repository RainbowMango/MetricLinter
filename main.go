package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/prometheus/prometheus/util/promlint"
)

type metrics struct {
	name     string
	endpoint string
	file     string
}

func main() {
	ms := []metrics{
		{
			name:     "kube-apiserver",
			endpoint: "/metrics",
			file:     "data/apimetrics",
		},
	}

	for _, m := range ms {
		fmt.Printf("Linting %s metrics from endpoint: %s\n", m.name, m.endpoint)

		// Read metrics from file
		content, err := ioutil.ReadFile(m.file)
		if err != nil {
			fmt.Printf("Read file from %s failed with error: %s\n", m.file, err)
		}

		linter := promlint.New(strings.NewReader(string(content)))
		problems, err := linter.Lint()
		if err != nil {
			fmt.Printf("Lint failed with error: %s\n", err)
			fmt.Printf("Metrics content: %s\n", string(content))
			return
		}

		if len(problems) == 0 {
			continue
		}

		fmt.Printf("The problems are: \n")
		for i := range problems {
			fmt.Printf("%v\n", problems[i])
		}
	}
}
