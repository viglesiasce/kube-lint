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
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/viglesiasce/kube-lint/pkg/check"
	"github.com/viglesiasce/kube-lint/pkg/rules"
)

var filename string
var kubeconfig string

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// load config
		var config rules.LinterConfig
		configFile, err := ioutil.ReadFile("example/config.yaml")
		if err != nil {
			panic("Unable to read config file")
		}
		err = yaml.Unmarshal(configFile, &config)
		if err != nil {
			panic("Unable to unmarshal config file")
		}

		resourceJSON := ""
		if kubeconfig != "" {
			resourceJSON = check.GetJSONFromKubernetes(kubeconfig)
		}

		if filename != "" {
			resourceJSON = check.GetJSONFromFile(filename)
		}

		if resourceJSON == "" {
			panic("NO RESOURCE FOUND")
		}

		table := check.CreateTable()
		check.EvaluateRules(table, config, resourceJSON)
		table.Render()
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")
	checkCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "example/pod.yaml", "Filename or directory of manifest(s)")
	checkCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file to use for requests")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
