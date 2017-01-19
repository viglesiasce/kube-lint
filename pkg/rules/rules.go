// Copyright © 2017 Vic Iglesias <viglesiasce@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
	Tags        []string    `json:"tags"`
}

type LinterConfig []KubernetesRule

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
	case "float64":
		var parsedValue float64
		// TODO what should you do if it is not found?
		if buf.String() != "" {
			floats := strings.Fields(buf.String())
			for _, float := range floats {
				parsedValue64, err := strconv.ParseFloat(float, 64)
				if err != nil {
					panic(err)
				}
				parsedValue = parsedValue64
				out, passed = kr.evaluateAsFloat(parsedValue)
				if passed {
					break
				}
			}
			actual = buf.String()
			value = strconv.FormatFloat(kr.Value.(float64), 'f', 4, 64)
		}
	case "bool":
		parsedValue := false
		// TODO what should you do if it is not found?
		if buf.String() != "" {
			bools := strings.Fields(buf.String())
			for _, boolean := range bools {
				parsedValue, err = strconv.ParseBool(boolean)
				if err != nil {
					panic(err)
				}
				out, passed = kr.evaluateAsBool(parsedValue)
				if passed {
					break
				}

			}
			actual = strconv.FormatBool(out.(bool))
			value = strconv.FormatBool(kr.Value.(bool))
		}

	// String is the default type
	default:
		// TODO need to figure out what to do with the evaluation of multiple
		// fields (ie .spec.containers[*].name )
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

func (kr KubernetesRule) evaluateAsFloat(value float64) (float64, bool) {
	var passed bool
	switch kr.Operator {
	case "null":
		passed = true
	case "equal":
		passed = value == kr.Value.(float64)
	case "notequal":
		passed = value != kr.Value.(float64)
	case "greaterthan":
		passed = value > kr.Value.(float64)
	case "lessthan":
		passed = value < kr.Value.(float64)
	case "set":
		passed = value != 0
	case "unset":
		passed = value == 0
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
