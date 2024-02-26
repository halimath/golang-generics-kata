# golang generics kata

A coding kata to learn generics with golang and practice building a library to
enable functional programming style with common collections.

# Exercise

Build a library with example tests that provides datastructures and functions to
represent ordered sequences of elements with common operations in a functional 
way.

The datastructures should be able to support different sources of input, i.e.
* pre-allocated sources - ordered lists of elements, i.e. 
    slices or arrays with varying element types
* generator sources - i.e. finite and infinit ranges of integral types 
    (those that have a "natural" successor)
* streaming sources - i.e. elements that are read from i/o sources (files,
    network, databases) or from `chan`nels - either one at a time or in bulks

The lib should provide implementations for common operations on those 
datastructures, that somehow transform, filter or discard elements, i.e.

* `map` - apply a mapping function converting elements from on datastructure
    returning a new datastructure holding the converted elements
* `filter` - return a new datatype that contains only those elements matching
    a given _predicate_ function
* `find` - find all elements matching a given _predicate_ function
* `first` - find the first element to match a _predicate_ function
* `take` - only take the first `n` elements from a datastructure returning a
    new one
* `skip` - skip the first `n` elements from a datastructure

These operations should be _combineable_ (i.e. via chaining, composition, ...).

The lib should provide implementations that serve as a _sink_ for the datasources,
i.e.

* `reduce` - aggregate elements from a datastructure by applying an aggregator
    function one element at a time and return the final aggregation result
    (i.e. calculate the `sum` or `product` of numbers)
* `collect` - collect elements from the datastructure into a go slice

# Challenges

1. Create the library with approprate structure and tests (use 
    [Go example tests] to showcase your lib)
1. Make sure to implement at least _two_ different _sources_ for datastructure
    (i.e. slices, ranges)
1. Make sure to implement at least two transformer operations (i.e. `map` and 
    `filter`) and at least two consumer operations (i.e. `collect` and `reduce`)
1. Write benchmark tests to compare the time and memory utilization of your
    lib in contrast to an _idiomatic_ go implementation of the same algorithm.

## Advanced Challenges

You may take the following challenges into account when designing the library
and especially the libray's API. 

* How to handle operations that might fail (i.e. `map` with the ability to
    return an `error` in addition to the regular transformation result)? How
    should the overall operation behave if the mapping function fails to transform
    a single element?
* How to partition transformation work (i.e. for `map`) onto separate goroutines?
    How to ensure the results are recollected "in order"? How to ensure goroutines
    do not leak?
* How to improve performance for memory allocating operations (i.e. `collect`ing
    elements into a slice)? How can you handle the case that the length of some
    datastructures is known in advance while it is unknown for others?

## Benchmarkung examples

You can use the following examples as benchmarks:

1. Calculation on generated sequence
    1. Produce an infinte range of natural numbers starting at `x`. 
    1. Multiply each element of that range by `y`. 
    1. Select only those numbers that are even. 
    1. Calculate the sum of the first `z` of such numbers. 
    1. Use values like `x == 19, y == 3, z == 17`.

1. Parse, filter and aggregate a logfile of JSON-lines 
    1. Read a logfile of JSON-lines (each line is a self-contained JSON object).
    1. Parse every line as JSON into an appropriate datastructure (i.e. `map` or `struct`).
    1. Only consider those messages, that have
        * a `app` field that is not empty
        * a `severity` field containing either `error` or `warn`
    1. Aggregate statistics on number of (filtered) log events per `app`.
    * To run benchmarks on this example use the [testdata generator](./cmd/testdata-generator/)
        to generate files of configurable sizes (these are not added to this repo)

[Go example tests]: https://go.dev/blog/examples