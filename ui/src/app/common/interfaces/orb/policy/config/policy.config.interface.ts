/**
 * Agent Policy / Tap Config Interface
 *
 * [Policies Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Policies-and-Datasets}
 * [Agent Taps](https://github.com/ns1labs/pktvisor/blob/develop/RFCs/2021-04-16-75-taps.md)
 */

/**
 * @interface PolicyConfig
 */
export interface PolicyConfig {
  /**
   * info reserved
   */
  kind?: string;

  /**
   * version {string}
   */
  version?: string;

  /**
   * tap name {string}
   */
  tap?: string;

  /**
   * input type
   */
  input_type?: string;

  /**
   * json object with configs
   */
  config?: {[propName: string]: any};

  /**
   * handlers object with configs
   */
  handlers?: {
    modules?: { [propName: string]: any },
    config?: { [propName: string]: any },
  };
}

