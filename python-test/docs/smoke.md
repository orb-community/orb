## Smoke tests

## Login

- [Request registration of a registered account using registered password username and company](login/request_registration_of_a_registered_account_using_registered_password_username_and_company.md)
- [Request registration of a registered account using registered password and username](login/request_registration_of_a_registered_account_using_registered_password_and_username.md)
- [Request registration of a registered account using registered password and company](login/request_registration_of_a_registered_account_using_registered_password_and_company.md)
- [Request registration of a registered account using registered password](login/request_registration_of_a_registered_account_using_registered_password.md)
- [Request registration of a registered account using unregistered password username and company](login/request_registration_of_a_registered_account_using_unregistered_password_username_and_company.md)
- [Request registration of a registered account using unregistered password and username](login/request_registration_of_a_registered_account_using_unregistered_password_and_username.md)
- [Request registration of a registered account using unregistered password and company](login/request_registration_of_a_registered_account_using_unregistered_password_and_company.md)
- [Request registration of a registered account using unregistered password](login/request_registration_of_a_registered_account_using_unregistered_password.md)
- [Request registration of an unregistered account with valid password and invalid email](login/request_registration_of_an_unregistered_account_with_valid_password_and_invalid_email.md)
- [Request registration of an unregistered account with valid password and valid email](login/request_registration_of_an_unregistered_account_with_valid_password_and_valid_email.md)
- [Request registration of an unregistered account with invalid password and valid email](login/request_registration_of_an_unregistered_account_with_invalid_password_and_valid_email.md)
- [Request registration of an unregistered account with invalid password and invalid email](login/request_registration_of_an_unregistered_account_with_invalid_password_and_invalid_email.md)
- [Check if email and password are required fields](login/check_if_email_and_password_are_required_fields.md)
- [Login with valid credentials](login/login_with_valid_credentials.md)
- [Login with invalid credentials](login/login_with_invalid_credentials.md)
- [Request password with registered email address](login/request_password_with_registered_email_address.md)


## Agents

- [Create agent with one tag](agents/create_agent_with_one_tag.md)
- [Edit agent name](agents/edit_agent_name.md)
- [Edit agent tag](agents/edit_agent_tag.md)
- [Save agent without tag](agents/save_agent_without_tag.md)
- [Insert tags in agents created without tags](agents/insert_tags_in_agents_created_without_tags.md)
- [Remove agent using correct name](agents/remove_agent_using_correct_name.md)
- [Run two orb agents on the same port](agents/run_two_orb_agents_on_the_same_port.md)
- [Run two orb agents on different ports](agents/run_two_orb_agents_on_different_ports.md)

## Agent Groups

- [Create agent group with description](agent_groups/create_agent_group_with_description.md)
- [Create agent group with one tag](agent_groups/create_agent_group_with_one_tag.md)
- [Edit agent group name](agent_groups/edit_agent_group_name.md)
- [Edit agent group tag](agent_groups/edit_agent_group_tag.md)
- [Remove agent group using correct name](agent_groups/remove_agent_group_using_correct_name.md)
- [Run two orb agents on the same port](agents/run_two_orb_agents_on_the_same_port.md)
- [Run two orb agents on different ports](agents/run_two_orb_agents_on_different_ports.md)
- [Edit Agent Group name removing name](agent_groups/edit_agent_group_name_removing_name.md)
- [Edit agent group name](agent_groups/edit_agent_group_name.md)
- [Edit agent group description](agent_groups/edit_agent_group_description.md)
- [Edit Agent Group description removing description](agent_groups/edit_agent_group_description_removing_description.md)
- [Edit Agent Group tags to subscribe agent](agent_groups/edit_agent_group_tags_to_subscribe_agent.md)
- [Edit Agent Group tags to unsubscribe agent](agent_groups/edit_agent_group_tags_to_unsubscribe_agent.md)
- [Edit Agent Group removing tags](agent_groups/edit_agent_group_removing_tags.md)
- [Edit Agent Group name, description and tags](agent_groups/edit_agent_group_name,_description_and_tags.md)


## Sinks

- [Create sink with description](sinks/create_sink_with_description.md)
- [Create sink without tags](sinks/create_sink_without_tags.md)
- [Remove sink using correct name](sinks/remove_sink_using_correct_name.md)

## Policies

- [Create policy with description](policies/create_policy_with_description.md)
- [Create policy with dhcp handler](policies/create_policy_with_dhcp_handler.md)
- [Create policy with dns handler](policies/create_policy_with_dns_handler.md)
- [Create policy with net handler](policies/create_policy_with_net_handler.md)
- [Edit policy handler](policies/edit_policy_handler.md)
- [Remove policy using correct name](policies/remove_policy_using_correct_name.md)


## Datasets

- [Create dataset](datasets/create_dataset.md)
- [Remove dataset using correct name](datasets/remove_dataset_using_correct_name.md)


## Integration tests

- [Check if sink is active while scraping metrics](integration/sink_active_while_scraping_metrics.md)
- [Check if sink with invalid credentials becomes active](integration/sink_error_invalid_credentials.md)
- [Provision agent before group (check if agent subscribes to the group)](integration/provision_agent_before_group.md)
- [Provision agent after group (check if agent subscribes to the group)](integration/provision_agent_after_group.md)
- [Provision agent with tag matching existing group linked to a valid dataset](integration/multiple_agents_subscribed_to_a_group.md)
- [Apply multiple policies to a group](integration/apply_multiple_policies.md)
- [Apply multiple policies to a group and remove one policy](integration/remove_one_of_multiple_policies.md)
- [Apply multiple policies to a group and remove one dataset](integration/remove_one_of_multiple_datasets.md)
- [Apply the same policy twice to the agent](integration/apply_policy_twice.md)
- [Remove group (invalid dataset, agent logs)](integration/remove_group.md)
- [Remove sink (invalid dataset, agent logs)](integration/remove_sink.md)
- [Remove policy (invalid dataset, agent logs, heartbeat)](integration/remove_policy.md)
- [Remove dataset (check agent logs, heartbeat)](integration/remove_dataset.md)
- [Remove agent container (logs, agent groups matches)](integration/remove_agent_container.md)
- [Remove agent container force (logs, agent groups matches)](integration/remove_agent_container_force.md)
- [Remove agent (logs, agent groups matches)](integration/remove_agent.md)
- [Subscribe an agent to multiple groups created before agent provisioning](integration/subscribe_an_agent_to_multiple_groups_created_before_agent_provisioning.md)
- [Subscribe an agent to multiple groups created after agent provisioning](integration/subscribe_an_agent_to_multiple_groups_created_after_agent_provisioning.md)
- [Agent subscription to group after editing agent's tags](integration/agent_subscription_to_group_after_editing_agent's_tags.md)
- [Agent subscription to group with policies after editing agent's tags](integration/agent_subscription_to_group_with_policies_after_editing_agent's_tags.md)
- [Edit agent name and apply policies to then](integration/edit_agent_name_and_apply_policies_to_then.md)
- [Insert tags in agents created without tags and apply policies to group matching new tags.md](integration/insert_tags_in_agents_created_without_tags_and_apply_policies_to_group_matching_new_tags.md)
- [Agent unsubscription to group with policies after editing agent group's tags (editing tags after agent provision)](integration/agent_unsubscription_to_group_with_policies_after_editing_agent_group's_tags_editing_tags_after_agent_provision.md)
- [Agent subscription to group with policies after editing agent group's tags (editing tags after agent provision)](integration/agent_subscription_to_group_with_policies_after_editing_agent_group's_tags_editing_tags_after_agent_provision.md)
- [Agent unsubscription to group with policies after editing agent group's tags (editing tags before agent provision)](integration/agent_unsubscription_to_group_with_policies_after_editing_agent_group's_tags_editing_tags_before_agent_provision.md)
- [Agent subscription to group with policies after editing agent group's tags (editing tags before agent provision)](integration/agent_subscription_to_group_with_policies_after_editing_agent_group's_tags_editing_tags_before_agent_provision.md)
- [Agent unsubscription to group with policies after editing agent's tags and agent group's tags (editing tags after agent provision)](integration/agent_unsubscription_to_group_with_policies_after_editing_agent's_tags_and_agent_group's_tags_editing_tags_after_agent_provision.md)
- [Agent subscription to group with policies after editing agent's tags and agent group's tags (editing tags after agent provision)](integration/agent_subscription_to_group_with_policies_after_editing_agent's_tags_and_agent_group's_tags_editing_tags_after_agent_provision.md)
- [Agent unsubscription to group with policies after editing agent's tags and agent group's tags (editing tags before agent provision)](integration/agent_unsubscription_to_group_with_policies_after_editing_agent's_tags_and_agent_group's_tags_editing_tags_before_agent_provision.md)
- [Agent subscription to group with policies after editing agent's tags and agent group's tags (editing tags before agent provision)](integration/agent_subscription_to_group_with_policies_after_editing_agent's_tags_and_agent_group's_tags_editing_tags_before_agent_provision.md)
- [Remotely reset an agent with policies](integration/reset_agent_remotely.md)
- [Remotely reset an agent without policies](integration/reset_agent_remotely.md)