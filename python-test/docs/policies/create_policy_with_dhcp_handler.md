## Scenario: Create policy with dhcp handler 
## 1 - Create a policy with dhcp handler, description, host specification, bpf filter and pcap source

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created