package operator

import (
	"context"
	"crypto/tls"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/rpc"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/yetanotherco/aligned_layer/core/types"
)

// AggregatorRpcClient is the client to communicate with the aggregator via RPC
type AggregatorRpcClient struct {
	rpcClient            *grpc.ClientConn
	aggregatorIpPortAddr string
	logger               logging.Logger
}

const (
	MaxRetries    = 10
	RetryInterval = 10 * time.Second
)

func NewAggregatorRpcClient(aggregatorIpPortAddr string, logger logging.Logger) (*AggregatorRpcClient, error) {
	client, err := newClient(aggregatorIpPortAddr)
	if err != nil {
		return nil, err
	}

	return &AggregatorRpcClient{
		rpcClient:            client,
		aggregatorIpPortAddr: aggregatorIpPortAddr,
		logger:               logger,
	}, nil
}

func newClient(aggregatorIpPortAddr string) (*grpc.ClientConn, error) {
	// TODO(juli): For mainnet, we should load the aggregator's certificate from a file.
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	return grpc.NewClient(aggregatorIpPortAddr, grpc.WithTransportCredentials(creds))
}

// SendSignedTaskResponseToAggregator is the method called by operators via RPC to send
// their signed task response.
func (c *AggregatorRpcClient) SendSignedTaskResponseToAggregator(signedTaskResponse *types.SignedTaskResponse) {
	var reply uint8
	for retries := 0; retries < MaxRetries; retries++ {
		err := c.rpcClient.Invoke(context.Background(), "Aggregator.ProcessOperatorSignedTaskResponse", signedTaskResponse, &reply)
		if err != nil {
			c.logger.Error("Received error from aggregator", "err", err)
			if errors.Is(err, rpc.ErrShutdown) {
				c.logger.Error("Aggregator is shutdown. Reconnecting...")
				client, err := newClient(c.aggregatorIpPortAddr)
				if err != nil {
					c.logger.Error("Could not reconnect to aggregator", "err", err)
					time.Sleep(RetryInterval)
				} else {
					c.rpcClient = client
					c.logger.Info("Reconnected to aggregator")
				}
			} else {
				c.logger.Infof("Received error from aggregator: %s. Retrying ProcessOperatorSignedTaskResponse RPC call...", err)
				time.Sleep(RetryInterval)
			}
		} else {
			c.logger.Info("Signed task response header accepted by aggregator.", "reply", reply)
			return
		}
	}
}
