export interface Sink {
    id?: string;
    name?: string;
    description?: string;
    tags?: any;
    status?: string;
    error?: string;
    backend?: string;
    config?: {
        remote_host: string;
        username: string;
    };
    ts_created?: string;
}
