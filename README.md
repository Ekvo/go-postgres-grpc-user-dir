# go-postgres-grpc-user-dir

---

This repository implemetn **proto** file: [go-postgres-grpc-apis/user/v1](https://github.com/Ekvo/go-postgres-grpc-apis/tree/main/user/v1 "https://github.com/Ekvo/go-postgres-grpc-apis/tree/main/user/v1")  

```protobuf
service UserService {
  rpc SignUp(SignUpRequest) returns (SignUpResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc UserData(UserDataRequest) returns (UserDataResponse);
}
```
* Create user with help - `SignUp`    
* Login with email address and password - `Login`  
* Get all user data without password - `UserData`

---

### Start with compose.yaml
```bash
docker compose --env-file ./init/.env up -d
```

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

