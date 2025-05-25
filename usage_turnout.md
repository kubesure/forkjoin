### GRPC server side result stream

GRPC server side streaming interface fans out multiple http requests to forkjoin base library using HTTP DispatchWorker as worker. DispatchWorker returns response or error as result to Result chan. GRPC server streams response to GRPC client as concurrent processing is completed for each HTTP URL requested. 



1. [httpMTLSServer.go](./cmd/mtls/httpMTLSServer.go) [httpTLSServer.go](./cmd/tls/httpTLSServer.go)
2. [certs.md](./certs/certs.md)