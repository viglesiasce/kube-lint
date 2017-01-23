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

package pods

import (
	"reflect"
	"testing"

	"encoding/json"

	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	"github.com/viglesiasce/kube-lint/pkg/rules"
	"k8s.io/client-go/pkg/api/v1"
)

func TestGetPodsFromServer(t *testing.T) {
	type args struct {
		kubeconfig string
	}
	tests := []struct {
		name string
		args args
		want []v1.Pod
	}{}
	for _, tt := range tests {
		if got := NewKubeServer(tt.args.kubeconfig).GetPods(""); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. GetPodsFromServer() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetPodsFromFile(t *testing.T) {
	type args struct {
		filename string
	}
	examplePodFile := "../../example/pod.yaml"
	examplePod := v1.Pod{}
	podBytes, err := ioutil.ReadFile(examplePodFile)
	if err != nil {
		t.Errorf("Unable to read example file")
	}
	podJSONBytes, err := yaml.YAMLToJSON(podBytes)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(podJSONBytes, &examplePod)
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name string
		args args
		want []v1.Pod
	}{
		{"Example file", args{filename: examplePodFile}, []v1.Pod{examplePod}},
	}
	for _, tt := range tests {
		if got := NewLocalFilesystem(tt.args.filename).GetPods(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. GetPodsFromFile() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCreateTable(t *testing.T) {
	tests := []struct {
		name string
		want *tablewriter.Table
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := CreateTable(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. CreateTable() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestEvaluateRules(t *testing.T) {
	type args struct {
		config  rules.LinterConfig
		pods    []v1.Pod
		tags    []string
		showAll bool
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		EvaluateRules(tt.args.config, tt.args.pods, tt.args.tags, tt.args.showAll)
	}
}
