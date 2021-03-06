// Code generated by protoc-gen-go. DO NOT EDIT.
// source: event.proto

package gravity_sdk_types_event

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Event_Type int32

const (
	Event_TYPE_EVENT    Event_Type = 0
	Event_TYPE_SNAPSHOT Event_Type = 1
	Event_TYPE_SYSTEM   Event_Type = 3
)

var Event_Type_name = map[int32]string{
	0: "TYPE_EVENT",
	1: "TYPE_SNAPSHOT",
	3: "TYPE_SYSTEM",
}

var Event_Type_value = map[string]int32{
	"TYPE_EVENT":    0,
	"TYPE_SNAPSHOT": 1,
	"TYPE_SYSTEM":   3,
}

func (x Event_Type) String() string {
	return proto.EnumName(Event_Type_name, int32(x))
}

func (Event_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{0, 0}
}

type EventPayload_State int32

const (
	EventPayload_STATE_NONE      EventPayload_State = 0
	EventPayload_STATE_CHUNK_END EventPayload_State = 1
)

var EventPayload_State_name = map[int32]string{
	0: "STATE_NONE",
	1: "STATE_CHUNK_END",
}

var EventPayload_State_value = map[string]int32{
	"STATE_NONE":      0,
	"STATE_CHUNK_END": 1,
}

func (x EventPayload_State) String() string {
	return proto.EnumName(EventPayload_State_name, int32(x))
}

func (EventPayload_State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{1, 0}
}

type SnapshotInfo_State int32

const (
	SnapshotInfo_STATE_NONE      SnapshotInfo_State = 0
	SnapshotInfo_STATE_CHUNK_END SnapshotInfo_State = 1
)

var SnapshotInfo_State_name = map[int32]string{
	0: "STATE_NONE",
	1: "STATE_CHUNK_END",
}

var SnapshotInfo_State_value = map[string]int32{
	"STATE_NONE":      0,
	"STATE_CHUNK_END": 1,
}

func (x SnapshotInfo_State) String() string {
	return proto.EnumName(SnapshotInfo_State_name, int32(x))
}

func (SnapshotInfo_State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{2, 0}
}

type SystemMessage_Type int32

const (
	SystemMessage_TYPE_WAKE SystemMessage_Type = 0
)

var SystemMessage_Type_name = map[int32]string{
	0: "TYPE_WAKE",
}

var SystemMessage_Type_value = map[string]int32{
	"TYPE_WAKE": 0,
}

func (x SystemMessage_Type) String() string {
	return proto.EnumName(SystemMessage_Type_name, int32(x))
}

func (SystemMessage_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{3, 0}
}

type Event struct {
	Type                 Event_Type     `protobuf:"varint,1,opt,name=type,proto3,enum=gravity.sdk.types.event.Event_Type" json:"type,omitempty"`
	EventPayload         *EventPayload  `protobuf:"bytes,2,opt,name=eventPayload,proto3" json:"eventPayload,omitempty"`
	SnapshotInfo         *SnapshotInfo  `protobuf:"bytes,3,opt,name=snapshotInfo,proto3" json:"snapshotInfo,omitempty"`
	SystemInfo           *SystemMessage `protobuf:"bytes,4,opt,name=systemInfo,proto3" json:"systemInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{0}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetType() Event_Type {
	if m != nil {
		return m.Type
	}
	return Event_TYPE_EVENT
}

func (m *Event) GetEventPayload() *EventPayload {
	if m != nil {
		return m.EventPayload
	}
	return nil
}

func (m *Event) GetSnapshotInfo() *SnapshotInfo {
	if m != nil {
		return m.SnapshotInfo
	}
	return nil
}

func (m *Event) GetSystemInfo() *SystemMessage {
	if m != nil {
		return m.SystemInfo
	}
	return nil
}

type EventPayload struct {
	PipelineID           uint64             `protobuf:"varint,1,opt,name=pipelineID,proto3" json:"pipelineID,omitempty"`
	Sequence             uint64             `protobuf:"varint,2,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Data                 []byte             `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	State                EventPayload_State `protobuf:"varint,4,opt,name=state,proto3,enum=gravity.sdk.types.event.EventPayload_State" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *EventPayload) Reset()         { *m = EventPayload{} }
func (m *EventPayload) String() string { return proto.CompactTextString(m) }
func (*EventPayload) ProtoMessage()    {}
func (*EventPayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{1}
}

func (m *EventPayload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventPayload.Unmarshal(m, b)
}
func (m *EventPayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventPayload.Marshal(b, m, deterministic)
}
func (m *EventPayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventPayload.Merge(m, src)
}
func (m *EventPayload) XXX_Size() int {
	return xxx_messageInfo_EventPayload.Size(m)
}
func (m *EventPayload) XXX_DiscardUnknown() {
	xxx_messageInfo_EventPayload.DiscardUnknown(m)
}

var xxx_messageInfo_EventPayload proto.InternalMessageInfo

func (m *EventPayload) GetPipelineID() uint64 {
	if m != nil {
		return m.PipelineID
	}
	return 0
}

func (m *EventPayload) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *EventPayload) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *EventPayload) GetState() EventPayload_State {
	if m != nil {
		return m.State
	}
	return EventPayload_STATE_NONE
}

type SnapshotInfo struct {
	PipelineID           uint64             `protobuf:"varint,1,opt,name=pipelineID,proto3" json:"pipelineID,omitempty"`
	SnapshotID           string             `protobuf:"bytes,2,opt,name=snapshotID,proto3" json:"snapshotID,omitempty"`
	Collection           string             `protobuf:"bytes,3,opt,name=collection,proto3" json:"collection,omitempty"`
	Key                  []byte             `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty"`
	Data                 []byte             `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	State                SnapshotInfo_State `protobuf:"varint,6,opt,name=state,proto3,enum=gravity.sdk.types.event.SnapshotInfo_State" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *SnapshotInfo) Reset()         { *m = SnapshotInfo{} }
func (m *SnapshotInfo) String() string { return proto.CompactTextString(m) }
func (*SnapshotInfo) ProtoMessage()    {}
func (*SnapshotInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{2}
}

func (m *SnapshotInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SnapshotInfo.Unmarshal(m, b)
}
func (m *SnapshotInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SnapshotInfo.Marshal(b, m, deterministic)
}
func (m *SnapshotInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SnapshotInfo.Merge(m, src)
}
func (m *SnapshotInfo) XXX_Size() int {
	return xxx_messageInfo_SnapshotInfo.Size(m)
}
func (m *SnapshotInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_SnapshotInfo.DiscardUnknown(m)
}

var xxx_messageInfo_SnapshotInfo proto.InternalMessageInfo

func (m *SnapshotInfo) GetPipelineID() uint64 {
	if m != nil {
		return m.PipelineID
	}
	return 0
}

func (m *SnapshotInfo) GetSnapshotID() string {
	if m != nil {
		return m.SnapshotID
	}
	return ""
}

func (m *SnapshotInfo) GetCollection() string {
	if m != nil {
		return m.Collection
	}
	return ""
}

func (m *SnapshotInfo) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *SnapshotInfo) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SnapshotInfo) GetState() SnapshotInfo_State {
	if m != nil {
		return m.State
	}
	return SnapshotInfo_STATE_NONE
}

type SystemMessage struct {
	Type                 SystemMessage_Type          `protobuf:"varint,1,opt,name=type,proto3,enum=gravity.sdk.types.event.SystemMessage_Type" json:"type,omitempty"`
	AwakeMessage         *SystemMessage_AwakeMessage `protobuf:"bytes,2,opt,name=awakeMessage,proto3" json:"awakeMessage,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *SystemMessage) Reset()         { *m = SystemMessage{} }
func (m *SystemMessage) String() string { return proto.CompactTextString(m) }
func (*SystemMessage) ProtoMessage()    {}
func (*SystemMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{3}
}

func (m *SystemMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemMessage.Unmarshal(m, b)
}
func (m *SystemMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemMessage.Marshal(b, m, deterministic)
}
func (m *SystemMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemMessage.Merge(m, src)
}
func (m *SystemMessage) XXX_Size() int {
	return xxx_messageInfo_SystemMessage.Size(m)
}
func (m *SystemMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemMessage.DiscardUnknown(m)
}

var xxx_messageInfo_SystemMessage proto.InternalMessageInfo

func (m *SystemMessage) GetType() SystemMessage_Type {
	if m != nil {
		return m.Type
	}
	return SystemMessage_TYPE_WAKE
}

func (m *SystemMessage) GetAwakeMessage() *SystemMessage_AwakeMessage {
	if m != nil {
		return m.AwakeMessage
	}
	return nil
}

type SystemMessage_AwakeMessage struct {
	PipelineID           uint64   `protobuf:"varint,1,opt,name=pipelineID,proto3" json:"pipelineID,omitempty"`
	Sequence             uint64   `protobuf:"varint,2,opt,name=Sequence,proto3" json:"Sequence,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SystemMessage_AwakeMessage) Reset()         { *m = SystemMessage_AwakeMessage{} }
func (m *SystemMessage_AwakeMessage) String() string { return proto.CompactTextString(m) }
func (*SystemMessage_AwakeMessage) ProtoMessage()    {}
func (*SystemMessage_AwakeMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{3, 0}
}

func (m *SystemMessage_AwakeMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemMessage_AwakeMessage.Unmarshal(m, b)
}
func (m *SystemMessage_AwakeMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemMessage_AwakeMessage.Marshal(b, m, deterministic)
}
func (m *SystemMessage_AwakeMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemMessage_AwakeMessage.Merge(m, src)
}
func (m *SystemMessage_AwakeMessage) XXX_Size() int {
	return xxx_messageInfo_SystemMessage_AwakeMessage.Size(m)
}
func (m *SystemMessage_AwakeMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemMessage_AwakeMessage.DiscardUnknown(m)
}

var xxx_messageInfo_SystemMessage_AwakeMessage proto.InternalMessageInfo

func (m *SystemMessage_AwakeMessage) GetPipelineID() uint64 {
	if m != nil {
		return m.PipelineID
	}
	return 0
}

func (m *SystemMessage_AwakeMessage) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func init() {
	proto.RegisterEnum("gravity.sdk.types.event.Event_Type", Event_Type_name, Event_Type_value)
	proto.RegisterEnum("gravity.sdk.types.event.EventPayload_State", EventPayload_State_name, EventPayload_State_value)
	proto.RegisterEnum("gravity.sdk.types.event.SnapshotInfo_State", SnapshotInfo_State_name, SnapshotInfo_State_value)
	proto.RegisterEnum("gravity.sdk.types.event.SystemMessage_Type", SystemMessage_Type_name, SystemMessage_Type_value)
	proto.RegisterType((*Event)(nil), "gravity.sdk.types.event.Event")
	proto.RegisterType((*EventPayload)(nil), "gravity.sdk.types.event.EventPayload")
	proto.RegisterType((*SnapshotInfo)(nil), "gravity.sdk.types.event.SnapshotInfo")
	proto.RegisterType((*SystemMessage)(nil), "gravity.sdk.types.event.SystemMessage")
	proto.RegisterType((*SystemMessage_AwakeMessage)(nil), "gravity.sdk.types.event.SystemMessage.AwakeMessage")
}

func init() { proto.RegisterFile("event.proto", fileDescriptor_2d17a9d3f0ddf27e) }

var fileDescriptor_2d17a9d3f0ddf27e = []byte{
	// 471 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x8d, 0x13, 0xa7, 0x22, 0x13, 0x27, 0x35, 0x8b, 0x10, 0x51, 0x0f, 0x55, 0x65, 0x04, 0xaa,
	0x54, 0xe4, 0x43, 0x7a, 0x40, 0xe2, 0x82, 0x2c, 0xb2, 0xa8, 0xa1, 0xaa, 0x1b, 0xad, 0x0d, 0x55,
	0x4f, 0xd1, 0x92, 0x0c, 0x25, 0x4a, 0xb0, 0x4d, 0x77, 0x29, 0xf2, 0x07, 0x70, 0xe3, 0xdf, 0xb8,
	0xf3, 0x35, 0x68, 0x27, 0x21, 0x5a, 0x23, 0x42, 0x22, 0x71, 0xdb, 0x99, 0x37, 0xef, 0x29, 0xef,
	0x65, 0xc6, 0xd0, 0xc6, 0x3b, 0xcc, 0x74, 0x58, 0xdc, 0xe6, 0x3a, 0x67, 0x8f, 0x6e, 0x6e, 0xe5,
	0xdd, 0x4c, 0x97, 0xa1, 0x9a, 0xce, 0x43, 0x5d, 0x16, 0xa8, 0x42, 0x82, 0x83, 0x9f, 0x75, 0x68,
	0x72, 0xf3, 0x62, 0xcf, 0xc1, 0x35, 0x40, 0xcf, 0x39, 0x72, 0x8e, 0xbb, 0xfd, 0xc7, 0xe1, 0x06,
	0x46, 0x48, 0xd3, 0x61, 0x5a, 0x16, 0x28, 0x88, 0xc0, 0x86, 0xe0, 0x11, 0x32, 0x92, 0xe5, 0x22,
	0x97, 0xd3, 0x5e, 0xfd, 0xc8, 0x39, 0x6e, 0xf7, 0x9f, 0xfc, 0x5b, 0x60, 0x35, 0x2c, 0x2a, 0x54,
	0x23, 0xa5, 0x32, 0x59, 0xa8, 0x8f, 0xb9, 0x1e, 0x66, 0x1f, 0xf2, 0x5e, 0x63, 0x8b, 0x54, 0x62,
	0x0d, 0x8b, 0x0a, 0x95, 0xbd, 0x06, 0x50, 0xa5, 0xd2, 0xf8, 0x89, 0x84, 0x5c, 0x12, 0x7a, 0xba,
	0x59, 0x88, 0x46, 0x2f, 0x50, 0x29, 0x79, 0x83, 0xc2, 0x62, 0x06, 0x2f, 0xc0, 0x35, 0x5e, 0x59,
	0x17, 0x20, 0xbd, 0x1e, 0xf1, 0x31, 0x7f, 0xc7, 0xe3, 0xd4, 0xaf, 0xb1, 0xfb, 0xd0, 0xa1, 0x3a,
	0x89, 0xa3, 0x51, 0x72, 0x76, 0x99, 0xfa, 0x0e, 0xdb, 0x87, 0xf6, 0xb2, 0x75, 0x9d, 0xa4, 0xfc,
	0xc2, 0x6f, 0x04, 0x3f, 0x1c, 0xf0, 0x6c, 0xb7, 0xec, 0x10, 0xa0, 0x98, 0x15, 0xb8, 0x98, 0x65,
	0x38, 0x1c, 0x50, 0xd2, 0xae, 0xb0, 0x3a, 0xec, 0x00, 0xee, 0x29, 0xfc, 0xfc, 0x05, 0xb3, 0x09,
	0x52, 0x8c, 0xae, 0x58, 0xd7, 0x8c, 0x81, 0x3b, 0x95, 0x5a, 0x52, 0x26, 0x9e, 0xa0, 0x37, 0x8b,
	0xa0, 0xa9, 0xb4, 0xd4, 0x48, 0xfe, 0xba, 0xfd, 0x93, 0x9d, 0x32, 0x0f, 0x13, 0x43, 0x11, 0x4b,
	0x66, 0xf0, 0x0c, 0x9a, 0x54, 0x1b, 0x83, 0x49, 0x1a, 0xa5, 0x7c, 0x1c, 0x5f, 0xc6, 0xdc, 0xaf,
	0xb1, 0x07, 0xb0, 0xbf, 0xac, 0x5f, 0x9d, 0xbd, 0x8d, 0xcf, 0xc7, 0x3c, 0x1e, 0xf8, 0x4e, 0xf0,
	0xad, 0x0e, 0x9e, 0x1d, 0xfa, 0x56, 0x47, 0x87, 0x00, 0xeb, 0xbf, 0x65, 0x40, 0x9e, 0x5a, 0xc2,
	0xea, 0x18, 0x7c, 0x92, 0x2f, 0x16, 0x38, 0xd1, 0xb3, 0x3c, 0x23, 0x6f, 0x2d, 0x61, 0x75, 0x98,
	0x0f, 0x8d, 0x39, 0x96, 0xe4, 0xcf, 0x13, 0xe6, 0xb9, 0xce, 0xa1, 0xf9, 0xb7, 0x1c, 0xf6, 0xb6,
	0xe4, 0x60, 0xff, 0xf6, 0xff, 0xc9, 0xe1, 0x7b, 0x1d, 0x3a, 0x95, 0x9d, 0x61, 0x2f, 0x2b, 0xe7,
	0x73, 0xb2, 0xdb, 0xa6, 0xd9, 0x67, 0x74, 0x05, 0x9e, 0xfc, 0x2a, 0xe7, 0xb8, 0x82, 0x56, 0x67,
	0x74, 0xba, 0xa3, 0x50, 0x64, 0x51, 0x45, 0x45, 0xe8, 0xe0, 0x0d, 0x78, 0x36, 0xba, 0xcb, 0x12,
	0x26, 0x7f, 0x2c, 0xe1, 0xef, 0x3a, 0x78, 0xb8, 0xba, 0x86, 0x0e, 0xb4, 0x68, 0xd5, 0xaf, 0xa2,
	0x73, 0xee, 0xd7, 0xde, 0xef, 0xd1, 0x57, 0xe6, 0xf4, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0x35,
	0xb2, 0x7c, 0x1c, 0x74, 0x04, 0x00, 0x00,
}
