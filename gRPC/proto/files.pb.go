// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.21.12
// source: gRPC/proto/files.proto

package proto

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

type Status int32

const (
	Status_SUCCESS Status = 0
	Status_FAILED  Status = 1
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0: "SUCCESS",
		1: "FAILED",
	}
	Status_value = map[string]int32{
		"SUCCESS": 0,
		"FAILED":  1,
	}
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status) Descriptor() protoreflect.EnumDescriptor {
	return file_gRPC_proto_files_proto_enumTypes[0].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_gRPC_proto_files_proto_enumTypes[0]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_gRPC_proto_files_proto_rawDescGZIP(), []int{0}
}

type FileSendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Request:
	//
	//	*FileSendRequest_FileName
	//	*FileSendRequest_Chunk
	Request isFileSendRequest_Request `protobuf_oneof:"request"`
}

func (x *FileSendRequest) Reset() {
	*x = FileSendRequest{}
	mi := &file_gRPC_proto_files_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileSendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileSendRequest) ProtoMessage() {}

func (x *FileSendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gRPC_proto_files_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileSendRequest.ProtoReflect.Descriptor instead.
func (*FileSendRequest) Descriptor() ([]byte, []int) {
	return file_gRPC_proto_files_proto_rawDescGZIP(), []int{0}
}

func (m *FileSendRequest) GetRequest() isFileSendRequest_Request {
	if m != nil {
		return m.Request
	}
	return nil
}

func (x *FileSendRequest) GetFileName() string {
	if x, ok := x.GetRequest().(*FileSendRequest_FileName); ok {
		return x.FileName
	}
	return ""
}

func (x *FileSendRequest) GetChunk() []byte {
	if x, ok := x.GetRequest().(*FileSendRequest_Chunk); ok {
		return x.Chunk
	}
	return nil
}

type isFileSendRequest_Request interface {
	isFileSendRequest_Request()
}

type FileSendRequest_FileName struct {
	FileName string `protobuf:"bytes,1,opt,name=fileName,proto3,oneof"`
}

type FileSendRequest_Chunk struct {
	Chunk []byte `protobuf:"bytes,2,opt,name=chunk,proto3,oneof"`
}

func (*FileSendRequest_FileName) isFileSendRequest_Request() {}

func (*FileSendRequest_Chunk) isFileSendRequest_Request() {}

type FileSendResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NewFileName string `protobuf:"bytes,1,opt,name=newFileName,proto3" json:"newFileName,omitempty"`
}

func (x *FileSendResponse) Reset() {
	*x = FileSendResponse{}
	mi := &file_gRPC_proto_files_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileSendResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileSendResponse) ProtoMessage() {}

func (x *FileSendResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gRPC_proto_files_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileSendResponse.ProtoReflect.Descriptor instead.
func (*FileSendResponse) Descriptor() ([]byte, []int) {
	return file_gRPC_proto_files_proto_rawDescGZIP(), []int{1}
}

func (x *FileSendResponse) GetNewFileName() string {
	if x != nil {
		return x.NewFileName
	}
	return ""
}

type FileDeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileName string `protobuf:"bytes,1,opt,name=fileName,proto3" json:"fileName,omitempty"`
}

func (x *FileDeleteRequest) Reset() {
	*x = FileDeleteRequest{}
	mi := &file_gRPC_proto_files_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileDeleteRequest) ProtoMessage() {}

func (x *FileDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gRPC_proto_files_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileDeleteRequest.ProtoReflect.Descriptor instead.
func (*FileDeleteRequest) Descriptor() ([]byte, []int) {
	return file_gRPC_proto_files_proto_rawDescGZIP(), []int{2}
}

func (x *FileDeleteRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

type FileDeleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FileDeleteResponse) Reset() {
	*x = FileDeleteResponse{}
	mi := &file_gRPC_proto_files_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileDeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileDeleteResponse) ProtoMessage() {}

func (x *FileDeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gRPC_proto_files_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileDeleteResponse.ProtoReflect.Descriptor instead.
func (*FileDeleteResponse) Descriptor() ([]byte, []int) {
	return file_gRPC_proto_files_proto_rawDescGZIP(), []int{3}
}

type FileGetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileName string `protobuf:"bytes,1,opt,name=fileName,proto3" json:"fileName,omitempty"`
	Start    int64  `protobuf:"varint,2,opt,name=start,proto3" json:"start,omitempty"`
	End      int64  `protobuf:"varint,3,opt,name=end,proto3" json:"end,omitempty"`
}

func (x *FileGetRequest) Reset() {
	*x = FileGetRequest{}
	mi := &file_gRPC_proto_files_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileGetRequest) ProtoMessage() {}

func (x *FileGetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gRPC_proto_files_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileGetRequest.ProtoReflect.Descriptor instead.
func (*FileGetRequest) Descriptor() ([]byte, []int) {
	return file_gRPC_proto_files_proto_rawDescGZIP(), []int{4}
}

func (x *FileGetRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *FileGetRequest) GetStart() int64 {
	if x != nil {
		return x.Start
	}
	return 0
}

func (x *FileGetRequest) GetEnd() int64 {
	if x != nil {
		return x.End
	}
	return 0
}

type FileGetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Response:
	//
	//	*FileGetResponse_FileSize
	//	*FileGetResponse_FileStream
	Response isFileGetResponse_Response `protobuf_oneof:"response"`
}

func (x *FileGetResponse) Reset() {
	*x = FileGetResponse{}
	mi := &file_gRPC_proto_files_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileGetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileGetResponse) ProtoMessage() {}

func (x *FileGetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gRPC_proto_files_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileGetResponse.ProtoReflect.Descriptor instead.
func (*FileGetResponse) Descriptor() ([]byte, []int) {
	return file_gRPC_proto_files_proto_rawDescGZIP(), []int{5}
}

func (m *FileGetResponse) GetResponse() isFileGetResponse_Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func (x *FileGetResponse) GetFileSize() int64 {
	if x, ok := x.GetResponse().(*FileGetResponse_FileSize); ok {
		return x.FileSize
	}
	return 0
}

func (x *FileGetResponse) GetFileStream() []byte {
	if x, ok := x.GetResponse().(*FileGetResponse_FileStream); ok {
		return x.FileStream
	}
	return nil
}

type isFileGetResponse_Response interface {
	isFileGetResponse_Response()
}

type FileGetResponse_FileSize struct {
	FileSize int64 `protobuf:"varint,1,opt,name=fileSize,proto3,oneof"`
}

type FileGetResponse_FileStream struct {
	FileStream []byte `protobuf:"bytes,2,opt,name=fileStream,proto3,oneof"`
}

func (*FileGetResponse_FileSize) isFileGetResponse_Response() {}

func (*FileGetResponse_FileStream) isFileGetResponse_Response() {}

var File_gRPC_proto_files_proto protoreflect.FileDescriptor

var file_gRPC_proto_files_proto_rawDesc = []byte{
	0x0a, 0x16, 0x67, 0x52, 0x50, 0x43, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x69, 0x6c,
	0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x52, 0x0a, 0x0f, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1c, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x16, 0x0a, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48,
	0x00, 0x52, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x42, 0x09, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0x34, 0x0a, 0x10, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x65, 0x77, 0x46, 0x69,
	0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6e, 0x65,
	0x77, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x2f, 0x0a, 0x11, 0x46, 0x69, 0x6c,
	0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x14, 0x0a, 0x12, 0x46, 0x69,
	0x6c, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x54, 0x0a, 0x0e, 0x46, 0x69, 0x6c, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x22, 0x5d, 0x0a, 0x0f, 0x46, 0x69, 0x6c, 0x65, 0x47, 0x65,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x08, 0x66, 0x69, 0x6c,
	0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x20, 0x0a, 0x0a, 0x66, 0x69, 0x6c, 0x65, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x0a, 0x66,
	0x69, 0x6c, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x42, 0x0a, 0x0a, 0x08, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2a, 0x21, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x0b, 0x0a, 0x07, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x01, 0x32, 0xe8, 0x01, 0x0a, 0x0c, 0x46, 0x69, 0x6c,
	0x65, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x45, 0x0a, 0x10, 0x53, 0x65, 0x6e,
	0x64, 0x54, 0x6f, 0x47, 0x52, 0x50, 0x43, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x16, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x69,
	0x6c, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01,
	0x12, 0x4b, 0x0a, 0x14, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x47, 0x52,
	0x50, 0x43, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a,
	0x11, 0x47, 0x65, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x47, 0x52, 0x50, 0x43, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x30, 0x01, 0x42, 0x12, 0x5a, 0x10, 0x67, 0x52, 0x50, 0x43, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gRPC_proto_files_proto_rawDescOnce sync.Once
	file_gRPC_proto_files_proto_rawDescData = file_gRPC_proto_files_proto_rawDesc
)

func file_gRPC_proto_files_proto_rawDescGZIP() []byte {
	file_gRPC_proto_files_proto_rawDescOnce.Do(func() {
		file_gRPC_proto_files_proto_rawDescData = protoimpl.X.CompressGZIP(file_gRPC_proto_files_proto_rawDescData)
	})
	return file_gRPC_proto_files_proto_rawDescData
}

var file_gRPC_proto_files_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_gRPC_proto_files_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_gRPC_proto_files_proto_goTypes = []any{
	(Status)(0),                // 0: proto.Status
	(*FileSendRequest)(nil),    // 1: proto.FileSendRequest
	(*FileSendResponse)(nil),   // 2: proto.FileSendResponse
	(*FileDeleteRequest)(nil),  // 3: proto.FileDeleteRequest
	(*FileDeleteResponse)(nil), // 4: proto.FileDeleteResponse
	(*FileGetRequest)(nil),     // 5: proto.FileGetRequest
	(*FileGetResponse)(nil),    // 6: proto.FileGetResponse
}
var file_gRPC_proto_files_proto_depIdxs = []int32{
	1, // 0: proto.FilesService.SendToGRPCServer:input_type -> proto.FileSendRequest
	3, // 1: proto.FilesService.DeleteFromGRPCServer:input_type -> proto.FileDeleteRequest
	5, // 2: proto.FilesService.GetFromGRPCServer:input_type -> proto.FileGetRequest
	2, // 3: proto.FilesService.SendToGRPCServer:output_type -> proto.FileSendResponse
	4, // 4: proto.FilesService.DeleteFromGRPCServer:output_type -> proto.FileDeleteResponse
	6, // 5: proto.FilesService.GetFromGRPCServer:output_type -> proto.FileGetResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gRPC_proto_files_proto_init() }
func file_gRPC_proto_files_proto_init() {
	if File_gRPC_proto_files_proto != nil {
		return
	}
	file_gRPC_proto_files_proto_msgTypes[0].OneofWrappers = []any{
		(*FileSendRequest_FileName)(nil),
		(*FileSendRequest_Chunk)(nil),
	}
	file_gRPC_proto_files_proto_msgTypes[5].OneofWrappers = []any{
		(*FileGetResponse_FileSize)(nil),
		(*FileGetResponse_FileStream)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gRPC_proto_files_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gRPC_proto_files_proto_goTypes,
		DependencyIndexes: file_gRPC_proto_files_proto_depIdxs,
		EnumInfos:         file_gRPC_proto_files_proto_enumTypes,
		MessageInfos:      file_gRPC_proto_files_proto_msgTypes,
	}.Build()
	File_gRPC_proto_files_proto = out.File
	file_gRPC_proto_files_proto_rawDesc = nil
	file_gRPC_proto_files_proto_goTypes = nil
	file_gRPC_proto_files_proto_depIdxs = nil
}
