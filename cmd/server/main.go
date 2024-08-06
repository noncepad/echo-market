package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	svr "github.com/noncepad/echo-market/myserver"
	"google.golang.org/grpc"
)

// ./server localhost:50063
func main() {
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go loopSignal(ctx, cancel, signalC)

	l, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	directionServer, err := svr.Run(ctx, cancel, s)
	if err != nil {
		panic(err)
	}
	go loopWait(directionServer.CloseSignal(), cancel)
	go loopClose(ctx, s, l)
	err = s.Serve(l)
	if err != nil {
		panic(err)
	}
}

func loopWait(errorC <-chan error, cancel context.CancelFunc) {
	<-errorC
	cancel()
}

func loopClose(ctx context.Context, s *grpc.Server, l net.Listener) {
	<-ctx.Done()
	log.Print("exiting echo server")
	//	s.GracefulStop()
	s.Stop()
	// log.Print("gracefully stopped echo server")
	l.Close()
	log.Print("closing listener")
}

func loopSignal(ctx context.Context, cancel context.CancelFunc, signalC <-chan os.Signal) {
	defer cancel()
	doneC := ctx.Done()
	select {
	case <-doneC:
	case s := <-signalC:
		os.Stderr.WriteString(fmt.Sprintf("%s\n", s.String()))
	}
}
