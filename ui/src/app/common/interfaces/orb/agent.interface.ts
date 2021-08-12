export interface Agent {
  /** id: UUIDv4 (read only) */
  id?: string;
  /** Name: string [a-zA-Z_:][a-zA-Z0-9_]* */
  name?: string;
  /** Description: string */
  description?: string;
  /**
   * ORB Tags: orb_tags string<JSON>
   * simple key/values - no recursive objects
   */
  tags?: { [propName: string]: string };
  /** Status: string ['active'|'error'] */
  status?: string;
  /** Error Message: string contains error message if status is 'error' (read only) */
  error?: string;
  /** ts_created: UUIDv4 (read only) */
  validate_only?: boolean;
  ts_created?: string;
  matching_agents?: { [propName: string]: Number };
  state?: string;
}
