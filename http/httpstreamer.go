package http

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

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
	UnimplementedHTTPForkJoinServiceServer
}

//FanoutFanin Fans out each http message to http dispatch works using the fork join interface
func (s *DispatchServer) FanoutFanin(request *HTTPRequest, stream HTTPForkJoinService_FanoutFaninServer) error {
	stream.Send(&HTTPResponse{})
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
	RegisterHTTPForkJoinServiceServer(s, server)
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
