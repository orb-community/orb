/**
 * Oltp Sink Config Interface
 * [Sinks Architecture]{@link https://github.com/orb-community/orb/wiki/Architecture:-Sinks}
 */
import { SinkConfig } from 'app/common/interfaces/orb/sink/sink.config.interface';

/**
 * @interface OtlpConfig
 */
export interface OtlpConfig extends SinkConfig<string> {
    name: 'Otlp';

    authentication: |any| {
        /**
         *  Authentication type, "type": "basicauth" Default, can be ommitted
         */
        type?: string;
        /**
         *  Password {string}
         */
        password?: string;
        /**
        *  Username|Email(?) {string}
        */
        username?: string;
    }
    exporter: |any| {
        /**
        *  Endpoint (Otlp sinks) or Remote Host (Prometheus sink) Link {string}
        */
        endpoint?: string;
        remote_host?: string;
    }
    
}