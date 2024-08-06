package myserver

import (
	"context"
	"errors"
	"fmt"

	pbs "github.com/noncepad/echo-market/proto/testecho"
	"github.com/noncepad/solpipe-market/go/cpumeter"
	"google.golang.org/grpc"
)

type myserverinfo struct {
	pbs.UnimplementedTestEchoServer
	ctx    context.Context
	cancel context.CancelFunc
}
type Base interface {
	Close() error
	CloseSignal() <-chan error
}

func Run(
	ctx context.Context,
	cancel context.CancelFunc,
	s *grpc.Server,
) (Base, error) {
	e1 := &myserverinfo{ctx: ctx, cancel: cancel}

	pbs.RegisterTestEchoServer(s, e1)
	// add a relay service so that the relay knows how much capacity is available
	cpumeter.Add(ctx, s)
	return e1, nil
}

func (e1 *myserverinfo) Close() error {
	signalC := e1.CloseSignal()
	e1.cancel()
	return <-signalC
}

func (e1 *myserverinfo) CloseSignal() <-chan error {
	signalC := make(chan error, 1)
	// TODO; missing internal{}
	go loopCtxWait(e1.ctx, signalC)
	return signalC
}

func loopCtxWait(ctx context.Context, signalC chan<- error) {
	<-ctx.Done()
	signalC <- ctx.Err()
}

// HandleEcho handles an incoming EchoRequest and sends back an EchoResponse.
/*
func (s *EchoServer) EchoRequest(ctx context.Context, req *echopb.EchoRequest) (*echopb.EchoResponse, error) {
	fmt.Printf("Received request: %s\n", req.Message)
	return &echopb.EchoResponse{Message: req.Message}, nil
}

*/

func (e1 *myserverinfo) Echo(ctx context.Context, req *pbs.EchoRequest) (*pbs.EchoResponse, error) {
	fmt.Printf("Received request: %s\n", req.Body)
	return &pbs.EchoResponse{Body: req.Body}, nil
}

func (e1 *myserverinfo) Feed(stream pbs.TestEcho_FeedServer) error {
	return errors.New("not implemented yet")
}
