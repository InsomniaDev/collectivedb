# Overview
[Grafana K6](https://k6.io/docs/) is used for load testing the CollectiveDB application

## Technical Overview

Pull the Grafana K6 docker image so that you can run tests using Docker without having to install K6 on the local development computer
```
docker pull grafana/k6
--- or ---
podman pull grafana/k6
```

Running it as a batch job within the Kubernetes cluster itself:
```
apiVersion: batch/v1
kind: Job
metadata:
  name: k6
spec:
  template:
    spec:
      containers:
      - name: k6
        image: grafana/k6
        command: ["run",  "- <script.js"]
      restartPolicy: Never
  backoffLimit: 4
```
There is a k8s job yml file in this same directory

Create a configmap that will use the scripts within the scripts directory
```
kubectl create configmap k6-scripts --from-file=scripts/
```

Describe the configmap 
```
kubectl describe configmap k6-scripts
```

Delete the configmap
```
kubectl delete configmap k6-scripts
```