## Scenario: Create policy with dns handler 

## 1 - Create a policy with dns handler, description, host specification, bpf filter, pcap source, only qname suffix and only rcode

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created


## 2 - Create a policy with dns handler, host specification, bpf filter, pcap source, only qname suffix and only rcode

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created


## 3 - Scenario: Create a policy with dns handler, bpf filter, pcap source, only qname suffix and only rcode

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created

## 4 - Scenario: Create a policy with dns handler, pcap source, only qname suffix and only rcode

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created


## 5 - Scenario: Create a policy with dns handler, only qname suffix

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created