---
title: Learning about Pod Security Policy
date: 2021-01-02
draft: false
summary: Beginner's guide for learning PSP in a kubernetes cluster.
draft: false
authors:
- admin
tags:
- kubernetes
- security
categories:
- kubernetes
- tech
---

{{% toc %}}

## What we will learn

- What is a PSP?
- How to enable and use PSP in a cluster(minikube)
- Create custom psp and use it in a pod

## What is Pod Security Policy?

Pod Security Policies are cluster-wide resources that control security sensitive aspects of pod specification. PSP objects define a set of conditions that a pod must run with in order to be accepted into the system, as well as defaults for their related fields. PodSecurityPolicy is an optional admission controller that is enabled by default through the API, thus policies can be deployed without the PSP admission plugin enabled. This functions as a validating and mutating controller simultaneously.

### Pod Security Policies allow you to control:

- The running of privileged containers
- Usage of host namespaces
- Usage of host networking and ports
- Usage of volume types
- Usage of the host filesystem
- A white list of Flexvolume drivers
- The allocation of an FSGroup that owns the pod’s volumes
- Requirments for use of a read only root file system
- The user and group IDs of the container
- Escalations of root privileges
- Linux capabilities, SELinux context, AppArmor, seccomp, sysctl profile

If you’re interested in more details, check out the [official Kubernetes documentation](https://kubernetes.io/docs/concepts/policy/pod-security-policy/).

## Enabling Pod Security Policies

For learning purpose we will be testing on our minikube cluster.

Enable `PodSecurityPolicy` on your minikube cluster by appening `PodSecurityPolicy` to the apiserver flag in minikube like this:

```bash
--extra-config=apiserver.enable-admission-plugins=PodSecurityPolicy
```

**Important**: As per kubernetes documentation

> Pod security policy control is implemented as an optional (but recommended) admission controller. PodSecurityPolicies are enforced by enabling the admission controller, but doing so without authorizing any policies will prevent any pods from being created in the cluster. Since the pod security policy API (policy/v1beta1/podsecuritypolicy) is enabled independently of the admission controller, for existing clusters it is recommended that policies are added and authorized before enabling the admission controller.

This means if your start your cluster without adding and authorizing policies it will fail to start any new pod. To see this in working, follow below steps:

```bash
# Start your minikube cluster without any policies defined
minikube start -p kautilya-cluster --kubernetes-version=v1.19.6 --feature-gates=EphemeralContainers=true --extra-config=apiserver.enable-admission-plugins=PodSecurityPolicy --addons=pod-security-policy

# Try running a simple pod
kubectl run nginx --image=nginx
```

The above command will fail with the follwing output:

`Error from server (Forbidden): pods "nginx" is forbidden: PodSecurityPolicy: no providers available to validate pod request`

To fix the above problem you have to create a `psp.yaml` file inside `~/.minikube/files/etc/kubernetes/addons` folder. This will fix the issue when the cluster is being bootstraped. Copy the contents of file below and run the command:

```yaml
mkdir -p ~/.minikube/files/etc/kubernetes/addons && \
cat <<EOF | tee ~/.minikube/files/etc/kubernetes/addons/psp.yaml | kubectl apply -f -
---
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: privileged
  annotations:
    seccomp.security.alpha.kubernetes.io/allowedProfileNames: "*"
  labels:
    addonmanager.kubernetes.io/mode: EnsureExists
spec:
  privileged: true
  allowPrivilegeEscalation: true
  allowedCapabilities:
  - "*"
  volumes:
  - "*"
  hostNetwork: true
  hostPorts:
  - min: 0
    max: 65535
  hostIPC: true
  hostPID: true
  runAsUser:
    rule: 'RunAsAny'
  seLinux:
    rule: 'RunAsAny'
  supplementalGroups:
    rule: 'RunAsAny'
  fsGroup:
    rule: 'RunAsAny'
---
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: restricted
  labels:
    addonmanager.kubernetes.io/mode: EnsureExists
spec:
  privileged: false
  allowPrivilegeEscalation: false
  requiredDropCapabilities:
    - ALL
  volumes:
    - 'configMap'
    - 'emptyDir'
    - 'projected'
    - 'secret'
    - 'downwardAPI'
    - 'persistentVolumeClaim'
  hostNetwork: false
  hostIPC: false
  hostPID: false
  runAsUser:
    rule: 'MustRunAsNonRoot'
  seLinux:
    rule: 'RunAsAny'
  supplementalGroups:
    rule: 'MustRunAs'
    ranges:
      # Forbid adding the root group.
      - min: 1
        max: 65535
  fsGroup:
    rule: 'MustRunAs'
    ranges:
      # Forbid adding the root group.
      - min: 1
        max: 65535
  readOnlyRootFilesystem: false
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: psp:privileged
  labels:
    addonmanager.kubernetes.io/mode: EnsureExists
rules:
- apiGroups: ['policy']
  resources: ['podsecuritypolicies']
  verbs:     ['use']
  resourceNames:
  - privileged
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: psp:restricted
  labels:
    addonmanager.kubernetes.io/mode: EnsureExists
rules:
- apiGroups: ['policy']
  resources: ['podsecuritypolicies']
  verbs:     ['use']
  resourceNames:
  - restricted
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: default:restricted
  labels:
    addonmanager.kubernetes.io/mode: EnsureExists
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: psp:restricted
subjects:
- kind: Group
  name: system:authenticated
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: default:privileged
  namespace: kube-system
  labels:
    addonmanager.kubernetes.io/mode: EnsureExists
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: psp:privileged
subjects:
- kind: Group
  name: system:masters
  apiGroup: rbac.authorization.k8s.io
- kind: Group
  name: system:nodes
  apiGroup: rbac.authorization.k8s.io
- kind: Group
  name: system:serviceaccounts:kube-system
  apiGroup: rbac.authorization.k8s.io
EOF
```

Let's check the PSPs created:

```bash
kubectl get psp
NAME         PRIV    CAPS   SELINUX    RUNASUSER          FSGROUP     SUPGROUP    READONLYROOTFS   VOLUMES
privileged   true    *      RunAsAny   RunAsAny           RunAsAny    RunAsAny    false            *
restricted   false          RunAsAny   MustRunAsNonRoot   MustRunAs   MustRunAs   false            configMap,emptyDir,projected,secret,downwardAPI,persistentVolumeClaim

```

Now if you try to create a pod, it will be created successfully.

```bash
kubectl run nginx --image=nginx
pod/nginx created
```

**Note**: One thing you need to see which psp is assigned to the pod. Run the command below:

```bash
kubectl describe po nginx | grep kubernetes.io/psp
Annotations:  kubernetes.io/psp: privileged
```

As you can see that `privileged` psp is assigned to it. Buy why privileged? That's where policy order come in action.

### Policy order

As per kubernetes documentation, it says:

In addition to restricting pod creation and update, pod security policies can also be used to provide default values for many of the fields that it controls. When multiple policies are available, the pod security policy controller selects policies according to the following criteria:

1. PodSecurityPolicies which allow the pod as-is, without changing defaults or mutating the pod, are preferred. The order of these non-mutating PodSecurityPolicies doesn't matter.

2. If the pod must be defaulted or mutated, the first PodSecurityPolicy (ordered by name) to allow the pod is selected.

> Understanding this policy order is a black hole. More your experiment more you go into the black hole. So one thing you need to keep in mind is that you need to match the `po.spec` and `psp.spec` and if more than 2 policies match the criterion, the policy which comes first will be selected. To understand more, we will try to create a custom policy and run it in our pod.

## Creating Custom Kubernetes Pod Security Policy

### Set up

Set up a namespace and a service account to act as for this example. We'll use this service account which will have access to our custom psp and use this inside our pod.

```bash
kubectl create namespace testing-psp
kubectl create serviceaccount -n testing-psp testing-psp
```

### Create a policy

Define the `minimal-psp-restricted` PodSecurityPolicy object. This is a policy that simply prevents the creation of privileged pods, and run user as non root. The name of a PodSecurityPolicy object must be a valid [DNS subdomain name](https://kubernetes.io/docs/concepts/overview/working-with-objects/names#dns-subdomain-names). Run the below command:

```yaml
cat <<EOF | kubectl apply -f -
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: minimal-psp-restricted
spec:
  privileged: false  # Don't allow privileged pods!
  # The rest fills in some required fields.
  seLinux:
    rule: RunAsAny
  supplementalGroups:
    rule: RunAsAny
  runAsUser:
    rule: MustRunAsNonRoot
  fsGroup:
    rule: RunAsAny
  volumes:
  - '*'
EOF
```

Check the psp created:

```bash
kubectl get psp
NAME                     PRIV    CAPS   SELINUX    RUNASUSER          FSGROUP     SUPGROUP    READONLYROOTFS   VOLUMES
minimal-psp-restricted   false          RunAsAny   MustRunAsNonRoot   RunAsAny    RunAsAny    false            *
privileged               true    *      RunAsAny   RunAsAny           RunAsAny    RunAsAny    false            *
restricted               false          RunAsAny   MustRunAsNonRoot   MustRunAs   MustRunAs   false            configMap,emptyDir,projected,secret,downwardAPI,persistentVolumeClaim

```

### Grant access to Service Account

Our Service Account `testing-psp` should  have permission to use the new policy. To check whether it has access run the command:

```bash
kubectl auth can-i use podsecuritypolicy/minimal-psp-restricted --as="system:serviceaccount:testing-psp:testing-psp"
no
```

Create the rolebinding to grant service account `testing-psp` the use `verb` on the `minimal-psp-restricted` policy:

```bash
# Create a role
kubectl create role psp:minimalunprivileged \
    --verb=use \
    --resource=podsecuritypolicy \
    --resource-name=minimal-psp-restricted \
    --namespace=testing-psp

role.rbac.authorization.k8s.io/psp:minimalunprivileged created

# Create a rolebinding
kubectl create rolebinding psp:minimalunprivilegedbinding \
    --role=psp:minimalunprivileged \
    --serviceaccount=testing-psp:testing-psp \
    --namespace=testing-psp

rolebinding.rbac.authorization.k8s.io/psp:minimalunprivilegedbinding created

# Check access to use psp
kubectl auth can-i use podsecuritypolicy/minimal-psp-restricted --as="system:serviceaccount:testing-psp:testing-psp"
yes
```

### Create a pod using custom psp

Let us create a pod using our service account. Run the file below:

```yaml
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: nginx
  name: nginx
spec:
  serviceAccountName: testing-psp
  containers:
  - image: nginx
    name: nginx
    resources: {}
    securityContext:
      runAsNonRoot: true # container should run as roon root
      runAsUser: 65534 # this means running as nobody
      runAsGroup: 65534 # this means running as nogroup
  dnsPolicy: ClusterFirst
  restartPolicy: Always
status: {}
EOF
```

Once the pod is created, let us check which psp is assigned to it

```bash
kubectl describe po nginx | grep kubernetes.io/psp
Annotations:  kubernetes.io/psp: minimal-psp-restricted
```

As you can see that it is using our custom psp `minimal-psp-restricted` because `po.spec` for our `nginx` pod matches `psp.spec` fields from `minimal-psp-restricted` policy.

That’s it. You’ve created and applied your first Kubernetes pod security policy. With the help of this technique, you can greatly enhance the security of your Kubernetes deployments.

### Pro Tip

When creating a custom policy, always start with privileged psp and then start changing fields for `po.spec` according to changes in your policy. This can help you in debugging.

## Did you find this page helpful? Consider sharing it 🙌
