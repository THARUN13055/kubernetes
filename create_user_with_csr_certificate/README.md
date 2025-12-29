# Kubernetes Certificate Signing Request (CSR) Guide

This guide walks you through creating and approving certificates for Kubernetes user authentication using Certificate Signing Requests.

## Reference Documentation
- [Kubernetes Certificate Signing Requests](https://kubernetes.io/docs/reference/access-authn-authz/certificate-signing-requests/)

## Prerequisites
- Access to a Kubernetes cluster
- `kubectl` configured with admin privileges
- `openssl` installed on your system

---

## Step 1: Generate Private Key

First, create a 2048-bit RSA private key:

```bash
openssl genrsa -out tharun.key 2048
```

**Note:** The original command `openssl key 2048 tharun.key` was incorrect. The correct command is `openssl genrsa`.

---

## Step 2: Create Certificate Signing Request (CSR)

Generate a CSR using the private key:

```bash
openssl req -new -key tharun.key -out tharun.csr -subj "/CN=tharun/O=developers"
```

**Changes made:**
- Corrected the command from `openssl req -new -key rsa:2048 -keyout tharun.key -out tharun.csr`
- Used `-key` to specify existing key (not `-keyout` which creates a new key)
- Added `-subj` flag to specify the Common Name (CN) and Organization (O) without interactive prompts

When prompted (if not using `-subj`), enter:
- **Common Name (CN)**: `tharun` (this will be the username)
- **Organization (O)**: `developers` (optional, used for group membership)

---

## Step 3: Encode CSR in Base64

Convert the CSR to base64 format (required for Kubernetes):

```bash
cat tharun.csr | base64 | tr -d '\n'
```

**Note:** Changed `tr "\n"` to `tr -d '\n'` to properly delete newlines.

Copy the output - you'll need it for the next step.

---

## Step 4: Create Kubernetes CSR YAML

Create a file named `csr.yaml`:

```yaml
apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: tharun-csr
spec:
  request: <PASTE_BASE64_CSR_HERE>
  signerName: kubernetes.io/kube-apiserver-client
  expirationSeconds: 86400  # 1 day (optional)
  usages:
  - client auth
```

**Important:** Replace `<PASTE_BASE64_CSR_HERE>` with the base64-encoded CSR from Step 3.

---

## Step 5: Apply the CSR

Submit the CSR to Kubernetes:

```bash
kubectl apply -f csr.yaml
```

Verify the CSR was created:

```bash
kubectl get csr
```

You should see `tharun-csr` with status `Pending`.

---

## Step 6: Approve the Certificate

Approve the certificate signing request:

```bash
kubectl certificate approve tharun-csr
```

**Note:** Corrected from `kubectl certificate approve csr tharuncsr` - the correct syntax is just the CSR name.

Verify approval:

```bash
kubectl get csr tharun-csr
```

Status should now show `Approved,Issued`.

**Troubleshooting:** If the certificate is not issued after approval, verify:
- The `signerName` is correct: `kubernetes.io/kube-apiserver-client`
- Check controller-manager logs: `kubectl logs -n kube-system <controller-manager-pod>`

---

## Step 7: Retrieve the Signed Certificate

Get the signed certificate:

```bash
kubectl get csr tharun-csr -o jsonpath='{.status.certificate}' | base64 -d > tharun.crt
```

**Note:** Changed from `kubectl get crt -o yaml` which was incorrect. The correct command retrieves the certificate from the CSR object.

---

## Step 8: Configure kubectl User

Now configure kubectl to use the new certificate. Three main components are needed:
- **Cluster**: The Kubernetes cluster endpoint
- **User**: Credentials (certificate and key)
- **Context**: Combination of cluster + user

### Set User Credentials

```bash
kubectl config set-credentials tharun \
  --client-certificate=tharun.crt \
  --client-key=tharun.key \
  --embed-certs=true
```

---

## Step 9: Create Context

Create a context that combines the user with the cluster:

```bash
kubectl config set-context tharun-context \
  --cluster=kind-kind \
  --user=tharun
```

**Note:** Changed context name to `tharun-context` for clarity. Replace `kind-kind` with your actual cluster name (check with `kubectl config get-clusters`).

---

## Step 10: Switch to New Context

Switch to the new user context:

```bash
kubectl config use-context tharun-context
```

Test access (note: the user has no permissions yet):

```bash
kubectl get pods
```

You'll likely see: `Error from server (Forbidden)` - this is expected as the user needs RBAC permissions.

---

## Step 11: Grant Permissions (Optional)

To grant permissions, switch back to admin context:

```bash
kubectl config use-context <admin-context>
```

Create a RoleBinding or ClusterRoleBinding:

```bash
# Example: Grant view access in default namespace
kubectl create rolebinding tharun-view \
  --clusterrole=view \
  --user=tharun \
  --namespace=default
```

---

## Summary of Changes Made

1. **Step 1**: Corrected `openssl key` to `openssl genrsa`
2. **Step 2**: Fixed CSR generation command and added `-subj` for non-interactive mode
3. **Step 3**: Changed `tr "\n"` to `tr -d '\n'`
4. **Step 4**: Clarified YAML structure and added required fields
5. **Step 6**: Corrected approve command syntax
6. **Step 7**: Fixed certificate retrieval command (was `kubectl get crt`, should be extracting from CSR)
7. **Step 8**: Added proper `kubectl config set-credentials` command
8. **Step 9**: Corrected context creation command and naming
9. **Overall**: Added proper structure, explanations, and troubleshooting tips

---

## Quick Reference Commands

```bash
# Generate key
openssl genrsa -out tharun.key 2048

# Generate CSR
openssl req -new -key tharun.key -out tharun.csr -subj "/CN=tharun/O=developers"

# Encode CSR
cat tharun.csr | base64 | tr -d '\n'

# Apply CSR
kubectl apply -f csr.yaml

# Approve CSR
kubectl certificate approve tharun-csr

# Get certificate
kubectl get csr tharun-csr -o jsonpath='{.status.certificate}' | base64 -d > tharun.crt

# Configure user
kubectl config set-credentials tharun --client-certificate=tharun.crt --client-key=tharun.key --embed-certs=true

# Create context
kubectl config set-context tharun-context --cluster=kind-kind --user=tharun

# Use context
kubectl config use-context tharun-context
```