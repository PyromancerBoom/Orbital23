namespace Go registy.proxy

struct ConnectRequest {
    1: string ApiKey;
    2: string ServiceName
    3: string ServerAddress
    4: string ServerPort
    5: i64 TTL
    6: i64 TTD
}

struct ConnectResponse {
    1: string Status;
    2: string Message;
    3: string ServerID;
}

struct HealtRequest {
    1: string ApiKey;
    2: string ServerID;
}

struct HealthResponse {
    1: string Status;
    2: string Message;
}

service RegistryProxy {
    ConnectResponse connectServer(1: ConnectRequest req);
    HealthResponse healthCheckServer(1: HealtRequest req);
}
