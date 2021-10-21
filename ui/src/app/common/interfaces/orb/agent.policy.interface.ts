/**
 * Agent Policy Data Model
 *
 * [Agent Policy Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Policies-and-Datasets}
 */

import { OrbEntity } from 'app/common/interfaces/orb/orb.entity.interface';
import { PolicyConfig } from 'app/common/interfaces/orb/policy/config/policy.config.interface';


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
   * A timestamp of creation {string}
   */
  ts_created?: string;

  /**
   * Agent backend specific policy data {{[propName: string]: string}}
   */
  policy?: PolicyConfig | any;

  /**
   * Tags {{[propName: string]: string}}
   * User defined tags
   */
  tags?: { [propName: string]: any } | any;

  /**
   * Policy Metadata {{[propName: string]: string}}
   */
  policy_metadata?: { [propName: string]: any } | any;

  window_config?: any;

}
