package seq_test

import (
	"testing"

	"github.com/halimath/coding-katas/seq"
	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func example1_seq() (int, error) {
	return seq.Reduce(
		seq.Take(
			seq.Filter(
				seq.Map(
					seq.InfiniteIntRange(19, 1),
					func(i int) int { return i * 3 },
				),
				func(i int) bool { return i%2 == 0 },
			),
			17,
		),
		func(sum, i int) int { return sum + i },
	)
}

func TestExample1_seq(t *testing.T) {
	c, err := example1_seq()
	expect.That(t,
		is.NoError(err),
		is.EqualTo(c, 1836),
	)
}

func BenchmarkExample1_seq(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := example1_seq()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func example1_idiomatic() int {
	count := 0
	i := 19
	sum := 0

	for {
		v := i * 3
		if v%2 == 0 {
			sum += v
			count++

			if count == 17 {
				return sum
			}
		}

		i++
	}
}

func TestExample1_idiomatic(t *testing.T) {
	c := example1_idiomatic()
	expect.That(t,
		is.EqualTo(c, 1836),
	)
}

func BenchmarkExample1_idiomatic(b *testing.B) {
	for n := 0; n < b.N; n++ {
		example1_idiomatic()
	}
}
