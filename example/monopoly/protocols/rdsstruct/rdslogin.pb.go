// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.5.1
// source: rdslogin.proto

package rdsstruct

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

type RdsPlayerData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UID             string            `protobuf:"bytes,1,opt,name=UID,proto3" json:"UID,omitempty"`
	DisplayName     string            `protobuf:"bytes,2,opt,name=DisplayName,proto3" json:"DisplayName,omitempty"`
	ClientIdAddress string            `protobuf:"bytes,3,opt,name=ClientIdAddress,proto3" json:"ClientIdAddress,omitempty"`
	ClientIdId      string            `protobuf:"bytes,4,opt,name=ClientIdId,proto3" json:"ClientIdId,omitempty"`
	Externs         map[string]string `protobuf:"bytes,5,rep,name=Externs,proto3" json:"Externs,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *RdsPlayerData) Reset() {
	*x = RdsPlayerData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rdslogin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RdsPlayerData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RdsPlayerData) ProtoMessage() {}

func (x *RdsPlayerData) ProtoReflect() protoreflect.Message {
	mi := &file_rdslogin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RdsPlayerData.ProtoReflect.Descriptor instead.
func (*RdsPlayerData) Descriptor() ([]byte, []int) {
	return file_rdslogin_proto_rawDescGZIP(), []int{0}
}

func (x *RdsPlayerData) GetUID() string {
	if x != nil {
		return x.UID
	}
	return ""
}

func (x *RdsPlayerData) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *RdsPlayerData) GetClientIdAddress() string {
	if x != nil {
		return x.ClientIdAddress
	}
	return ""
}

func (x *RdsPlayerData) GetClientIdId() string {
	if x != nil {
		return x.ClientIdId
	}
	return ""
}

func (x *RdsPlayerData) GetExterns() map[string]string {
	if x != nil {
		return x.Externs
	}
	return nil
}

var File_rdslogin_proto protoreflect.FileDescriptor

var file_rdslogin_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x64, 0x73, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x09, 0x72, 0x64, 0x73, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x22, 0x8a, 0x02, 0x0a, 0x0d,
	0x52, 0x64, 0x73, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a,
	0x03, 0x55, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x55, 0x49, 0x44, 0x12,
	0x20, 0x0a, 0x0b, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x28, 0x0a, 0x0f, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x43, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x49, 0x64, 0x12, 0x3f, 0x0a, 0x07, 0x45,
	0x78, 0x74, 0x65, 0x72, 0x6e, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x72,
	0x64, 0x73, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x52, 0x64, 0x73, 0x50, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x07, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x73, 0x1a, 0x3a, 0x0a, 0x0c,
	0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x3b, 0x72, 0x64,
	0x73, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rdslogin_proto_rawDescOnce sync.Once
	file_rdslogin_proto_rawDescData = file_rdslogin_proto_rawDesc
)

func file_rdslogin_proto_rawDescGZIP() []byte {
	file_rdslogin_proto_rawDescOnce.Do(func() {
		file_rdslogin_proto_rawDescData = protoimpl.X.CompressGZIP(file_rdslogin_proto_rawDescData)
	})
	return file_rdslogin_proto_rawDescData
}

var file_rdslogin_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rdslogin_proto_goTypes = []interface{}{
	(*RdsPlayerData)(nil), // 0: rdsstruct.RdsPlayerData
	nil,                   // 1: rdsstruct.RdsPlayerData.ExternsEntry
}
var file_rdslogin_proto_depIdxs = []int32{
	1, // 0: rdsstruct.RdsPlayerData.Externs:type_name -> rdsstruct.RdsPlayerData.ExternsEntry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rdslogin_proto_init() }
func file_rdslogin_proto_init() {
	if File_rdslogin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rdslogin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RdsPlayerData); i {
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
			RawDescriptor: file_rdslogin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rdslogin_proto_goTypes,
		DependencyIndexes: file_rdslogin_proto_depIdxs,
		MessageInfos:      file_rdslogin_proto_msgTypes,
	}.Build()
	File_rdslogin_proto = out.File
	file_rdslogin_proto_rawDesc = nil
	file_rdslogin_proto_goTypes = nil
	file_rdslogin_proto_depIdxs = nil
}
