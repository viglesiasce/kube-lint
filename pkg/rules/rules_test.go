package rules

import (
	"log"
	"testing"

	"github.com/ghodss/yaml"
)

var pod = `---
apiVersion: v1
kind: Pod
metadata:
  name: redis-django
  labels:
    app: web
spec:
  containers:
    - name: key-value-store
      image: redis
      ports:
        - containerPort: 6379
      resources:
        requests:
          cpu: 100m
          memory: 100Mi
        limits:
          cpu: 100m
          memory: 100Mi
`

func TestRulesOperators(t *testing.T) {
	var ruletests = []struct {
		rule     KubernetesRule
		resource string
	}{
		{KubernetesRule{Operator: "equal", Field: ".metadata.name", Kind: "Pod", Value: "redis-django"}, pod},
		{KubernetesRule{Operator: "equal", Field: ".spec.containers[0].image", Kind: "Pod", Value: "redis"}, pod},
		{KubernetesRule{Operator: "set", Field: ".spec.containers[0].resources.requests", Kind: "Pod"}, pod},
		{KubernetesRule{Operator: "unset", Field: ".spec.containers[0].privileged", Kind: "Pod"}, pod},
	}
	for _, tt := range ruletests {
		t.Run(tt.rule.GetName(), func(t *testing.T) {
			resourceJSONBytes, err := yaml.YAMLToJSON([]byte(tt.resource))
			if err != nil {
				log.Panicf("Unable to convert resource to JSON: %v", err)
			}
			result := tt.rule.Evaluate(resourceJSONBytes)
			if !result.Passed {
				t.Errorf("Test did not pass\nRule: %v\nResult %v", tt.rule, tt.resource)
			}
		})
	}
}
