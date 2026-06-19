# Copyright 2026 Walter Schulze
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     https://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

generate:
	go run ./cmd/generator -seed=123456 -num=10000 ../schemas/

# draft-04-rmUniqueItems-invalid: limit generation to 500 instances to avoid issues with go-kaptinlin. TODO report issue and when it is fixed, undo the limit.
# jsconfig-rmUniqueItems-invalid: limit generation to 400 instances to avoid issues with networknt. TODO report issue and when it is fixed, undo the limit.
# vercel-invalid: blaze (https://github.com/sourcemeta/blaze/issues/878)
# ui5-invalid: jsu-c, jsu-java, jsu-js, jsu-py (https://github.com/clairey-zx81/json-model/issues/8)

mutate:
	((rm -rf ../schemas/ansible-meta-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/ansible-meta)
	((rm -rf ../schemas/aws-cdk-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/aws-cdk)
	((rm -rf ../schemas/babelrc-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/babelrc)
	((rm -rf ../schemas/clang-format-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/clang-format)
	((rm -rf ../schemas/cmake-presets-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/cmake-presets)
	((rm -rf ../schemas/code-climate-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/code-climate)
	((rm -rf ../schemas/cypress-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/cypress)
	((rm -rf ../schemas/deno-rmUniqueItems-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/deno-rmUniqueItems)
	((rm -rf ../schemas/dependabot-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/dependabot)
	((rm -rf ../schemas/draft-04-rmUniqueItems-invalid || true) && go run ./cmd/mutator -num=500 -seed=123456 ../schemas/draft-04-rmUniqueItems)
	((rm -rf ../schemas/fabric-mod-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/fabric-mod)
	((rm -rf ../schemas/geojson-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/geojson)
	((rm -rf ../schemas/gitpod-configuration-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/gitpod-configuration)
	((rm -rf ../schemas/importmap-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/importmap)
	((rm -rf ../schemas/jasmine-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/jasmine)
	((rm -rf ../schemas/jsconfig-rmUniqueItems-invalid || true) && go run ./cmd/mutator -num=400 -seed=123456 ../schemas/jsconfig-rmUniqueItems)
	((rm -rf ../schemas/jshintrc-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/jshintrc)
	((rm -rf ../schemas/lazygit-rmUniqueItems-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/lazygit-rmUniqueItems)
	((rm -rf ../schemas/lerna-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/lerna)
	((rm -rf ../schemas/nest-cli-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/nest-cli)
	((rm -rf ../schemas/omnisharp-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/omnisharp)
	((rm -rf ../schemas/pre-commit-hooks-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/pre-commit-hooks)
	((rm -rf ../schemas/pulumi-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/pulumi)
	((rm -rf ../schemas/semantic-release-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/semantic-release)
	((rm -rf ../schemas/stale-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/stale)
	((rm -rf ../schemas/stylecop-rmUniqueItems-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/stylecop-rmUniqueItems)
	((rm -rf ../schemas/tmuxinator-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/tmuxinator)
	((rm -rf ../schemas/ui5-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/ui5)
	((rm -rf ../schemas/unreal-engine-uproject-rmUniqueItems-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/unreal-engine-uproject-rmUniqueItems)
	((rm -rf ../schemas/vercel-invalid || true) && go run ./cmd/mutator -seed=1 -num=500 ../schemas/vercel)
	((rm -rf ../schemas/yamllint-invalid || true) && go run ./cmd/mutator -seed=123456 ../schemas/yamllint)