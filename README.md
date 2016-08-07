Roller Coaster
==============

[![Go Report Card](https://goreportcard.com/badge/github.com/pavel-paulau/rollercoaster)](https://goreportcard.com/report/github.com/pavel-paulau/rollercoaster)
[![Travis CI](https://travis-ci.org/pavel-paulau/rollercoaster.svg?branch=master)](https://travis-ci.org/pavel-paulau/rollercoaster)
[![Coverage Status](https://coveralls.io/repos/github/pavel-paulau/rollercoaster/badge.svg?branch=master)](https://coveralls.io/github/pavel-paulau/rollercoaster?branch=master)

Roller Coaster is a standalone web application for visualization of performance trends.

It provides the following features:

* Built-in chart plotter (based on [Google Chart](https://developers.google.com/chart/))
* Embedded data storage ([Bolt](https://github.com/boltdb/bolt))
* Simple [RESTful API](https://github.com/gin-gonic/gin) for data manipulation
* No external dependencies

API
===

Currently, the application supports these endpoints: 

| Endpoint                                | Method | Payload   | Description                                      |
|-----------------------------------------|--------|-----------|--------------------------------------------------|
| http://127.0.0.1:8080/api/v1/benchmarks | GET    | N/A       | Gets a list of all "benchmark" objects           |
| http://127.0.0.1:8080/api/v1/benchmarks | POST   | benchmark | Adds a new "benchmark" object to the data bucket |
| http://127.0.0.1:8080/api/v1/benchmarks | DELETE | id        | Deletes an existing "benchmark" object by id     |

The following status codes are used in API:

| Methods           | Code | Description     |
|-------------------|------|-----------------|
| GET, DELETE       | 200  | Success         |
| POST              | 201  | Benchmark added |
| DELETE, POST      | 400  | Bad payload     |
| DELETE, GET, POST | 500  | Internal error  |

"benchmark" object can be described using this schema:

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

"id" object can be described using this schema:

```
{
  "type": "object",
  "properties": {
    "id": {
      "type": "integer"
    }
  },
  "required": [
    "id""
  ]
}
```

Please notice that Unix timestamps and incremental IDs are automatically added to the documents upon successful POST request.

Examples:

```
> curl -XPOST -d '{"group":"ForestDB, Write-heavy workload","metric":"Read throughput, ops/sec","value":25000}' http://127.0.0.1:8080/api/v1/benchmarks
{"message":"ok"}
```

```
> curl -XGET http://127.0.0.1:8080/api/v1/benchmarks
[{"group":"ForestDB, Write-heavy workload","id":1,"metric":"Read throughput, ops/sec","timestamp":1470424035931557290,"value":25000}]
```

```
> curl -XDELETE -d '{"id":1}' http://127.0.0.1:8080/api/v1/benchmarks
{"message":"ok"}
```

Docker image
============

A small Docker image (7.3MB) is available for this project:

```
> docker pull perflab/rollercoaster

> docker run -t -d -p 8080:8080 perflab/rollercoaster
```
