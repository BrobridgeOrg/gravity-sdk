// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: types/product_event/product_event.proto

package product_event

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

type Method int32

const (
	Method_INSERT   Method = 0
	Method_UPDATE   Method = 1
	Method_DELETE   Method = 2
	Method_TRUNCATE Method = 3
)

// Enum value maps for Method.
var (
	Method_name = map[int32]string{
		0: "INSERT",
		1: "UPDATE",
		2: "DELETE",
		3: "TRUNCATE",
	}
	Method_value = map[string]int32{
		"INSERT":   0,
		"UPDATE":   1,
		"DELETE":   2,
		"TRUNCATE": 3,
	}
)

func (x Method) Enum() *Method {
	p := new(Method)
	*p = x
	return p
}

func (x Method) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Method) Descriptor() protoreflect.EnumDescriptor {
	return file_types_product_event_product_event_proto_enumTypes[0].Descriptor()
}

func (Method) Type() protoreflect.EnumType {
	return &file_types_product_event_product_event_proto_enumTypes[0]
}

func (x Method) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Method.Descriptor instead.
func (Method) EnumDescriptor() ([]byte, []int) {
	return file_types_product_event_product_event_proto_rawDescGZIP(), []int{0}
}

type ProductEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventName   string   `protobuf:"bytes,1,opt,name=eventName,proto3" json:"eventName,omitempty"`
	Table       string   `protobuf:"bytes,2,opt,name=table,proto3" json:"table,omitempty"`
	Method      Method   `protobuf:"varint,3,opt,name=method,proto3,enum=gravity.sdk.types.product_event.Method" json:"method,omitempty"`
	PrimaryKeys []string `protobuf:"bytes,4,rep,name=primaryKeys,proto3" json:"primaryKeys,omitempty"`
	PrimaryKey  []byte   `protobuf:"bytes,5,opt,name=primaryKey,proto3" json:"primaryKey,omitempty"`
	Data        []byte   `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ProductEvent) Reset() {
	*x = ProductEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_types_product_event_product_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProductEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProductEvent) ProtoMessage() {}

func (x *ProductEvent) ProtoReflect() protoreflect.Message {
	mi := &file_types_product_event_product_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProductEvent.ProtoReflect.Descriptor instead.
func (*ProductEvent) Descriptor() ([]byte, []int) {
	return file_types_product_event_product_event_proto_rawDescGZIP(), []int{0}
}

func (x *ProductEvent) GetEventName() string {
	if x != nil {
		return x.EventName
	}
	return ""
}

func (x *ProductEvent) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *ProductEvent) GetMethod() Method {
	if x != nil {
		return x.Method
	}
	return Method_INSERT
}

func (x *ProductEvent) GetPrimaryKeys() []string {
	if x != nil {
		return x.PrimaryKeys
	}
	return nil
}

func (x *ProductEvent) GetPrimaryKey() []byte {
	if x != nil {
		return x.PrimaryKey
	}
	return nil
}

func (x *ProductEvent) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_types_product_event_product_event_proto protoreflect.FileDescriptor

var file_types_product_event_product_event_proto_rawDesc = []byte{
	0x0a, 0x27, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1f, 0x67, 0x72, 0x61, 0x76, 0x69,
	0x74, 0x79, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0xd9, 0x01, 0x0a, 0x0c, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x61, 0x62,
	0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x12,
	0x3f, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x27, 0x2e, 0x67, 0x72, 0x61, 0x76, 0x69, 0x74, 0x79, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x4b, 0x65, 0x79, 0x73, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x4b, 0x65,
	0x79, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x4b, 0x65, 0x79,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x4b,
	0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x2a, 0x3a, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x12, 0x0a, 0x0a, 0x06, 0x49, 0x4e, 0x53, 0x45, 0x52, 0x54, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06,
	0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x45, 0x4c, 0x45,
	0x54, 0x45, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x54, 0x52, 0x55, 0x4e, 0x43, 0x41, 0x54, 0x45,
	0x10, 0x03, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x42, 0x72, 0x6f, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x4f, 0x72, 0x67, 0x2f, 0x67, 0x72,
	0x61, 0x76, 0x69, 0x74, 0x79, 0x2d, 0x73, 0x64, 0x6b, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_types_product_event_product_event_proto_rawDescOnce sync.Once
	file_types_product_event_product_event_proto_rawDescData = file_types_product_event_product_event_proto_rawDesc
)

func file_types_product_event_product_event_proto_rawDescGZIP() []byte {
	file_types_product_event_product_event_proto_rawDescOnce.Do(func() {
		file_types_product_event_product_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_types_product_event_product_event_proto_rawDescData)
	})
	return file_types_product_event_product_event_proto_rawDescData
}

var file_types_product_event_product_event_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_types_product_event_product_event_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_types_product_event_product_event_proto_goTypes = []interface{}{
	(Method)(0),          // 0: gravity.sdk.types.product_event.Method
	(*ProductEvent)(nil), // 1: gravity.sdk.types.product_event.ProductEvent
}
var file_types_product_event_product_event_proto_depIdxs = []int32{
	0, // 0: gravity.sdk.types.product_event.ProductEvent.method:type_name -> gravity.sdk.types.product_event.Method
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_types_product_event_product_event_proto_init() }
func file_types_product_event_product_event_proto_init() {
	if File_types_product_event_product_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_types_product_event_product_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProductEvent); i {
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
			RawDescriptor: file_types_product_event_product_event_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_types_product_event_product_event_proto_goTypes,
		DependencyIndexes: file_types_product_event_product_event_proto_depIdxs,
		EnumInfos:         file_types_product_event_product_event_proto_enumTypes,
		MessageInfos:      file_types_product_event_product_event_proto_msgTypes,
	}.Build()
	File_types_product_event_product_event_proto = out.File
	file_types_product_event_product_event_proto_rawDesc = nil
	file_types_product_event_product_event_proto_goTypes = nil
	file_types_product_event_product_event_proto_depIdxs = nil
}
