module github.com/threefoldtech/grid3-go

go 1.18

require (
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/centrifuge/go-substrate-rpc-client/v4 v4.0.5
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/threefoldtech/go-rmb v0.2.0
	github.com/threefoldtech/grid_proxy_server v1.5.6
	github.com/threefoldtech/substrate-client v0.0.0-20220822132933-d0d75781793c
	github.com/threefoldtech/zos v0.5.6-0.20220823125932-7df5043ab018
	golang.zx2c4.com/wireguard/wgctrl v0.0.0-20200609130330-bd2cb7843e1b
)

require github.com/goombaio/namegenerator v0.0.0-20181006234301-989e774b106e

require (
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/tklauser/numcpus v0.3.0 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/vishvananda/netlink v1.1.1-0.20201029203352-d40f9887b852 // indirect
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	github.com/whs/nacl-sealed-box v0.0.0-20180930164530-92b9ba845d8d // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
)

require (
	github.com/ChainSafe/go-schnorrkel v1.0.0 // indirect
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cosmos/go-bip39 v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/decred/base58 v1.0.3 // indirect
	github.com/decred/dcrd/crypto/blake256 v1.0.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/ethereum/go-ethereum v1.10.17 // indirect
	github.com/go-redis/redis/v8 v8.11.1 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang/mock v1.6.0
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gtank/merlin v0.1.1 // indirect
	github.com/gtank/ristretto255 v0.1.2 // indirect
	github.com/jbenet/go-base58 v0.0.0-20150317085156-6237cf65f3a6 // indirect
	github.com/mimoo/StrobeGo v0.0.0-20210601165009-122bf33a46e0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pierrec/xxHash v0.1.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/cors v1.8.2 // indirect
	github.com/rs/zerolog v1.26.0 // indirect
	github.com/vedhavyas/go-subkey v1.0.3 // indirect
	golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898
	golang.org/x/sys v0.0.0-20220403020550-483a9cbc67c0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/centrifuge/go-substrate-rpc-client/v4 v4.0.5 => github.com/threefoldtech/go-substrate-rpc-client/v4 v4.0.3-0.20220629145942-1ef6a654b4b5
