package client

import (
	"github.com/529Crew/blade/internal/config"
	"github.com/gagliardetto/solana-go/rpc"
)

var client *rpc.Client

func Get() *rpc.Client {
	return client
}

var utilClient *rpc.Client

func GetUtil() *rpc.Client {
	return utilClient
}

func init() {
	client = rpc.New(config.Get().RpcUrl)
	utilClient = rpc.New(config.Get().UtilRpcUrl)
}
