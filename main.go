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
			endpoint: "/metrics", // curl localhost:8080/metrics
			file:     "data/apimetrics",
			report:   "report/kube-apiserver.log",
		},
		{
			name:     "kube-scheduler",
			endpoint: "/metrics", // curl localhost:10251/metrics
			file:     "data/kubescheduler",
			report:   "report/kube-scheduler.log",
		},
		{
			name:     "kube-proxy",
			endpoint: "/metrics", // curl localhost:10249/metrics
			file:     "data/kubeproxy",
			report:   "report/kube-proxy.log",
		},
		{
			name:     "kubelet-resource",
			endpoint: "/metrics/resource/v1alpha1", // curl localhost:10255/metrics/resource/v1alpha1
			file:     "data/kubeletresource",
			report:   "report/kubelet-resource.log",
		},
		{
			name:     "kubelet-probes",
			endpoint: "/metrics/probes", // curl localhost:10255/metrics/probes
			file:     "data/kubeletprobes",
			report:   "report/kubelet-probes.log",
		},
		{
			name:     "kube-controller-manager",
			endpoint: "/metrics", // curl localhost:10252/metrics
			file:     "data/kubecontrollermanager",
			report:   "report/kube-controller-manager.log",
		},
		{
			name:     "cloud-controller-manager",
			endpoint: "/metrics", // curl localhost:10253/metrics (not available on local cluster started by hack/local-cluster.sh)
			file:     "data/cloudcontrollermanager",
			report:   "report/cloud-controller-manager.log",
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
