/**
 * Agent Data Model Interface
 *
 * [Fleet Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Fleet}
 */

/**
 * @interface Agent
 */
export interface Agent {
  /**
   *  ID {string} UUIDv4 (read only)
   */
  id?: string;

  /**
   * Name {string} [a-zA-Z_:][a-zA-Z0-9_]*
   */
  name?: string;

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

  // TODO why not go with status as Sink?
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
  ts_lst_hb?: string;
}
