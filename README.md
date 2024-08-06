# Echo Example

This repository gives an example grpc server and client that echo messages.

## Compile

Compile the go libraries:

```bash
go build -o ./server github.com/noncepad/echo-market/cmd/server
go build -o ./client github.com/noncepad/echo-market/cmd/client
```

## Run

### Run Server

Select an available TCP port that is not currently in use:

```bash
./server localhost:40064
```

### Run Client

Open a new terminal and execute the following:

```bash
./client localhost:40064
```

* Ensure that the TCP port numbers match with the server

## Output

The response should look something like:

### Output Client

```log
INFO[0000] hello world                                  
INFO[0000] request: body:"helo", response body:"helo"   
INFO[0000] request: body:"helo", response body:"helo"   
INFO[0000] request: body:"helo", response body:"helo"   
INFO[0000] request: body:"helo", response body:"helo"   
INFO[0000] request: body:"helo", response body:"helo"  
```

### Output Server

```log
Received request: helo
Received request: helo
Received request: helo
Received request: helo
Received request: helo
```
