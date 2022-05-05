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

// FIXME: Validate Input
//FanoutFanin Fans out each http message to http dispatch works using the fork join interface
func (s *DispatchServer) FanoutFanin(request *HTTPRequest, stream HTTPForkJoinService_FanoutFaninServer) error {
	ctx := context.WithValue(context.Background(), CtxRequestID, request.Id)
	log := fj.NewLogger()
	mtplx := fj.NewMultiplexer()
	for i, m := range request.Messages {
		value, _ := Message_Method_name[int32(m.Method)]
		msg := fj.HTTPMessage{
			Method:         fj.METHOD(value),
			URL:            m.URL,
			Payload:        m.Payload,
			Headers:        m.Headers,
			ActiveDeadLine: m.ActiveDeadLineSeconds,
			ID:             fmt.Sprint(i + 1),
		}
		reqMsg := fj.HTTPRequest{Message: msg}
		mtplx.AddWorker(&DispatchWorker{Request: reqMsg, activeDeadLineSeconds: msg.ActiveDeadLine})
	}

	resultStream := mtplx.Multiplex(ctx, nil)
	log.LogInfo(RequestID(ctx), "Forked")
	for result := range resultStream {
		response, _ := result.X.(fj.HTTPResponse)
		if result.Err != nil {
			err := stream.Send(makeErrRes(result.Err.Code, result.Err.Message))
			if err != nil {
				log.LogResponseError(result.ID, response.Message.ID, fmt.Sprintf("Error while writing to stream: %v", err.Error()))
				return status.Errorf(codes.Internal, fmt.Sprintf("Error while request id: %s message id: %s writing to stream", result.ID, response.Message.ID), err)
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
			r := HTTPResponse{Id: result.ID, Message: m}
			err := stream.Send(&r)
			if err != nil {
				return err
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
