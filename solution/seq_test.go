package seq_test

import (
	"fmt"

	"github.com/halimath/coding-katas/seq"
)

func ExampleForEach() {
	err := seq.ForEach(
		seq.Take(
			seq.InfiniteIntRange(1, 1),
			10,
		),
		func(i int) { fmt.Println(i) },
	)
	fmt.Println(err)

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
	// 10
	// <nil>
}

func ExampleReduce() {
	sum, err := seq.Reduce(
		seq.Take(
			seq.InfiniteIntRange(1, 1),
			10,
		),
		func(sum, v int) int { return sum + v },
	)
	fmt.Println(sum, err)

	// Output:
	// 55 <nil>
}

func Example() {
	r, err := seq.Reduce(
		seq.Take(
			seq.Filter(
				seq.Map(
					seq.InfiniteIntRange(19, 3),
					func(i int) int { return i * 3 },
				),
				func(i int) bool { return i%2 == 0 },
			),
			17,
		),
		func(sum, i int) int { return sum + i },
	)
	fmt.Println(r, err)

	// Output: 3570 <nil>
}
