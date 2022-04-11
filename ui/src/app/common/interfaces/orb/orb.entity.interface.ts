/**
 * Base Orb Entity Data Model Interface
 *
 * [Orb Base Entity]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Common-Patterns}
 */

/**
 * @interface OrbEntity
 */
export interface OrbEntity {
  /**
   *  ID {string} UUIDv4 (read only)
   */
  id?: string;

  /**
   * Name {string} [a-zA-Z_:][a-zA-Z0-9_]*
   */
  name?: string;

  /**
   *  Tenant owner {string} UUIDv4 (read only)
   */
  mf_owner_id?: string;

  /**
   * Error - dict of errors
   */
  error?: any;

  /**
   * Error Message
   */
  message?: string;

  /**
   * Error Status
   */
  status?: string;

  /**
   * Error Status Message
   */
  statusText?: string;

}
