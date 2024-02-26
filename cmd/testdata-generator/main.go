package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"time"
)

var (
	minSize = flag.Int("min-size", 1024*1024, "Min file size to produce")
	outFile = flag.String("out", "", "Name of the output file. If empty output will be written to STDOUT")
)

func main() {
	flag.Parse()

	out := os.Stdout
	if *outFile != "" {
		of, err := os.OpenFile(*outFile, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer of.Close()
		out = of
	}

	var count int
	start := time.Now()

	for count < *minSize {
		b, err := json.Marshal(logEvent{
			Time:     start,
			App:      app(),
			Severity: severity(),
			Message:  message(),
		})

		if err != nil {
			panic(err)
		}

		s := string(b)
		fmt.Fprintln(out, s)
		count += len(s)

		start = start.Add(duration())
	}
}

const loremIpsum = `Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et e`

func message() string {
	return loremIpsum[:rand.IntN(len(loremIpsum)-50)+50]
}

var severities = []string{"debug", "info", "warn", "error"}

func severity() string {
	return severities[rand.IntN(len(severities))]
}

var apps = []string{"", "service", "api_gateway", "auth", "repository"}

func app() string {
	return apps[rand.IntN(len(apps))]
}

type logEvent struct {
	Time     time.Time `json:"timestamp"`
	App      string    `json:"app"`
	Severity string    `json:"severity"`
	Message  string    `json:"msg"`
}

func duration() time.Duration {
	return time.Millisecond * time.Duration(rand.IntN(20))
}
