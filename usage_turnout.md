### GRPC server side result stream

GRPC server side streaming interface fans out multiple http requests to forkjoin base library using HTTP DispatchWorker as worker. DispatchWorker returns response or error as result to Result chan. GRPC server streams response to GRPC client as concurrent processing is completed for each HTTP URL requested. 

set ENV var TLS and MLTS server

1. SERVER_CRT path to server certificate. [httpMTLSServer.go](./cmd/mtls/httpMTLSServer.go) [httpTLSServer.go](./cmd/tls/httpTLSServer.go)
2. SERVER_KEY path to server key. [httpMTLSServer.go](./cmd/mtls/httpMTLSServer.go) [httpTLSServer.go](./cmd/tls/httpTLSServer.go) 
4. CA_CRT path ca certificate. required by MTLS connection refer [httpMTLSServer.go](./cmd/mtls/httpMTLSServer.go)
5. SERVER_HOSTNAME DNS name specified in certificate file. Required by clients. refer [certs.md](./certs/certs.md)