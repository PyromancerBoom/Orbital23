namespace Go ApiGateway

struct GatewayRequest {
    1: string serviceName (api.path = "serviceName");
    2: string serviceMethod (api.path = "serviceMethod");
    3: string requestData (api.form = "*");
    4: string serviceId (api.path = "serviceId");
}

struct GatewayResponse {
    1: i32 statusCode;
    2: map<string, string> responseData;
}

// Assuming service name is unique
service ApiGateway {
    GatewayResponse processGetRequest(1: GatewayRequest request) (api.get = '/:serviceName/:serviceMethod');
    GatewayResponse processPostRequest(1: GatewayRequest request) (api.post = '/:serviceName/:serviceMethod');
}
