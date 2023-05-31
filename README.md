# collectivedb
Collective Intelligence Database


- In-memory database
- Self-healing with up to two distributed replicas
- Data is distributed across multiple instances

- Each node has map of all instances
- When instance is added then all maps update
- Hashing on the key to discover which instance has all of the data
    - Hash to number
    - Number is percentage of instances or something and rounded to the nearest depending on number of instances

- Should the instance mapping be distributed as well?
- Performance?
    - Search algorithm....
    - Network should be extremely fast if in a kubernetes cluster 

### Thoughts of how to distribute
- Database can be added to any container -> is a go package
- Or it can just be a docker image
    - Recognizes local docker and kubernetes

## Thoughts on Connecting
- How should it go through and connect between instances to start replication?
  - Can have a few ip addresses added, then it pings and requests for all ips
  - Can be setup to work with a k8s dns or cidr range, then searches across that range automatically for all ips
    - Can determine automatically if it is in a k8s cluster and then automatically search within the deployment that it is part of
      - have this setup as a flag

