//go:generate go run ../cmd/testdata-generator/main.go -min-size 7340032 -out ./testdata/benchmark.s.jsonl
//go:generate go run ../cmd/testdata-generator/main.go -min-size 73400320 -out ./testdata/benchmark.m.jsonl
//go:generate go run ../cmd/testdata-generator/main.go -min-size 734003200 -out ./testdata/benchmark.l.jsonl
package seq_test

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/halimath/coding-katas/seq"
	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

type logEntry struct {
	App      string `json:"app"`
	Severity string `json:"severity"`
}

func calcMessageStats_seq(filename string) (map[string]int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return seq.Reduce(
		seq.Filter(
			seq.MapE(
				seq.Scanner(bufio.NewScanner(f)),
				func(b []byte) (e logEntry, err error) {
					err = json.Unmarshal(b, &e)
					return
				},
			),
			func(e logEntry) bool {
				if e.App == "" {
					return false
				}

				if e.Severity == "error" {
					return true
				}

				if e.Severity == "warn" {
					return true
				}

				return false
			},
		),
		func(stats map[string]int, e logEntry) map[string]int {
			stats[e.App] = stats[e.App] + 1
			return stats
		},
		make(map[string]int),
	)
}

func TestCalcMessageStats_seq(t *testing.T) {
	s, err := calcMessageStats_seq("./testdata/test.jsonl")
	expect.That(t,
		is.NoError(err),
		is.DeepEqualTo(s, map[string]int{
			"auth":        132,
			"service":     130,
			"repository":  159,
			"api_gateway": 137,
		}),
	)
}

func BenchmarkCalcMessageStats_seq_small(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := calcMessageStats_seq("./testdata/benchmark.s.jsonl")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCalcMessageStats_seq_medium(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := calcMessageStats_seq("./testdata/benchmark.m.jsonl")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCalcMessageStats_seq_large(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := calcMessageStats_seq("./testdata/benchmark.l.jsonl")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func calcMessageStats_idiomatic(filename string) (map[string]int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stats := make(map[string]int)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var e logEntry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			return nil, err
		}

		if e.App == "" {
			continue
		}

		if e.Severity != "error" && e.Severity != "warn" {
			continue
		}

		stats[e.App] = stats[e.App] + 1
	}

	return stats, nil
}

func TestCalcMessageStats_idiomatic(t *testing.T) {
	s, err := calcMessageStats_idiomatic("./testdata/test.jsonl")
	expect.That(t,
		is.NoError(err),
		is.DeepEqualTo(s, map[string]int{
			"auth":        132,
			"service":     130,
			"repository":  159,
			"api_gateway": 137,
		}),
	)
}

func BenchmarkCalcMessageStats_idiomatic_small(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := calcMessageStats_idiomatic("./testdata/benchmark.s.jsonl")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCalcMessageStats_idiomatic_medium(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := calcMessageStats_idiomatic("./testdata/benchmark.m.jsonl")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCalcMessageStats_idiomatic_large(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := calcMessageStats_idiomatic("./testdata/benchmark.l.jsonl")
		if err != nil {
			b.Fatal(err)
		}
	}
}
