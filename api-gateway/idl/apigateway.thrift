namespace Go ApiGateway

struct GatewayRequest {
    1: string serviceName (api.form = "servicename");
    2: string serviceMethod (api.form = "servicenethod");
    3: string requestData (api.form = "requestdata");
    4: string serviceId (api.form = "serviceid");
}

struct GatewayResponse {
    1: i32 statusCode;
    2: string responseData;
}

service ApiGateway {
    GatewayResponse processGetRequest(1: GatewayRequest REQ) (api.post = "hertzgateway/get");
    GatewayResponse processPostRequest(1: GatewayRequest REQ) (api.post = "hertzgateway/post");
}
