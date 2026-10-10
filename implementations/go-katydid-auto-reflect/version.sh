
#!/bin/bash

set -o errexit
set -o nounset

# Extract the version of the jsonschema module from go.mod
jsonschema_version=$(grep 'katydid.org.za/go/validator-go-jsonschema' implementations/go-katydid-auto-reflect/go.mod | sed -E 's/.*\/validator-go-jsonschema v([0-9\.]+)-.*/\1/')

# Output the version
echo "$jsonschema_version"
