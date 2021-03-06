package http

import (
	"context"
	"log"

	fj "github.com/kubesure/forkjoin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//DispatchServer implements GRPC server interface FannoutFannin
type DispatchServer struct {
	UnimplementedHTTPForkJoinServiceServer
}

//FanoutFanin Fans out each http message to http dispatch works using the fork join interface
func (s *DispatchServer) FanoutFanin(request *HTTPRequest, stream HTTPForkJoinService_FanoutFaninServer) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mtplx := fj.NewMultiplexer()
	for i, m := range request.Messages {
		value, _ := Message_Method_name[int32(m.Method)]
		msg := fj.HTTPMessage{
			Method:  fj.METHOD(value),
			URL:     m.URL,
			Payload: m.Payload,
			Headers: m.Headers,
			ID:      i,
		}
		reqMsg := fj.HTTPRequest{Message: msg}
		mtplx.AddWorker(&DispatchWorker{Request: reqMsg})
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
				method, _ := Message_Method_value[string(response.Message.Method)]
				m := &Message{
					URL:        response.Message.URL,
					Method:     Message_Method(method),
					Headers:    response.Message.Headers,
					Payload:    response.Message.Payload,
					StatusCode: uint32(response.Message.StatusCode),
				}
				r := HTTPResponse{Message: m}
				err := stream.Send(&r)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func makeErrRes(code fj.ErrorCode, msg string) *HTTPResponse {
	fjerrors := []*Error{}
	fjerr := Error{Code: ErrorCode(code), Message: msg}
	fjerrors = append(fjerrors, &fjerr)
	r := HTTPResponse{Errors: fjerrors}
	return &r
}
