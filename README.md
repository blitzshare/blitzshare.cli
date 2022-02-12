[![CircleCI](https://circleci.com/gh/blitzshare/blitzshare.cli/tree/main.svg?style=svg&circle-token=2d8209870f559c209f3167d0f28404d05339975e)](https://circleci.com/gh/blitzshare/blitzshare.cli/tree/main)

![logo](./assets/logo.png)

# blitzshare.cli
Blitzshare API client and P2p peer in libp2p network.

[api](https://github.com/blitzshare/blitzshare.api)

[bootstrap node](https://github.com/blitzshare/blitzshare.bootstrap.node)

# Build executable
```bash
make build
```

# Usage

## P2p chat connection
```bash
# start init peer
$ blitz --start
```
Send OTP to connecting peer
```
$ blitz  --connect
```

## P2p file share

```bash
# start init peer
$ blitz --start --file <Local FILE PATH>
```
Send OTP to connecting peer
```
$ blitz  --connect
```
Notice local file created with `blitzshare-<OTP>.txt` name format
 

## Tools & Libraries used
[libp2p](https://docs.ipfs.io/concepts/libp2p/)

[websequencediagrams](./websequencediagrams)