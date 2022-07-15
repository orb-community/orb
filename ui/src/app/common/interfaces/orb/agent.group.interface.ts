/**
 * Agent Group Data Model Interface
 *
 * [Fleet Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Fleet}
 */

import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { TagMatch } from 'app/common/interfaces/orb/tag.match.interface';
import { OrbEntity } from 'app/common/interfaces/orb/orb.entity.interface';

export interface AgentGroupState {
  name?: string;
  channel?: string;
}

/**
 * @interface AgentGroup
 */
export interface AgentGroup extends OrbEntity {
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

  /**
   * save summary of matching agents
   * @ignore
   */
  matching_agents?: TagMatch;

  /**
   * keep track of matched agents
   * @ignore
   */
  agents?: Agent[];
}
