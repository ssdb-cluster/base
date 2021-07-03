benchmark

```
go test -v -run TestServer

./redis-benchmark -p 9000 -d 100 -t set -c 100 -n 10000000
```
