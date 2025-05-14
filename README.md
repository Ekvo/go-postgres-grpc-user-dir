# go-postgres-grpc-user-dir

---

This repository implemetn **proto** file: [go-grpc-apis/user/v1](https://github.com/Ekvo/go-grpc-apis/tree/main/user/v1 "https://github.com/Ekvo/go-grpc-apis/tree/main/user/v1")  

```protobuf
service UserService {
  rpc UserRegister(UserRegisterRequest) returns (UserRegisterResponse);

  rpc UserLogin(UserLoginRequest) returns (UserLoginResponse);

  // UserData, UserUpdate, UserDelete - get 'user_id' from metadata -H "authorization"

  // by 'user_id' find 'User' in 'db'
  rpc UserData(UserDataRequest) returns (UserDataResponse);

  // by 'user_id' find 'User' in 'db'
  // if found -> set empty 'User' fields with old User data from 'db' -> Update
  rpc UserUpdate(UserUpdateRequest) returns (UserUpdateResponse);

  rpc UserDelete(UserDeleteRequest) returns (UserDeleteResponse);
}
```

### Start with compose.yaml
```bash
docker compose --env-file ./init/.env up -d
```

### grpcurl

for grpcurl need load user.proto

```protobuf
git clone https://github.com/Ekvo/go-grpc-apis.git
```

**Use grpcurl directly from the directory where you clone `https://github.com/Ekvo/go-grpc-apis.git`, not from go-grpc-apis itself**
this rules for work with Docker container


* Create user with help - `UserRegister`
```http request
grpcurl -plaintext -d '{"login":"ekvo", "first_name": "Alex", "email": "alex@example.com", "password": "somepass","created_at": "2024-10-05T15:34:56Z" }' -proto=go-grpc-apis/user/v1/user.proto localhost:50051 user.v1.UserService/UserRegister
```
* Login with email address and password - `UserLogin` 
```http request
grpcurl -plaintext -d '{ "email": "alex@example.com", "password": "somepass" }' -proto=go-grpc-apis/user/v1/user.proto localhost:50051 user.v1.UserService/UserLogin
```
* Get all user data without password - `UserData`

**next grpcurl - change `JWT_TOKEN` to token from response `UserLoginResponse` after `UserLogin`**
```http request
grpcurl -plaintext -H "authorization: bearer JWT_TOKEN" -proto=go-grpc-apis/user/v1/user.proto localhost:50051 user.v1.UserService/UserData
```
* Update user data - `UserUpdate`
```http request
grpcurl -plaintext -H "authorization: bearer JWT_TOKEN" -d '{ "last_name": "Ekvo", "updated_at": "2024-10-05T16:34:56Z"}' -proto=go-grpc-apis/user/v1/user.proto localhost:50051 user.v1.UserService/UserUpdate
```
* Remove user  - `UserDelete`
```http request
grpcurl -plaintext -H "authorization: bearer JWT_TOKEN" -proto=go-grpc-apis/user/v1/user.proto localhost:50051 user.v1.UserService/UserDelete
```

---



### Basic principles:
 * DTO
 * Solid
 
### Stuff 
 
* postgresql use pgx
```bash
go get github.com/jackc/pgx/v5
```

* for parse config
```bash
go get github.com/spf13/viper
```

* read data from .env
```bash
go get github.com/joho/godotenv
```

* use jwt
```bash
go get github.com/golang-jwt/jwt/v5
```
