# go-postgres-grpc-user-dir

---
This repository implemetn **proto** file: [go-grpc-apis/user/v1](https://github.com/Ekvo/go-grpc-apis/tree/main/user/v1 "https://github.com/Ekvo/go-grpc-apis/tree/main/user/v1")  

The main idea is to implement a service for the user data handler. We can `create`, `read`, `update` and `delete` (`CRUD`) and check authorization if we need to access the user data store.
Also deploy this service in a container using SQL (`postgresql`) as the storage. 


### Structure of applicationgodo

```txt
├── cmd/app
│   └──── main.go  
├── init
│   └──── .env // have .env file - because it's not a commercial service 
├── internal
|   ├── app             // heart of application
|   │   └──── app.go    // run and stop
|   ├── config
|   │   └──── config.go   
|   ├── db
|   │   ├──── migration
|   │   │     └──── migration.go 
|   │   ├──── mock
|   │   │     └──── db_mock.go // for test service
|   │   ├──── db.go   
|   │   └──── query.go  
|   ├── model            // data models define
|   │   ├──── login.go    
|   │   └──── user.go    
|   ├── lib            
|   │   └──── jwtsign     // work with jwt.Token  
|   │         └──── jwtsign.go    
|   ├── listen  
|   │   └──── listen.go   // listen for server
|   └── servises 
|       ├── deserializer  // entities to get data from query or ctx
|       │   ├── deserializer.go      
|       │   ├── login_decode.go     
|       │   ├── token_decode.go      
|       │   ├── user_deocde.go   
|       │   ├── user_id_decode.go     
|       │   └── user_update_decode.go 
|       ├── serializer    // entities - create objects for response
|       │   ├── login_encode.go      
|       │   └── user_encode.go  
|       ├── middleware.go // authorization 
|       ├── service.go    // biz logic
|       ├── user_data.go  
|       ├── user_delete.go 
|       ├── user_login.go       
|       ├── user_register.go 
|       └── user_update.go   
├── pkg/utils 
│   └──── utils.go         // general helper functions
├── script        
│   └──── start.sh
└── sql
    ├──── migrations // contain num_files.up.sql            
    └──── init_compose.sql // use in compose.yaml for create data base  
 .gitignore       
 compose.yaml
 Dockerfile
 README.md
```

### Tech stack: 
- golang 1.24.1, sql, PostgreSQL, /migrate/v4, pgx/v5, net/http, caarlos0/env/v11, testify, jwt, git, Dockerfile, compose.yaml, linux, shell

### Main 'service' from protofile

```protobuf
service UserService {
  rpc UserRegister(UserRegisterRequest) returns (UserRegisterResponse);

  rpc UserLogin(UserLoginRequest) returns (UserLoginResponse);

  // UserData, UserUpdate, UserDelete - get 'user_id' from metadata -H "authorization"
  
  rpc UserData(UserDataRequest) returns (UserDataResponse);
 
  rpc UserUpdate(UserUpdateRequest) returns (UserUpdateResponse);

  rpc UserDelete(UserDeleteRequest) returns (UserDeleteResponse);
}
```


### Start with compose.yaml
```bash
# have .env file 
docker compose --env-file ./init/.env up -d
```

### Local start
To run locally, you need to start Docker, see above, stop the server for service (port 50051:50001)
```bash
# after start docker and close server
go run cmd/app/main.go
```

### grpcurl

For `grpcurl` need load `user.proto`

```bash
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
grpcurl -plaintext -H "authorization: bearer JWT_TOKEN" -d '{"login": "linxy","first_name": "Dmitry", "last_name": "Tai","email": "linxybest@gmail.com", "updated_at": "2024-10-05T16:34:56Z"}' -proto=go-grpc-apis/user/v1/user.proto localhost:50051 user.v1.UserService/UserUpdate
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

* migration tools
```bash
go get github.com/golang-migrate/migrate/v4
```
 
* Use pgx driver for work with postgresql in golang
```bash
go get github.com/jackc/pgx/v5
```

* For parse config
```bash
go get github.com/caarlos0/env/v11
```

* Read data from .env
```bash
go get github.com/joho/godotenv
```

* Use jwt for authorization
```bash
go get github.com/golang-jwt/jwt/v5
```

### Test 

* For comfortable testing is used
```bash
go get github.com/stretchr/testify
```

* Start test from main directory
```bash
go test ./...
```

we can also find out the test **coverage** of specific packages of an application.
```bash
# . - set direct for testing
go test . -coverprofile=coverage.out
```

```bash
# after 'go test ./some_direct/_test.go -coverprofile=coverage.out'
go tool cover -html=coverage
```

##### Сoverage of packages

| file                                      | percent % |
|:------------------------------------------|----------:|
| internal/db/db.go                         |      75.9 |
| internal/db/query.go                      |     100.0 |
| internal/db/schema.go                     |     100.0 |
|                                           |           |
| internal/service/service.go               |     100.0 |
| internal/service/user_data.go             |      60.0 |
| internal/service/user_delete.go           |      75.0 |
| internal/service/user_login.go            |      88.2 |
| internal/service/user_register.go         |      85.7 |
| internal/service/user_update.go           |      74.2 |
| internal/service/middleware.go            |      79.2 |

p.s. Thanks for your time:)

