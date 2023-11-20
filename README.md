# meta-oci-sync


## Setup

Install Metacontroller: https://metacontroller.github.io/metacontroller/guide/install.html

Apply CRD:
```
k apply -f manifests/crds.yaml 
```
Apply custom controller:
```
k apply -f manifests/metacontroller.yaml
```

Build image:
```
~/.rd/bin/nerdctl --namespace k8s.io build -t mirror-controller:v1.0.0 .
```

Create ConfigMap
```
kubectl create configmap oci-sync-controller --from-file=main.go -o yaml --dry-run > manifests/controller-configmap.yml
```

Deploy Controller
```
k apply -f manifests/controller-configmap.yml
k apply -f manifests/controller-deployment.yml
```

Test 