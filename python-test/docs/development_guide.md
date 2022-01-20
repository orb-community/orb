| Integration Scenario                                                                | Automated via API  | Automated via UI | Smoke               | Sanity              | 
|-------------------------------------------------------------------------------------|--------------------|------------------|---------------------|---------------------|
| Check if sink is active while scraping metrics                                      | <center>✅</center> |                  | <center>👍</center> | <center>👍</center> |
| Check if sink with invalid credentials becomes active                               |                    |                  |                     |
| Check if after 30 minutes without data sink becomes idle                            |                    |                  |                     |                     |
| Provision agent before group (check if agent subscribes to the group)               | <center>✅</center> |                  | <center>👍</center> | <center>👍</center> |
| Provision agent after group (check if agent subscribes to the group)                |                    |                  | <center>👍</center> | <center>👍</center> |
| Provision agent with tag matching existing group linked to a valid dataset          |                    |                  | <center>👍</center> | <center>👍</center> |
| Apply multiple policies to a group                                                  | <center>✅</center> |                  | <center>👍</center> | <center>👍</center> |
| Apply multiple policies to a group and remove one policy                            | <center>✅</center> |                  | <center>👍</center> | <center>👍</center> |
| Apply multiple policies to a group and remove all of them                           |
| Apply multiple policies to a group and remove one dataset                           | <center>✅</center> |                  | <center>👍</center> | <center>👍</center> |
| Apply multiple policies to a group and remove all datasets                          |
| Apply the same policy twice to the agent                                            |
| Delete sink linked to a dataset, create another one and edit dataset using new sink |
| Remove one of multiples datasets that apply the same policy to the agent            |                    |                  ||
| Remove group (invalid dataset, agent logs)                                          |                    |                  | <center>👍</center> | <center>👍</center> |
| Remove sink (invalid dataset, agent logs)                                           |                    |                  | <center>👍</center> | <center>👍</center> |
| Remove policy (invalid dataset, agent logs, heartbeat)                              | <center>✅</center> |                  | <center>👍</center> | <center>👍</center> |
| Remove dataset (check agent logs, heartbeat)                                        | <center>✅</center> |                  | <center>👍</center> | <center>👍</center> |
| Remove agent container (logs, agent groups matches)                                 |                    |                  | <center>👍</center> | <center>👍</center> |
| Remove agent container force (logs, agent groups matches)                           |                    |                  | <center>👍</center> | <center>👍</center> |
| Remove agent (logs, agent groups matches)                                           |                    |                  | <center>👍</center> | <center>👍</center> |