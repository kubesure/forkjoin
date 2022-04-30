package http

import (
	"context"
	"fmt"

	fj "github.com/kubesure/forkjoin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//Request ID key passsed as data in context
type ctxKey int

const (
	CtxRequestID ctxKey = iota
)

//DispatchServer implements GRPC server interface FannoutFannin
type DispatchServer struct {
	UnimplementedHTTPForkJoinServiceServer
}

//FanoutFanin Fans out each http message to http dispatch works using the fork join interface
func (s *DispatchServer) FanoutFanin(request *HTTPRequest, stream HTTPForkJoinService_FanoutFaninServer) error {
	ctx := context.WithValue(context.Background(), CtxRequestID, request.Id)
	log := fj.NewLogger()
	mtplx := fj.NewMultiplexer()
	for i, m := range request.Messages {
		value, _ := Message_Method_name[int32(m.Method)]
		msg := fj.HTTPMessage{
			Method:  fj.METHOD(value),
			URL:     m.URL,
			Payload: m.Payload,
			Headers: m.Headers,
			ID:      fmt.Sprint(i),
		}
		reqMsg := fj.HTTPRequest{Message: msg}
		mtplx.AddWorker(&DispatchWorker{Request: reqMsg})
	}

	resultStream := mtplx.Multiplex(ctx, nil)
	log.LogInfo(RequestID(ctx), "Forked")
	for result := range resultStream {
		if result.Err != nil {
			err := stream.Send(makeErrRes(result.Err.Code, result.Err.Message))
			if err != nil {
				//TODO: add message id to log if possible
				log.LogResponseError(result.ID, "nil", fmt.Sprintf("Error while writing to stream: %v", err.Error()))
				return status.Errorf(codes.Internal, "Error while writing to stream", err)
			}
		} else {
			response, ok := result.X.(fj.HTTPResponse)
			if !ok {
				log.LogResponseError(result.ID, "nil", "type assertion error http.Response not found")
				log.Println("type assertion err http.Response not found")
				err := stream.Send(makeErrRes(fj.InternalError, "type assertion error http.Response not found"))
				if err != nil {
					//TODO: add message id to log if possible
					log.LogResponseError(result.ID, "nil", fmt.Sprintf("Error while writing to stream: %v", err.Error()))
					return status.Errorf(codes.Internal, "Error while writing to stream", err)
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
	log.LogInfo(RequestID(ctx), "Joined")
	return nil
}

func makeErrRes(code fj.EventCode, msg string) *HTTPResponse {
	fjerrors := []*Error{}
	fjerr := Error{Code: ErrorCode(code), Message: msg}
	fjerrors = append(fjerrors, &fjerr)
	r := HTTPResponse{Errors: fjerrors}
	return &r
}

func RequestID(ctx context.Context) string {
	return ctx.Value(CtxRequestID).(string)
}
