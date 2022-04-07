## Scenario: Reset agent remotely 

### Agent with policies applied:
## Steps: 
- Provision an orb-agent and apply 2 policies to them
- Restart the agent through a POST request on `/agents/{agent_id}/rpc/reset` endpoint
- Check the logs and agent's view page


## Expected Result: 
- this agent's heartbeat shows that 2 policies are successfully applied and has status running
- the container logs should contain the message "restarting all backends" within 5 seconds
- the container logs that were output after reset the agent contain the message "removing policies" within 5 seconds
- the container logs that were output after reset the agent contain the message "resetting backend" within 5 seconds
- the container logs that were output after reset the agent contain the message "pktvisor process stopped" within 5 seconds
- the container logs that were output after reset the agent contain the message "reapplying policies" within 5 seconds
- the container logs that were output after reset the agent contain the message "all backends were restarted" within 5 seconds
- the container logs that were output after reset the agent contain the message "completed RPC subscription to group" within 10 seconds
- the container logs that were output after reset the agent contain the message "policy applied successfully" referred to each applied policy within 10 seconds
- the container logs that were output after reset the agent contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds

____
### Agent without policies applied:
## Steps:
- Provision an agent and subscribe them to a group
- Restart the agent through a POST request on `/agents/{agent_id}/rpc/reset` endpoint
- Apply 2 policies to this agent
- Check the logs and agent's view page

## Expected Result:
- the container logs should contain the message "restarting all backends" within 5 seconds
- the container logs that were output after reset the agent contain the message "resetting backend" within 5 seconds
- the container logs that were output after reset the agent contain the message "pktvisor process stopped" within 5 seconds
- the container logs that were output after reset the agent contain the message "reapplying policies" within 5 seconds
- the container logs that were output after reset the agent contain the message "all backends were restarted" within 5 seconds
- the container logs that were output after reset the agent contain the message "completed RPC subscription to group" within 10 seconds
- the container logs that were output after reset the agent contain the message "policy applied successfully" referred to each applied policy within 10 seconds
- the container logs that were output after reset the agent contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
- this agent's heartbeat shows that 2 policies are successfully applied and has status running