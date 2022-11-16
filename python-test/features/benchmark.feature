@benchmark
Feature: Integrated Benchmark Tests

  @benchmark
  Scenario Outline: Benchmark memory usage multiple policies
    Given the Orb user has a registered account
    And the Orb user logs in
    And that an agent with 1 orb tag(s) already exists and is online
    And pktvisor state is running
    And referred agent is subscribed to 1 group
    And this agent's heartbeat shows that 1 groups are matching the agent
    And that a sink already exists
    When "<amount>" mixed policies are applied to the group
    Then this agent's heartbeat shows that "<amount>" policies are applied and all has status running
    And the container logs contain the message "policy applied successfully" referred to each policy within "<waiting_time>" seconds
    And "<amount>" dataset(s) have validity valid and 0 have validity invalid in "<waiting_time>" seconds
    Examples:
      | amount | waiting_time |
      | 10     | 30           |
      | 20     | 60           |
      | 50     | 120          |
      | 100    | 120          |
      | 200    | 120          |
