package client

import (
	"context"

	"github.com/529Crew/blade/internal/config"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	jito_go "github.com/weeaa/jito-go"
	"github.com/weeaa/jito-go/clients/searcher_client"
)

var jitoClient *searcher_client.Client

func GetJito() *searcher_client.Client {
	return jitoClient
}

func init() {
	jito_wallet, err := solana.PrivateKeyFromBase58(config.Get().JitoPrivateKey)
	if err != nil {
		panic(err)
	}

	client, err := searcher_client.New(
		context.Background(),
		jito_go.NewYork.BlockEngineURL,
		rpc.New(config.Get().JitoRpcUrl),
		rpc.New(config.Get().RpcUrl),
		jito_wallet,
		nil,
	)
	if err != nil {
		panic(err)
	}
	jitoClient = client
}
