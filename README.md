Roller Coaster
==============

Roller Coaster is a standalone web application for visualization of performance trends.

It has the following features:

* Built-in chart plotter (based on [Google Chart](https://developers.google.com/chart/))
* Embedded data storage ([Bolt](https://github.com/boltdb/bolt))
* Simple [RESTful API](https://github.com/gin-gonic/gin) for data manipulations
* No external dependencies

API
===

Currently, the application supports these endpoints: 

| Endpoint                                | Method | Payload   | Description                                        |
|-----------------------------------------|--------|-----------|----------------------------------------------------|
| http://127.0.0.1:8080/api/v1/benchmarks | GET    | N/A       | Getting a list of all "benchmark" objects          |
| http://127.0.0.1:8080/api/v1/benchmarks | POST   | benchmark | Adding a new "benchmark" object to the data bucket |

The following status codes are used in API:

| Methods   | Code | Description     |
|-----------|------|-----------------|
| GET       | 200  | Success         |
| POST      | 201  | Benchmark added |
| POST      | 400  | Bad payload     |
| GET, POST | 500  | Internal error  |

"benchmark" objects can be described using this schema:

```
{
  "type": "object",
  "properties": {
    "group": {
      "type": "string"
    },
    "metric": {
      "type": "string"
    },
    "value": {
      "type": "number"
    }
  },
  "required": [
    "group",
    "metric",
    "value"
  ]
}
```

Please notice that Unix timestamps are automatically added to the documents upon successful POST request.

Examples:

```
> curl -XPOST -d '{"group":"ForestDB, Write-heavy workload","metric":"Read throughput, ops/sec","value":25000}' http://127.0.0.1:8080/api/v1/benchmarks
{"message":"ok"}
```

```
> curl -XGET http://127.0.0.1:8080/api/v1/benchmarks
[{"group":"ForestDB, Write-heavy workload","metric":"Read throughput, ops/sec","timestamp":1470268675328944907,"value":25000}]

```
