GOOS=linux GOARCH=amd64 go build .

podman build --platform=linux/amd64 -t collective . 
podman tag collective 192.168.1.19:30500/collective:latest
podman push 192.168.1.19:30500/collective:latest --tls-verify=false

COLLECTIVE_CONTAINER=`kubectl get pods | grep collective | awk '{ print $1 }'`
kubectl delete pod $COLLECTIVE_CONTAINER