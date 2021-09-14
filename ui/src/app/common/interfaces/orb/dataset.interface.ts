/**
 * Dataset Data Model
 *
 * [Dataset Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Policies-and-Datasets}
 */

/**
 * @interface Dataset
 */
export interface Dataset {
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
   *  Agent Group ID {string} UUIDv4 (read only)
   */
  agent_group_id?: string;

  /**
   *  Agent Policy ID {string} UUIDv4 (read only)
   */
  agent_policy_id?: string;

  /**
   *  Array of Sink ID {<string>[]} UUIDv4 (read only)
   */
  sink_id?: string[];

  /**
   * Indicates whether dataset is valid or not {boolean}
   */
  valid?: boolean;

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
   * Dataset Metadata {{[propName: string]: string}}
   */
  dataset_metadata?: any;
}
