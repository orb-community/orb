## Scenario: Create policy with net handler 

## 1 - Create a policy with net handler, description, host specification, bpf filter and pcap source

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created

