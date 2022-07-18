/**
 * Agent Policy Data Model
 *
 * [Agent Policy Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Policies-and-Datasets}
 */

import { OrbEntity } from 'app/common/interfaces/orb/orb.entity.interface';
import { PolicyConfig } from 'app/common/interfaces/orb/policy/config/policy.config.interface';

export enum AgentPolicyStates {
  running = 'running',
  failedToApply = 'failed_to_apply',
}

export interface AgentPolicyState {
  id?: string;
  name?: string;
  state?: AgentPolicyStates;
  error?: string;
  datasets?: string[];
}

/**
 * @interface AgentPolicy
 */
export interface AgentPolicy extends OrbEntity {
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
   * Schema Version
   */
  schema_version?: string;

  /**
   * A timestamp of creation {string}
   */
  ts_created?: string;

  /**
   * Agent backend specific policy data {{[propName: string]: string}}
   */
  policy?: PolicyConfig;

  /**
   * Tags {{[propName: string]: string}}
   * User defined tags
   */
  tags?: { [propName: string]: any } | any;

  /**
   * Format {string}
   * When specified, determines language of
   * policy_data content
   */
  format?: string;

  /**
   * Policy Data {string}
   * Either JSON or YAML compatible string;
   */
  policy_data?: string;
}
