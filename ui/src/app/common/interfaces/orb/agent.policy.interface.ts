/**
 * Agent Policy Data Model
 *
 * [Agent Policy Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Policies-and-Datasets}
 */

/**
 * @interface AgentPolicy
 */
export interface AgentPolicy {
  /**
   *  Tenant owner {string} UUIDv4 (read only)
   */
  mf_owner_id?: string;

  /**
   *  ID {string} UUIDv4 (read only)
   */
  id?: string;

  /**
   * Name {string} [a-zA-Z_:][a-zA-Z0-9_]*
   */
  name?: string;

  /**
   * Description {string}
   */
  description?: string;

  /**
   * Agent Backend this policy is for: {string}
   * Cannot change once created (read only)
   */
  backend?: string;

  /**
   * Version
   * monotonically increases on each update
   * starts at 0
   */
  version?: number;

  /**
   * A timestamp of creation {string}
   */
  ts_created?: string;

  /**
   * Tags {{[propName: string]: string}}
   * User defined tags
   */
  tags?: any;

  /**
   * Policy Metadata {{[propName: string]: string}}
   */
  policy_metadata?: any;
}
