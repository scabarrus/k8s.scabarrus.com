# cd /root/kubernetes/golang/src

# docker build -t authorization-webhook:1.0 -f k8s.scabarrus.com/webhook/kubernetes/authz/Dockerfile .

# docker run -d -p 8082:8080 authorization-webhook:1.0

# docker save -o authorization-webhook.tar authorization-webhook:1.0


[root@rec-apache-1 kubernetes]# cd authz
[root@rec-apache-1 authz]# cat extfile.cnf
subjectAltName = IP:192.168.169.129
[root@rec-apache-1 authz]# openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -extfile extfile.cnf

