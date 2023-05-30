# proto definition

We can define protocol data types with `.proto` file.

Follow https://protobuf.dev/programming-guides/proto3/.

## Defining struct

The struct in [user.go](../api/user.go)

```go
type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"fullName"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}
```

defined in grpc format is 

```proto
syntax = "proto3";

package pb;

import  "google/protobuf/timestamp.proto";

option go_package = "github.com/tgfukuda/be-master/pb";

message User {
    string username = 1;
    string full_name = 2;
    string email = 3;
    google.protobuf.Timestamp password_changed_at = 4;
    google.protobuf.Timestamp created_at = 5;
}
```

Each number of parameter defines a serialization of the parameter (Important).

To import local file and make vscode recognize,

```proto
import "user.proto";

option go_package = "github.com/tgfukuda/be-master/pb";

message CreateUserRequest {
    string username = 1;
    string full_name = 2;
    string email = 3;
    string password = 4;
}

message CreateUserResponse {
    User user = 1;
}
```

and add below to settings.json.

```json
"protoc": {
    "options": [
        "--proto_path=proto"
    ]
}
```

### Deprecation of optional and required

`optional` and `required` keywords are deprecated and all of the parameter is optional now.

## Define RPC

Define service to extract rpcs.

```proto
service SimpleBank {
    rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {}
    rpc LoginUser(LoginUserRequest) returns (LoginUserResponse) {}
}
```

## Build

```bash
	$ rm -f pb/*.go # clean up
	$ protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto
```
