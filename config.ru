require 'msgpack-rpc-over-http'
class MyHandler
  def add(x,y) return x+y end
end

run MessagePack::RPCOverHTTP::Server.app(MyHandler.new)
