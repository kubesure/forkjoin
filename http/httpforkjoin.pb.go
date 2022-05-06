// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: api/httpforkjoin.proto

package http

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ErrorCode int32

const (
	ErrorCode_InternalError           ErrorCode = 0
	ErrorCode_RequestError            ErrorCode = 1
	ErrorCode_ResponseError           ErrorCode = 2
	ErrorCode_ConnectionError         ErrorCode = 3
	ErrorCode_ConcurrencyContextError ErrorCode = 4
	ErrorCode_RequestAborted          ErrorCode = 5
)

// Enum value maps for ErrorCode.
var (
	ErrorCode_name = map[int32]string{
		0: "InternalError",
		1: "RequestError",
		2: "ResponseError",
		3: "ConnectionError",
		4: "ConcurrencyContextError",
		5: "RequestAborted",
	}
	ErrorCode_value = map[string]int32{
		"InternalError":           0,
		"RequestError":            1,
		"ResponseError":           2,
		"ConnectionError":         3,
		"ConcurrencyContextError": 4,
		"RequestAborted":          5,
	}
)

func (x ErrorCode) Enum() *ErrorCode {
	p := new(ErrorCode)
	*p = x
	return p
}

func (x ErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_api_httpforkjoin_proto_enumTypes[0].Descriptor()
}

func (ErrorCode) Type() protoreflect.EnumType {
	return &file_api_httpforkjoin_proto_enumTypes[0]
}

func (x ErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorCode.Descriptor instead.
func (ErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_api_httpforkjoin_proto_rawDescGZIP(), []int{0}
}

type Message_Method int32

const (
	Message_NIL   Message_Method = 0
	Message_GET   Message_Method = 1
	Message_POST  Message_Method = 2
	Message_PUT   Message_Method = 3
	Message_PATCH Message_Method = 4
)

// Enum value maps for Message_Method.
var (
	Message_Method_name = map[int32]string{
		0: "NIL",
		1: "GET",
		2: "POST",
		3: "PUT",
		4: "PATCH",
	}
	Message_Method_value = map[string]int32{
		"NIL":   0,
		"GET":   1,
		"POST":  2,
		"PUT":   3,
		"PATCH": 4,
	}
)

func (x Message_Method) Enum() *Message_Method {
	p := new(Message_Method)
	*p = x
	return p
}

func (x Message_Method) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Message_Method) Descriptor() protoreflect.EnumDescriptor {
	return file_api_httpforkjoin_proto_enumTypes[1].Descriptor()
}

func (Message_Method) Type() protoreflect.EnumType {
	return &file_api_httpforkjoin_proto_enumTypes[1]
}

func (x Message_Method) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Message_Method.Descriptor instead.
func (Message_Method) EnumDescriptor() ([]byte, []int) {
	return file_api_httpforkjoin_proto_rawDescGZIP(), []int{0, 0}
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	URL                   string            `protobuf:"bytes,1,opt,name=URL,proto3" json:"URL,omitempty"`
	Method                Message_Method    `protobuf:"varint,3,opt,name=method,proto3,enum=http.Message_Method" json:"method,omitempty"`
	Headers               map[string]string `protobuf:"bytes,4,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Payload               string            `protobuf:"bytes,5,opt,name=payload,proto3" json:"payload,omitempty"`
	StatusCode            uint32            `protobuf:"varint,6,opt,name=statusCode,proto3" json:"statusCode,omitempty"`
	Id                    string            `protobuf:"bytes,7,opt,name=Id,proto3" json:"Id,omitempty"`
	ActiveDeadLineSeconds uint32            `protobuf:"varint,8,opt,name=activeDeadLineSeconds,proto3" json:"activeDeadLineSeconds,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_httpforkjoin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_api_httpforkjoin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_api_httpforkjoin_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetURL() string {
	if x != nil {
		return x.URL
	}
	return ""
}

func (x *Message) GetMethod() Message_Method {
	if x != nil {
		return x.Method
	}
	return Message_NIL
}

func (x *Message) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *Message) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

func (x *Message) GetStatusCode() uint32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *Message) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Message) GetActiveDeadLineSeconds() uint32 {
	if x != nil {
		return x.ActiveDeadLineSeconds
	}
	return 0
}

type HTTPRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Api      string     `protobuf:"bytes,1,opt,name=api,proto3" json:"api,omitempty"`
	Messages []*Message `protobuf:"bytes,2,rep,name=messages,proto3" json:"messages,omitempty"`
	Id       string     `protobuf:"bytes,3,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *HTTPRequest) Reset() {
	*x = HTTPRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_httpforkjoin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HTTPRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HTTPRequest) ProtoMessage() {}

func (x *HTTPRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_httpforkjoin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HTTPRequest.ProtoReflect.Descriptor instead.
func (*HTTPRequest) Descriptor() ([]byte, []int) {
	return file_api_httpforkjoin_proto_rawDescGZIP(), []int{1}
}

func (x *HTTPRequest) GetApi() string {
	if x != nil {
		return x.Api
	}
	return ""
}

func (x *HTTPRequest) GetMessages() []*Message {
	if x != nil {
		return x.Messages
	}
	return nil
}

func (x *HTTPRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type HTTPResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message *Message `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Errors  []*Error `protobuf:"bytes,2,rep,name=errors,proto3" json:"errors,omitempty"`
	Id      string   `protobuf:"bytes,3,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *HTTPResponse) Reset() {
	*x = HTTPResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_httpforkjoin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HTTPResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HTTPResponse) ProtoMessage() {}

func (x *HTTPResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_httpforkjoin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HTTPResponse.ProtoReflect.Descriptor instead.
func (*HTTPResponse) Descriptor() ([]byte, []int) {
	return file_api_httpforkjoin_proto_rawDescGZIP(), []int{2}
}

func (x *HTTPResponse) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *HTTPResponse) GetErrors() []*Error {
	if x != nil {
		return x.Errors
	}
	return nil
}

func (x *HTTPResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    ErrorCode `protobuf:"varint,1,opt,name=code,proto3,enum=http.ErrorCode" json:"code,omitempty"`
	Message string    `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_httpforkjoin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_api_httpforkjoin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_api_httpforkjoin_proto_rawDescGZIP(), []int{3}
}

func (x *Error) GetCode() ErrorCode {
	if x != nil {
		return x.Code
	}
	return ErrorCode_InternalError
}

func (x *Error) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_api_httpforkjoin_proto protoreflect.FileDescriptor

var file_api_httpforkjoin_proto_rawDesc = []byte{
	0x0a, 0x16, 0x61, 0x70, 0x69, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x66, 0x6f, 0x72, 0x6b, 0x6a, 0x6f,
	0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x68, 0x74, 0x74, 0x70, 0x22, 0xf5,
	0x02, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x52,
	0x4c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x55, 0x52, 0x4c, 0x12, 0x2c, 0x0a, 0x06,
	0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x68,
	0x74, 0x74, 0x70, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x34, 0x0a, 0x07, 0x68, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x68, 0x74,
	0x74, 0x70, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73,
	0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x15, 0x61, 0x63,
	0x74, 0x69, 0x76, 0x65, 0x44, 0x65, 0x61, 0x64, 0x4c, 0x69, 0x6e, 0x65, 0x53, 0x65, 0x63, 0x6f,
	0x6e, 0x64, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x15, 0x61, 0x63, 0x74, 0x69, 0x76,
	0x65, 0x44, 0x65, 0x61, 0x64, 0x4c, 0x69, 0x6e, 0x65, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73,
	0x1a, 0x3a, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x38, 0x0a, 0x06,
	0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x07, 0x0a, 0x03, 0x4e, 0x49, 0x4c, 0x10, 0x00, 0x12,
	0x07, 0x0a, 0x03, 0x47, 0x45, 0x54, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x4f, 0x53, 0x54,
	0x10, 0x02, 0x12, 0x07, 0x0a, 0x03, 0x50, 0x55, 0x54, 0x10, 0x03, 0x12, 0x09, 0x0a, 0x05, 0x50,
	0x41, 0x54, 0x43, 0x48, 0x10, 0x04, 0x22, 0x5a, 0x0a, 0x0b, 0x48, 0x54, 0x54, 0x50, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x70, 0x69, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x61, 0x70, 0x69, 0x12, 0x29, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x68, 0x74, 0x74, 0x70,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x49, 0x64, 0x22, 0x6c, 0x0a, 0x0c, 0x48, 0x54, 0x54, 0x50, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x27, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x68, 0x74, 0x74, 0x70, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x23, 0x0a, 0x06, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x68, 0x74,
	0x74, 0x70, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x64,
	0x22, 0x46, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x23, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x68, 0x74, 0x74, 0x70, 0x2e, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2a, 0x89, 0x01, 0x0a, 0x09, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x11, 0x0a, 0x0d, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x01, 0x12, 0x11, 0x0a, 0x0d, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x02, 0x12, 0x13,
	0x0a, 0x0f, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x10, 0x03, 0x12, 0x1b, 0x0a, 0x17, 0x43, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x63, 0x79, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x04,
	0x12, 0x12, 0x0a, 0x0e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x62, 0x6f, 0x72, 0x74,
	0x65, 0x64, 0x10, 0x05, 0x32, 0x4d, 0x0a, 0x13, 0x48, 0x54, 0x54, 0x50, 0x46, 0x6f, 0x72, 0x6b,
	0x4a, 0x6f, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x36, 0x0a, 0x0b, 0x46,
	0x61, 0x6e, 0x6f, 0x75, 0x74, 0x46, 0x61, 0x6e, 0x69, 0x6e, 0x12, 0x11, 0x2e, 0x68, 0x74, 0x74,
	0x70, 0x2e, 0x48, 0x54, 0x54, 0x50, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e,
	0x68, 0x74, 0x74, 0x70, 0x2e, 0x48, 0x54, 0x54, 0x50, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x30, 0x01, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x3b, 0x68, 0x74, 0x74, 0x70, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_httpforkjoin_proto_rawDescOnce sync.Once
	file_api_httpforkjoin_proto_rawDescData = file_api_httpforkjoin_proto_rawDesc
)

func file_api_httpforkjoin_proto_rawDescGZIP() []byte {
	file_api_httpforkjoin_proto_rawDescOnce.Do(func() {
		file_api_httpforkjoin_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_httpforkjoin_proto_rawDescData)
	})
	return file_api_httpforkjoin_proto_rawDescData
}

var file_api_httpforkjoin_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_api_httpforkjoin_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_api_httpforkjoin_proto_goTypes = []interface{}{
	(ErrorCode)(0),       // 0: http.ErrorCode
	(Message_Method)(0),  // 1: http.Message.Method
	(*Message)(nil),      // 2: http.Message
	(*HTTPRequest)(nil),  // 3: http.HTTPRequest
	(*HTTPResponse)(nil), // 4: http.HTTPResponse
	(*Error)(nil),        // 5: http.Error
	nil,                  // 6: http.Message.HeadersEntry
}
var file_api_httpforkjoin_proto_depIdxs = []int32{
	1, // 0: http.Message.method:type_name -> http.Message.Method
	6, // 1: http.Message.headers:type_name -> http.Message.HeadersEntry
	2, // 2: http.HTTPRequest.messages:type_name -> http.Message
	2, // 3: http.HTTPResponse.message:type_name -> http.Message
	5, // 4: http.HTTPResponse.errors:type_name -> http.Error
	0, // 5: http.Error.code:type_name -> http.ErrorCode
	3, // 6: http.HTTPForkJoinService.FanoutFanin:input_type -> http.HTTPRequest
	4, // 7: http.HTTPForkJoinService.FanoutFanin:output_type -> http.HTTPResponse
	7, // [7:8] is the sub-list for method output_type
	6, // [6:7] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_api_httpforkjoin_proto_init() }
func file_api_httpforkjoin_proto_init() {
	if File_api_httpforkjoin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_httpforkjoin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_httpforkjoin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HTTPRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_httpforkjoin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HTTPResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_httpforkjoin_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_httpforkjoin_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_httpforkjoin_proto_goTypes,
		DependencyIndexes: file_api_httpforkjoin_proto_depIdxs,
		EnumInfos:         file_api_httpforkjoin_proto_enumTypes,
		MessageInfos:      file_api_httpforkjoin_proto_msgTypes,
	}.Build()
	File_api_httpforkjoin_proto = out.File
	file_api_httpforkjoin_proto_rawDesc = nil
	file_api_httpforkjoin_proto_goTypes = nil
	file_api_httpforkjoin_proto_depIdxs = nil
}
