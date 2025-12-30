# Kubernetes Network Policy - Multi-Namespace Access Control

This project demonstrates how to implement Kubernetes Network Policies to control traffic between namespaces. We'll create three isolated namespaces and configure a network policy that allows only `project-1` to access `project-2`, while blocking all other namespaces including `project-3`.

## 📋 Project Overview

We have three separate projects deployed across three namespaces:

- **project-1**: Contains `website-1` with environment label `testa`
- **project-2**: Contains `website-2` with environment label `testb` (secured namespace)
- **project-3**: Contains `website-3` with environment label `testc`

## 🎯 Network Policy Goal

**Secure `project-2` namespace** by:
1. Only allowing incoming traffic from `project-1` namespace
2. Blocking all other namespaces (including `project-3`)
3. Demonstrating namespace-level isolation in Kubernetes

## 📁 Project Structure

```
.
├── namespace.yaml       # Creates three namespaces with environment labels
├── deployment.yaml      # Deploys three nginx websites across namespaces
└── networkpolicy.yaml   # Network policy to secure project-2
```

## 🚀 Deployment Steps

### Step 1: Create Namespaces

First, create the three namespaces with their respective environment labels:

```bash
kubectl apply -f namespace.yaml
```

This creates:
- `project-1` with label `env: a`
- `project-2` with label `env: b`
- `project-3` with label `env: c`

**Verify namespaces:**
```bash
kubectl get namespaces --show-labels
```

### Step 2: Deploy Applications

Deploy the three nginx websites to their respective namespaces:

```bash
kubectl apply -f deployment.yaml
```

This creates:
- `website-1` in `project-1` (labeled with `env: testa`)
- `website-2` in `project-2` (labeled with `env: testb`)
- `website-3` in `project-3` (labeled with `env: testc`)

**Verify deployments:**
```bash
kubectl get deployments -A
kubectl get pods -A
```

### Step 3: Apply Network Policy

Apply the network policy to secure `project-2`:

```bash
kubectl apply -f networkpolicy.yaml
```

**Verify network policy:**
```bash
kubectl get networkpolicy -n project-2
kubectl describe networkpolicy allow-project1-testa-to-project2-testb -n project-2
```

## 🔒 Network Policy Explanation

The network policy is applied to the `project-2` namespace and works as follows:

```yaml
spec:
  podSelector:
    matchLabels:
      env: testb          # Applies to pods with label env: testb
  policyTypes:
    - Ingress            # Controls incoming traffic only
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              env: a      # Allow from namespace with env: a (project-1)
          podSelector:
            matchLabels:
              env: testa  # Allow from pods with env: testa
```

### Key Points:

1. **Target**: Pods in `project-2` with label `env: testb`
2. **Allow**: Traffic from `project-1` namespace (labeled `env: a`) AND pods labeled `env: testa`
3. **Block**: All other traffic, including from `project-3`

## 🧪 Testing the Network Policy

### Create Services (Optional)

To test connectivity, you may want to create services:

```bash
# Create service for website-2 in project-2
kubectl expose deployment website-2 -n project-2 --port=80 --target-port=80 --name=website-2-service

# Create service for website-1 in project-1
kubectl expose deployment website-1 -n project-1 --port=80 --target-port=80 --name=website-1-service
```

### Test Access from project-1 (Should Work ✅)

```bash
# Get a shell in project-1 pod
kubectl exec -it -n project-1 $(kubectl get pod -n project-1 -l app=website-1 -o jsonpath='{.items[0].metadata.name}') -- bash

# Try to access website-2
curl website-2-service.project-2.svc.cluster.local
# Expected: Success - HTML response from nginx
```

### Test Access from project-3 (Should Fail ❌)

```bash
# Get a shell in project-3 pod
kubectl exec -it -n project-3 $(kubectl get pod -n project-3 -l app=website-3 -o jsonpath='{.items[0].metadata.name}') -- bash

# Try to access website-2
curl website-2-service.project-2.svc.cluster.local --max-time 5
# Expected: Timeout or connection refused
```

## 📊 Resource Specifications

Each deployment includes:
- **Replicas**: 1 pod
- **Container**: nginx:latest
- **CPU Requests**: 50m
- **CPU Limits**: 70m
- **Memory Requests**: 256Mi
- **Memory Limits**: 512Mi

## 🔍 Troubleshooting

### Check if Network Policy is Applied

```bash
kubectl get networkpolicy -n project-2
kubectl describe networkpolicy -n project-2
```

### View Pod Labels

```bash
kubectl get pods -n project-1 --show-labels
kubectl get pods -n project-2 --show-labels
kubectl get pods -n project-3 --show-labels
```

### Check Namespace Labels

```bash
kubectl get namespace project-1 --show-labels
kubectl get namespace project-2 --show-labels
kubectl get namespace project-3 --show-labels
```

### Debug Connectivity Issues

```bash
# Install curl in nginx container if needed
kubectl exec -it -n project-1 <pod-name> -- apt-get update && apt-get install -y curl

# Check DNS resolution
kubectl exec -it -n project-1 <pod-name> -- nslookup website-2-service.project-2.svc.cluster.local
```

## 🧹 Cleanup

To remove all resources:

```bash
kubectl delete -f networkpolicy.yaml
kubectl delete -f deployment.yaml
kubectl delete -f namespace.yaml
```

Or delete namespaces (which cascades to all resources):

```bash
kubectl delete namespace project-1 project-2 project-3
```

## 📚 Additional Notes

- Network Policies require a CNI plugin that supports them (e.g., Calico, Cilium, Weave Net)
- By default, if no Network Policy is applied, all traffic is allowed
- Once a Network Policy is applied to a namespace, only explicitly allowed traffic is permitted
- Network Policies are additive - multiple policies can apply to the same pod

## 🔐 Security Best Practices

1. **Default Deny**: Consider creating a default deny-all policy first, then whitelist specific traffic
2. **Principle of Least Privilege**: Only allow the minimum necessary access
3. **Label Strategy**: Use consistent labeling strategy across namespaces and pods
4. **Regular Audits**: Periodically review network policies to ensure they match current requirements
5. **Documentation**: Keep network policy documentation up-to-date as your architecture evolves

## 📖 References

- [Kubernetes Network Policies Documentation](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
- [Network Policy Recipes](https://github.com/ahmetb/kubernetes-network-policy-recipes)