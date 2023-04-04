// Copyright 2023 NJWS Inc.

package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/foliagecp/inventory-bmc-app/internal/bmc"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbbmc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// set max 100 MB
const maxMsgSize = 100 * 1024 * 1024

var (
	bmcAppAddr = "127.0.0.1"
	bmcAppPort = "32415"

	log = logrus.New()
)

func init() {
	if value, ok := os.LookupEnv("BMC_APP_ADDR"); ok {
		bmcAppAddr = value
	}
	if value, ok := os.LookupEnv("BMC_APP_PORT"); ok {
		bmcAppPort = value
	}
}

func Run(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)

	go func() {
		for {
			select {
			case sig := <-sigChan:
				switch sig {
				case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
					cancel()
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return serve(ctx)
}

func serve(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	port := fmt.Sprintf(":%s", bmcAppPort)
	pc, err := net.Listen("tcp4", port)
	if err != nil {
		return
	}
	defer pc.Close()

	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)
	defer server.Stop()

	log.Info("register discovery")
	bmcApp, err := bmc.New(ctx)
	if err != nil {
		return
	}
	pbbmc.RegisterBmcServiceServer(server, pbbmc.BmcServiceServer(bmcApp))

	log.Info("serving ", port)

	go func() {
		if err = server.Serve(pc); err != nil {
			cancel()
		}

	}()

	<-ctx.Done()

	return
}

func Client() (conn *grpc.ClientConn, err error) {
	return grpc.Dial(fmt.Sprintf("%s:%s", bmcAppAddr, bmcAppPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize), grpc.MaxCallSendMsgSize(maxMsgSize)))
}
