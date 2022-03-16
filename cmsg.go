package baticli

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

type ClientMsgType int32

const (
	ClientMsgType_Unused   ClientMsgType = 0
	ClientMsgType_Init     ClientMsgType = 1
	ClientMsgType_InitResp ClientMsgType = 2
	ClientMsgType_Biz      ClientMsgType = 3
	ClientMsgType_Ack      ClientMsgType = 4
	ClientMsgType_Echo     ClientMsgType = 100
)

// Enum value maps for ClientMsgType.
var (
	ClientMsgType_name = map[int32]string{
		0:   "Unused",
		1:   "Init",
		2:   "InitResp",
		3:   "Biz",
		4:   "Ack",
		100: "Echo",
	}
	ClientMsgType_value = map[string]int32{
		"Unused":   0,
		"Init":     1,
		"InitResp": 2,
		"Biz":      3,
		"Ack":      4,
		"Echo":     100,
	}
)

func (m ClientMsgType) Enum() *ClientMsgType {
	p := new(ClientMsgType)
	*p = m
	return p
}

func (m ClientMsgType) String() string {
	return protoimpl.X.EnumStringOf(m.Descriptor(), protoreflect.EnumNumber(m))
}

func (ClientMsgType) Descriptor() protoreflect.EnumDescriptor {
	return file_cmsg_proto_enumTypes[0].Descriptor()
}

func (ClientMsgType) Type() protoreflect.EnumType {
	return &file_cmsg_proto_enumTypes[0]
}

func (m ClientMsgType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(m)
}

// Deprecated: Use ClientMsgType.Descriptor instead.
func (ClientMsgType) EnumDescriptor() ([]byte, []int) {
	return file_cmsg_proto_rawDescGZIP(), []int{0}
}

type Compressor int32

const (
	Compressor_Null    Compressor = 0
	Compressor_Deflate Compressor = 1
	Compressor_Zstd    Compressor = 2
)

// Enum value maps for Compressor.
var (
	Compressor_name = map[int32]string{
		0: "Null",
		1: "Deflate",
		2: "Zstd",
	}
	Compressor_value = map[string]int32{
		"Null":    0,
		"Deflate": 1,
		"Zstd":    2,
	}
)

func (x Compressor) Enum() *Compressor {
	p := new(Compressor)
	*p = x
	return p
}

func (x Compressor) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Compressor) Descriptor() protoreflect.EnumDescriptor {
	return file_cmsg_proto_enumTypes[1].Descriptor()
}

func (Compressor) Type() protoreflect.EnumType {
	return &file_cmsg_proto_enumTypes[1]
}

func (x Compressor) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Compressor.Descriptor instead.
func (Compressor) EnumDescriptor() ([]byte, []int) {
	return file_cmsg_proto_rawDescGZIP(), []int{1}
}

type ClientMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type       ClientMsgType `protobuf:"varint,2,opt,name=type,proto3,enum=cmsg.ClientMsgType" json:"type,omitempty"`
	Ack        int32         `protobuf:"varint,3,opt,name=ack,proto3" json:"ack,omitempty"`
	ServiceId  *string       `protobuf:"bytes,4,opt,name=service_id,json=serviceId,proto3,oneof" json:"service_id,omitempty"`
	Compressor *Compressor   `protobuf:"varint,5,opt,name=compressor,proto3,enum=cmsg.Compressor,oneof" json:"compressor,omitempty"`
	BizData    []byte        `protobuf:"bytes,6,opt,name=biz_data,json=bizData,proto3,oneof" json:"biz_data,omitempty"`
	InitData   *InitData     `protobuf:"bytes,7,opt,name=init_data,json=initData,proto3,oneof" json:"init_data,omitempty"`
}

func (msg *ClientMsg) Reset() {
	*msg = ClientMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cmsg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(msg))
		ms.StoreMessageInfo(mi)
	}
}

func (msg *ClientMsg) String() string {
	return protoimpl.X.MessageStringOf(msg)
}

func (*ClientMsg) ProtoMessage() {}

func (msg *ClientMsg) ProtoReflect() protoreflect.Message {
	mi := &file_cmsg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && msg != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(msg))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(msg)
}

// Deprecated: Use ClientMsg.ProtoReflect.Descriptor instead.
func (*ClientMsg) Descriptor() ([]byte, []int) {
	return file_cmsg_proto_rawDescGZIP(), []int{0}
}

func (msg *ClientMsg) GetId() string {
	if msg != nil {
		return msg.Id
	}
	return ""
}

func (msg *ClientMsg) GetType() ClientMsgType {
	if msg != nil {
		return msg.Type
	}
	return ClientMsgType_Unused
}

func (msg *ClientMsg) GetAck() int32 {
	if msg != nil {
		return msg.Ack
	}
	return 0
}

func (msg *ClientMsg) GetServiceId() string {
	if msg != nil && msg.ServiceId != nil {
		return *msg.ServiceId
	}
	return ""
}

func (msg *ClientMsg) GetCompressor() Compressor {
	if msg != nil && msg.Compressor != nil {
		return *msg.Compressor
	}
	return Compressor_Null
}

func (msg *ClientMsg) GetBizData() []byte {
	if msg != nil {
		return msg.BizData
	}
	return nil
}

func (msg *ClientMsg) GetInitData() *InitData {
	if msg != nil {
		return msg.InitData
	}
	return nil
}

type InitData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AcceptCompressor Compressor `protobuf:"varint,1,opt,name=accept_compressor,json=acceptCompressor,proto3,enum=cmsg.Compressor" json:"accept_compressor,omitempty"`
	PingInterval     uint32     `protobuf:"varint,2,opt,name=ping_interval,json=pingInterval,proto3" json:"ping_interval,omitempty"`
}

func (x *InitData) Reset() {
	*x = InitData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cmsg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitData) ProtoMessage() {}

func (x *InitData) ProtoReflect() protoreflect.Message {
	mi := &file_cmsg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitData.ProtoReflect.Descriptor instead.
func (*InitData) Descriptor() ([]byte, []int) {
	return file_cmsg_proto_rawDescGZIP(), []int{1}
}

func (x *InitData) GetAcceptCompressor() Compressor {
	if x != nil {
		return x.AcceptCompressor
	}
	return Compressor_Null
}

func (x *InitData) GetPingInterval() uint32 {
	if x != nil {
		return x.PingInterval
	}
	return 0
}

var File_cmsg_proto protoreflect.FileDescriptor

var file_cmsg_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x63, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x6d,
	0x73, 0x67, 0x22, 0xbc, 0x02, 0x0a, 0x09, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x73, 0x67,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x27, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13,
	0x2e, 0x63, 0x6d, 0x73, 0x67, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x73, 0x67, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x63, 0x6b,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x61, 0x63, 0x6b, 0x12, 0x22, 0x0a, 0x0a, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x09, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12,
	0x35, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x63, 0x6d, 0x73, 0x67, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x72,
	0x65, 0x73, 0x73, 0x6f, 0x72, 0x48, 0x01, 0x52, 0x0a, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73,
	0x73, 0x6f, 0x72, 0x88, 0x01, 0x01, 0x12, 0x1e, 0x0a, 0x08, 0x62, 0x69, 0x7a, 0x5f, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x02, 0x52, 0x07, 0x62, 0x69, 0x7a, 0x44,
	0x61, 0x74, 0x61, 0x88, 0x01, 0x01, 0x12, 0x30, 0x0a, 0x09, 0x69, 0x6e, 0x69, 0x74, 0x5f, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x6d, 0x73, 0x67,
	0x2e, 0x49, 0x6e, 0x69, 0x74, 0x44, 0x61, 0x74, 0x61, 0x48, 0x03, 0x52, 0x08, 0x69, 0x6e, 0x69,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x6f, 0x6d, 0x70,
	0x72, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x62, 0x69, 0x7a, 0x5f, 0x64,
	0x61, 0x74, 0x61, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x69, 0x6e, 0x69, 0x74, 0x5f, 0x64, 0x61, 0x74,
	0x61, 0x22, 0x6e, 0x0a, 0x08, 0x49, 0x6e, 0x69, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x3d, 0x0a,
	0x11, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73,
	0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x63, 0x6d, 0x73, 0x67, 0x2e,
	0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x52, 0x10, 0x61, 0x63, 0x63, 0x65,
	0x70, 0x74, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x23, 0x0a, 0x0d,
	0x70, 0x69, 0x6e, 0x67, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0c, 0x70, 0x69, 0x6e, 0x67, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61,
	0x6c, 0x2a, 0x4f, 0x0a, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x73, 0x67, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x55, 0x6e, 0x75, 0x73, 0x65, 0x64, 0x10, 0x00, 0x12, 0x08,
	0x0a, 0x04, 0x49, 0x6e, 0x69, 0x74, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x49, 0x6e, 0x69, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x10, 0x02, 0x12, 0x07, 0x0a, 0x03, 0x42, 0x69, 0x7a, 0x10, 0x03, 0x12,
	0x07, 0x0a, 0x03, 0x41, 0x63, 0x6b, 0x10, 0x04, 0x12, 0x08, 0x0a, 0x04, 0x45, 0x63, 0x68, 0x6f,
	0x10, 0x64, 0x2a, 0x2d, 0x0a, 0x0a, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x6f, 0x72,
	0x12, 0x08, 0x0a, 0x04, 0x4e, 0x75, 0x6c, 0x6c, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x65,
	0x66, 0x6c, 0x61, 0x74, 0x65, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x5a, 0x73, 0x74, 0x64, 0x10,
	0x02, 0x42, 0x18, 0x5a, 0x16, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x62, 0x61, 0x74, 0x69, 0x67, 0x6f, 0x2f, 0x63, 0x6d, 0x73, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_cmsg_proto_rawDescOnce sync.Once
	file_cmsg_proto_rawDescData = file_cmsg_proto_rawDesc
)

func file_cmsg_proto_rawDescGZIP() []byte {
	file_cmsg_proto_rawDescOnce.Do(func() {
		file_cmsg_proto_rawDescData = protoimpl.X.CompressGZIP(file_cmsg_proto_rawDescData)
	})
	return file_cmsg_proto_rawDescData
}

var file_cmsg_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_cmsg_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_cmsg_proto_goTypes = []interface{}{
	(ClientMsgType)(0), // 0: cmsg.ClientMsgType
	(Compressor)(0),    // 1: cmsg.Compressor
	(*ClientMsg)(nil),  // 2: cmsg.ClientMsg
	(*InitData)(nil),   // 3: cmsg.InitData
}
var file_cmsg_proto_depIdxs = []int32{
	0, // 0: cmsg.ClientMsg.type:type_name -> cmsg.ClientMsgType
	1, // 1: cmsg.ClientMsg.compressor:type_name -> cmsg.Compressor
	3, // 2: cmsg.ClientMsg.init_data:type_name -> cmsg.InitData
	1, // 3: cmsg.InitData.accept_compressor:type_name -> cmsg.Compressor
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_cmsg_proto_init() }
func file_cmsg_proto_init() {
	if File_cmsg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cmsg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientMsg); i {
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
		file_cmsg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitData); i {
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
	file_cmsg_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cmsg_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cmsg_proto_goTypes,
		DependencyIndexes: file_cmsg_proto_depIdxs,
		EnumInfos:         file_cmsg_proto_enumTypes,
		MessageInfos:      file_cmsg_proto_msgTypes,
	}.Build()
	File_cmsg_proto = out.File
	file_cmsg_proto_rawDesc = nil
	file_cmsg_proto_goTypes = nil
	file_cmsg_proto_depIdxs = nil
}
