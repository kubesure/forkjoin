package http

//DispatchServer implements GRPC server interface FannoutFannin
type DispatchServer struct {
}

//FanoutFanin Fans out each http message to http dispatch works using the fork join interface
func (s *DispatchServer) FanoutFanin(request *HTTPRequest, stream HTTPForJoinService_FanoutFaninServer) error {
	stream.Send(&HTTPResponse{})
	return nil
}
