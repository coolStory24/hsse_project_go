// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.29.1
// source: proto/user_service.proto

package gen

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetUserContactDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *GetUserContactDataRequest) Reset() {
	*x = GetUserContactDataRequest{}
	mi := &file_proto_user_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserContactDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserContactDataRequest) ProtoMessage() {}

func (x *GetUserContactDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_user_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserContactDataRequest.ProtoReflect.Descriptor instead.
func (*GetUserContactDataRequest) Descriptor() ([]byte, []int) {
	return file_proto_user_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetUserContactDataRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type GetUserContactDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Phone string `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty"`
}

func (x *GetUserContactDataResponse) Reset() {
	*x = GetUserContactDataResponse{}
	mi := &file_proto_user_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserContactDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserContactDataResponse) ProtoMessage() {}

func (x *GetUserContactDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_user_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserContactDataResponse.ProtoReflect.Descriptor instead.
func (*GetUserContactDataResponse) Descriptor() ([]byte, []int) {
	return file_proto_user_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetUserContactDataResponse) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *GetUserContactDataResponse) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

var File_proto_user_service_proto protoreflect.FileDescriptor

var file_proto_user_service_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x1a,
	0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x34, 0x0a, 0x19, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x48, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f,
	0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x32,
	0x84, 0x01, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x75, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63,
	0x74, 0x44, 0x61, 0x74, 0x65, 0x12, 0x2e, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x42, 0x5a, 0x40, 0x62, 0x6f, 0x6f, 0x6b, 0x69, 0x6e,
	0x67, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x3b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_user_service_proto_rawDescOnce sync.Once
	file_proto_user_service_proto_rawDescData = file_proto_user_service_proto_rawDesc
)

func file_proto_user_service_proto_rawDescGZIP() []byte {
	file_proto_user_service_proto_rawDescOnce.Do(func() {
		file_proto_user_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_user_service_proto_rawDescData)
	})
	return file_proto_user_service_proto_rawDescData
}

var file_proto_user_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_user_service_proto_goTypes = []any{
	(*GetUserContactDataRequest)(nil),  // 0: service_interaction.GetUserContactDataRequest
	(*GetUserContactDataResponse)(nil), // 1: service_interaction.GetUserContactDataResponse
}
var file_proto_user_service_proto_depIdxs = []int32{
	0, // 0: service_interaction.UserService.GetUserContactDate:input_type -> service_interaction.GetUserContactDataRequest
	1, // 1: service_interaction.UserService.GetUserContactDate:output_type -> service_interaction.GetUserContactDataResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_user_service_proto_init() }
func file_proto_user_service_proto_init() {
	if File_proto_user_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_user_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_user_service_proto_goTypes,
		DependencyIndexes: file_proto_user_service_proto_depIdxs,
		MessageInfos:      file_proto_user_service_proto_msgTypes,
	}.Build()
	File_proto_user_service_proto = out.File
	file_proto_user_service_proto_rawDesc = nil
	file_proto_user_service_proto_goTypes = nil
	file_proto_user_service_proto_depIdxs = nil
}
