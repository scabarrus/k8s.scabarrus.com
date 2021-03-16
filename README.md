# k8s.scabarrus.com
This repository contains services to understand how kubernetes webhook can be implemented

Currently, module contains or will contains following project:
* [authn-webhook] 
* authz-webhook
* validate-webhook (not yet)
* mutate-webhook (not yet)
* audit-webhook (not yet)

## Lab environment
My lab environment is a minimal k8s cluster with one master and one worker (poor cluster :-)).

## Project layout 
```
.
├── img
└── webhook
    ├── cmd
    │   ├── authn
    │   ├── authz
    │   └── mutate
    ├── kubernetes
    │   ├── authn
    │   └── authz
    └── pkg
        ├── authn
        └── authz
```

| Folder		| Description					|
|-----------------------|-----------------------------------------------|
| wehook		| module					|
| webhook/cmd		| main per project (authn| authz | mutate	|
| webhook/kubernetes	| config file per project			|
| webhook/pkg		| package use per project			|


## Authn

### Description

The authn webhook is just an example or a starting project to understand how authentication webhook can be implemented (Not use it in production).
It not check the user in repository but take your token and split it to find which user and group of the user (I said not use it in production :)).

The token format is <user>@<secret>:<group1,group2...>

The request send by kubectl for user authentication is :
```
{
  "apiVersion": "authentication.k8s.io/v1beta1",
  "kind": "TokenReview",
  "spec": {
    "token": "user1@mysecret:admin,dev"
  }
}
```

The response should be:
```
{
    "metadata": {
        "creationTimestamp": null
    },
    "spec": {},
    "status": {
        "authenticated": true,
        "user": {
            "username": "user1",
            "groups": [
                "admin",
                "dev"
            ]
        }
    }
}
```


### Create image
```
# docker build -t authentication-webhook:1.0 -f k8s.scabarrus.com/webhook/kubernetes/authn/Dockerfile .
```

### Create certificate
I have use the same certificate of my kubernetes cluster, but it's not mandatory
```
# openssl genrsa -out server.key 2048 #Generate key
# openssl req -new -key server.key -out server.csr  #Create Certificate Service Request
# echo "subjectAltName = IP:${WOKRER_IP}" > extfile.cnf #Add IP of Worker in SAN List
# openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -extfile extfile.cnf
```


### Deploy image
Create secret :
```
# mysecret=$(echo -n "mysecret" | base64)
# kubectl create secret authn-secret --from-literal=SECRET=${mysecret} 
```

Create configmap with certificate and key:
```
# kubectl create configmap authn-configmap --from-file=server.crt --from-file=server.key
```
Create deployment and service :
```
# kubectl apply -f k8s.scabarrus.com/webhook/kubernetes/authn/authn-deploy.yaml #Deploy pod
# kubectl expose deployment authn-webhook --type=NodePort --port=8080 --target-port=8080 # Expose as nodeport service
```

### Configure api server

Copy  authn-config file in a folder on master node 

Add option to api-server in file /etc/kubernetes/manifests/kube-apiserver.yaml

```
- --authentication-token-webhook-config-file=/etc/authn-config.yaml
- --authentication-token-webhook-version=v1
```
NB: option v1 is important to match with struct managed by this project.

Mount this file in kube-apiserver:
``` 
 volumes:
  - hostPath:
      path: /root/authn-config.yaml
    name: authn-config
```

## Authz
The authz webhook just check if a user has right to perform an action on following condition:
- if namespace is webhook
- if user's group is admin he can perform all action on resources
- if user's group is dev he can get or list resource(s)

As said before don't use it in production.

When kubectl is used the step after authentication is authorization step.
The request sent should look like:
```
{
  "apiVersion": "authorization.k8s.io/v1beta1",
  "kind": "SubjectAccessReview",
  "spec": {
    "resourceAttributes": {
      "namespace": "webhook",
      "verb": "get",
      "group": "",
      "version": "v1",
      "resource": "pods",
       "subresource": ",name:authn-webhook-6dcb76c688-n94zk,"
      
    },
    "user": "bob",
    "groups": ["utut", "admin"]
  }
}
```

The response should look like:
```
{
    "metadata": {
        "creationTimestamp": null
    },
    "spec": {},
    "status": {
        "allowed": true
    }
}
```
### Deploy image

Create configmap with certificate and key:
```
# kubectl create configmap authn-configmap --from-file=server.crt --from-file=server.key
```
Create deployment and service :
```
# kubectl apply -f k8s.scabarrus.com/webhook/kubernetes/authz/authz-deploy.yaml #Deploy pod
# kubectl expose deployment authz-webhook --type=NodePort --port=8080 --target-port=8080 # Expose as nodeport service
```

### Configure api server

Copy  authn-config file in a folder on master node

Add option to api-server in file /etc/kubernetes/manifests/kube-apiserver.yaml

```
- --authorization-webhook-version=v1
- --authorization-mode=Node,RBAC,Webhook
- --authorization-webhook-config-file=/etc/authz-config.yaml
```
NB: option v1 is important to match with struct managed by this project.

Mount this file in kube-apiserver:
```
 volumes:
  - hostPath:
      path: /root/authz-config.yaml
    name: authz-config
```

### Configure kubeconfig 
```
apiVersion: v1
clusters:
- cluster:
    certificate-authority: ./k8s.scabarrus.com/webhook/kubernetes/authz/ca.crt
    server: https://192.168.169.128:6443
  name: authn-webhook
contexts:
- context:
    cluster: authn-webhook
    namespace: webhook
    user: marvel
  name: authn-webhook
current-context: authn-webhook
kind: Config
preferences: {}
users:
- name: bob
  user:
    token: bob@mysecret:admin,tutu
```
NB: the certificate used for the webhook is configured to be able to validate the response from the webhook

### Test
Specify the user and kubeconfig file
```
# kubectl  --kubeconfig=admin.conf get pods  --user bob
kubectl  --kubeconfig=admin.conf get pods  --user bob
NAME                             READY   STATUS    RESTARTS   AGE
authn-webhook-6dcb76c688-n94zk   1/1     Running   0          17h
authz-webhook-79455f668f-h4xlc   1/1     Running   0          5h53m

```
