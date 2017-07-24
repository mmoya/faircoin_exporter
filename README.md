# FairCoin Exporter

[![Build Status](https://travis-ci.org/mmoya/faircoin_exporter.svg)](https://travis-ci.org/mmoya/faircoin_exporter)

This is a prometheus exporter for monitoring the health of [CVN][1] in the
FairCoin network. A dashboard can be seen [here][2].

## Building

```
git clone https://github.com/mmoya/faircoin_exporter.git
cd faircoin_exporter
sudo apt-get install libzmq3-dev
dep ensure
make
```

`dep` can be installed with `go get -u github.com/golang/dep/cmd/dep`.

## Running

The exporter needs access to both the RPC server (with proper authentication)
and a ZeroMQ socket where hashblock are announced.

You can run everything with docker by following these steps:

1. Create a credential for the RPC server
   ```
   mkdir -p faircoin/blockchain
   cd faircoin
   echo "$(openssl rand -hex 4):$(openssl rand -hex 32)" >faircoin.cookie
   chmod 400 faircoin.cookie
   ```

1. Create a dedicated docker network
   ```
   docker network create faircoin
   ```

1. Launch wallet
   ```
   docker run -d --restart=always \
       --name faircoin \
       --net faircoin \
       -v $PWD/blockchain:/root/.faircoin2 \
       -p 40404:40404 \
       mmoya/faircoin \
           -rpcuser=$(cut -d: -f 1 faircoin.cookie) \
           -rpcpassword=$(cut -d: -f 2 faircoin.cookie) \
           -zmqpubhashblock=tcp://0.0.0.0:28332 \
   ```

1. Launch exporter
   ```
   docker run -d --restart=always \
       --name faircoin-exporter \
       --net faircoin \
       mmoya/faircoin-exporter \
         -rpc.url="http://faircoin:40405" \
         -rpc.user=$(cut -d: -f 1 faircoin.cookie) \
         -rpc.password=$(cut -d: -f 2 faircoin.cookie) \
         -zmq.url="tcp://faircoin:28332"
   ```

the metrics can be queried with:

```
curl -s `docker inspect faircoin-exporter -f '{{ .NetworkSettings.Networks.faircoin.IPAddress }}'`:9132/metrics | grep -E faircoin_
```


[1]: https://github.com/faircoin/faircoin/blob/master/doc/CVN-operators-guide.md
[2]: https://dashboard.faircoin.io/dashboard/db/faircoin-cvn
