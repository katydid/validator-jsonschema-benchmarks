# JSON Schema Benchmark

This benchmark suite builds on the amazing work of [SourceMeta's JSONSchema Benchmarks](https://github.com/sourcemeta-research/jsonschema-benchmark/) for the purpose of benchmarking Katydid's JSONSchema implementation against other JSONSchema validators.
Each validator is run with multiple schemas and a collection of documents that might be valid or invalid.

## Results

See the [most recent recorded results](./recorded/2026-09-02/Readme.md).

We record results in the [recorded](./recorded/Readme.md) folder.

Results can also be found [via GitHub Actions](https://github.com/katydid/validator-jsonschema-benchmarks/actions/workflows/ci.yml).

## Setup

The benchmarks require:

* posix tools: make, sed, printf
* [Docker](https://www.docker.com/)

## Running

There are several ways to run the benchmarks:

* Run all benchmarks and produce a report: `rm -rf dist; RUNS=3 make run`.
* Run only specific implementations: `make IMPLEMENTATIONS='blaze jsoncons' RUNS=5`
* Run only specific schemas: `make SCHEMAS='example-address-valid'`

## Analytics

Analytics currently requires Go to be installed:

* `make analytics-impls-latex` creates a latex file `impls.latex` with analytics on the results after running per implementation.
* `make analytics-schemas-md` creates a markdown file `schemas.md` with analytics on the schemas.

## Implementations

Outside of Katydid, the benchmarks include 18 implementations in various programming languages:
* Ajv: The popular Javascript implementation Ajv, including an implementation executed on the BUN~\cite{javascript_engine_bun} Javascript runtime for extra speed;
* Blaze: A highly optimized C++ implementation that uses a custom field name hash function for faster string comparison;
* JSON Schema Utils (JSU) is a tool that generates code in several target languages: C, Java, Javascript or Python to validate a specific schema (we left out the Perl version);
* Corvus.JsonSchema: A tool that generates C# for a specific schema;
* and libraries in other languages, including: Go, Ruby, Elixir, Kotlin, and PHP.

We exclude libraries that:
* do not support `"format"`.
* fails to validate three or more schemas, do not respond to error reports, and are not in contention for winning the benchmark.
See ignore below, to see exactly which implementations are currently excluded.

All implementations can be found in the [implementations](./implementations/) subdirectory.
A summary of these implementations is given below:
- [ajv](./implementations/ajv/) (Javascript) (returns bool)
- [ajv-bun](./implementations/ajv-bun/) (Javascript with BUN runtime) (returns bool)
- [blaze](./implementations/blaze/) (C++) (returns bool)
- [boon](./implementations/boon/) (Rust) (returns Result)
- [corvus](./implementations/corvus/) (generated C#) (returns Result)
- [go-google](./implementations/go-google/) (Go) (returns error) ([ignored](./implementations/go-google/.benchmark-ignore))
- [go-json-schema-spec](./implementations/go-json-schema-spec/) (Go) (returns error) ([ignored](./implementations/go-json-schema-spec/.benchmark-ignore))
- [go-kaptinlin](./implementations/go-kaptinlin/) (Go) (returns error)
- [go-katydid-auto-json](./implementations/go-katydid-auto-json/) (Go)  (returns bool)
- [go-katydid-auto-reflect](./implementations/go-katydid-auto-reflect/) (Go)  (returns bool)
- [go-katydid-mem-json](./implementations/go-katydid-mem-json/) (Go) (returns bool)
- [go-katydid-mem-reflect](./implementations/go-katydid-mem-reflect/) (Go) (returns bool)
- [go-santhosh-tekuri](./implementations/go-santhosh-tekuri/) (Go) (returns error)
- [hyperjump](./implementations/hyperjump/) (Javascript) (returns Result)
- [jsdotnet](./implementations/jsdotnet/) (C#) (returns Result)
- [json_schemer](./implementations/json_schemer/) (Ruby) (returns bool)
- [jsoncons](./implementations/jsoncons/) (C++) (throws Exception)
- [jsu-c](./implementations/jsu-c/) (generated C) (returns bool, optionally can report errors)
- [jsu-java](./implementations/jsu-java/) (generated Java) (returns bool)
- [jsu-js](./implementations/jsu-js/) (generated Javascript) (returns bool, optionally can report errors)
- [jsu-pl](./implementations/jsu-pl/) (generated Perl) ([ignored](./implementations/jsu-pl/.benchmark-ignore))
- [jsu-py](./implementations/jsu-py/) (generated Python) (returns bool, optionally can report errors)
- [JSV](./implementations/jsv) (Elixir) (returns Result)
- [kmp](./implementations/kmp) (Kotlin) (returns bool, optionally can return report errors)
- [networknt](./implementations/networknt/) (Java) (returns bool, optionally can return result)
- [opis](./implementations/opis/) (PHP) (returns Result) ([ignored](./implementations/opis/.benchmark-ignore))
- [py-jsonschema](./implementations/py-jsonschema/) (Python) (returns bool)
- [rapidjson](./implementations/rapidjson/) (C++) (returns bool) ([ignored](./implementations/rapidjson/.benchmark-ignore))
- [schemasafe](./implementations/schemasafe/) (Javascript) (returns bool) ([ignored](./implementations/schemasafe/.benchmark-ignore))

Each implementation is warmed up by first validating all the documents for over a hundred iterations.
This implies the CPU caches are warmed up, JIT has been triggered, and also that our memoized solution has a populated table.
This is representative of the real-world situations of an API gateway or filtering through petabytes of records.

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

[SourceMeta's JSONSchema Benchmarks](https://github.com/sourcemeta-research/jsonschema-benchmark/) originally collected schemas from the JSONSchema Store and then crawled GitHub for valid JSON documents for each of these schemas.
We developed a mutator that uses these valid instances as input to generate invalid instances, where one object field or array item is mutated, added, or removed, and the [go-santhosh-tekuri](./implementations/go-santhosh-tekuri/) library checks that it is truly invalid according to the schema.
We added schemas from the official [JSONSchema examples](https://json-schema.org/learn/json-schema-examples\#blog-post): (address, blogpost, calendar, devicetype, health-record, job-posting, movie, userprofile) and Ajv's benchmarks (cosmicrealms, complex, medium, advanced, basic) and also include our own [conference filtering example](./schemas/katydid-conf-valid/).
To evaluate the extended schema collection, we developed a random generator capable of producing both valid instances and non-conforming (invalid) documents containing at least one schema violation.

We run a curated list of [curated_schemas.txt](./curated_schemas.txt) where:
* Schemas with `uniqueItems` have been removed and replaced with schemes where `uniqueItems` have been removed from the schema, because katydid cannot handle uniqueItems, this allows for isolated benchmarking of the remaining schema logic;
* We have totally excluded schemas that use `dynamicRef`, which are the `cql2` and `openapi` schemas.
* We have removed `cspell`, since the following implementations all have problems with it: boon, go-kaptinlin, go-santhosh-tekuri, json_schemer and kmp. One problem is for using Perl syntax regexes `(?=`.
* We have removed `ui5-manifest`, since the following implementations all have problems with it: ajv, ajv-bun, boon, go-kaptinlin, go-santhosh-tekuri, hyperjump, networknt. One problem is for using Perl syntax regexes `(?=`.
* We have removed `krakend`, since the following implementations all have problems with it: ajv, ajv-bun, hyperjump, networknt.  One problem is for `Invalid regular expression: /^\/[^\*\?\&\%]*(\/\*)?$/u`.
* We have removed `helm-chart-lock` since the following implementations all have problems with it: [jsu-py, jsu-java](https://github.com/clairey-zx81/json-model/issues/5), [opis](https://github.com/opis/json-schema/issues/166), which might be because it uses invalid uri's that are missing an authority `file://../../`, see explanation [here](https://github.com/json-everything/json-everything/issues/1043#issuecomment-5038687635).

In total, we have 45 schemas, each of which has a valid and an invalid data set.
All schemas are found in the [schemas](./schemas/) folder.

|#|Dataset name|Source|Valid/InValid|uniqueItems|# Docs|Schema Size (KB)|Mean Doc. Size (B)|
|---|---|---|---|---|---|---|---|
| 1 | ajv-cosmicrealms-rmUniqueItems-invalid | generated | invalid | removed | 10000 | 1583 | 718 |
| 2 | ajv-cosmicrealms-rmUniqueItems-valid | generated | valid | removed | 10000 | 1583 | 492 |
| 3 | ansible-meta | collected | valid | none | 333 | 37015 | 312 |
| 4 | ansible-meta-invalid | mutated | invalid | none | 333 | 37015 | 353 |
| 5 | aws-cdk | collected | valid | none | 483 | 740 | 1149 |
| 6 | aws-cdk-invalid | mutated | invalid | none | 483 | 740 | 517 |
| 7 | babelrc | collected | valid | none | 794 | 6717 | 140 |
| 8 | babelrc-invalid | mutated | invalid | none | 794 | 6717 | 241 |
| 9 | clang-format | collected | valid | none | 133 | 55557 | 334 |
| 10 | clang-format-invalid | mutated | invalid | none | 133 | 55557 | 441 |
| 11 | cmake-presets | collected | valid | none | 967 | 86084 | 2719 |
| 12 | cmake-presets-invalid | mutated | invalid | none | 967 | 86084 | 2637 |
| 13 | code-climate | collected | valid | none | 2484 | 6069 | 281 |
| 14 | code-climate-invalid | mutated | invalid | none | 2484 | 6069 | 381 |
| 15 | cypress | collected | valid | none | 981 | 16457 | 403 |
| 16 | cypress-invalid | mutated | invalid | none | 981 | 16457 | 250 |
| 17 | deno-rmUniqueItems | collected | valid | removed | 987 | 22822 | 1022 |
| 18 | deno-rmUniqueItems-invalid | mutated | invalid | removed | 987 | 22822 | 540 |
| 19 | dependabot | collected | valid | none | 967 | 9631 | 403 |
| 20 | dependabot-invalid | mutated | invalid | none | 967 | 9631 | 435 |
| 21 | draft-04-rmUniqueItems | collected | valid | removed | 561 | 4020 | 12657 |
| 22 | draft-04-rmUniqueItems-invalid | mutated | invalid | removed | 561 | 4020 | 11430 |
| 23 | example-address-invalid | generated | invalid | none | 10000 | 761 | 135 |
| 24 | example-address-valid | generated | valid | none | 10000 | 761 | 142 |
| 25 | example-blogpost-invalid | generated | invalid | none | 10000 | 1265 | 511 |
| 26 | example-blogpost-valid | generated | valid | none | 10000 | 1265 | 258 |
| 27 | example-calendar-invalid | generated | invalid | none | 10000 | 1499 | 658 |
| 28 | example-calendar-valid | generated | valid | none | 10000 | 1499 | 218 |
| 29 | example-devicetype-invalid | generated | invalid | none | 10000 | 1457 | 428 |
| 30 | example-devicetype-valid | generated | valid | none | 10000 | 1457 | 95 |
| 31 | example-health-record-invalid | generated | invalid | none | 10000 | 1509 | 678 |
| 32 | example-health-record-valid | generated | valid | none | 10000 | 1509 | 324 |
| 33 | example-job-posting-invalid | generated | invalid | none | 10000 | 674 | 424 |
| 34 | example-job-posting-valid | generated | valid | none | 10000 | 674 | 145 |
| 35 | example-movie-invalid | generated | invalid | none | 10000 | 708 | 507 |
| 36 | example-movie-valid | generated | valid | none | 10000 | 708 | 140 |
| 37 | example-userprofile-invalid | generated | invalid | none | 10000 | 631 | 451 |
| 38 | example-userprofile-valid | generated | valid | none | 10000 | 631 | 135 |
| 39 | fabric-mod | collected | valid | none | 911 | 11446 | 691 |
| 40 | fabric-mod-invalid | mutated | invalid | none | 911 | 11446 | 744 |
| 41 | geojson | collected | valid | none | 500 | 46177 | 52329 |
| 42 | geojson-invalid | mutated | invalid | none | 500 | 46177 | 52401 |
| 43 | gitpod-configuration | collected | valid | none | 986 | 13489 | 357 |
| 44 | gitpod-configuration-invalid | mutated | invalid | none | 986 | 13489 | 436 |
| 45 | importmap | collected | valid | none | 964 | 618 | 629 |
| 46 | importmap-invalid | mutated | invalid | none | 964 | 618 | 724 |
| 47 | jasmine | collected | valid | none | 980 | 3690 | 133 |
| 48 | jasmine-invalid | mutated | invalid | none | 980 | 3690 | 203 |
| 49 | jsck-complex-invalid | generated | invalid | none | 100 | 4060 | 25892 |
| 50 | jsck-complex-valid | generated | valid | none | 100 | 4060 | 29874 |
| 51 | jsck-medium-invalid | generated | invalid | none | 10000 | 1887 | 476 |
| 52 | jsck-medium-valid | generated | valid | none | 10000 | 1887 | 313 |
| 53 | jsconfig-rmUniqueItems | collected | valid | removed | 981 | 60482 | 177 |
| 54 | jsconfig-rmUniqueItems-invalid | mutated | invalid | removed | 981 | 60482 | 272 |
| 55 | jshintrc | collected | valid | none | 966 | 12142 | 428 |
| 56 | jshintrc-invalid | mutated | invalid | none | 966 | 12142 | 490 |
| 57 | katydid-conf-invalid | generated | invalid | none | 10000 | 1053 | 197 |
| 58 | katydid-conf-valid | generated | valid | none | 10000 | 1053 | 233 |
| 59 | lazygit-rmUniqueItems | collected | valid | removed | 280 | 89462 | 276 |
| 60 | lazygit-rmUniqueItems-invalid | mutated | invalid | removed | 280 | 89462 | 375 |
| 61 | lerna | collected | valid | none | 985 | 4717 | 173 |
| 62 | lerna-invalid | mutated | invalid | none | 985 | 4717 | 228 |
| 63 | nest-cli | collected | valid | none | 1025 | 19450 | 290 |
| 64 | nest-cli-invalid | mutated | invalid | none | 1025 | 19450 | 355 |
| 65 | omnisharp | collected | valid | none | 987 | 13924 | 595 |
| 66 | omnisharp-invalid | mutated | invalid | none | 987 | 13924 | 683 |
| 67 | pre-commit-hooks | collected | valid | none | 985 | 9920 | 551 |
| 68 | pre-commit-hooks-invalid | mutated | invalid | none | 985 | 9920 | 610 |
| 69 | pulumi | collected | valid | none | 3807 | 7983 | 251 |
| 70 | pulumi-invalid | mutated | invalid | none | 3807 | 7983 | 290 |
| 71 | semantic-release | collected | valid | none | 794 | 3422 | 461 |
| 72 | semantic-release-invalid | mutated | invalid | none | 794 | 3422 | 531 |
| 73 | stale | collected | valid | none | 961 | 3702 | 469 |
| 74 | stale-invalid | mutated | invalid | none | 961 | 3702 | 519 |
| 75 | stylecop-rmUniqueItems | collected | valid | removed | 983 | 11732 | 568 |
| 76 | stylecop-rmUniqueItems-invalid | mutated | invalid | removed | 983 | 11732 | 590 |
| 77 | tmuxinator | collected | valid | none | 382 | 4545 | 629 |
| 78 | tmuxinator-invalid | mutated | invalid | none | 382 | 4545 | 669 |
| 79 | ui5 | collected | valid | none | 942 | 96365 | 487 |
| 80 | ui5-invalid | mutated | invalid | none | 942 | 96365 | 488 |
| 81 | unreal-engine-uproject-rmUniqueItems | collected | valid | removed | 859 | 10402 | 394 |
| 82 | unreal-engine-uproject-rmUniqueItems-invalid | mutated | invalid | removed | 859 | 10402 | 444 |
| 83 | vercel | collected | valid | none | 710 | 38187 | 406 |
| 84 | vercel-invalid | mutated | invalid | none | 710 | 38187 | 467 |
| 85 | yamllint | collected | valid | none | 984 | 890 | 352 |
| 86 | yamllint-invalid | mutated | invalid | none | 984 | 890 | 376 |
| 87 | zschema-advanced-rmUniqueItems-invalid | generated | invalid | removed | 10000 | 2618 | 620 |
| 88 | zschema-advanced-rmUniqueItems-valid | generated | valid | removed | 10000 | 2618 | 282 |
| 89 | zschema-basic-rmUniqueItems-invalid | generated | invalid | removed | 200 | 975 | 2894 |
| 90 | zschema-basic-rmUniqueItems-valid | generated | valid | removed | 200 | 975 | 2731 |

### Generation

Some schemas had a collection of instances gathered from github, which are mutated to create invalid instances via `make mutate`.
The rest we can regenerate valid and invalid instances for by running: `make generate`.

Shortcomings of document generator:
* Lots of libraries do not support unicode strings properly, so random generation of documents has been limited to ascii strings.
* Since we using JSONL for documents, we do not generate newlines in strings.
* Number generation has been limited to 64 bit floats for the generator and ints for the mutator, to avoid issues some implementations we having.
