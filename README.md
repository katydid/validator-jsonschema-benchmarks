# JSON Schema Benchmark

This benchmark suite builds on the amazing work of [SourceMeta's JSONSchema Benchmarks](https://github.com/sourcemeta-research/jsonschema-benchmark/) for the purpose of benchmarking Katydid's JSONSchema implementation against other JSONSchema validators.

Each validator is run with multiple schemas and a collection of documents that might be valid or invalid, which is expressed in the suffix of the folder name.

## Results

The most recent results can be seen [via GitHub Actions](https://github.com/katydid/validator-jsonschema-benchmarks/actions/workflows/ci.yml).

## Setup

The benchmarks require:

* posix tools: make, sed, printf
* [Docker](https://www.docker.com/)

## Running

There are several ways to run the benchmarks:

* Run all benchmarks and produce a report: `make run` or `make dist/report.csv`
* Run only specific implementations: `make IMPLEMENTATIONS='blaze jsoncons' RUNS=5`
* Run only specific schemas: `make SCHEMAS='example-address-valid'`

## Analytics

Analytics currently requires Go to be installed:

* `make analytics-schemas-latex` creates a latex file `schemas.latex` with analytics on the schemas.
* `make analytics-results-latex` creates a latex file `results.latex` with analytics on the results after running per schema and implementation.
* `make analytics-impls-md` creates a markdown file `impls.md` with analytics on the results after running per implementation.

## Implementations

All implementations can be found in the `implementations/` subdirectory.
A summary of these implementations is given below:

- [ajv](./implementations/ajv/) (JS) (returns bool)
- [ajv-bun](./implementations/ajv-bun/) (JS with BUN runtime) (returns bool)
- [blaze](./implementations/blaze/) (C++) (returns bool)
- [boon](./implementations/boon/) (Rust) (returns Result)
- [corvus](./implementations/corvus/) (generated C#) (returns Result)
- [go-google](./implementations/go-google/) (Go) (returns error) (ignored)
- [go-json-schema-spec](./implementations/go-json-schema-spec/) (Go) (returns error) (ignored)
- [go-kaptinlin](./implementations/go-kaptinlin/) (Go) (returns error)
- [go-katydid-auto-json](./implementations/go-katydid-auto-json/) (Go)  (returns bool)
- [go-katydid-auto-reflect](./implementations/go-katydid-auto-reflect/) (Go)  (returns bool)
- [go-katydid-mem-json](./implementations/go-katydid-mem-json/) (Go) (returns bool)
- [go-katydid-mem-reflect](./implementations/go-katydid-mem-reflect/) (Go) (returns bool)
- [go-santhosh-tekuri](./implementations/go-santhosh-tekuri/) (Go) (returns error)
- [hyperjump](./implementations/hyperjump/) (JS) (returns Result)
- [jsdotnet](./implementations/jsdotnet/) (C#) (returns Result)
- [json_schemer](./implementations/json_schemer/) (Ruby) (returns bool)
- [jsoncons](./implementations/jsoncons/) (C++) (throws Exception)
- [jsu-c](./implementations/jsu-c/) (generated C) (returns bool, optionally can report errors)
- [jsu-java](./implementations/jsu-java/) (generated Java) (returns bool)
- [jsu-js](./implementations/jsu-js/) (generated JS) (returns bool, optionally can report errors)
- [jsu-pl](./implementations/jsu-pl/) (generated Perl) (ignored)
- [jsu-py](./implementations/jsu-py/) (generated Python) (returns bool, optionally can report errors)
- [JSV](./implementations/jsv) (Elixir) (returns Result)
- [kmp](./implementations/kmp) (Kotlin) (returns bool, optionally can return report errors)
- [networknt](./implementations/networknt/) (Java) (returns bool, optionally can return result)
- [opis](./implementations/opis/) (PHP) (returns Result)
- [py-jsonschema](./implementations/py-jsonschema/) (Python) (returns bool)
- [rapidjson](./implementations/rapidjson/) (C++) (returns bool) (ignored)
- [schemasafe](./implementations/schemasafe/) (JS) (returns bool) (ignored)

Compared to the original [SourceMeta's JSONSchema Benchmarks](https://github.com/sourcemeta-research/jsonschema-benchmark/) the following libraries were added: 

- [go-google](./implementations/go-google/) (Go) (ignored)
- [go-json-schema-spec](./implementations/go-json-schema-spec/) (Go) (ignored)
- [go-kaptinlin](./implementations/go-kaptinlin/) (Go)
- [go-katydid-auto-json](./implementations/go-katydid-auto-json/) (Go)
- [go-katydid-mem-json](./implementations/go-katydid-mem-json/) (Go)
- [go-katydid-auto-reflect](./implementations/go-katydid-auto-reflect/) (Go)
- [go-katydid-mem-reflect](./implementations/go-katydid-mem-reflect/) (Go)
- [rapidjson](./implementations/rapidjson/) (C++) (ignored)

Each implementation is run via Docker.
First, a Docker container is built with all the necessary dependencies.
Then, at runtime, a folder containing the schema and the necessary dependencies is mounted and the time to validate all documents is measured.

Implementations can be ignored by adding a `.benchmark-ignore` file in the implementation subdirectory.
It also worth noting that some implementations compile schemas ahead of time into a more efficient representation, while others interpret the entire schema at runtime.

### Adding a new implementation

First, each implementation must have a `Dockerfile` that copies in any necessary scripts and installs dependencies.
There is also a `version.sh` script that must output the version of the implementation (often extracted from whatever dependency management tool is used).
Finally, appropriate targets must be added to the `Makefile` to build the Docker container and run the benchmark.

## Schemas

All schemas are found in the schemas folder.
Some schemas had a collection of instances gathered from github, which are mutated to create invalid instances via `make mutate`.
The rest we can regenerate valid and invalid instances for by running: `make generate`.

We run a curated list of [curated_schemas.txt](./curated_schemas.txt) where:
* schemas with `uniqueItems` have been removed and replaced with schemes where `uniqueItems` have been removed from the schema,
* we have totally excluded schemas that use `dynamicRef`: `cql2` and `openapi`.
* we have removed `cspell`, since the following implementations all have problems with it: boon, go-kaptinlin, go-santhosh-tekuri, json_schemer and kmp. One problem is for using Perl syntax regexes `(?=`.
* we have removed `ui5-manifest`, since the following implementations all have problems with it: ajv, ajv-bun, boon, go-kaptinlin, go-santhosh-tekuri, hyperjump, networknt. One problem is for using Perl syntax regexes `(?=`.
* we have removed `krakend`, since the following implementations all have problems with it: ajv, ajv-bun, hyperjump, networknt.  One problem is for `Invalid regular expression: /^\/[^\*\?\&\%]*(\/\*)?$/u`.
* we have removed `helm-chart-lock` since the following implementations all have problems with it: [jsu-py, jsu-java](https://github.com/clairey-zx81/json-model/issues/5), [opis](https://github.com/opis/json-schema/issues/166), [jsdotnet](https://github.com/json-everything/json-everything/issues/1043).

### Shortcomings of document generator

* Lots of libraries do not support unicode strings properly, so random generation of documents has been limited to ascii strings.
* Since we using JSONL for documents, we do not generate newlines in strings.
* Number generation has been limited to 64 bit floats.