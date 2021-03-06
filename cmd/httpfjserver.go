package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	fj "github.com/kubesure/forkjoin"
	h "github.com/kubesure/forkjoin/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime)
}

//DispatchServer implements GRPC server interface FannoutFannin
type DispatchServer struct {
	h.UnimplementedHTTPForkJoinServiceServer
}

//FanoutFanin Fans out each http message to http dispatch works using the fork join interface
func (s *DispatchServer) FanoutFanin(request *h.HTTPRequest, stream h.HTTPForkJoinService_FanoutFaninServer) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mtplx := fj.NewMultiplexer()
	for i, m := range request.Messages {
		value, _ := h.Message_Method_name[int32(m.Method)]
		msg := fj.HTTPMessage{
			Method:  fj.METHOD(value),
			URL:     m.URL,
			Payload: m.Payload,
			Headers: m.Headers,
			ID:      i,
		}
		reqMsg := fj.HTTPRequest{Message: msg}
		mtplx.AddWorker(&h.DispatchWorker{Request: reqMsg})
	}

	resultStream := mtplx.Multiplex(ctx, nil)
	for result := range resultStream {
		if result.Err != nil {
			log.Printf("Error for id: %v %v %v\n", result.ID, result.Err.Code, result.Err.Message)
			err := stream.Send(makeErrRes(result.Err.Code, result.Err.Message))
			if err != nil {
				return status.Errorf(codes.Internal, "Error while sending to stream", err)
			}
		} else {
			response, ok := result.X.(fj.HTTPResponse)
			if !ok {
				log.Println("type assertion err http.Response not found")
				err := stream.Send(makeErrRes(fj.InternalError, "type assertion err http.Response not found"))
				if err != nil {
					return status.Errorf(codes.Internal, "Error while sending to stream", err)
				}
			} else {
				method, _ := h.Message_Method_value[string(response.Message.Method)]
				m := &h.Message{
					URL:        response.Message.URL,
					Method:     h.Message_Method(method),
					Headers:    response.Message.Headers,
					Payload:    response.Message.Payload,
					StatusCode: int32(response.Message.StatusCode),
				}
				r := h.HTTPResponse{Message: m}
				err := stream.Send(&r)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func makeErrRes(code fj.ErrorCode, msg string) *h.HTTPResponse {
	fjerrors := []*h.Error{}
	fjerr := h.Error{Code: h.ErrorCode(code), Message: msg}
	fjerrors = append(fjerrors, &fjerr)
	r := h.HTTPResponse{Errors: fjerrors}
	return &r
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
	h.RegisterHTTPForkJoinServiceServer(s, server)
	reflection.Register(s)

	h := health.NewServer()
	h.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, h)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Println("shutting down fork join service...")
			s.GracefulStop()
			<-ctx.Done()
		}
	}()
	s.Serve(lis)
}
