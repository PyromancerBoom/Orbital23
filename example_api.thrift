namespace go API

struct EchoReq{
    1: string Msg (api.form = "msg");
}

struct EchoResp{
    1: string Msg;
}

service EchoService{
    EchoResp echo(1: EchoReq REQ) (api.post = "/echo");
}