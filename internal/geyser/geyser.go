package geyser

import (
	"crypto/x509"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"net/url"
	"solanaindexer/internal/logger"
	"time"
)

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

func StartGeyser(grpc string) (*grpc.ClientConn, error) {
	u, err := url.Parse(grpc)
	if err != nil {
		logger.Errorf("Error while parsing grpc %v", err)
		return nil, err
	}

	plaintext := u.Scheme == "http"

	port := u.Port()
	if port == "" {
		port = "443"
	}

	hostname := u.Hostname()
	if hostname == "" {
		return nil, err
	}

	address := hostname + ":" + port
	logger.Infof("%v", address)

	conn, err := grpcConnect(address, plaintext)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func grpcConnect(address string, plaintext bool) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if plaintext {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		pool, _ := x509.SystemCertPool()
		creds := credentials.NewClientTLSFromCert(pool, "")
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	opts = append(opts, grpc.WithKeepaliveParams(kacp))

	logger.Infof("Starting grpc client, connecting to %v", address)
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, errors.New("fail to dial")
	}
	return conn, nil
}
