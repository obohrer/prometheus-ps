package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	ps "github.com/mitchellh/go-ps"
)

// globals
var cfg Configuration

func sanitizeName(s string) string {
	return strings.Replace(strings.Replace(strings.ToLower(s), "-", "_", -1), " ", "_", -1)
}

func buildMetricName(executable string) string {
	var buffer bytes.Buffer
	buffer.WriteString("process_")
	buffer.WriteString(sanitizeName(executable))
	buffer.WriteString("_up")
	return buffer.String()
}

func writeProcessesMetrics(w http.ResponseWriter, ps map[string]int) {
	for process, c := range ps {
		metricName := buildMetricName(process)
		fmt.Fprintf(w, "# HELP %s The process %s is up\n", metricName, metricName)
		fmt.Fprintf(w, "# TYPE %s gauge\n", metricName)
		fmt.Fprintf(w, "%s %d\n", metricName, c)
	}
}

func groupByName(pl []ps.Process, watchList map[string]Void) map[string]int {
	processFreq := make(map[string]int)
	for _, p := range pl {
		exec := p.Executable()
		_, present := watchList[exec]
		if present {
			processFreq[exec]++
		}
	}
	return processFreq
}

func startServer(config Configuration) {
	var listen = fmt.Sprintf(":%d", config.Port)
	s := &http.Server{
		Addr:           listen,
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		ps, err := ps.Processes()
		if err != nil {
			log.Fatal(err, "Unexpected error")
		} else {
			writeProcessesMetrics(w, groupByName(ps, config.Wl))
		}
	})

	fmt.Printf("Starting server on %s ...\n", listen)
	s.ListenAndServe()
}

func showUsage() {
	fmt.Printf("Usage : prometheus_ps --config <path-to-conf.json>\n")
}

func main() {
	args := os.Args[1:]
	if len(args) <= 1 || args[0] != "--config" {
		showUsage()
	}
	confLoc := args[1]

	config := ReadConfig(confLoc)
	startServer(config)
}
