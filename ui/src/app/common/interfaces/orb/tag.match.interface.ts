/**
 * Tag Match Model Interface
 *
 * [Fleet Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Fleet}
 */

/**
 * @interface TagMatch
 * Define some common information fields for validation
 * response when matching tags. Add as needed
 */
export interface TagMatch {
  /**
   * Total #matches
   */
  total?: number;

  /**
   * Total #online of #matches
   */
  online?: number;
}
