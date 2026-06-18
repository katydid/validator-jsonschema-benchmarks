warmup_iterations = 100
# 10 seconds in ns
max_warmup_time = 10_000_000_000

[schema_path | _] = System.argv()
want = !(schema_path =~ "-invalid")
schema_string = File.read!(Path.join(schema_path, "schema.json"))
schema = JSON.decode!(schema_string)

# Compile the schema
{compile_duration, schema} =
  :timer.tc(fn -> JSV.build!(schema, formats: true) end, :nanosecond)

# Load instances
stream = File.stream!(Path.join(schema_path, "instances.jsonl"), [:read, :utf8])
lines = Enum.to_list(stream)

{parse_duration, instances} =
  :timer.tc(fn -> Enum.map(lines, &JSON.decode!/1) end, :nanosecond)

# Validate the data
{cold_duration, _} =
  :timer.tc(
    fn ->
      Enum.each(instances, fn instance ->
        case JSV.validate(instance, schema) do
          {:ok, _} when want -> :ok
          {:error, _} when not want -> :ok
          _ -> exit(:invalid)
        end
      end)
    end,
    :nanosecond
  )

iterations = trunc(Float.ceil(max_warmup_time / cold_duration))

Enum.each(0..min(iterations, warmup_iterations), fn _ ->
  Enum.each(instances, fn instance ->
    JSV.validate(instance, schema)
  end)
end)

# Validate the data
{warm_duration, _} =
  :timer.tc(
    fn ->
      Enum.each(instances, fn instance ->
        JSV.validate(instance, schema)
      end)
    end,
    :nanosecond
  )

IO.puts("#{cold_duration},#{warm_duration},#{parse_duration},#{compile_duration}")
