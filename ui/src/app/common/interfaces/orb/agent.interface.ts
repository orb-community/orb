/**
 * Agent Data Model Interface
 *
 * [Fleet Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Fleet}
 */

import { OrbEntity } from 'app/common/interfaces/orb/orb.entity.interface';
import { AgentGroupState } from './agent.group.interface';
import { AgentPolicyState } from './agent.policy.interface';

/**
 * @enum AgentStates
 */
export enum AgentStates {
  new = 'new',
  online = 'online',
  offline = 'offline',
  stale = 'stale',
  removed = 'removed',
}

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
  agent_metadata?: { [propname: string]: any };

  /**
   * State {string} = 'new'|'online'|'offline'|'stale'|'removed'
   * Current Status of the Agent's Connection
   */
  state?: string;

  /**
   * Last Heartbeat Data {{[propName: string]: string}}
   */
  last_hb_data?:
    | any
    | {
        backend_state?: any;
        group_state?: { [id: string]: AgentGroupState };
        policy_state?: { [id: string]: AgentPolicyState };
      };

  /**
   * Last Heartbeat timestamp {string}
   */
  ts_last_hb?: string;

  /**
   * Boolean which indicates whether the agent is in an error state or not.
   * Heartbeat data contains error information.
   */
  error_state?: boolean;

  /**
   * MQTT KEY
   */
  key?: string;

  /**
   * Combines tags for display in UI
   * Internal use
   * See
   */
  combined_tags?: any;
}
