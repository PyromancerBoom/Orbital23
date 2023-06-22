namespace Go ApiGateway

struct GatewayRequest {
    1: string serviceName (api.path = "serviceName");
    // 2: string serviceMethod (api.path = "serviceMethod");
    2: string serviceId (api.path = "serviceId");
    3: string path (api.path = "path");
}

struct GatewayResponse {
    1: i32 statusCode;
    2: map<string, string> responseData;
}

// Assuming service name is unique
service ApiGateway {
    GatewayResponse processGetRequest(1: GatewayRequest request) (api.get = '/:serviceName/:path');
    GatewayResponse processPostRequest(1: GatewayRequest request) (api.post = '/:serviceName/:path');
}
