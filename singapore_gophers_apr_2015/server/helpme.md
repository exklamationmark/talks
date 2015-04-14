## Showing RSS (real memory usage) & SZ (Virtual memory usage)
while true; do ps -p 29920 -o command,rss,sz; sleep 1; done

## Build and show escape analysis
go build -gcflags=-m

## Send lots of request
wrk -t20 -c400 -d30s 'http://localhost:9090'
