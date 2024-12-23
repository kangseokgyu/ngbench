package reporter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	pb "github.com/kangseokgyu/ngbench/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Reporter struct {
	conn   *grpc.ClientConn
	client pb.NGBenchServiceClient
}

func NewReporter(server_ip string, port string) (*Reporter, error) {
	conn, err := grpc.Dial(
		server_ip+":"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, err
	}
	return &Reporter{
		conn:   conn,
		client: pb.NewNGBenchServiceClient(conn),
	}, nil
}

func (r *Reporter) Close() {
	r.conn.Close()
}

func SendResult() {
	conn, err := grpc.Dial(
		"localhost:15536",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
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

func (r *Reporter) SendDeauthTimestamp() {
	stream, err := r.client.ReportDeauthTimestampResult(context.Background())
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	pd, err := pcap.OpenLive("ath2", 65535, true, time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}
	defer pd.Close()

	err = pd.SetBPFFilter("wlan type mgt subtype deauth and wlan addr2 00:00:00:66:00:00")
	if err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(pd, pd.LinkType())

	for packet := range packetSource.Packets() {
		if packet == nil {
			log.Println("End of Test")
			break
		}

		log.Printf("Capture a deauth frame. (timestamp: %s)\n", packet.Metadata().Timestamp.String())
		result := &pb.DeauthTimestampResult{Timestamp: uint64(packet.Metadata().Timestamp.UnixNano())}
		stream.Send(result)
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response: %v", err)
	}
}
