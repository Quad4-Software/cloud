package checks

import (
	"context"
	"fmt"

	"github.com/Quad4-Software/cloud/pkg/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewGRPCCheck checks the reachability of a grpc server.
func NewGRPCCheck(address string) func(context.Context) error {
	return func(_ context.Context) error {
		address, err := handlers.FailSaveAddress(address)
		if err != nil {
			return err
		}

		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("could not connect to grpc server: %v", err)
		}
		_ = conn.Close()
		return nil
	}
}
