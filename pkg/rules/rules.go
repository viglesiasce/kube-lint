package rules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"k8s.io/client-go/pkg/util/jsonpath"
)

//
type Rule interface {
	Evaluate(resource []byte) Result
	GetName() string
}

//
type Result struct {
	Passed   bool
	Expected string
	Actual   string
}

//
type Operator interface {
	Evaluate(resource []byte, field string)
}

// Represents a single policy for the linting of a resource
type KubernetesRule struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Kind        string      `json:"kind"`
	Operator    string      `json:"operator"`
	Field       string      `json:"field"`
	Value       interface{} `json:"value"`
	ValueType   string      `json:"valueType"`
}

type LinterConfig map[string][]KubernetesRule

// Evaluate rule
func (lr KubernetesRule) Evaluate(resource []byte) Result {
	j := jsonpath.New(lr.Name)
	j.AllowMissingKeys(true)
	// insert brackets for the user so that the YAML gets easier to write
	err := j.Parse("{" + lr.Field + "}")
	if err != nil {
		panic(fmt.Errorf("in %s, parse %s error %v", lr.Name, resource, err))
	}
	buf := new(bytes.Buffer)
	var data interface{}
	err = json.Unmarshal(resource, &data)
	if err != nil {
		panic(err)
	}
	err = j.Execute(buf, data)
	if err != nil {
		panic(fmt.Errorf("in %s, execute error %v", lr.Name, err))
	}
	out := buf.String()
	var passed bool
	switch lr.Operator {
	case "null":
		passed = true
	case "equal":
		passed = out == lr.Value.(string)
	case "notequal":
		passed = out != lr.Value.(string)
	case "set":
		passed = out != ""
	case "unset":
		passed = out == ""
	case "matches":
		regex := regexp.MustCompile(lr.Value.(string))
		passed = regex.MatchString(out)
	default:
		panic("Operator not implemented: " + lr.Operator)
	}
	var value string
	switch lr.ValueType {
	case "bool":
		value = strconv.FormatBool(lr.Value.(bool))
	default:
		if lr.Value != nil {
			value = lr.Value.(string)
		} else {
			value = ""
		}
	}
	return Result{Passed: passed, Expected: value, Actual: out}
}

// GetName returns the name of a KubernetesRule
func (lr KubernetesRule) GetName() string {
	return lr.Name
}

// NewKubernetesRule returns a KubernetesRule object
func NewKubernetesRule(operator string, field string, value interface{}, valueType string) Rule {
	return KubernetesRule{Operator: operator, Field: field, Value: value, ValueType: valueType}
}
