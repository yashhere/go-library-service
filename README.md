# A GRPC service in Go

To compile proto files, use the two commands -

```bash
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:. \
  api/library.proto
```

```bash
protoc -I/usr/local/include -I. \
-I$GOPATH/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--go_out=plugins=grpc:. \
api/library.proto
```

## How to run the service

Run gRPC server in one terminal -

```bash
go run cmd/main.go
```

This will start two server, one gRPC at port :50051 and other HTTP on port :8181.

Now run OPA in other terminal -

```bash
opa run -a "localhost:8080" -s OPA/policy.rego OPA/users.json
```

Now add the books:

```bash
bash add_books.sh
```

After that you can try calling all the functions on following urls -

```bash
http://localhost:8181/listBooks
http://localhost:8181/searchBook
http://localhost:8181/addBook
```

Do not forget to add the input data in the body field of the POST/PUT request.

A sample input data is below -

```json
{
    "user": {
        "userType": 2
    },
    "book": {
      "isbn": "9488900377"
    }
}
```