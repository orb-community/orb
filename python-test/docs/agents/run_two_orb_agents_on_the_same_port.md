## Scenario: Run two orb agents on the same port 

## Steps:
1 - Provision an agent
2 - Provision another agent on same port

## Expected Result:
- Second container must be exited
- the container logs should contain the message "agent startup error"