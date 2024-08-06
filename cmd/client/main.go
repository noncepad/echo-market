package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	pbs "github.com/noncepad/echo-market/proto/testecho"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ./client localhost:50061
func main() {
	log.Info("hello world")
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go loopSignal(ctx, cancel, signalC)
	conn, err := grpc.Dial(
		// localhost:50063
		os.Args[1],
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 5; i++ {
		client := pbs.NewTestEchoClient(conn)
		rq := &pbs.EchoRequest{
			Body: "helo",
		}
		resp, err := client.Echo(ctx, rq)
		log.Infof("request: %+v, response %+v", rq, resp)
		if err != nil {
			panic(err)
		}
	}

	conn.Close()

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
