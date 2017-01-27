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

package cmd

import (
	"io/ioutil"
	"os/user"

	"k8s.io/client-go/pkg/api/v1"

	"fmt"
	"os"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/viglesiasce/kube-lint/pkg/pods"
	"github.com/viglesiasce/kube-lint/pkg/rules"
)

var filename string
var configFile string
var kubeconfig string
var namespace string
var tags []string
var showAll bool

// podsCmd represents the pods command
var podsCmd = &cobra.Command{
	Use:   "pods",
	Short: "Evaluate rules for pods",
	Long:  `Evaluate all rules marked as kind "Pod"`,
	Run: func(cmd *cobra.Command, args []string) {
		// load config
		if configFile == "" {
			fmt.Println("[ERROR] Pass your linter config file using --config or -c")
			os.Exit(1)
		}
		var config rules.LinterConfig
		configFile, err := ioutil.ReadFile(configFile)
		if err != nil {
			panic("Unable to read config file")
		}
		err = yaml.Unmarshal(configFile, &config)
		if err != nil {
			panic("Unable to unmarshal config file")
		}

		inputPods := []v1.Pod{}
		if filename != "" {
			fmt.Println("Getting pods for", filename)
			inputPods = pods.NewLocalFilesystem(filename).GetPods()
		} else {
			// kubeconfig has a default value so will always be populated
			fmt.Println("Getting pods for", kubeconfig)
			inputPods = pods.NewKubeServer(kubeconfig).GetPods(namespace)
		}

		if len(inputPods) == 0 {
			fmt.Println("NO PODS FOUND")
			os.Exit(0)
		}
		pods.EvaluateRules(config, inputPods, tags, showAll)
	},
}

func init() {
	user, err := user.Current()
	if err != nil {
		fmt.Println("Unable to determine current user")
		os.Exit(1)
	}
	defaultKubeConfig := fmt.Sprintf("%s/%s", user.HomeDir, ".kube/config")
	RootCmd.AddCommand(podsCmd)
	podsCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "", "Filename or directory of manifest(s)")
	podsCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Filename of config file with linter rules")
	podsCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", defaultKubeConfig, "Path to the kubeconfig file to use for requests")
	podsCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace to use for requests")
	podsCmd.PersistentFlags().StringSliceVarP(&tags, "tags", "t", []string{}, "Tags used to filter rules (all by default)")
	podsCmd.PersistentFlags().BoolVar(&showAll, "show-all", false, "Show passing rules and failing rules")
}
