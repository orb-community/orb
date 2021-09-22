/**
 * Agent Data Model Interface
 *
 * [Fleet Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Fleet}
 */

import { OrbEntity } from 'app/common/interfaces/orb/orb.entity.interface';

/**
 * @interface Agent
 */
export interface Agent extends OrbEntity {
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
   * Agent Tags {{[propName: string]: string}}
   * Sent in by the agent when it connects
   */
  agent_tags?: any;

  /**
   * Orb Tags {{[propName: string]: string}}
   * User defined tags
   */
  orb_tags?: any;

  /**
   * Agent Metadata {{[propName: string]: string}}
   * Sent in by agent, defining its capabilities.
   */
  agent_metadata?: any;

  /**
   * State {string} = 'new'|'online'|'offline'|'stale'|'removed'
   * Current Status of the Agent's Connection
   */
  state?: string;

  /**
   * Last Heartbeat Data {{[propName: string]: string}}
   */
  last_hb_data?: any;

  /**
   * Last Heartbeat timestamp {string}
   */
  ts_last_hb?: string;

  /**
   * Boolean which indicates whether the agent is in an error state or not.
   * Heartbeat data contains error information.
   */
  error_state?: boolean;

  key?: string;

}
