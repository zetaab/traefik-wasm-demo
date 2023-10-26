
Run server:

```
go run main.go logger.go
```

Make request:

```
% curl localhost:8090/hello
User-Agent: curl/8.1.2
Accept: */*
X-Foo: Hello, World!
Hello, World
```

In server logs you should see:

```
{"level":"debug","time":"2023-10-26T20:53:21+03:00","message":"hello from handleRequest debug"}
{"level":"debug","time":"2023-10-26T20:53:21+03:00","message":"time is 2022-01-01 00:00:00.001 +0000 UTC m=+1640995200.001000001"}
```