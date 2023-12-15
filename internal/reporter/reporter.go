package reporter

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/kangseokgyu/ngbench/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SendResult() {
	conn, err := grpc.Dial("localhost:19895", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNGBenchServiceClient(conn)

	r, err := c.ReportResult(context.Background())
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	for i := 0; i < 5; i++ {
		err := r.Send(&pb.Result{Data: "true"})
		if err != nil {
			fmt.Println(err)
			break
		}
		time.Sleep(time.Second)
	}

	resp, err := r.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response: %v", err)
	}
	log.Printf("Received response: %v", resp.Message)
}
