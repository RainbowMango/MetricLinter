package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/prometheus/prometheus/util/promlint"
)

type metrics struct {
	name     string
	endpoint string
	file     string
	report   string
}

func RecordReport(report string, problems []promlint.Problem) {
	f, err := os.Create(report)
	if err != nil {
		panic(fmt.Sprintf("create file %s failed with error: %v", report, err))
	}

	for i := range problems {
		_, _ = f.WriteString(fmt.Sprintf("%s\n", problems[i].Text))
	}
}

func main() {
	ms := []metrics{
		{
			name:     "kube-apiserver",
			endpoint: "/metrics",
			file:     "data/apimetrics",
			report:   "report/kube-apiserver.log",
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

		RecordReport(m.report, problems)
		fmt.Printf("The problems number: %d\n", len(problems))
	}
}
