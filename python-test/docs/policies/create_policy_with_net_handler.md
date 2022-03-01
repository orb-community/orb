## Scenario: Create policy with net handler 

## 1 - Create a policy net with description, host specification, bpf filter, pcap source, only qname suffix and only rcode

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created


## 2 - Create a policy net with host specification, bpf filter, pcap source, only qname suffix and only rcode

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created


## 3 - Scenario: Create a policy net with bpf filter, pcap source, only qname suffix and only rcode

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created

## 4 - Scenario: Create a policy net with pcap source, only qname suffix and only rcode

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created


## 5 - Scenario: Create a policy net with only qname suffix

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created