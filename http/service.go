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

// TODO: Validate Input
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
			// TODO: should it return or process next message
			err := stream.Send(makeGrpcErrRes(result.ID, *result.Err, response))
			if err != nil {
				log.LogResponseError(result.ID, response.Message.ID, fmt.Sprintf("Error while writing to stream: %v", err.Error()))
				return status.Errorf(codes.Internal, fmt.Sprintf("Error while request id: %s message id: %s writing to stream", result.ID, response.Message.ID), err)
			}
		} else {
			m := makeMessage(response)
			r := HTTPResponse{Id: result.ID, Message: m}
			err := stream.Send(&r)
			if err != nil {
				log.LogResponseError(result.ID, response.Message.ID, fmt.Sprintf("Error while writing to stream: %v", err.Error()))
				return status.Errorf(codes.Internal, fmt.Sprintf("Error while request id: %s message id: %s writing to stream", result.ID, response.Message.ID), err)
			}
		}
	}
	log.LogInfo(RequestID(ctx), "Joined")
	return nil
}

func makeMessage(response fj.HTTPResponse) *Message {
	method, _ := Message_Method_value[string(response.Message.Method)]
	m := &Message{
		URL:        response.Message.URL,
		Method:     Message_Method(method),
		Headers:    response.Message.Headers,
		Payload:    response.Message.Payload,
		StatusCode: uint32(response.Message.StatusCode),
		Id:         response.Message.ID,
	}
	return m
}

//func makeGrpcErrRes(code fj.EventCode, msg string) *HTTPResponse {
func makeGrpcErrRes(reqID string, err fj.FJError, response fj.HTTPResponse) *HTTPResponse {
	fjerrors := []*Error{}
	fjerr := Error{Code: ErrorCode(err.Code), Message: err.Message}
	fjerrors = append(fjerrors, &fjerr)
	m := makeMessage(response)
	r := HTTPResponse{Id: reqID, Message: m, Errors: fjerrors}
	return &r
}

func RequestID(ctx context.Context) string {
	return ctx.Value(CtxRequestID).(string)
}
