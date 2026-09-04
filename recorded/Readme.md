# Recorded Results

* [2026-06-24](./2026-06-24/Readme.md)
* [2026-09-02](./2026-09-02/Readme.md)
* [2026-09-04](./2026-09-04/Readme.md)

## How to produce a report

1. Startup Docker

2. Do a preliminary run to check that all the dockers are working:

```
$ rm -rf dist; make SCHEMAS='example-address-valid'
```

This will run one schema accross all implementations and makes sure that all the Dockers are prepared.

3. Do the actual run:

```
$ rm -rf dist; RUNS=3 make run
```

Be sure not to do anything else on your computer and that your power settings are set so that the computer uses its full power.

4. Analyse

Now that all the results have been gather they need to be analysed.
* Creating `results_per_impl.md` requires running `$ make analytics-impls-md` and renaming the output file `impls.md`
* Creating `results_per_impl.tex` requires running `$ make analytics-impls-latex` and renaming the output file `impls.md`
* Creating `schema_summary.md` requires running `$ make analytics-schemas-md` and renaming the output file `schemas.md`
* Creating `results_per_schema.md` requires running `$ make analytics-results-md` and renaming the output file `results.md`
