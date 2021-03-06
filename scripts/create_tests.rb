# (c) 2020, AXIA Systems, Inc. All rights reserved.
# See the file LICENSE for licensing terms.

require 'net/http'
require 'json'

Username = "ASDFasdf1245asdfASDFasdf234"
Password = "ASDFasdf1245asdfASDFasdf234"

keystoreURI = "http://127.0.0.1:9650/ext/keystore"
SwapChainURI = "http://127.0.0.1:9650/ext/bc/X"
coreChainURI = "http://127.0.0.1:9650/ext/bc/P"


def send(uri, data)
    uri = URI(uri)

    req = Net::HTTP::Post.new(uri.request_uri, {'Content-Type': 'application/json'})
    req.body = data.to_json

    http = Net::HTTP.new(uri.host, uri.port)
    resp = http.request(req).body
    p resp
    resp
end

def newSwapChainAddr()
    resp = JSON.parse(send(SwapChainURI, {
        "jsonrpc": "2.0",
        "id": nextID(),
        "method": "avm.createAddress",
        "params":{
            "username": Username,
            "password": Password
        }
    }))
    addr2 = resp["result"]["address"]
end

@id = 0
def nextID()
    @id = @id + 1
    @id
end

send(keystoreURI, {
    "jsonrpc": "2.0",
    "id": nextID(),
    "method": "keystore.createUser",
    "params": {
        "username": Username,
        "password": Password
    }
})

send(SwapChainURI, {
    "jsonrpc": "2.0",
    "id": nextID(),
    "method": "avm.importKey",
    "params":{
        "username": Username,
        "password": Password,
        "privateKey": "ewoqjP7PxY4yr3iLTpLisriqt94hdyDFNgchSxGGztUrTXtNN"
    }
})

addr1 = newSwapChainAddr()
addr2 = newSwapChainAddr()
addr3 = newSwapChainAddr()

send(SwapChainURI, {
    "jsonrpc": "2.0",
    "id": nextID(),
    "method": "avm.send",
    "params":{
        "username": Username,
        "password": Password,
        "assetID": "AVA",
        "amount": 100_000,
        "to": addr1
    }
})

sleep 2

send(SwapChainURI, {
    "jsonrpc": "2.0",
    "id": 3,
    "method": "avm.send",
    "params":{
        "username": Username,
        "password": Password,
        "assetID": "AVA",
        "amount": 10_000,
        "to": addr2
    }
})

sleep 2
send(SwapChainURI, {
    "jsonrpc": "2.0",
    "id": nextID(),
    "method": "avm.send",
    "params":{
        "username": Username,
        "password": Password,
        "assetID": "AVA",
        "amount": 10_001,
        "to": addr2
    }
})

sleep 2
send(SwapChainURI, {
    "jsonrpc": "2.0",
    "id": nextID(),
    "method": "avm.send",
    "params":{
        "username": Username,
        "password": Password,
        "assetID": "AVA",
        "amount": 20_002,
        "to": addr3
    }
})

