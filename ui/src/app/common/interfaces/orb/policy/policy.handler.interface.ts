/**
 * Agent Policy / Handler Module Interface
 *
 * [Policies Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Policies-and-Datasets}
 */

/**
 * @interface PolicyHandler
 */
export interface PolicyHandler {
  /**
   * name {string}
   */
  name?: string;

  /**
   * version {string}
   */
  version?: string;

  /**
   * type {string}
   */
  type?: string;

  /**
   * config {}
   */
  config?: { [propName: string]: {} | any };

  /**
   * filter {}
   */
  filter?: { [propName: string]: {} | any };

  /**
   * metrics {}
   */
  metrics?: { [propName: string]: {} | any };

  /**
   * metrics_groups {}
   */
  metrics_groups?: { [propName: string]: {} | any };

  /**
   * content
   */
  content?: { [propName: string]: {} | any };
}

