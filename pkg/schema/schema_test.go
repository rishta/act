package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestAdditionalFunctions(t *testing.T) {
	var node yaml.Node
	err := yaml.Unmarshal([]byte(`
on: push
jobs:
  job-with-condition:
    runs-on: self-hosted
    if: success() || success('joba', 'jobb') || failure() || failure('joba', 'jobb') || always() || cancelled()
    steps:
    - run: exit 0
`), &node)
	if !assert.NoError(t, err) {
		return
	}
	err = (&Node{
		Definition: "workflow-root-strict",
		Schema:     GetWorkflowSchema(),
	}).UnmarshalYAML(&node)
	assert.NoError(t, err)
}

func TestAdditionalFunctionsFailure(t *testing.T) {
	var node yaml.Node
	err := yaml.Unmarshal([]byte(`
on: push
jobs:
  job-with-condition:
    runs-on: self-hosted
    if: success() || success('joba', 'jobb') || failure() || failure('joba', 'jobb') || always('error')
    steps:
    - run: exit 0
`), &node)
	if !assert.NoError(t, err) {
		return
	}
	err = (&Node{
		Definition: "workflow-root-strict",
		Schema:     GetWorkflowSchema(),
	}).UnmarshalYAML(&node)
	assert.Error(t, err)
}

func TestAdditionalFunctionsSteps(t *testing.T) {
	var node yaml.Node
	err := yaml.Unmarshal([]byte(`
on: push
jobs:
  job-with-condition:
    runs-on: self-hosted
    steps:
    - run: exit 0
      if: success() || failure() || always()
`), &node)
	if !assert.NoError(t, err) {
		return
	}
	err = (&Node{
		Definition: "workflow-root-strict",
		Schema:     GetWorkflowSchema(),
	}).UnmarshalYAML(&node)
	assert.NoError(t, err)
}

func TestAdditionalFunctionsStepsExprSyntax(t *testing.T) {
	var node yaml.Node
	err := yaml.Unmarshal([]byte(`
on: push
jobs:
  job-with-condition:
    runs-on: self-hosted
    steps:
    - run: exit 0
      if: ${{ success() || failure() || always() }}
`), &node)
	if !assert.NoError(t, err) {
		return
	}
	err = (&Node{
		Definition: "workflow-root-strict",
		Schema:     GetWorkflowSchema(),
	}).UnmarshalYAML(&node)
	assert.NoError(t, err)
}

func TestJobEnvironmentDeployment(t *testing.T) {
	for _, tc := range []struct {
		value string
		ok    bool
	}{
		{"false", true},
		{"true", true},
		{"maybe", false},
	} {
		var node yaml.Node
		err := yaml.Unmarshal([]byte(`
on: push
jobs:
  job-with-environment:
    runs-on: ubuntu-latest
    environment:
      name: production
      deployment: `+tc.value+`
    steps:
    - run: exit 0
`), &node)
		if !assert.NoError(t, err) {
			return
		}
		err = (&Node{
			Definition: "workflow-root-strict",
			Schema:     GetWorkflowSchema(),
		}).UnmarshalYAML(&node)
		if tc.ok {
			assert.NoError(t, err, "deployment: %s", tc.value)
		} else {
			assert.Error(t, err, "deployment: %s", tc.value)
		}
	}
}
