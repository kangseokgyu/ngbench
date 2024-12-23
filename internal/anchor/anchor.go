package anchor

import (
	"fmt"
	"log"
	"net"
	"time"

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

func (s *server) ReportDeauthTimestampResult(stream pb.NGBenchService_ReportDeauthTimestampResultServer) error {
	for {
		data, err := stream.Recv()
		if err != nil {
			log.Printf("error while receiving stream: %v", err)
			break
		}
		log.Printf("Received data: %s", time.Unix(0, int64(data.Timestamp)).String())
	}
	resp := &pb.DeauthTimestampResultReply{}
	return stream.SendAndClose(resp)
}

func RecvResult(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNGBenchServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
