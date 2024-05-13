package geyser

import (
	"context"
	"time"

	"github.com/529Crew/blade/internal/config"
	"github.com/529Crew/blade/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	pb "github.com/529Crew/blade/internal/geyser/proto"
)

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, /* send pings every 10 seconds if there is no activity */
	Timeout:             time.Second,      /* wait 1 second for ping ack before considering the connection dead */
	PermitWithoutStream: true,             /* send pings even without active streams */
}

const retryInterval = 5 * time.Second

var geyserStream pb.Geyser_SubscribeClient
var geyserConnected bool

func Connect(connChans ...chan struct{}) error {
	for {
		logger.Log.Println("[GEYSER]: connecting")
		conn, err := grpc.Dial(
			config.Get().GeyserGrpcUrl,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithKeepaliveParams(kacp),
		)
		if err != nil {
			logger.Log.Println("[GEYSER]: failed to connect, retrying...")
			time.Sleep(retryInterval)
			continue
		}

		client := pb.NewGeyserClient(conn)

		stream, err := client.Subscribe(context.Background())
		if err != nil {
			logger.Log.Println("[GEYSER]: failed to establish stream, re-connecting...")
			time.Sleep(retryInterval)
			continue
		}
		geyserStream = stream

		logger.Log.Println("[GEYSER]: connected")
		geyserConnected = true

		/* alert monitors that new connection was established */
		for _, conn := range connChans {
			go func(c chan struct{}) {
				c <- struct{}{}
			}(conn)
		}

		for {
			resp, err := stream.Recv()
			if err != nil {
				logger.Log.Println("[GEYSER]: failed to recv, re-connecting...")
				time.Sleep(retryInterval)
				continue
			}
			go AlertListeners(resp)
		}
	}
}
