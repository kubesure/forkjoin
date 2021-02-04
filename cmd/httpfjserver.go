package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	fj "github.com/kubesure/forkjoin/http"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime)
}

//DispatchServer implements GRPC server interface FannoutFannin
type DispatchServer struct {
	fj.UnimplementedHTTPForkJoinServiceServer
}

//FanoutFanin Fans out each http message to http dispatch works using the fork join interface
func (s *DispatchServer) FanoutFanin(request *fj.HTTPRequest, stream fj.HTTPForkJoinService_FanoutFaninServer) error {
	for i, v := range [5]int{1, 2, 3, 4, 5} {
		time.Sleep(2 * time.Second)
		m := &fj.Message{StatusCode: int32(i), Payload: strconv.Itoa(v)}
		stream.Send(&fj.HTTPResponse{Message: m})
	}

	return nil
}

const (
	port = ":50051"
)

func main() {
	log.Println("starting http fork join service...")
	ctx := context.Background()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	server := &DispatchServer{}
	fj.RegisterHTTPForkJoinServiceServer(s, server)
	reflection.Register(s)

	h := health.NewServer()
	h.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, h)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Println("shutting down party server...")
			s.GracefulStop()
			<-ctx.Done()
		}
	}()
	s.Serve(lis)
}
