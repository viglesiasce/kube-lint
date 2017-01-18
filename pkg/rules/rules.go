package rules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"k8s.io/client-go/pkg/util/jsonpath"
)

type Rule interface {
	Evaluate(resource []byte) Result
	GetName() string
}

type Result struct {
	Passed   bool
	Expected string
	Actual   string
}

type Operator interface {
	Evaluate(resource []byte, field string)
}

// KubernetesRule represents a single policy for the linting of a resource
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
func (kr KubernetesRule) Evaluate(resource []byte) Result {
	j := jsonpath.New(kr.Name)
	j.AllowMissingKeys(true)
	// insert brackets for the user so that the YAML gets easier to write
	err := j.Parse("{" + kr.Field + "}")
	if err != nil {
		panic(fmt.Errorf("in %s, parse %s error %v", kr.Name, resource, err))
	}
	buf := new(bytes.Buffer)
	var data interface{}
	err = json.Unmarshal(resource, &data)
	if err != nil {
		panic(err)
	}
	err = j.Execute(buf, data)
	if err != nil {
		panic(fmt.Errorf("in %s, execute error %v", kr.Name, err))
	}

	var value string
	var passed bool
	var out interface{}
	var actual string
	switch kr.ValueType {
	case "bool":
		parsedValue := false
		// TODO this is the case where it is unset, what to do?
		if buf.String() != "" {
			parsedValue, err = strconv.ParseBool(buf.String())
			if err != nil {
				panic(err)
			}
		}
		out, passed = kr.evaluateAsBool(parsedValue)
		actual = strconv.FormatBool(out.(bool))
		value = strconv.FormatBool(kr.Value.(bool))
	// String is the default type
	default:
		out, passed = kr.evaluateAsString(buf.String())
		actual = out.(string)
		if kr.Value != nil {
			value = kr.Value.(string)
		} else {
			value = ""
		}
	}
	return Result{Passed: passed, Expected: value, Actual: actual}
}

func (kr KubernetesRule) evaluateAsBool(value bool) (bool, bool) {
	var passed bool
	switch kr.Operator {
	case "null":
		passed = true
	case "equal":
		passed = value == kr.Value.(bool)
	case "notequal":
		passed = value != kr.Value.(bool)
	case "set":
		passed = value == true
	case "unset":
		passed = value == false
	default:
		panic("Operator not implemented for boolean type: " + kr.Operator)
	}
	return value, passed
}

func (kr KubernetesRule) evaluateAsString(value string) (string, bool) {
	var passed bool
	switch kr.Operator {
	case "null":
		passed = true
	case "equal":
		passed = value == kr.Value.(string)
	case "notequal":
		passed = value != kr.Value.(string)
	case "set":
		passed = value != ""
	case "unset":
		passed = value == ""
	case "matches":
		regex := regexp.MustCompile(kr.Value.(string))
		passed = regex.MatchString(value)
	default:
		panic("Operator not implemented for string type: " + kr.Operator)
	}
	return value, passed
}

// GetName returns the name of a KubernetesRule
func (kr KubernetesRule) GetName() string {
	return kr.Name
}

// NewKubernetesRule returns a KubernetesRule object
func NewKubernetesRule(operator string, field string, value interface{}, valueType string) Rule {
	return KubernetesRule{Operator: operator, Field: field, Value: value, ValueType: valueType}
}
