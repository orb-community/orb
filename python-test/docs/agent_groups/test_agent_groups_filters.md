## Scenario: Test agent groups filters 
## Steps:
1 - Create multiple agent groups

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}

2 - On agent groups' page (`orb.live/pages/fleet/groups`) use the filter:

* Name
* Description
* Agents
* Tags
* Search by


## Expected Result:

- All filters must be working properly

 
