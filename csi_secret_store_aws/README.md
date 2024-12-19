# AWS Secret Manager Setup Guide 🔐

## Step 1: Create AWS Secrets Manager Secret 🛠️

1. Navigate to AWS Secrets Manager console
2. Click "Store a new secret"
3. Choose secret type (e.g., "Other type of secret")
4. Enter your secret key-value pairs
5. Set a secret name (e.g., "my-application-secrets")
6. Add description (optional)
7. Click "Next" and complete the creation

**Important**: Note down the ARN of your secret, you'll need it for the IAM policy

## Step 2: Create IAM Policy for Secrets Access 👮

1. Go to IAM Console
2. Navigate to Policies
3. Click "Create Policy"
4. Choose JSON and use this policy:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "secretsmanager:GetSecretValue",
                "secretsmanager:DescribeSecret"
            ],
            "Resource": "YOUR-SECRET-ARN-HERE"
        }
    ]
}
```

5. Replace `YOUR-SECRET-ARN-HERE` with the ARN from Step 1
6. Click "Next"
7. Add policy name (e.g., "SecretsManagerReadOnly")
8. Add description (optional)
9. Click "Create Policy"

**Note**: This policy provides minimum required permissions - only read and get value access 🔒

## Same time we need to create the Role for this Policy

**before creating that we need to create OIDC VALIDATE provider**

1. Create Identify provier
2. Select OIDC provider
3. for url Go to you eks cluster there you see the OIDC url copy and past that
4. and for auditon you need to name as **sts.amazonaws.com**
5. create it

### Now we create the Role

1. Create role
2. Select Webidentify
3. in web identify choose the identify provider which we before create
4. Click next and add your policy which we created in **Step 2**
5. And Create the role is create. Use this urn for **service account**

## Step 3: Install Secrets Store CSI Driver 🔌

### A. Add Helm Repository
```bash
helm repo add secrets-store-csi-driver https://kubernetes-sigs.github.io/secrets-store-csi-driver/charts
```

### B. Install CSI Driver with Secret Rotation
```bash
helm install csi-secrets-store secrets-store-csi-driver/secrets-store-csi-driver \
  --namespace kube-system \
  --set enableSecretRotation=true
```

### C. Install AWS Provider
1. Visit the AWS Provider GitHub repository:
   - [AWS Secret Store CSI Driver Provider](https://github.com/aws/secrets-store-csi-driver-provider-aws)
2. Follow provider installation instructions from the repository
3. Verify the installation is complete


### D. Verify Installation ✅
Run the following command to verify the installation:
```bash
kubectl get crd
```

Expected output:
```
secretproviderclasses.secrets-store.csi.x-k8s.io
```

**Important Notes**:
- Ensure your Kubernetes cluster has Helm initialized 🎡
- The CSI driver installation requires cluster admin privileges
- Check [official documentation](https://secrets-store-csi-driver.sigs.k8s.io/getting-started/usage.html) for detailed configuration options
- Monitor logs after installation to ensure proper functionality

## Step 5: Create Secret Provider Class 📝

### A. Verify API Resources
First, verify the available secret-related resources:
```bash
kubectl api-resources | grep secret
```

Expected output:
```
secrets                                                v1                                true         Secret
secretproviderclasses                                  secrets-store.csi.x-k8s.io/v1     true         SecretProviderClass
secretproviderclasspodstatuses                         secrets-store.csi.x-k8s.io/v1     true         SecretProviderClassPodStatus
```

### B. Create Secret Provider Class
Create a YAML file named `secret-provider-class.yaml`:

```yaml
apiVersion: secrets-store.csi.x-k8s.io/v1
kind: SecretProviderClass
metadata:
  name: aws-test-secret-provider
  namespace: test
spec:
  provider: aws
  parameters:
    objects: |
      - objectName: "db-creds"
        objectType: "secretsmanager"
```

Apply the configuration:
```bash
kubectl apply -f secret-provider-class.yaml
```

### C. Verify Secret Provider Class
Check if the Secret Provider Class was created successfully:
```bash
kubectl describe secretproviderclass -n test aws-test-secret-provider
```

Expected output:
```yaml
Name:         aws-test-secret-provider
Namespace:    test
Labels:       
Annotations:  
API Version:  secrets-store.csi.x-k8s.io/v1
Kind:         SecretProviderClass
Metadata:
  Creation Timestamp:  2024-12-19T08:32:21Z
  Generation:         1
  Resource Version:   6912439
  UID:               6e4febe9-cb9c-485b-8a84-28bd5c0813a1
Spec:
  Parameters:
    Objects:  - objectName: "db-creds"
              objectType: "secretsmanager"
  Provider:  aws
Events:     
```

**Important Notes**: 
- The Secret Provider Class stores data in JSON format, not as environment variables 🔄
- Additional environment variable configurations can be added in the Secret Provider Class (covered in future steps)
- Ensure the namespace matches your application's namespace
- The `objectName` should match your AWS Secret name
- Monitor Events field for any issues during creation

---