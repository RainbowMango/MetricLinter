package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/prometheus/client_golang/prometheus/testutil/promlint"
)

type metrics struct {
	name     string
	endpoint string
	file     string
	report   string
	url      string
}

func RecordReport(report string, problems []promlint.Problem) {
	problemMap := make(map[string][]string)

	f, err := os.Create(report)
	if err != nil {
		panic(fmt.Sprintf("create file %s failed with error: %v", report, err))
	}

	for i := range problems {
		_, _ = f.WriteString(fmt.Sprintf("%s: %s\n", problems[i].Metric, problems[i].Text))
		problemMap[problems[i].Metric] = append(problemMap[problems[i].Metric], problems[i].Text)
	}

	_, _ = f.WriteString(fmt.Sprintf("\n\nTotal number of metrics with problems: %d\n", len(problemMap)))
	metricsNames := make([]string, 0, len(problemMap))
	for name, _ := range problemMap {
		metricsNames = append(metricsNames, name)
	}
	sort.SliceStable(metricsNames, func(i, j int) bool {
		return metricsNames[i] < metricsNames[j]
	})
	for _, m := range metricsNames {
		_, _ = f.WriteString(fmt.Sprintf("%s\n", m))
	}
}

func main() {
	ms := []metrics{
		{
			name:     "kube-apiserver",
			endpoint: "/metrics", // curl localhost:8080/metrics
			// url:      "http://localhost:8080/metrics",
			file:   "data/apimetrics",
			report: "report/kube-apiserver.log",
		},
		{
			name:     "kube-scheduler",
			endpoint: "/metrics", // curl localhost:10251/metrics
			// url:      "http://localhost:10251/metrics",
			file:   "data/kubescheduler",
			report: "report/kube-scheduler.log",
		},
		{
			name:     "kube-proxy",
			endpoint: "/metrics", // curl localhost:10249/metrics
			// url:      "http://localhost:10249/metrics",
			file:   "data/kubeproxy",
			report: "report/kube-proxy.log",
		},
		{
			name:     "kubelet-resource-v1alpha1",
			endpoint: "/metrics/resource/v1alpha1", // curl localhost:10255/metrics/resource/v1alpha1
			// url:      "http://localhost:10255/metrics/resource/v1alpha1",
			file:   "data/kubeletresourcev1alpha1",
			report: "report/kubelet-resource-v1alpha1.log",
		},
		{
			name:     "kubelet-resource",
			endpoint: "/metrics/resource", // curl localhost:10255/metrics/resource
			// url:      "http://localhost:10255/metrics/resource",
			file:   "data/kubeletresource",
			report: "report/kubelet-resource.log",
		},
		{
			name:     "kubelet-probes",
			endpoint: "/metrics/probes", // curl localhost:10255/metrics/probes
			// url:      "http://localhost:10255/metrics/probes",
			file:   "data/kubeletprobes",
			report: "report/kubelet-probes.log",
		},
		{
			name:     "kube-controller-manager",
			endpoint: "/metrics", // curl localhost:10252/metrics
			// url:      "http://localhost:10252/metrics",
			file:   "data/kubecontrollermanager",
			report: "report/kube-controller-manager.log",
		},
		{
			name:     "cloud-controller-manager",
			endpoint: "/metrics", // curl localhost:10253/metrics (not available on local cluster started by hack/local-cluster.sh)
			file:     "data/cloudcontrollermanager",
			report:   "report/cloud-controller-manager.log",
		},
		{
			name:     "karmada-scheduler",
			endpoint: "/metrics", // curl localhost:10253/metrics (not available on local cluster started by hack/local-cluster.sh)
			file:     "data/karmada-scheduler",
			report:   "report/karmada-cheduler.log",
		},
	}

	for _, m := range ms {
		if len(m.url) > 0 {
			fmt.Printf("grabbing metrics from %s\n", m.url)
			resp, err := http.Get(m.url)
			if err != nil {
				fmt.Printf("grabbing metrics from %s failed as %v\n", m.url, err)
				continue
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Read response body from %s failed as %v\n", m.url, err)
				continue
			}

			if err = ioutil.WriteFile(m.file, body, os.ModePerm); err != nil {
				fmt.Printf("Update data file failed for %s as %v\n", m.name, err)
			}
		}

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
	}
}
