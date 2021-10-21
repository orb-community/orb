/**
 * Agent Policy / Tap Config Interface
 *
 * [Policies Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Policies-and-Datasets}
 * [Agent Taps](https://github.com/ns1labs/pktvisor/blob/develop/RFCs/2021-04-16-75-taps.md)
 */

/**
 * @interface PolicyTap
 */
export interface PolicyTap {
  /**
   * info reserved
   */
  name?: string;

  /**
   * agents
   */
  agents?: {total?: number};

  /**
   * input type
   */
  input_type?: string;

  /**
   * array with configs fields that are predefined,
   * without their predefined value
   */
  config_predefined?: string[];
}

