# collector-opentsdb
Qframe collector for OpenTSDB JSON API (v2.x)


## Example

To run the collector, start `main.go`.

```bash
$ go run main.go
  2017/07/02 20:24:00 [II] Dispatch broadcast for Back, Data and Tick
  2017/07/02 20:24:00.299409 [NOTICE]        opentsdb Name:opentsdb   >> Start collector v0.0.0
  [negroni] listening on 0.0.0.0:8070
```

### plain json

By using predefined json blob, one or two metrics are send.

```bash
$ cat resources/metric.json
{
    "metric": "sys.cpu.nice",
    "timestamp": 1346846400,
    "value": 18.0,
    "tags": {
       "host": "web01",
       "dc": "lga"
    }
}
$ curl -H "Content-Type: application/json" -X POST -d @resources/metric.json \
       http://localhost:8070/api/put
```

Output of the collector:

```bash
2017/07/02 20:25:38.042472 sys.cpu.nice 18.000000 1346846400 host=web01,dc=lga
[negroni] 2017-07-02T20:25:38.040237462+02:00 | 204 | 	 4.522366ms | localhost:8070 | POST /api/put
```

Sending two metrics...

```bash
$ cat resources/metrics.json
  [{
      "metric": "sys.cpu.nice",
      "timestamp": 1346846400,
      "value": 18.0,
      "tags": {
         "host": "web01",
         "dc": "lga"
      }
  },{
      "metric": "sys.cpu.nice",
      "timestamp": 1346846400,
      "value": 18.0,
      "tags": {
         "host": "web02",
         "dc": "lga"
      }
  }]
$ curl -H "Content-Type: application/json" -X POST -d @resources/metrics.json \
       http://localhost:8070/api/put
```

Output:

```bash
2017/07/02 20:27:40.883746 sys.cpu.nice 18.000000 1346846400 host=web01,dc=lga
2017/07/02 20:27:40.883765 sys.cpu.nice 18.000000 1346846400 host=web02,dc=lga
[negroni] 2017-07-02T20:27:40.882531713+02:00 | 204 | 	 1.257217ms | localhost:8070 | POST /api/put
```

### GZIP encoded JSON

Using scollector of the bosun project...

```bash
$ ~/bin/scollector -f c_dfstat_darwin -h http://localhost:8070
```

...the JSON is encoded.

```bash
$ nc -l 8070                                                                                                                                                                                  git:(master|●1✚4…
  POST /api/put HTTP/1.1
  Host: localhost:8070
  User-Agent: Scollector/0.6.0-beta1
  Content-Length: 272
  Content-Encoding: gzip
  Content-Type: application/json
  Accept-Encoding: gzip
```

But now worries, the `Content-Encoding` is parsed and handled.

```bash
2017/07/02 20:28:59.288550 darwin.disk.fs.total 487374848.000000 1499020138 mount=/,host=kniebook
2017/07/02 20:28:59.288572 darwin.disk.fs.used 464089196.000000 1499020138 host=kniebook,mount=/
2017/07/02 20:28:59.288582 darwin.disk.fs.free 23029652.000000 1499020138 host=kniebook,mount=/
2017/07/02 20:28:59.288605 darwin.disk.fs.inodes.total 4294967279.000000 1499020138 mount=/,host=kniebook
2017/07/02 20:28:59.288620 darwin.disk.fs.inodes.used 3613419.000000 1499020138 host=kniebook,mount=/
2017/07/02 20:28:59.288631 darwin.disk.fs.inodes.free 4291353860.000000 1499020138 host=kniebook,mount=/
2017/07/02 20:28:59.288641 scollector.collector.duration 0.006347 1499020138 collector=bosun.org/cmd/scollector/collectors.c_dfstat_darwin,host=kniebook,os=darwin
2017/07/02 20:28:59.288668 scollector.collector.error 0.000000 1499020138 collector=bosun.org/cmd/scollector/collectors.c_dfstat_darwin,host=kniebook,os=darwin
[negroni] 2017-07-02T20:28:59.287241996+02:00 | 204 | 	 1.43224ms | localhost:8070 | POST /api/put
```
