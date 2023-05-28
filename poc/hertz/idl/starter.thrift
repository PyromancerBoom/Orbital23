namespace go gateway

struct HelloReq {
     1: string Name
}

struct HelloResp {
     1: string RespBody;
}

service HelloService {
     HelloResp HelloMethod(1: HelloReq request);
     HelloResp HelloMethod1(1: HelloReq request);
     HelloResp HelloMethod2(1: HelloReq request);
}
