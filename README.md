# collectivedb
Collective Intelligence Database

This in-memory database is tightly integrated with the application, providing instant access to data and enabling fast retrieval. Moreover, its data is automatically replicated across the entire cluster, ensuring high availability and resilience.

Imagine if Redis and Kafka had a baby, that's how this database is designed to perform. There are still a few kinks and pieces that need to be ironed out, but the database as a whole is pretty close to going live.

### Features
- In-memory database
- Replicated through a hashing algorithm on multiple nodes within database cluster
  - Self-healing and distributed
- gRPC and API layers
- Each node is capable of acting as a leader and contains a map of all instances and location of data across cluster
- Database can autoscale alongside application
