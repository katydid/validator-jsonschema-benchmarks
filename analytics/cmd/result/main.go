// Copyright 2026 Walter Schulze
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/katydid/validator-jsonschema-benchmarks/analytics/analytics"
)

func main() {
	log.SetFlags(log.Lshortfile)
	format := flag.String("format", "latex", "output format (md|latex)")
	rmUniqueItems := flag.Bool("rmUniqueItems", false, "if there is an rmUniqueItems version replace the original with it")
	latexPifont := flag.Bool("latex.pifont", false, "replace yes/no in with \\cmark/\\xmark, requires adding the following to your latex: \\usepackage{pifont}\\newcommand{\\cmark}{\\ding{51}}\\newcommand{\\xmark}{\\ding{55}}")
	shortSchemaName := flag.Bool("shortSchemaName", false, "write short schema name")
	rmSchema1 := flag.Bool("filterSchema1", false, "filter all schemas where one implementation had an non zero exit code")
	rootFolder := flag.String("rootFolder", "./schemas", "location of schemas folder")
	impl := flag.String("impls", "", "space separated list of implementations to filter")
	flag.Parse()
	reportFilename := "./dist/report.csv"
	if len(flag.Args()) == 0 {
		log.Printf("expected location of report.csv as first argument to command, defaulting to ./dist/report.csv")
	} else {
		reportFilename = flag.Args()[0]
	}

	log.Printf("analysing schemas folder at %s\n", *rootFolder)
	curated := analytics.GetCuratedSchemas(*rootFolder)
	schemasFolder := filepath.Join(*rootFolder, "schemas")
	schemas, err := analytics.CollectSchemas(schemasFolder, curated)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("analysing report file %s\n", reportFilename)
	lines, err := analytics.ReadLines(schemas, reportFilename)
	if err != nil {
		log.Fatal(err)
	}

	if *rmUniqueItems {
		lines = analytics.FilterNoUniqueItems(lines)
	}

	if impl != nil && *impl != "" {
		impls := strings.Split(*impl, " ")
		log.Printf("filtering implementations: %#v", impls)
		lines = analytics.FilterImplementations(lines, impls)
	}

	if *rmSchema1 {
		lines = analytics.FilterSchemasExitStatus0(lines)
	} else {
		lines = analytics.FilterExitStatus0(lines)
	}

	lines = analytics.MeanRuns(lines)

	groups := analytics.GroupBySchema(lines)
	scores := analytics.ScoreGroups(groups)

	sprintName := func(name string) string {
		if *rmUniqueItems {
			name = strings.Replace(name, "-rmUniqueItems", "", 1)
		}
		if *shortSchemaName {
			schemaName, err := analytics.ParseSchemaName(name)
			if err != nil {
				panic(err)
			}
			name = schemaName.ShortName
		}
		esc := func(s string) string {
			return strings.Replace(s, "_", "\\_", -1)
		}
		return esc(name)
	}

	switch *format {
	case "latex":
		sprintBool := func(b bool) string {
			return analytics.SprintLatexBool(*latexPifont, b)
		}
		fprintLatex(os.Stdout, sprintName, sprintBool, scores)
	case "md":
		sprintBool := func(b bool) string {
			if b {
				return ":white_check_mark:"
			} else {
				return ":x:"
			}
		}
		fprintMarkdown(os.Stdout, sprintName, sprintBool, scores)
	}
}

func fprintLatex(
	w io.Writer,
	sprintName func(string) string,
	sprintBool func(bool) string,
	scores []*analytics.ScoredLine,
) {
	p := func(format string, a ...any) {
		fmt.Fprintf(w, format, a...)
	}

	p("%% BEGIN Generated tabular\n")
	p("\\begin{tabular}{lll|ll|ll}\n")
	p(`name & valid & impl & rank (with parsing) & rank (without parsing) \\`)
	p("\n")
	currentSchemaName := ""
	for _, score := range scores {
		if currentSchemaName != score.Line.Schema.Name {
			currentSchemaName = score.Line.Schema.Name
			p("\\hline")
			p("\n")
		}
		p("%s", sprintName(score.Line.Schema.Name))
		p(" & ")
		p("%s", sprintBool(score.Line.Schema.ValidationKind == "valid" || score.Line.Schema.ValidationKind == ""))
		p(" & ")
		p("%d", score.ParsePlusWarmRank)
		p(" & ")
		p("%s", score.Line.Implementation)
		p(" & ")
		p("%d", score.WarmRank)
		p(" \\\\\n")

	}
	p("\\end{tabular}\n")
	p("%% END Generated tabular\n")
}

func fprintMarkdown(
	w io.Writer,
	sprintName func(string) string,
	sprintBool func(bool) string,
	scores []*analytics.ScoredLine,
) {
	p := func(format string, a ...any) {
		fmt.Fprintf(w, format, a...)
	}

	p(`| name | valid | impl | ns/op (with parsing) | ns/op (without parsing) |`)
	p("\n")
	p(`| ---  | ---   | ---  | ---                  | ---                     |`)
	p("\n")
	for _, score := range scores {
		p("| ")
		p("%s", sprintName(score.Line.Schema.ShortName))
		p(" | ")
		p("%s", sprintBool(score.Line.Schema.ValidationKind == "valid" || score.Line.Schema.ValidationKind == ""))
		p(" | ")
		p("%s", score.Line.Implementation)
		p(" | ")
		p("%.0f", score.ParsePlusWarmNsPerDoc)
		p(" | ")
		p("%.0f", score.WarmNsPerDoc)
		p(" |")
		p("\n")

	}
}
