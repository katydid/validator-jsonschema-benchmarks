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
	"slices"
	"strings"

	"github.com/katydid/validator-jsonschema-benchmarks/analytics/analytics"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	log.SetFlags(log.Lshortfile)
	format := flag.String("format", "latex", "output format (md|latex)")
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

	if impl != nil && *impl != "" {
		impls := strings.Split(*impl, " ")
		log.Printf("filtering implementations: %#v", impls)
		lines = analytics.FilterImplementations(lines, impls)
	}

	totalSchemas := analytics.ListSchemas(lines)
	completedPerImpl := analytics.ListCompletedSchemasPerImplementation(lines)

	if *rmSchema1 {
		lines = analytics.FilterSchemasExitStatus0(lines)
	} else {
		lines = analytics.FilterExitStatus0(lines)
	}

	lines = analytics.MeanRuns(lines)

	kindGroups := analytics.GroupByKind(lines)
	analyseds := [][]*analytics.Implementation{}
	for i := range kindGroups {
		kind := kindGroups[i][0].Schema.ValidationKind
		if kind == "" {
			kind = "valid"
		}
		log.Printf("analysing kind group %d = %v with %d lines", i, kind, len(kindGroups[i]))
		groups := analytics.GroupBySchema(kindGroups[i])
		log.Printf("got %d groups", len(groups))
		scores := analytics.ScoreGroups(groups)
		log.Printf("got %d scores", len(scores))
		analysed := analytics.AnalyseImplementations(scores, totalSchemas, completedPerImpl)
		if len(analysed[0].Scores) == 0 {
			log.Fatalf("no scores for kind %v", kind)
		}
		analyseds = append(analyseds, analysed)
	}

	switch *format {
	case "latex":
		fprintLatex(os.Stdout, analyseds)
	case "md":
		fprintMarkdown(os.Stdout, analyseds)
	}
}

func getKind(impls []*analytics.Implementation) string {
	for _, impl := range impls {
		if len(impl.Scores) == 0 {
			log.Fatalf("no scores for %s", impl.Name)
		}
		for _, score := range impl.Scores {
			kind := score.Line.Schema.ValidationKind
			if kind == "" {
				return "valid"
			}
			return kind
		}
	}
	log.Fatalf("could not find kind in %v", impls)
	return ""
}

type details struct {
	Name string
	Lang string
	Link string
}

var implDetails = map[string]*details{
	"ajv":                              {Name: "Ajv", Lang: "Javascript (gen)", Link: "https://ajv.js.org/"},
	"ajv-bun":                          {Name: "Ajv-Bun", Lang: "Javascript (gen)", Link: "https://ajv.js.org/"},
	"blaze":                            {Name: "Blaze", Lang: "C++", Link: "https://github.com/sourcemeta/blaze"},
	"boon":                             {Name: "boon", Lang: "Rust", Link: "https://github.com/santhosh-tekuri/boon"},
	"corvus":                           {Name: "Corvus", Lang: "C# (gen)", Link: "https://github.com/corvus-dotnet/Corvus.JsonSchema"},
	"go-google":                        {Name: "Google", Lang: "Golang", Link: "https://github.com/google/jsonschema-go"},
	"go-json-schema-spec":              {Name: "json-schema-spec", Lang: "Golang", Link: "https://github.com/json-schema-spec/json-schema-go"},
	"go-kaptinlin":                     {Name: "kaptinlin", Lang: "Golang", Link: "https://github.com/kaptinlin/jsonschema"},
	"go-katydid-auto-json":             {Name: "Katydid-Comp-Fused", Lang: "Golang", Link: "https://git.katydid.org.za/validator-go-jsonschema"},
	"go-katydid-auto-reflect":          {Name: "Katydid-Comp-Steps", Lang: "Golang", Link: "https://git.katydid.org.za/validator-go-jsonschema"},
	"go-katydid-mem-json":              {Name: "Katydid-Memo-Fused", Lang: "Golang", Link: "https://git.katydid.org.za/validator-go-jsonschema"},
	"go-katydid-mem-reflect":           {Name: "Katydid-Memo-Steps", Lang: "Golang", Link: "https://git.katydid.org.za/validator-go-jsonschema"},
	"go-katydid-auto-json-min-parse":   {Name: "Katydid-Comp-Fused-min-Parse", Lang: "Golang", Link: "https://git.katydid.org.za/validator-go-jsonschema"},
	"go-katydid-mem-reflect-min-parse": {Name: "Katydid-Memo-Steps-min-Parse", Lang: "Golang", Link: "https://git.katydid.org.za/validator-go-jsonschema"},
	"go-santhosh-tekuri":               {Name: "santhosh-tekuri", Lang: "Golang", Link: "https://github.com/santhosh-tekuri/jsonschema/"},
	"hyperjump":                        {Name: "Hyperjump", Lang: "Javascript", Link: "https://github.com/hyperjump-io/json-schema"},
	"jsdotnet":                         {Name: "json-everything", Lang: "C#", Link: "https://github.com/json-everything/json-everything"},
	"json_schemer":                     {Name: "JSONSchemer", Lang: "Ruby", Link: "https://github.com/davishmcclurg/json_schemer"},
	"jsoncons":                         {Name: "jsoncons", Lang: "C++", Link: "https://github.com/danielaparker/jsoncons"},
	"jsu-c":                            {Name: "JSU", Lang: "C++ (gen)", Link: "https://github.com/zx80/json-schema-utils"},
	"jsu-java":                         {Name: "JSU", Lang: "Java (gen)", Link: "https://github.com/zx80/json-schema-utils"},
	"jsu-js":                           {Name: "JSU", Lang: "Javascript (gen)", Link: "https://github.com/zx80/json-schema-utils"},
	"jsu-pl":                           {Name: "JSU", Lang: "Perl (gen)", Link: "https://github.com/zx80/json-schema-utils"},
	"jsu-py":                           {Name: "JSU", Lang: "Python (gen)", Link: "https://github.com/zx80/json-schema-utils"},
	"jsv":                              {Name: "JSV", Lang: "Elixir", Link: "https://github.com/lud/jsv"},
	"kmp":                              {Name: "OptimumCode", Lang: "Kotlin", Link: "https://github.com/OptimumCode/json-schema-validator"},
	"networknt":                        {Name: "networknt", Lang: "Java", Link: "https://github.com/networknt/json-schema-validator"},
	"opis":                             {Name: "Opis", Lang: "PHP", Link: "https://opis.io/json-schema"},
	"py-jsonschema":                    {Name: "py-jsonschema", Lang: "Python", Link: "https://github.com/python-jsonschema/jsonschema/"},
	"rapidjson":                        {Name: "RapidJSON", Lang: "C++", Link: "https://github.com/Tencent/rapidjson/"},
	"schemasafe":                       {Name: "schemasafe", Lang: "Javascript", Link: "https://github.com/ExodusMovement/schemasafe"},
}

func fprintMarkdown(
	w io.Writer,
	implss [][]*analytics.Implementation,
) {
	p := func(format string, a ...any) {
		fmt.Fprintf(w, format, a...)
	}

	p(`## Filtering use case (invalid docs and parsing time included)`)
	p("\n")
	p("\n")
	p(`| impl | language  | completed | median/doc | mean/doc | median rank | mean rank | not completed |`)
	p("\n")
	p(`| ---  | ---       | ---       | ---      | ---        | ---       | ---         | ---           | `)
	p("\n")
	for i, impls := range implss {
		if getKind(implss[i]) != "invalid" {
			continue
		}
		slices.SortFunc(impls, func(x, y *analytics.Implementation) int {
			return analytics.FloatCompare(x.MedianParsePlusWarmNsPerDoc, y.MedianParsePlusWarmNsPerDoc)
		})
		for _, impl := range impls {
			p("| ")
			p("[%s](%s)", implDetails[impl.Name].Name, implDetails[impl.Name].Link)
			p(" | ")
			p("%s", implDetails[impl.Name].Lang)
			p(" | ")
			p("%d/%d", len(impl.CompletedSchemas), len(impl.AllSchemas))
			p(" | ")
			p("%.0f", impl.MedianParsePlusWarmNsPerDoc)
			p(" | ")
			p("%.0f", impl.MeanParsePlusWarmNsPerDoc)
			p(" | ")
			p("%.0f", impl.MedianParsePlusWarmRank)
			p(" | ")
			p("%.0f", impl.MeanParsePlusWarmRank)
			p(" | ")
			p("%v", impl.NotCompletedSchemas)
			p(" |")
			p("\n")
		}
	}

	p("\n")
	p(`## API Gateway with no reuse use case (valid docs and parsing time included)`)
	for i, impls := range implss {
		if getKind(implss[i]) != "valid" {
			continue
		}
		p("\n")
		p("\n")
		p(`| impl | language  | completed | median/doc | mean/doc | median rank | mean rank | not completed |`)
		p("\n")
		p(`| ---  | ---       | ---       | ---      | ---        | ---       | ---         | ---           | `)
		p("\n")
		slices.SortFunc(impls, func(x, y *analytics.Implementation) int {
			return analytics.FloatCompare(x.MedianParsePlusWarmNsPerDoc, y.MedianParsePlusWarmNsPerDoc)
		})
		for _, impl := range impls {
			p("| ")
			p("[%s](%s)", implDetails[impl.Name].Name, implDetails[impl.Name].Link)
			p(" | ")
			p("%s", implDetails[impl.Name].Lang)
			p(" | ")
			p("%d/%d", len(impl.CompletedSchemas), len(impl.AllSchemas))
			p(" | ")
			p("%.0f", impl.MedianParsePlusWarmNsPerDoc)
			p(" | ")
			p("%.0f", impl.MeanParsePlusWarmNsPerDoc)
			p(" | ")
			p("%.0f", impl.MedianParsePlusWarmRank)
			p(" | ")
			p("%.0f", impl.MeanParsePlusWarmRank)
			p(" | ")
			p("%v", impl.NotCompletedSchemas)
			p(" |")
			p("\n")
		}
	}

	// remove katydid-auto-json and katydid-mem-json, since they both do not create any structures that can be reused
	for i := range implss {
		implss[i] = slices.DeleteFunc(implss[i], func(impl *analytics.Implementation) bool { return impl.Name == "go-katydid-auto-json" })
		implss[i] = slices.DeleteFunc(implss[i], func(impl *analytics.Implementation) bool { return impl.Name == "go-katydid-mem-json" })
	}
	p("\n")
	p(`## API Gateway with generic reuse use case (valid docs and parsing time excluded)`)
	for i, impls := range implss {
		if getKind(implss[i]) != "valid" {
			continue
		}
		p("\n")
		p("\n")
		p(`| impl | language  | completed | median/doc | mean/doc | median rank | mean rank | not completed |`)
		p("\n")
		p(`| ---  | ---       | ---       | ---      | ---        | ---       | ---         | ---           | `)
		p("\n")
		slices.SortFunc(impls, func(x, y *analytics.Implementation) int {
			return analytics.FloatCompare(x.MedianWarmNsPerDoc, y.MedianWarmNsPerDoc)
		})
		for _, impl := range impls {
			p("| ")
			p("[%s](%s)", implDetails[impl.Name].Name, implDetails[impl.Name].Link)
			p(" | ")
			p("%s", implDetails[impl.Name].Lang)
			p(" | ")
			p("%d/%d", len(impl.CompletedSchemas), len(impl.AllSchemas))
			p(" | ")
			p("%.0f", impl.MedianWarmNsPerDoc)
			p(" | ")
			p("%.0f", impl.MeanWarmNsPerDoc)
			p(" | ")
			p("%.0f", impl.MedianWarmRank)
			p(" | ")
			p("%.0f", impl.MeanWarmRank)
			p(" | ")
			p("%v", impl.NotCompletedSchemas)
			p(" |")
			p("\n")
		}
	}
}

func fprintLatex(
	w io.Writer,
	implss [][]*analytics.Implementation,
) {
	var messagePrinter *message.Printer = message.NewPrinter(language.English)
	p := func(format string, a ...any) {
		messagePrinter.Fprintf(w, format, a...)
	}
	sprint := func(name string) string {
		name = strings.Replace(name, "#", "\\#", -1)
		return strings.Replace(name, "_", "\\_", -1)
	}

	for i, impls := range implss {
		p("%% BEGIN Generated tabular for kind: %s and parsing included\n", getKind(implss[i]))
		p("\\begin{tabular}{ll|l|ll}\n")
		p(`\# & impl & language & median ns/doc & mean ns/doc \\`)
		p("\n")
		p(`\hline`)
		p("\n")
		slices.SortFunc(impls, func(x, y *analytics.Implementation) int {
			return analytics.FloatCompare(x.MedianParsePlusWarmNsPerDoc, y.MedianParsePlusWarmNsPerDoc)
		})
		j := 0
		for _, impl := range impls {
			if impl.Name == "go-katydid-auto-reflect" ||
				impl.Name == "go-katydid-mem-reflect" ||
				impl.Name == "go-katydid-mem-json" ||
				impl.Name == "go-katydid-auto-json-min-parse" ||
				impl.Name == "go-katydid-mem-json-min-parse" {
				// remove katydid-auto-reflect and katydid-mem-reflect the objects they create will be discard and ...-json alternatives are faster.
				// remove go-katydid-mem-json since auto is faster.
				// remove ...-min-parse, since these have parsing times included.
				continue
			}
			j++
			p("%d", j)
			p(" & ")
			p("%s", sprint(implDetails[impl.Name].Name))
			p(" & ")
			p("%s", sprint(implDetails[impl.Name].Lang))
			p(" & ")
			p("%.0f", impl.MedianParsePlusWarmNsPerDoc)
			p(" & ")
			p("%.0f", impl.MeanParsePlusWarmNsPerDoc)
			p(" \\\\\n")

		}
		p("\\end{tabular}\n")
		p("%% END Generated tabular\n")
	}

	for i, impls := range implss {
		p("%% BEGIN Generated tabular for kind: %s and parsing excluded\n", getKind(implss[i]))
		p("\\begin{tabular}{ll|l|ll}\n")
		p(`\# & impl & median ns/doc & mean ns/doc \\`)
		p("\n")
		p(`\hline`)
		p("\n")
		slices.SortFunc(impls, func(x, y *analytics.Implementation) int {
			return analytics.FloatCompare(x.MedianWarmNsPerDoc, y.MedianWarmNsPerDoc)
		})
		j := 0
		for _, impl := range impls {
			if impl.Name == "go-katydid-auto-json" || impl.Name == "go-katydid-mem-json" || impl.Name == "go-katydid-mem-reflect" {
				// remove katydid-auto-json and katydid-mem-json, since they both do not create any structures that can be reused
				// remove go-katydid-mem-reflect, since auto is faster.
				continue
			}
			j++
			p("%d", j)
			p(" & ")
			p("%s", sprint(implDetails[impl.Name].Name))
			p(" & ")
			p("%s", sprint(implDetails[impl.Name].Lang))
			p(" & ")
			p("%.0f", impl.MedianWarmNsPerDoc)
			p(" & ")
			p("%.0f", impl.MeanWarmNsPerDoc)
			p(" \\\\\n")

		}
		p("\\end{tabular}\n")
		p("%% END Generated tabular\n")
	}
}
