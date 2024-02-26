# testdata generator

This directory contains a small testdata generator used to generate "log files"
for running benchmarks for the example given in the repo's [README](../README.md).

# Usage

```shell
Usage testdata-generator:
  -min-size int
        Min file size to produce (default 1048576)
  -out string
        Name of the output file. If empty output will be written to STDOUT
```