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
   * tap name {string}
   */
  tap?: string;

  /**
   * agents
   */
  agents?: {total?: number};

  /**
   * input type
   */
  input_type?: string;

  /**
   * json object with configs
   */
  config?: { [propName: string]: any };

  /**
   * json object with configs
   */
  filter?: { [propName: string]: any };

  /**
   * array with configs fields that are predefined,
   * without their predefined value
   */
  config_predefined?: { [propName: string]: string | number | boolean | any };

  /**
   * array with filter fields that are predefined,
   * without their predefined value
   */
  filter_predefined?: { [propName: string]: string | number | boolean | any };
}

