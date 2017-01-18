package pods

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

type PodList interface {
	GetPods() []v1.Pod
}

type KubeServer struct {
	kubeconfig string
}

func NewKubeServer(kubeconfig string) KubeServer {
	return KubeServer{kubeconfig: kubeconfig}
}

func (ks KubeServer) GetPods(namespace string) []v1.Pod {
	k8sClientConfig, err := clientcmd.BuildConfigFromFlags("", ks.kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(k8sClientConfig)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.Core().Pods(namespace).List(v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return pods.Items
}

type LocalFilesystem struct {
	path string
}

func NewLocalFilesystem(path string) LocalFilesystem {
	return LocalFilesystem{path: path}
}

func (lfs LocalFilesystem) GetPods() []v1.Pod {
	resourceFile, err := ioutil.ReadFile(lfs.path)

	if err != nil {
		panic("Unable to read pod file")
	}

	resourceJSONBytes, err := yaml.YAMLToJSON(resourceFile)
	if err != nil {
		panic("Unable to convert resource to JSON")
	}

	if err != nil {
		panic("Unable to Unmarshal json")
	}
	pod := v1.Pod{}
	err = json.Unmarshal(resourceJSONBytes, &pod)
	pods := []v1.Pod{pod}
	return pods
}

func CreateTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	//  "PROFILE", "OPERATOR", "EXPECTED", "ACTUAL",
	table.SetHeader([]string{"POD", "DESCRIPTION", "RESULT"})
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	return table
}

func EvaluateRules(config rules.LinterConfig, pods []v1.Pod, profiles []string, showAll bool) {
	table := CreateTable()
	activeProfiles := rules.LinterConfig{}
	if len(profiles) == 0 {
		// Evaluate all rules by default
		activeProfiles = config
	} else {
		// Filter out active profiles
		for _, profile := range profiles {
			activeProfiles[profile] = config[profile]
		}
	}
	for _, pod := range pods {
		for _, value := range activeProfiles {
			for _, rule := range value {
				// TODO need to be able to provide a list of resources to test against
				if rule.Kind != "Pod" && rule.Kind != "" {
					continue
				}
				k8sRule := rules.NewKubernetesRule(rule.Operator, rule.Field, rule.Value, rule.ValueType)
				resourceJSON, err := json.Marshal(pod)
				if err != nil {
					panic(err)
				}
				result := k8sRule.Evaluate(resourceJSON)
				if err != nil {
					panic(err)
				}
				var colorizedResult string
				if result.Passed {
					if showAll {
						colorizedResult = color.GreenString("passed")
						table.Append([]string{pod.Name, rule.Description, colorizedResult})
					}
				} else {
					colorizedResult = color.RedString("failed")
					table.Append([]string{pod.Name, rule.Description, colorizedResult})
				}
			}
		}
	}
	table.Render()
}
