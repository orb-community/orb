/**
 * Prometheus Sink Config Interface
 * for more details:
 * /src/cmd/prom-sink/main.go
 * https://github.com/ns1labs/orb/wiki/Architecture:-Sinks
 */
export interface Prometheus {
  /** Remote Host Name: string */
  remote_host?: string;
  /** Username: string */
  username?: string;
  /** Password: string */
  password?: string;
}
