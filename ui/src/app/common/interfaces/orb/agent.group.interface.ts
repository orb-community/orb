import { Agent } from 'app/common/interfaces/orb/agent.interface';

export interface AgentGroup {
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
  tags?: any;
  ts_created?: string;
  matching_agents?: {
    total: number;
    online: number;
  };
  agents?: Agent[];
  validate_only?: boolean;
}
