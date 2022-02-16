/**
 * Agent Policy / Input Config Interface
 *
 * [Policies Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Policies-and-Datasets}
 * [Agent Taps](https://github.com/ns1labs/pktvisor/blob/develop/RFCs/2021-04-16-75-taps.md)
 */

/**
 * @interface PolicyInput
 */
export interface PolicyInput {
  version?: string;

  /**
   * json object with configs
   */
  config?: { [propName: string]: any };

  /**
   * json object with configs
   */
  filter?: { [propName: string]: any };
}

