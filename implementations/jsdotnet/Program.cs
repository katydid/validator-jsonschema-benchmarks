using Json.Schema;
using System.Diagnostics;
using System.Linq;
using System.Text.Json;
using System.Text.Json.Serialization;
using Json.More;

const int WarmupIterations = 1000;
const long MaxWarmupTime = 10_000_000_000;

var evaluationOptions = new EvaluationOptions
	{
		RequireFormatValidation = true,
    // OutputFormat = OutputFormat.Hierarchical, // uncomment for more detailed error reporting.
	};

bool ValidateAll(JsonSchema schema, JsonElement[] docs, bool want) {
  foreach (var doc in docs) {
    try
		{
      var result = schema.Evaluate(doc, evaluationOptions);
      if (result.IsValid != want) {
        Console.Error.WriteLine(doc);
        Console.Error.WriteLine(JsonSerializer.Serialize(result, new JsonSerializerOptions{}));
        return false; 
      }
    } catch (Exception e) {
      if (want) {
        return false;
      }
    }
  }

  return true;
}

var want = !args[0].Contains("-invalid");
Stopwatch stopWatch = new Stopwatch();

// Load the schema
stopWatch.Start();
var schema = JsonSchema.FromFile(args[0]);
stopWatch.Stop();
TimeSpan compileTs = stopWatch.Elapsed;

// Read and parse all instances
var lines = File.ReadLines(args[1]);
stopWatch.Start();
var docs = lines.Select(l => JsonDocument.Parse(l).RootElement).ToArray();
stopWatch.Stop();
TimeSpan parseTs = stopWatch.Elapsed;

// Loop and validate all instances
stopWatch.Start();
var valid = ValidateAll(schema, docs, want);
stopWatch.Stop();
TimeSpan coldTs = stopWatch.Elapsed;

var iterations = (int) Math.Ceiling(((double) MaxWarmupTime) / coldTs.TotalNanoseconds);
for (int i = 0; i < Math.Min(iterations, WarmupIterations); i++) {
  ValidateAll(schema, docs, want);
}

stopWatch.Restart();
ValidateAll(schema, docs, want);
stopWatch.Stop();
TimeSpan warmTs = stopWatch.Elapsed;

// Output file time and exit
Console.WriteLine(coldTs.TotalNanoseconds + "," + warmTs.TotalNanoseconds + "," + parseTs.TotalNanoseconds + "," + compileTs.TotalNanoseconds);
Environment.Exit(valid ? 0 : 1);
