# msgpack-rpc-over-http-go

Client library for [msgpack-rpc-over-http](https://github.com/authorNari/msgpack-rpc-over-http) for Go.

## Install

```
$ go get github.com/hakobera/msgpack-rpc-over-http-go
```

## Usage

If server code is followings:

```ruby
require 'msgpack-rpc-over-http'
class MyHandler
  def add(x,y) return x+y end
end

run MessagePack::RPCOverHTTP::Server.app(MyHandler.new)
```

Client code is like this:

```go
package main

import (
  "fmt"
  "os"
  "github.com/hakobera/msgpack-rpc-over-http-go/overhttp"
)

func main() {
  url := "http://localhost:9000"
  opts := make(map[string]int32)
  client := overhttp.NewMsgpackClient(url, &opts)

  ret, err := client.Call("add", 1, 2)
  if err != nil {
    fmt.Printf("got %v", err)
    os.Exit(1)
  }

  fmt.Printf("result is %v", ret)
}
```
