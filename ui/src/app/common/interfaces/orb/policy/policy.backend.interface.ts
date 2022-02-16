/**
 * Agent Policy / Backend Interface
 *
 * [Policies Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Policies-and-Datasets}
 */

/**
 * @interface PolicyBackend
 */
export interface PolicyBackend {
  /**
   * backend denomination {string}
   */
  backend?: string;

  /**
   * description {string}
   */
  description?: string;

  /**
   * schema version {string}
   */
  schema_version?: string;
}

