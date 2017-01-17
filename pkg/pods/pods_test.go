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
		if got := GetPodsFromServer(tt.args.kubeconfig); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. GetPodsFromServer() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetPodsFromFile(t *testing.T) {
	type args struct {
		filename string
	}
	examplePodFile := "/Users/viglesias/Dropbox/go/src/github.com/viglesiasce/kube-lint/example/pod.yaml"
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
		if got := GetPodsFromFile(tt.args.filename); !reflect.DeepEqual(got, tt.want) {
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
		table   *tablewriter.Table
		config  rules.LinterConfig
		pods    []v1.Pod
		showAll bool
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		EvaluateRules(tt.args.table, tt.args.config, tt.args.pods, tt.args.showAll)
	}
}
