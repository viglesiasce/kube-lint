package check

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	"github.com/viglesiasce/kube-lint/pkg/rules"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func GetJSONFromKubernetes(kubeconfig string) string {
	// TODO Make this take a reader interface or the like
	k8sClientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(k8sClientConfig)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.Core().Pods("").List(v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var resourceJSONBytes []byte
	for _, pod := range pods.Items {
		resourceJSONBytes, err = json.Marshal(pod)
		if err != nil {
			panic("Unable to convert resource to JSON")
		}
	}
	return string(resourceJSONBytes)
}

func GetJSONFromFile(filename string) string {
	resourceFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic("Unable to read config file")
	}

	resourceJSONBytes, err := yaml.YAMLToJSON(resourceFile)
	if err != nil {
		panic("Unable to convert resource to JSON")
	}
	return string(resourceJSONBytes)
}

func CreateTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	//  "PROFILE", "OPERATOR", "EXPECTED", "ACTUAL",
	table.SetHeader([]string{"RULE", "POD", "RESULT"})
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	return table
}

func EvaluateRules(table *tablewriter.Table, config rules.LinterConfig, resourceJSON string) {
	for _, value := range config {
		for _, rule := range value {
			k8sRule := rules.NewKubernetesRule(rule.Operator, rule.Field, rule.Value, rule.ValueType)
			result := k8sRule.Evaluate(resourceJSON)
			var colorizedResult string
			if result.Passed {
				colorizedResult = color.GreenString("passed")
			} else {
				colorizedResult = color.RedString("failed")
			}
			// profile, rule.Operator, result.Expected, result.Actual,
			table.Append([]string{rule.Name, "podname", colorizedResult})
		}
	}
}
