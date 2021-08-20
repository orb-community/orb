/**
 * Agent Group Data Model Interface
 *
 * [Fleet Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Fleet}
 */

import { Agent } from 'app/common/interfaces/orb/agent.interface';

/**
 * @interface AgentGroup
 */
export interface AgentGroup {
// #region public API
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
   * A timestamp of creation {string}
   */
  ts_created?: string;

  /**
   * Channel ID {string}
   * Comm. Ch. ID
   * Unique to this agent
   */
  channel_id?: string;

  /**
   * Tags {{[propName: string]: string}}
   */
  tags?: any;
// #endregion

// #region internal API
  /**
   * validate rule only?
   * @ignore
   */
  validate_only?: boolean;

  /**
   * save summary of matching agents
   * @ignore
   */
  matching_agents?: {
    total: number;
    online: number;
  };

  /**
   * keep track of matched agents
   * @ignore
   */
  agents?: Agent[];
// #endregion
}
