# kube-lint
A linter for Kubernetes resources with a customizable rule set.

## Introduction
`kube-lint` hopes to make it easy to validate that your Kubernetes configuration files and your running resources
adhere to a standard that you define. You define a list of rules that you would like to validate against your resources
and `kube-lint` will evaluate those rules against them.

In many organizations you will want to have a standard for what is considered "correct" enough to be deployed into 
your Kubernetes clusters. You may have conventions for labels or restrictions on certain types of services being created.
You can use `kube-lint` during your CI/CD pipeline to gate resources being created that do not adhere to your standards.
Additionally you can use kube-lint to audit against a running set of resources in your cluster. 

***CONSIDER THIS A PROTOTYPE. PLEASE PROVIDE FEEDBACK IN THE [ISSUES](https://github.com/viglesiasce/kube-lint/issues)***

***Only Pod linting is currently implemented***

## Installation

- Download a release from the [releases page](https://github.com/viglesiasce/kube-lint/releases/) that matches your platform.
- Extract the archive

### For MacOS
```
wget https://github.com/viglesiasce/kube-lint/releases/download/v0.0.1-prototype/kube-lint-prototype-darwin.tgz
tar zxfv kube-lint-prototype-darwin.tgz
./darwin/kube-lint -h
```

### For Linux
```
wget https://github.com/viglesiasce/kube-lint/releases/download/v0.0.1-prototype/kube-lint-prototype-linux.tgz
tar zxfv kube-lint-prototype-linux.tgz
./linux/kube-lint -h
```

## Rule configuration
The rule configuration file is a YAML formatted list of [KubernetesRules](https://github.com/viglesiasce/kube-lint/blob/master/pkg/rules/rules.go#L44). An example config file is 
available at `example/config.yaml` in this repository.

A KubernetesRule has the following format:
```
name: app-label
description: Includes a label with key "app"
kind: Pod
field: .metadata.labels.app
operator: set
valueType: string
tags:
- operations
- security
```

`name` is an identifier for this rule.

`description` provides details about what the rule is checking for.

`kind` is the type of resource this check should be done against.

`field` is a [jsonpath](https://kubernetes.io/docs/user-guide/jsonpath/) used to get the value you want to evaluate against.

`operator` is the check that youd like to do against your expected vs actual values (ie equal, matches, lessthan). 
For `string` type the available operators are `equal`, `notequal`, `set`, `unset`, `matches`. For `bool` type the available
operators are `equal`, `notequal`, `set`, `unset`. For `float64` type, the available operators are `equal`, `notequal`,
`set`, `unset`, `greaterthan`, `lessthan`.


`valueType` is the type of the value that needs to be evaluated. `string` is the default. `bool` and `float64` are also implemented. 

`tags` is a list of strings that can be used to decide whether to run this rule or not via the CLI. 

## Running kube-lint
Once installed you can run kube-lint from this directory as follows:
```
kube-lint pods --config example/config.yaml
```

To change the rules edit `example/config.yaml`. You rulebender you.

## TODO if this seems like a reasonable approach to pursue
- Replace `panic` everywhere with proper error handling
- Add tests. Lots of tests.
- Add docstrings to all exported functions/types/methods
- Make -f be able to load a directories of yaml files (like kubectl)
- Decide on how to deal with unset parameters
- Choose a logging framework and use it
- Add more resources (services/deployments/etc.)
- Use ${HOME}/.kube-lint for config params
- Develop standardized baseline of rules that are useful
- Vendor dependencies using glide

## Contributing
Add an issue to talk about what youd like to see changed. Lets talk about it then come up with a plan of action. 
