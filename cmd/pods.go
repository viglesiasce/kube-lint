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

	"k8s.io/client-go/pkg/api/v1"

	"fmt"
	"os"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/viglesiasce/kube-lint/pkg/pods"
	"github.com/viglesiasce/kube-lint/pkg/rules"
)

var filename string
var kubeconfig string
var namespace string
var profiles []string
var showAll bool

// podsCmd represents the pods command
var podsCmd = &cobra.Command{
	Use:   "pods",
	Short: "Evaluate rules for pods",
	Long:  `Evaluate all rules marked as kind "Pod"`,
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

		inputPods := []v1.Pod{}
		if kubeconfig != "" {
			inputPods = pods.NewKubeServer(kubeconfig).GetPods(namespace)
		} else if filename != "" {
			inputPods = pods.NewLocalFilesystem(filename).GetPods()
		} else {
			fmt.Println("--filename or --kubeconfig required")
			os.Exit(1)
		}

		if len(inputPods) == 0 {
			fmt.Println("NO PODS FOUND")
			os.Exit(0)
		}
		pods.EvaluateRules(config, inputPods, profiles, showAll)
	},
}

func init() {
	RootCmd.AddCommand(podsCmd)
	podsCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "", "Filename or directory of manifest(s)")
	podsCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file to use for requests")
	podsCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace to use for requests")
	podsCmd.PersistentFlags().StringSliceVarP(&profiles, "profiles", "p", []string{}, "Profiles to check against (all by default)")
	podsCmd.PersistentFlags().BoolVar(&showAll, "show-all", false, "Show passing rules and failing rules")
}
