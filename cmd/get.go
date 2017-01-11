// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/viglesiasce/kube-lint/pkg/rules"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "evaluate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// load yaml file
		var config rules.LinterConfig
		configFile, err := ioutil.ReadFile("config.yaml")
		if err != nil {
			panic("Unable to read config file")
		}
		// TODO Make this take a reader interface or the like
		k8sClientConfig, err := clientcmd.BuildConfigFromFlags("", "/Users/viglesias/.kube/config")
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

		// resourceFile, err := ioutil.ReadFile("pod.yaml")

		if err != nil {
			panic("Unable to read config file")
		}
		err = yaml.Unmarshal(configFile, &config)
		if err != nil {
			panic("Unable to unmarshal config file")
		}
		// resourceJSONBytes, err := yaml.YAMLToJSON(resourceFile)
		// resourceJSON := string(resourceJSONBytes)
		// if err != nil {
		// 	panic("Unable to convert resource to JSON")
		// }
		table := tablewriter.NewWriter(os.Stdout)
		//  "PROFILE", "OPERATOR", "EXPECTED", "ACTUAL",
		table.SetHeader([]string{"RULE", "POD", "RESULT"})
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		for _, pod := range pods.Items {
			resourceJSONBytes, err := json.Marshal(pod)
			if err != nil {
				panic("Unable to convert resource to JSON")
			}
			resourceJSON := string(resourceJSONBytes)

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
					table.Append([]string{rule.Name, pod.Name, colorizedResult})
				}
			}

		}
		table.Render()
	},
}

func init() {
	RootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("rules", "", "List rules")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
