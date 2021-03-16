## Authn webhook

### Build image from the source 
```
# cd /root/kubernetes/golang/src

# docker build -t authentication-webhook:1.0 -f k8s.scabarrus.com/webhook/kubernetes/authn/Dockerfile .
```

### Test image 
```
# docker run -d -p 8081:8080 authentication-webhook:1.0
curl -X POST -H "Content-Type: application/json" \
--data '{ 
  "apiVersion": "authentication.k8s.io/v1beta1",
  "kind": "TokenReview",
  "spec": {
    "token": "user1@mysecret:admin,dev"
  }
}' --insecure  https://localhost:8080/
```


### Save image and copy it on the worker node (Better to use a docker registry inside of your cluster)
```
# docker save -o authentication-webhook.tar authentication-webhook:1.0
```

### Create your certificate
```
#  openssl genrsa -out server.key 2048
# openssl req -new -key server.key -out server.csr
#  openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt
#  openssl x509 -in server.crt -text
```

### Create your deployment
```
# kubectl create deployment authn-webhook --replica=1 --image authentication-webhook:1.0.0

```

### Expose your deployment
# kubectl expose deploy authn-webhook --name=authn-webhook-deploy --port=8080 --target-port=8080 --type=NodePort
