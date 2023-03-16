# Overview
Business logic behind the replication between nodes.

## Main Thinking
- There is a replicaCount environment variable set that will dictate what level of replication between nodes should occur. _This needs to be set at the cluster level when the first node is created rather than as an environment variable_
- When adding nodes to the collective, it will first try to add the new node to a replica group that isn't full. 
- If all replica groups are full, then it will create a new group

### Not Full Replica Group
- We never want data stored on less nodes than requested through the replica group, so when there are less nodes than what the replicaCount is set at then there is a secondaryGroup used.
- The secondaryGroup will store the data on all of it's nodes, but it will not be used for retrieving the data itself. _or should it be?_

### Becoming a Full Replica Group
- At this point, the secondaryGroup will be removed from the replicaGroup collective data entry and all of the data for the replicaGroup will be removed from the secondaryGroup

### Becoming Part of a Replica Group
- When a node becomes part of a replica group, it will pull all of the data from the replica group nodes and it should immediately set a secondaryGroup

### Deleting a Replica Group
- When the last node from a replica group is removed, then all of the data will fail over to where it has been stored on the secondaryGroup
  - All of the collective data entries will be updated to reflect their location as being on the new replicaGroup (what once was the secondaryGroup)