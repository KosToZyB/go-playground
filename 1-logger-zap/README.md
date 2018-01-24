[zap](https://github.com/uber-go/zap) logger with rest api

/changelevel
---
Change log level:
* info
* warning
* error
* dpanic
* panic
* fatal

```
curl -H 'content-type: application/json' http://localhost:8080/changelevel -d '{"level": "error"}'
```

/hello
---
Method for test log on server
```
curl -H 'content-type: application/json' http://localhost:8080/hello
```
