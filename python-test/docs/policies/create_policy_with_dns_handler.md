## Scenario: Create policy with dns handler 

## 1 - Create a policy dns with description, host specification, bpf filter, pcap source, only qname suffix and only rcode

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created


## 2 - Create a policy dns with host specification, bpf filter, pcap source, only qname suffix and only rcode

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created


## 3 - Scenario: Create a policy dns with bpf filter, pcap source, only qname suffix and only rcode

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created

## 4 - Scenario: Create a policy dns with pcap source, only qname suffix and only rcode

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created


## 5 - Scenario: Create a policy dns with only qname suffix

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


### Expected Result:
- Request must have status code 201 (created) and the policy must be created