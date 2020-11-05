# collective-db
Collective Intelligence Database


- In-memory database
- Self-healing with up to two distributed replicas
- Data is distributed across multiple instances

- Each node has map of all instances
- When instance is added then all mapus update
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

