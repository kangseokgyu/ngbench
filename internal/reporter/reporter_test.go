package reporter_test

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/kangseokgyu/ngbench/internal/reporter"
	pb "github.com/kangseokgyu/ngbench/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.NGBenchServiceServer
}

func (s *server) ReportResult(stream pb.NGBenchService_ReportResultServer) error {
	for {
		data, err := stream.Recv()
		if err != nil {
			log.Printf("error while receiving stream: %v", err)
			break
		}
		log.Printf("Received data: %s", data)
	}
	resp := &pb.ResultReply{
		Message: "Data received successfully!",
	}
	return stream.SendAndClose(resp)
}

func TestSendResult(t *testing.T) {
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 15536))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterNGBenchServiceServer(s, &server{})
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	time.Sleep(time.Second)

	r, err := reporter.NewReporter("localhost", "15536")
	if err != nil {
		log.Fatalf("failed to create reporter: %v", err)
	}
	defer r.Close()
	r.SendDeauthTimestamp()
}
