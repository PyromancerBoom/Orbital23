namespace go UserService

struct QueryUser {
  1: string ID;
}

struct QueryUserResponse {
  1: bool Exist;
  2: string ID;
  3: string Name;
  4: string Email;
  5: i32 Age;
}

struct InsertUser {
  1: string ID;
  2: string Name;
  3: string Email;
  4: i32 Age;
}

struct InsertUserResponse {
  1: bool Ok;
  2: string Msg;
}

service UserService {
  QueryUserResponse queryUser(1: QueryUser req);
  InsertUserResponse insertUser(1: InsertUser req);
}