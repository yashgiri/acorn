---
title: "acorn init"
---
## acorn init

Install and configure acorn in the cluster

```
acorn init [flags]
```

### Examples

```

acorn init
```

### Options

```
      --cluster-domains strings               The externally addressable cluster domain (default .local.on-acorn.io)
  -h, --help                                  help for init
      --image string                          Override the default image used for the deployment
      --ingress-class-name string             The ingress class name to assign to all created ingress resources (default '')
      --mode string                           Initialize only 'config', 'resources', or 'both' (default 'both')
  -o, --output string                         Output manifests instead of applying them (json, yaml)
      --pod-security-enforce-profile string   The name of the PodSecurity profile to set (default baseline)
      --set-pod-security-enforce-profile      Set the PodSecurity profile on created namespaces (default true)
      --tls-enabled                           If true HTTPS URLs will be rendered for HTTP endpoint URLs (default false)
```

### Options inherited from parent commands

```
  -A, --all-namespaces      Namespace to work in
      --context string      Context to use in the kubeconfig file
      --kubeconfig string   Location of a kubeconfig file
      --namespace string    Namespace to work in (default "acorn")
```

### SEE ALSO

* [acorn](acorn.md)	 - 
