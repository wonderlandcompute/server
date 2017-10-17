// Code generated by protoc-gen-go. DO NOT EDIT.
// source: disneyland.proto

/*
Package disneyland is a generated protocol buffer package.

It is generated from these files:
	disneyland.proto

It has these top-level messages:
	Job
	ListOfJobs
	RequestWithId
	ListJobsRequest
	ListJobStatusResponse
	ListJobsKindRequest
*/
package disneyland

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Job_Status int32

const (
	Job_PENDING   Job_Status = 0
	Job_PULLED    Job_Status = 1
	Job_RUNNING   Job_Status = 2
	Job_FAILED    Job_Status = 3
	Job_COMPLETED Job_Status = 4
)

var Job_Status_name = map[int32]string{
	0: "PENDING",
	1: "PULLED",
	2: "RUNNING",
	3: "FAILED",
	4: "COMPLETED",
}
var Job_Status_value = map[string]int32{
	"PENDING":   0,
	"PULLED":    1,
	"RUNNING":   2,
	"FAILED":    3,
	"COMPLETED": 4,
}

func (x Job_Status) String() string {
	return proto.EnumName(Job_Status_name, int32(x))
}
func (Job_Status) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Job struct {
	Project  string     `protobuf:"bytes,1,opt,name=project" json:"project,omitempty"`
	Id       uint64     `protobuf:"varint,2,opt,name=id" json:"id,omitempty"`
	Kind     string     `protobuf:"bytes,3,opt,name=kind" json:"kind,omitempty"`
	Status   Job_Status `protobuf:"varint,4,opt,name=status,enum=Job_Status" json:"status,omitempty"`
	Input    string     `protobuf:"bytes,5,opt,name=input" json:"input,omitempty"`
	Output   string     `protobuf:"bytes,6,opt,name=output" json:"output,omitempty"`
	Metadata string     `protobuf:"bytes,7,opt,name=metadata" json:"metadata,omitempty"`
}

func (m *Job) Reset()                    { *m = Job{} }
func (m *Job) String() string            { return proto.CompactTextString(m) }
func (*Job) ProtoMessage()               {}
func (*Job) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Job) GetProject() string {
	if m != nil {
		return m.Project
	}
	return ""
}

func (m *Job) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Job) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func (m *Job) GetStatus() Job_Status {
	if m != nil {
		return m.Status
	}
	return Job_PENDING
}

func (m *Job) GetInput() string {
	if m != nil {
		return m.Input
	}
	return ""
}

func (m *Job) GetOutput() string {
	if m != nil {
		return m.Output
	}
	return ""
}

func (m *Job) GetMetadata() string {
	if m != nil {
		return m.Metadata
	}
	return ""
}

type ListOfJobs struct {
	Jobs []*Job `protobuf:"bytes,1,rep,name=jobs" json:"jobs,omitempty"`
}

func (m *ListOfJobs) Reset()                    { *m = ListOfJobs{} }
func (m *ListOfJobs) String() string            { return proto.CompactTextString(m) }
func (*ListOfJobs) ProtoMessage()               {}
func (*ListOfJobs) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ListOfJobs) GetJobs() []*Job {
	if m != nil {
		return m.Jobs
	}
	return nil
}

type RequestWithId struct {
	Id uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
}

func (m *RequestWithId) Reset()                    { *m = RequestWithId{} }
func (m *RequestWithId) String() string            { return proto.CompactTextString(m) }
func (*RequestWithId) ProtoMessage()               {}
func (*RequestWithId) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *RequestWithId) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type ListJobsRequest struct {
	HowMany uint32 `protobuf:"varint,1,opt,name=how_many,json=howMany" json:"how_many,omitempty"`
}

func (m *ListJobsRequest) Reset()                    { *m = ListJobsRequest{} }
func (m *ListJobsRequest) String() string            { return proto.CompactTextString(m) }
func (*ListJobsRequest) ProtoMessage()               {}
func (*ListJobsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ListJobsRequest) GetHowMany() uint32 {
	if m != nil {
		return m.HowMany
	}
	return 0
}

type ListJobStatusResponse struct {
	Status []Job_Status `protobuf:"varint,1,rep,packed,name=status,enum=Job_Status" json:"status,omitempty"`
}

func (m *ListJobStatusResponse) Reset()                    { *m = ListJobStatusResponse{} }
func (m *ListJobStatusResponse) String() string            { return proto.CompactTextString(m) }
func (*ListJobStatusResponse) ProtoMessage()               {}
func (*ListJobStatusResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ListJobStatusResponse) GetStatus() []Job_Status {
	if m != nil {
		return m.Status
	}
	return nil
}

type ListJobsKindRequest struct {
	HowMany uint32 `protobuf:"varint,1,opt,name=how_many,json=howMany" json:"how_many,omitempty"`
	Kind    string `protobuf:"bytes,2,opt,name=kind" json:"kind,omitempty"`
}

func (m *ListJobsKindRequest) Reset()                    { *m = ListJobsKindRequest{} }
func (m *ListJobsKindRequest) String() string            { return proto.CompactTextString(m) }
func (*ListJobsKindRequest) ProtoMessage()               {}
func (*ListJobsKindRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ListJobsKindRequest) GetHowMany() uint32 {
	if m != nil {
		return m.HowMany
	}
	return 0
}

func (m *ListJobsKindRequest) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func init() {
	proto.RegisterType((*Job)(nil), "Job")
	proto.RegisterType((*ListOfJobs)(nil), "ListOfJobs")
	proto.RegisterType((*RequestWithId)(nil), "RequestWithId")
	proto.RegisterType((*ListJobsRequest)(nil), "ListJobsRequest")
	proto.RegisterType((*ListJobStatusResponse)(nil), "ListJobStatusResponse")
	proto.RegisterType((*ListJobsKindRequest)(nil), "ListJobsKindRequest")
	proto.RegisterEnum("Job_Status", Job_Status_name, Job_Status_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Disneyland service

type DisneylandClient interface {
	CreateJob(ctx context.Context, in *Job, opts ...grpc.CallOption) (*Job, error)
	GetJob(ctx context.Context, in *RequestWithId, opts ...grpc.CallOption) (*Job, error)
	ListJobs(ctx context.Context, in *ListJobsRequest, opts ...grpc.CallOption) (*ListOfJobs, error)
	ModifyJob(ctx context.Context, in *Job, opts ...grpc.CallOption) (*Job, error)
	PullPendingJobs(ctx context.Context, in *ListJobsKindRequest, opts ...grpc.CallOption) (*ListOfJobs, error)
}

type disneylandClient struct {
	cc *grpc.ClientConn
}

func NewDisneylandClient(cc *grpc.ClientConn) DisneylandClient {
	return &disneylandClient{cc}
}

func (c *disneylandClient) CreateJob(ctx context.Context, in *Job, opts ...grpc.CallOption) (*Job, error) {
	out := new(Job)
	err := grpc.Invoke(ctx, "/Disneyland/CreateJob", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *disneylandClient) GetJob(ctx context.Context, in *RequestWithId, opts ...grpc.CallOption) (*Job, error) {
	out := new(Job)
	err := grpc.Invoke(ctx, "/Disneyland/GetJob", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *disneylandClient) ListJobs(ctx context.Context, in *ListJobsRequest, opts ...grpc.CallOption) (*ListOfJobs, error) {
	out := new(ListOfJobs)
	err := grpc.Invoke(ctx, "/Disneyland/ListJobs", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *disneylandClient) ModifyJob(ctx context.Context, in *Job, opts ...grpc.CallOption) (*Job, error) {
	out := new(Job)
	err := grpc.Invoke(ctx, "/Disneyland/ModifyJob", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *disneylandClient) PullPendingJobs(ctx context.Context, in *ListJobsKindRequest, opts ...grpc.CallOption) (*ListOfJobs, error) {
	out := new(ListOfJobs)
	err := grpc.Invoke(ctx, "/Disneyland/PullPendingJobs", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Disneyland service

type DisneylandServer interface {
	CreateJob(context.Context, *Job) (*Job, error)
	GetJob(context.Context, *RequestWithId) (*Job, error)
	ListJobs(context.Context, *ListJobsRequest) (*ListOfJobs, error)
	ModifyJob(context.Context, *Job) (*Job, error)
	PullPendingJobs(context.Context, *ListJobsKindRequest) (*ListOfJobs, error)
}

func RegisterDisneylandServer(s *grpc.Server, srv DisneylandServer) {
	s.RegisterService(&_Disneyland_serviceDesc, srv)
}

func _Disneyland_CreateJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Job)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DisneylandServer).CreateJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Disneyland/CreateJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DisneylandServer).CreateJob(ctx, req.(*Job))
	}
	return interceptor(ctx, in, info, handler)
}

func _Disneyland_GetJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestWithId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DisneylandServer).GetJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Disneyland/GetJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DisneylandServer).GetJob(ctx, req.(*RequestWithId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Disneyland_ListJobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListJobsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DisneylandServer).ListJobs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Disneyland/ListJobs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DisneylandServer).ListJobs(ctx, req.(*ListJobsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Disneyland_ModifyJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Job)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DisneylandServer).ModifyJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Disneyland/ModifyJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DisneylandServer).ModifyJob(ctx, req.(*Job))
	}
	return interceptor(ctx, in, info, handler)
}

func _Disneyland_PullPendingJobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListJobsKindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DisneylandServer).PullPendingJobs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Disneyland/PullPendingJobs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DisneylandServer).PullPendingJobs(ctx, req.(*ListJobsKindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Disneyland_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Disneyland",
	HandlerType: (*DisneylandServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateJob",
			Handler:    _Disneyland_CreateJob_Handler,
		},
		{
			MethodName: "GetJob",
			Handler:    _Disneyland_GetJob_Handler,
		},
		{
			MethodName: "ListJobs",
			Handler:    _Disneyland_ListJobs_Handler,
		},
		{
			MethodName: "ModifyJob",
			Handler:    _Disneyland_ModifyJob_Handler,
		},
		{
			MethodName: "PullPendingJobs",
			Handler:    _Disneyland_PullPendingJobs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "disneyland.proto",
}

func init() { proto.RegisterFile("disneyland.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 439 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0xd1, 0x6e, 0xd3, 0x30,
	0x14, 0x6d, 0xd2, 0x2c, 0x6d, 0x6f, 0xd5, 0x2e, 0xba, 0x0c, 0xe4, 0xf5, 0x85, 0xc8, 0x48, 0xa8,
	0x12, 0x28, 0x0f, 0x45, 0xe2, 0x89, 0x17, 0xb4, 0x94, 0xa9, 0xa5, 0xed, 0xa2, 0xc0, 0xc4, 0x23,
	0x4a, 0xb0, 0x47, 0x3d, 0x3a, 0xbb, 0xd4, 0x8e, 0xa6, 0x7e, 0x22, 0x1f, 0x85, 0x84, 0xec, 0xa5,
	0x1d, 0x8c, 0x4a, 0x7b, 0x89, 0x72, 0xef, 0x39, 0x3e, 0xc7, 0xf7, 0xf8, 0x42, 0xc4, 0x84, 0x96,
	0x7c, 0xbb, 0x2a, 0x24, 0x4b, 0xd6, 0x1b, 0x65, 0x14, 0xfd, 0xed, 0x41, 0x73, 0xaa, 0x4a, 0x24,
	0xd0, 0x5a, 0x6f, 0xd4, 0x35, 0xff, 0x66, 0x88, 0x17, 0x7b, 0xc3, 0x4e, 0xbe, 0x2b, 0xb1, 0x0f,
	0xbe, 0x60, 0xc4, 0x8f, 0xbd, 0x61, 0x90, 0xfb, 0x82, 0x21, 0x42, 0xf0, 0x43, 0x48, 0x46, 0x9a,
	0x8e, 0xe6, 0xfe, 0xf1, 0x05, 0x84, 0xda, 0x14, 0xa6, 0xd2, 0x24, 0x88, 0xbd, 0x61, 0x7f, 0xd4,
	0x4d, 0xa6, 0xaa, 0x4c, 0x3e, 0xb9, 0x56, 0x5e, 0x43, 0x78, 0x02, 0x47, 0x42, 0xae, 0x2b, 0x43,
	0x8e, 0xdc, 0xc9, 0xbb, 0x02, 0x9f, 0x41, 0xa8, 0x2a, 0x63, 0xdb, 0xa1, 0x6b, 0xd7, 0x15, 0x0e,
	0xa0, 0x7d, 0xc3, 0x4d, 0xc1, 0x0a, 0x53, 0x90, 0x96, 0x43, 0xf6, 0x35, 0x9d, 0x40, 0x78, 0xa7,
	0x8d, 0x5d, 0x68, 0x65, 0xe3, 0x45, 0x3a, 0x59, 0x9c, 0x47, 0x0d, 0x04, 0x08, 0xb3, 0xcb, 0xd9,
	0x6c, 0x9c, 0x46, 0x9e, 0x05, 0xf2, 0xcb, 0xc5, 0xc2, 0x02, 0xbe, 0x05, 0x3e, 0xbc, 0x9f, 0x58,
	0xa0, 0x89, 0x3d, 0xe8, 0x9c, 0x5d, 0xcc, 0xb3, 0xd9, 0xf8, 0xf3, 0x38, 0x8d, 0x02, 0xfa, 0x12,
	0x60, 0x26, 0xb4, 0xb9, 0xb8, 0x9a, 0xaa, 0x52, 0x23, 0x81, 0xe0, 0x5a, 0x95, 0x9a, 0x78, 0x71,
	0x73, 0xd8, 0x1d, 0x05, 0x76, 0x8a, 0xdc, 0x75, 0xe8, 0x73, 0xe8, 0xe5, 0xfc, 0x67, 0xc5, 0xb5,
	0xf9, 0x22, 0xcc, 0x72, 0xc2, 0xea, 0x58, 0xbc, 0x5d, 0x2c, 0xf4, 0x35, 0x1c, 0x5b, 0x21, 0x2b,
	0x53, 0x13, 0xf1, 0x14, 0xda, 0x4b, 0x75, 0xfb, 0xf5, 0xa6, 0x90, 0x5b, 0x47, 0xec, 0xe5, 0xad,
	0xa5, 0xba, 0x9d, 0x17, 0x72, 0x4b, 0xdf, 0xc1, 0xd3, 0x9a, 0x5d, 0x87, 0xc4, 0xf5, 0x5a, 0x49,
	0xcd, 0xff, 0x4a, 0xd2, 0xde, 0xe1, 0x70, 0x92, 0x34, 0x85, 0x27, 0x3b, 0xaf, 0x8f, 0x42, 0xb2,
	0xc7, 0xfd, 0xf6, 0x8f, 0xe6, 0xdf, 0x3f, 0xda, 0xe8, 0x97, 0x07, 0x90, 0xee, 0xf7, 0x01, 0x4f,
	0xa1, 0x73, 0xb6, 0xe1, 0x85, 0xe1, 0x76, 0x1d, 0xdc, 0xe8, 0x03, 0xf7, 0xa5, 0x0d, 0x8c, 0x21,
	0x3c, 0xe7, 0xd6, 0x0e, 0xfb, 0xc9, 0x3f, 0x29, 0xec, 0x19, 0xaf, 0xa0, 0xbd, 0xbb, 0x11, 0x46,
	0xc9, 0x83, 0x20, 0x06, 0xdd, 0xe4, 0x3e, 0x63, 0xda, 0xb0, 0x4e, 0x73, 0xc5, 0xc4, 0xd5, 0xf6,
	0x7f, 0xa7, 0xb7, 0x70, 0x9c, 0x55, 0xab, 0x55, 0xc6, 0x25, 0x13, 0xf2, 0xbb, 0x93, 0x3b, 0x49,
	0x0e, 0xcc, 0xfa, 0x40, 0xb2, 0x0c, 0xdd, 0x36, 0xbf, 0xf9, 0x13, 0x00, 0x00, 0xff, 0xff, 0xe3,
	0x52, 0xf9, 0x30, 0xe1, 0x02, 0x00, 0x00,
}
