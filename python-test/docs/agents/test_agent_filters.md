## Scenario: Test agent filters 
## Steps: 
1 - Create multiple agents

- REST API Method: POST
- endpoint: /agents
- header: {authorization:token}

2 - On agents' page (`orb.live/pages/fleet/agents`) use the filter:

   * Name
   * Status
   * Tags
   * Search by

 
## Expected Result: 

- All filters must be working properly
