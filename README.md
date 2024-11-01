# Infra Load Balancer

A lightweight load balancer supporting HTTP API and gRPC request handling. The load balancer configuration includes two primary configuration files: one for server settings and one for target servers. It provides three working modes: 

1. **Default Mode**: Standard Go HTTP server without concurrency optimizations.
2. **Goroutine Mode**: Adds concurrency for handling requests.
3. **FastHTTP Mode**: Uses the `fasthttp` library for high-performance request handling.

## Configuration Files

1. **Server Configuration**: Defines listening address and timeout settings.
    ```yaml
    ListenAddr: "127.0.0.1:8080"
    ReadTimeout: 30
    ReadHeaderTimeout: 5
    WriteTimeout: 10
    IdleTimeout: 60
    ```

2. **Target Server Configuration**: Specifies target groups, each with its own path, balancing strategy, and list of URLs.
    ```yaml
    - name: "API group"
      handlePath: "/api"
      balanceStrategy: "roundRobin"
      targets:
        urls:
          - "http://localhost:5000/api"
          - "http://localhost:5001/api"
          - "http://localhost:5002/api"
          
    - name: "GRPC group"
      handlePath: "/grpc"
      balanceStrategy: "roundRobin"
      targets:
        urls:
          - "http://localhost:6000/grpc"
          - "http://localhost:6001/grpc"
    ```

## Usage

### Running the Load Balancer

To start the load balancer, specify the configuration files with the balancing strategy:

```bash
go run main.go -choose fast(deafult / defaultG ) -serverConf server.yaml -tragFile trg.yaml
```

##Example Console Output

```
INFRA: 
choose: fast
___
    Name:              API group
    Handle path:       /api
    BalanceStrategy:   roundRobin
    URLs:              [http://localhost:5000/api, http://localhost:5001/api, http://localhost:5002/api]
___
    Name:              GRPC group
    Handle path:       /grpc
    BalanceStrategy:   roundRobin
    URLs:              [http://localhost:6000/grpc, http://localhost:6001/grpc]

created router:       /api with: [http://localhost:5000/api, http://localhost:5001/api, http://localhost:5002/api]
created router:       /grpc with: [http://localhost:6000/grpc, http://localhost:6001/grpc]
server started:       127.0.0.1:8080

```

##Sample Request and Response Log

Below is an example of a request being processed by the load balancer, with both request and response details.

```
--------------------
request:     GET /api HTTP/1.1
User-Agent:  PostmanRuntime/7.42.0
Host:        localhost:5000
Content-Type: application/json
Content-Length: 42
Additional:  header
Accept:      */*
Postman-Token: c51cce4e-95d2-4965-a649-470c96539bb1
Accept-Encoding: gzip, deflate, br
Connection:   keep-alive

{
    "lalala": 1344,
    "some": "data"
}
--------------------
response:    HTTP/1.1 200 OK
Date:        Thu, 31 Oct 2024 23:50:30 GMT
Content-Type: text/plain; charset=utf-8
Content-Length: 12

API response
--------------------
resp 2:  200
```


###Target Server

To run a target server, specify the listening port as shown below. The server will handle incoming requests routed by the load balancer.

```
go run main.go -port 5000

listen: 5000
-------------------
accept request: /api method: GET
body length: 42
___________________
```

Example Request via Postman

A sample request showing a load-balanced API response:


![image](https://github.com/user-attachments/assets/95ba1ffa-bf8f-46aa-a9bd-0c4bc3151457)


