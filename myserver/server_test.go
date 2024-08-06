package myserver_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	svr "github.com/noncepad/echo-market/myserver"
	pbs "github.com/noncepad/echo-market/proto/testecho"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestEcho(t *testing.T) {
	// comment about grpc......
	port := uint16(10051)
	// use the context to kill the server (pipeline) and client (bidder)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// set up server
	errorC := make(chan error, 1)
	go loopServer(ctx, cancel, errorC, port)
	t.Cleanup(func() {
		// wait for the server to exit so we do not end up with a hung process
		<-errorC
	})

	time.Sleep(2 * time.Second)

	// set up client
	conn, err := grpc.Dial(
		fmt.Sprintf("127.0.0.1:%d", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	client := pbs.NewTestEchoClient(conn)

	// run test suite'
	suite := []*testCase{
		{
			expectedError: nil,
			input: &pbs.EchoRequest{
				Body: "Hello World!",
			},
			expected: &pbs.EchoResponse{
				Body: "Hello World!",
			},
		},
		{
			expectedError: nil,
			input: &pbs.EchoRequest{
				Body: "",
			},
			expected: &pbs.EchoResponse{
				Body: "",
			},
		},
	}

	for _, s := range suite {
		actualResponse, actualError := client.Echo(
			ctx,
			s.input,
		)

		if s.expectedError != actualError {
			t.Fatalf("errors do not match: got %+v vs expected %+v", actualError, s.expectedError)
		}

		// actual same as expected
		if actualError != nil {
			if s.expected.Body != actualResponse.Body {
				t.Fatalf("got %s, but expected %s", actualResponse.Body, s.expected.Body)
			}
		}
	}

	/*for _, s := range Suite {

	resp, err := client.Echo(ctx, &pbs.DirectionRequest{
		Segment: s.req,S
	})
		}
	}
	   message := "Hello World!"
	   	result := Echo(message)
	   	if result != message {
	   		t.Errorf("Expected '%s', but actual '%s'", message, result)
	   	}
	*/
}

func loopServer(
	ctx context.Context,
	cancel context.CancelFunc,
	errorC chan<- error,
	port uint16,
) {
	defer cancel()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		errorC <- err
		return
	}
	s := grpc.NewServer()
	directionServer, err := svr.Run(ctx, cancel, s)
	if err != nil {
		errorC <- err
		return
	}
	go func() {
		<-directionServer.CloseSignal()
		cancel()
	}()
	go loopClose(ctx, s, l)
	err = s.Serve(l)
	if err != nil {
		errorC <- err
		return
	}
	errorC <- nil
}

func loopClose(ctx context.Context, s *grpc.Server, l net.Listener) {
	<-ctx.Done()
	s.GracefulStop()
	l.Close()
}

type testCase struct {
	input         *pbs.EchoRequest
	expected      *pbs.EchoResponse
	expectedError error
}
