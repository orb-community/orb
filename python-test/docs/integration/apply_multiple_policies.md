## 1- Scenario: apply multiple advanced policies to agents subscribed to a group

Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create a sink
4. Create multiple advanced policies (with filters, source pcap)
5. Create a dataset linking the group, the sink and one of the policies
6. Create another dataset linking the same group, sink and the other policy

Expected result:
-
- All the policies must be applied to the agent (orb-agent API response)
- The container logs contain the message "policy applied successfully" referred to each policy
- The container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy
- Referred sink must have active state on response
- Datasets related to all existing policies have validity valid


## 2- Scenario: apply multiple simple policies to agents subscribed to a group

Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create a sink
4. Create 2 policies
5. Create a dataset linking the group, the sink and one of the policies
6. Create another dataset linking the same group, sink and the other policy

Expected result:
-
- All the policies must be applied to the agent (orb-agent API response)
- The container logs contain the message "policy applied successfully" referred to each policy
- The container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy
- Referred sink must have active state on response
- Datasets related to all existing policies have validity valid