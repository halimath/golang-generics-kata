# Testdata

This directory will contain testdata used for benchmarking the package.

This directory intentionally does not contain _all_ testdata when cloned from git
to save space in git. Use the [testdata generator](./cmd/testdata-generator) to
generate file.

Running `go generate .` in the repo's root will also generate testdata used for
benchmark tests.