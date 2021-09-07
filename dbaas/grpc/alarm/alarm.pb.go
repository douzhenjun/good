// Code generated by protoc-gen-go. DO NOT EDIT.
// source: alarm.proto

package alarmpb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type SendRequest struct {
	CollectId            int32    `protobuf:"varint,1,opt,name=collect_id,json=collectId,proto3" json:"collect_id,omitempty"`
	Interval             int32    `protobuf:"varint,2,opt,name=interval,proto3" json:"interval,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendRequest) Reset()         { *m = SendRequest{} }
func (m *SendRequest) String() string { return proto.CompactTextString(m) }
func (*SendRequest) ProtoMessage()    {}
func (*SendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a4142572412ce8e, []int{0}
}

func (m *SendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendRequest.Unmarshal(m, b)
}
func (m *SendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendRequest.Marshal(b, m, deterministic)
}
func (m *SendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendRequest.Merge(m, src)
}
func (m *SendRequest) XXX_Size() int {
	return xxx_messageInfo_SendRequest.Size(m)
}
func (m *SendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendRequest proto.InternalMessageInfo

func (m *SendRequest) GetCollectId() int32 {
	if m != nil {
		return m.CollectId
	}
	return 0
}

func (m *SendRequest) GetInterval() int32 {
	if m != nil {
		return m.Interval
	}
	return 0
}

type GetAlarmStatusRequest struct {
	ModelObject          string   `protobuf:"bytes,1,opt,name=model_object,json=modelObject,proto3" json:"model_object,omitempty"`
	InstId               string   `protobuf:"bytes,2,opt,name=inst_id,json=instId,proto3" json:"inst_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAlarmStatusRequest) Reset()         { *m = GetAlarmStatusRequest{} }
func (m *GetAlarmStatusRequest) String() string { return proto.CompactTextString(m) }
func (*GetAlarmStatusRequest) ProtoMessage()    {}
func (*GetAlarmStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a4142572412ce8e, []int{1}
}

func (m *GetAlarmStatusRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAlarmStatusRequest.Unmarshal(m, b)
}
func (m *GetAlarmStatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAlarmStatusRequest.Marshal(b, m, deterministic)
}
func (m *GetAlarmStatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAlarmStatusRequest.Merge(m, src)
}
func (m *GetAlarmStatusRequest) XXX_Size() int {
	return xxx_messageInfo_GetAlarmStatusRequest.Size(m)
}
func (m *GetAlarmStatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAlarmStatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetAlarmStatusRequest proto.InternalMessageInfo

func (m *GetAlarmStatusRequest) GetModelObject() string {
	if m != nil {
		return m.ModelObject
	}
	return ""
}

func (m *GetAlarmStatusRequest) GetInstId() string {
	if m != nil {
		return m.InstId
	}
	return ""
}

type GetAlarmItemRequest struct {
	ModelId              string   `protobuf:"bytes,1,opt,name=model_id,json=modelId,proto3" json:"model_id,omitempty"`
	InstId               string   `protobuf:"bytes,2,opt,name=inst_id,json=instId,proto3" json:"inst_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAlarmItemRequest) Reset()         { *m = GetAlarmItemRequest{} }
func (m *GetAlarmItemRequest) String() string { return proto.CompactTextString(m) }
func (*GetAlarmItemRequest) ProtoMessage()    {}
func (*GetAlarmItemRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a4142572412ce8e, []int{2}
}

func (m *GetAlarmItemRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAlarmItemRequest.Unmarshal(m, b)
}
func (m *GetAlarmItemRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAlarmItemRequest.Marshal(b, m, deterministic)
}
func (m *GetAlarmItemRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAlarmItemRequest.Merge(m, src)
}
func (m *GetAlarmItemRequest) XXX_Size() int {
	return xxx_messageInfo_GetAlarmItemRequest.Size(m)
}
func (m *GetAlarmItemRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAlarmItemRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetAlarmItemRequest proto.InternalMessageInfo

func (m *GetAlarmItemRequest) GetModelId() string {
	if m != nil {
		return m.ModelId
	}
	return ""
}

func (m *GetAlarmItemRequest) GetInstId() string {
	if m != nil {
		return m.InstId
	}
	return ""
}

type SendRespond struct {
	Errorno              int32    `protobuf:"varint,1,opt,name=errorno,proto3" json:"errorno,omitempty"`
	Data                 string   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendRespond) Reset()         { *m = SendRespond{} }
func (m *SendRespond) String() string { return proto.CompactTextString(m) }
func (*SendRespond) ProtoMessage()    {}
func (*SendRespond) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a4142572412ce8e, []int{3}
}

func (m *SendRespond) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendRespond.Unmarshal(m, b)
}
func (m *SendRespond) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendRespond.Marshal(b, m, deterministic)
}
func (m *SendRespond) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendRespond.Merge(m, src)
}
func (m *SendRespond) XXX_Size() int {
	return xxx_messageInfo_SendRespond.Size(m)
}
func (m *SendRespond) XXX_DiscardUnknown() {
	xxx_messageInfo_SendRespond.DiscardUnknown(m)
}

var xxx_messageInfo_SendRespond proto.InternalMessageInfo

func (m *SendRespond) GetErrorno() int32 {
	if m != nil {
		return m.Errorno
	}
	return 0
}

func (m *SendRespond) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

type GetAlarmStatusRespond struct {
	Errorno              int32    `protobuf:"varint,1,opt,name=errorno,proto3" json:"errorno,omitempty"`
	Data                 string   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAlarmStatusRespond) Reset()         { *m = GetAlarmStatusRespond{} }
func (m *GetAlarmStatusRespond) String() string { return proto.CompactTextString(m) }
func (*GetAlarmStatusRespond) ProtoMessage()    {}
func (*GetAlarmStatusRespond) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a4142572412ce8e, []int{4}
}

func (m *GetAlarmStatusRespond) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAlarmStatusRespond.Unmarshal(m, b)
}
func (m *GetAlarmStatusRespond) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAlarmStatusRespond.Marshal(b, m, deterministic)
}
func (m *GetAlarmStatusRespond) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAlarmStatusRespond.Merge(m, src)
}
func (m *GetAlarmStatusRespond) XXX_Size() int {
	return xxx_messageInfo_GetAlarmStatusRespond.Size(m)
}
func (m *GetAlarmStatusRespond) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAlarmStatusRespond.DiscardUnknown(m)
}

var xxx_messageInfo_GetAlarmStatusRespond proto.InternalMessageInfo

func (m *GetAlarmStatusRespond) GetErrorno() int32 {
	if m != nil {
		return m.Errorno
	}
	return 0
}

func (m *GetAlarmStatusRespond) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

type GetAlarmItemRespond struct {
	Errorno              int32    `protobuf:"varint,1,opt,name=errorno,proto3" json:"errorno,omitempty"`
	Data                 string   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAlarmItemRespond) Reset()         { *m = GetAlarmItemRespond{} }
func (m *GetAlarmItemRespond) String() string { return proto.CompactTextString(m) }
func (*GetAlarmItemRespond) ProtoMessage()    {}
func (*GetAlarmItemRespond) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a4142572412ce8e, []int{5}
}

func (m *GetAlarmItemRespond) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAlarmItemRespond.Unmarshal(m, b)
}
func (m *GetAlarmItemRespond) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAlarmItemRespond.Marshal(b, m, deterministic)
}
func (m *GetAlarmItemRespond) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAlarmItemRespond.Merge(m, src)
}
func (m *GetAlarmItemRespond) XXX_Size() int {
	return xxx_messageInfo_GetAlarmItemRespond.Size(m)
}
func (m *GetAlarmItemRespond) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAlarmItemRespond.DiscardUnknown(m)
}

var xxx_messageInfo_GetAlarmItemRespond proto.InternalMessageInfo

func (m *GetAlarmItemRespond) GetErrorno() int32 {
	if m != nil {
		return m.Errorno
	}
	return 0
}

func (m *GetAlarmItemRespond) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

type GetAlarmResultByInstidRequest struct {
	InstId               string   `protobuf:"bytes,1,opt,name=inst_id,json=instId,proto3" json:"inst_id,omitempty"`
	UserTag              string   `protobuf:"bytes,2,opt,name=user_tag,json=userTag,proto3" json:"user_tag,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAlarmResultByInstidRequest) Reset()         { *m = GetAlarmResultByInstidRequest{} }
func (m *GetAlarmResultByInstidRequest) String() string { return proto.CompactTextString(m) }
func (*GetAlarmResultByInstidRequest) ProtoMessage()    {}
func (*GetAlarmResultByInstidRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a4142572412ce8e, []int{6}
}

func (m *GetAlarmResultByInstidRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAlarmResultByInstidRequest.Unmarshal(m, b)
}
func (m *GetAlarmResultByInstidRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAlarmResultByInstidRequest.Marshal(b, m, deterministic)
}
func (m *GetAlarmResultByInstidRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAlarmResultByInstidRequest.Merge(m, src)
}
func (m *GetAlarmResultByInstidRequest) XXX_Size() int {
	return xxx_messageInfo_GetAlarmResultByInstidRequest.Size(m)
}
func (m *GetAlarmResultByInstidRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAlarmResultByInstidRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetAlarmResultByInstidRequest proto.InternalMessageInfo

func (m *GetAlarmResultByInstidRequest) GetInstId() string {
	if m != nil {
		return m.InstId
	}
	return ""
}

func (m *GetAlarmResultByInstidRequest) GetUserTag() string {
	if m != nil {
		return m.UserTag
	}
	return ""
}

type GetAlarmResultByInstidRespond struct {
	Errorno              int32    `protobuf:"varint,1,opt,name=errorno,proto3" json:"errorno,omitempty"`
	ErrorMsgZh           string   `protobuf:"bytes,2,opt,name=error_msg_zh,json=errorMsgZh,proto3" json:"error_msg_zh,omitempty"`
	ErrorMsgEn           string   `protobuf:"bytes,3,opt,name=error_msg_en,json=errorMsgEn,proto3" json:"error_msg_en,omitempty"`
	Data                 string   `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAlarmResultByInstidRespond) Reset()         { *m = GetAlarmResultByInstidRespond{} }
func (m *GetAlarmResultByInstidRespond) String() string { return proto.CompactTextString(m) }
func (*GetAlarmResultByInstidRespond) ProtoMessage()    {}
func (*GetAlarmResultByInstidRespond) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a4142572412ce8e, []int{7}
}

func (m *GetAlarmResultByInstidRespond) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAlarmResultByInstidRespond.Unmarshal(m, b)
}
func (m *GetAlarmResultByInstidRespond) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAlarmResultByInstidRespond.Marshal(b, m, deterministic)
}
func (m *GetAlarmResultByInstidRespond) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAlarmResultByInstidRespond.Merge(m, src)
}
func (m *GetAlarmResultByInstidRespond) XXX_Size() int {
	return xxx_messageInfo_GetAlarmResultByInstidRespond.Size(m)
}
func (m *GetAlarmResultByInstidRespond) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAlarmResultByInstidRespond.DiscardUnknown(m)
}

var xxx_messageInfo_GetAlarmResultByInstidRespond proto.InternalMessageInfo

func (m *GetAlarmResultByInstidRespond) GetErrorno() int32 {
	if m != nil {
		return m.Errorno
	}
	return 0
}

func (m *GetAlarmResultByInstidRespond) GetErrorMsgZh() string {
	if m != nil {
		return m.ErrorMsgZh
	}
	return ""
}

func (m *GetAlarmResultByInstidRespond) GetErrorMsgEn() string {
	if m != nil {
		return m.ErrorMsgEn
	}
	return ""
}

func (m *GetAlarmResultByInstidRespond) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

type GetAlarmResultRequest struct {
	InstIds              string   `protobuf:"bytes,1,opt,name=inst_ids,json=instIds,proto3" json:"inst_ids,omitempty"`
	UserTag              string   `protobuf:"bytes,2,opt,name=user_tag,json=userTag,proto3" json:"user_tag,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAlarmResultRequest) Reset()         { *m = GetAlarmResultRequest{} }
func (m *GetAlarmResultRequest) String() string { return proto.CompactTextString(m) }
func (*GetAlarmResultRequest) ProtoMessage()    {}
func (*GetAlarmResultRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a4142572412ce8e, []int{8}
}

func (m *GetAlarmResultRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAlarmResultRequest.Unmarshal(m, b)
}
func (m *GetAlarmResultRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAlarmResultRequest.Marshal(b, m, deterministic)
}
func (m *GetAlarmResultRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAlarmResultRequest.Merge(m, src)
}
func (m *GetAlarmResultRequest) XXX_Size() int {
	return xxx_messageInfo_GetAlarmResultRequest.Size(m)
}
func (m *GetAlarmResultRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAlarmResultRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetAlarmResultRequest proto.InternalMessageInfo

func (m *GetAlarmResultRequest) GetInstIds() string {
	if m != nil {
		return m.InstIds
	}
	return ""
}

func (m *GetAlarmResultRequest) GetUserTag() string {
	if m != nil {
		return m.UserTag
	}
	return ""
}

type GetAlarmResultRespond struct {
	Errorno              int32    `protobuf:"varint,1,opt,name=errorno,proto3" json:"errorno,omitempty"`
	ErrorMsgZh           string   `protobuf:"bytes,2,opt,name=error_msg_zh,json=errorMsgZh,proto3" json:"error_msg_zh,omitempty"`
	ErrorMsgEn           string   `protobuf:"bytes,3,opt,name=error_msg_en,json=errorMsgEn,proto3" json:"error_msg_en,omitempty"`
	Data                 string   `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAlarmResultRespond) Reset()         { *m = GetAlarmResultRespond{} }
func (m *GetAlarmResultRespond) String() string { return proto.CompactTextString(m) }
func (*GetAlarmResultRespond) ProtoMessage()    {}
func (*GetAlarmResultRespond) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a4142572412ce8e, []int{9}
}

func (m *GetAlarmResultRespond) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAlarmResultRespond.Unmarshal(m, b)
}
func (m *GetAlarmResultRespond) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAlarmResultRespond.Marshal(b, m, deterministic)
}
func (m *GetAlarmResultRespond) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAlarmResultRespond.Merge(m, src)
}
func (m *GetAlarmResultRespond) XXX_Size() int {
	return xxx_messageInfo_GetAlarmResultRespond.Size(m)
}
func (m *GetAlarmResultRespond) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAlarmResultRespond.DiscardUnknown(m)
}

var xxx_messageInfo_GetAlarmResultRespond proto.InternalMessageInfo

func (m *GetAlarmResultRespond) GetErrorno() int32 {
	if m != nil {
		return m.Errorno
	}
	return 0
}

func (m *GetAlarmResultRespond) GetErrorMsgZh() string {
	if m != nil {
		return m.ErrorMsgZh
	}
	return ""
}

func (m *GetAlarmResultRespond) GetErrorMsgEn() string {
	if m != nil {
		return m.ErrorMsgEn
	}
	return ""
}

func (m *GetAlarmResultRespond) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*SendRequest)(nil), "alarmpb.SendRequest")
	proto.RegisterType((*GetAlarmStatusRequest)(nil), "alarmpb.GetAlarmStatusRequest")
	proto.RegisterType((*GetAlarmItemRequest)(nil), "alarmpb.GetAlarmItemRequest")
	proto.RegisterType((*SendRespond)(nil), "alarmpb.SendRespond")
	proto.RegisterType((*GetAlarmStatusRespond)(nil), "alarmpb.GetAlarmStatusRespond")
	proto.RegisterType((*GetAlarmItemRespond)(nil), "alarmpb.GetAlarmItemRespond")
	proto.RegisterType((*GetAlarmResultByInstidRequest)(nil), "alarmpb.GetAlarmResultByInstidRequest")
	proto.RegisterType((*GetAlarmResultByInstidRespond)(nil), "alarmpb.GetAlarmResultByInstidRespond")
	proto.RegisterType((*GetAlarmResultRequest)(nil), "alarmpb.GetAlarmResultRequest")
	proto.RegisterType((*GetAlarmResultRespond)(nil), "alarmpb.GetAlarmResultRespond")
}

func init() { proto.RegisterFile("alarm.proto", fileDescriptor_4a4142572412ce8e) }

var fileDescriptor_4a4142572412ce8e = []byte{
	// 438 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x94, 0xcd, 0xae, 0x93, 0x40,
	0x14, 0xc7, 0xe5, 0xde, 0xde, 0xd2, 0x9e, 0x36, 0x2e, 0xc6, 0x2f, 0x24, 0xd6, 0x54, 0x16, 0xc6,
	0x55, 0x17, 0xea, 0xce, 0x95, 0x9a, 0x46, 0x49, 0x6c, 0x4c, 0xc0, 0x95, 0x1b, 0x32, 0xed, 0x4c,
	0x28, 0x0d, 0xcc, 0x54, 0x66, 0x30, 0xd1, 0x37, 0x70, 0xe9, 0x13, 0xf8, 0xaa, 0x86, 0x61, 0x06,
	0xa1, 0x80, 0x4d, 0xba, 0xb9, 0xbb, 0x39, 0x1f, 0xf3, 0xe3, 0xfc, 0xcf, 0x9c, 0x03, 0xcc, 0x70,
	0x8a, 0xf3, 0x6c, 0x75, 0xcc, 0xb9, 0xe4, 0xc8, 0x56, 0xc6, 0x71, 0xeb, 0x7d, 0x84, 0x59, 0x48,
	0x19, 0x09, 0xe8, 0xb7, 0x82, 0x0a, 0x89, 0x16, 0x00, 0x3b, 0x9e, 0xa6, 0x74, 0x27, 0xa3, 0x84,
	0x38, 0xd6, 0xd2, 0x7a, 0x71, 0x13, 0x4c, 0xb5, 0xc7, 0x27, 0xc8, 0x85, 0x49, 0xc2, 0x24, 0xcd,
	0xbf, 0xe3, 0xd4, 0xb9, 0x52, 0xc1, 0xda, 0xf6, 0x42, 0x78, 0xf0, 0x81, 0xca, 0xb7, 0x25, 0x37,
	0x94, 0x58, 0x16, 0xc2, 0x30, 0x9f, 0xc1, 0x3c, 0xe3, 0x84, 0xa6, 0x11, 0xdf, 0x1e, 0xe8, 0x4e,
	0x2a, 0xea, 0x34, 0x98, 0x29, 0xdf, 0x67, 0xe5, 0x42, 0x8f, 0xc0, 0x4e, 0x98, 0x50, 0xdf, 0xbc,
	0x52, 0xd1, 0x71, 0x69, 0xfa, 0xc4, 0xf3, 0xe1, 0x9e, 0x81, 0xfa, 0x92, 0x66, 0x06, 0xf9, 0x18,
	0x26, 0x15, 0x52, 0x17, 0x39, 0x0d, 0x6c, 0x65, 0xfb, 0x64, 0x18, 0xf5, 0xc6, 0x28, 0x15, 0x47,
	0xce, 0x08, 0x72, 0xc0, 0xa6, 0x79, 0xce, 0x73, 0xc6, 0xb5, 0x4c, 0x63, 0x22, 0x04, 0x23, 0x82,
	0x25, 0xd6, 0xd7, 0xd5, 0xd9, 0x5b, 0x77, 0xc5, 0x5d, 0x82, 0x79, 0x7f, 0x2a, 0xe7, 0x12, 0x48,
	0x08, 0x0b, 0x03, 0x09, 0xa8, 0x28, 0x52, 0xf9, 0xee, 0x87, 0xcf, 0x84, 0x4c, 0xea, 0x47, 0x6c,
	0xb4, 0xc0, 0x6a, 0xb6, 0xa0, 0x6c, 0x5b, 0x21, 0x68, 0x1e, 0x49, 0x1c, 0x6b, 0xa2, 0x5d, 0xda,
	0x5f, 0x70, 0xec, 0xfd, 0xb6, 0x86, 0xa9, 0xe7, 0x8a, 0x5c, 0xc2, 0x5c, 0x1d, 0xa3, 0x4c, 0xc4,
	0xd1, 0xcf, 0xbd, 0x46, 0x83, 0xf2, 0x6d, 0x44, 0xfc, 0x75, 0xdf, 0xce, 0xa0, 0xcc, 0xb9, 0x6e,
	0x67, 0xac, 0x59, 0x2d, 0x74, 0xd4, 0x10, 0xba, 0xf9, 0xd7, 0xf4, 0xaa, 0xa4, 0xc6, 0xf3, 0x6b,
	0x81, 0xc2, 0x3c, 0x7f, 0xa5, 0x50, 0xfc, 0x4f, 0xe2, 0x2f, 0xab, 0xcb, 0xbb, 0x25, 0x69, 0x2f,
	0xff, 0x5c, 0xc3, 0x8d, 0x2a, 0x04, 0xbd, 0x86, 0x51, 0x39, 0x96, 0xe8, 0xfe, 0x4a, 0xaf, 0xe4,
	0xaa, 0xb1, 0x8f, 0xee, 0xa9, 0x57, 0xd5, 0xeb, 0xdd, 0x41, 0x01, 0xdc, 0x6d, 0xcf, 0x23, 0x7a,
	0x5a, 0x67, 0xf6, 0x6e, 0xa1, 0x3b, 0x1c, 0x37, 0xcc, 0x4f, 0x30, 0x6f, 0x0e, 0x27, 0x7a, 0xd2,
	0xb9, 0xd1, 0x58, 0x41, 0x77, 0x28, 0x6a, 0x68, 0x07, 0x78, 0xd8, 0x3f, 0x4f, 0xe8, 0x79, 0xe7,
	0x66, 0xef, 0x18, 0xbb, 0xe7, 0xf3, 0x7a, 0xba, 0x51, 0xa5, 0xf4, 0x74, 0xa3, 0x35, 0x41, 0xee,
	0x70, 0x5c, 0x33, 0xb7, 0x63, 0xf5, 0xa3, 0x7c, 0xf5, 0x37, 0x00, 0x00, 0xff, 0xff, 0x2a, 0x3a,
	0x38, 0x6c, 0x37, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AlarmClient is the client API for Alarm service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AlarmClient interface {
	Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*SendRespond, error)
	GetAlarmStatus(ctx context.Context, in *GetAlarmStatusRequest, opts ...grpc.CallOption) (*GetAlarmStatusRespond, error)
	GetAlarmItem(ctx context.Context, in *GetAlarmItemRequest, opts ...grpc.CallOption) (*GetAlarmItemRespond, error)
	GetAlarmResultByInstid(ctx context.Context, in *GetAlarmResultByInstidRequest, opts ...grpc.CallOption) (*GetAlarmResultByInstidRespond, error)
	GetAlarmResult(ctx context.Context, in *GetAlarmResultRequest, opts ...grpc.CallOption) (*GetAlarmResultRespond, error)
}

type alarmClient struct {
	cc grpc.ClientConnInterface
}

func NewAlarmClient(cc grpc.ClientConnInterface) AlarmClient {
	return &alarmClient{cc}
}

func (c *alarmClient) Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*SendRespond, error) {
	out := new(SendRespond)
	err := c.cc.Invoke(ctx, "/alarmpb.Alarm/Send", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alarmClient) GetAlarmStatus(ctx context.Context, in *GetAlarmStatusRequest, opts ...grpc.CallOption) (*GetAlarmStatusRespond, error) {
	out := new(GetAlarmStatusRespond)
	err := c.cc.Invoke(ctx, "/alarmpb.Alarm/GetAlarmStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alarmClient) GetAlarmItem(ctx context.Context, in *GetAlarmItemRequest, opts ...grpc.CallOption) (*GetAlarmItemRespond, error) {
	out := new(GetAlarmItemRespond)
	err := c.cc.Invoke(ctx, "/alarmpb.Alarm/GetAlarmItem", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alarmClient) GetAlarmResultByInstid(ctx context.Context, in *GetAlarmResultByInstidRequest, opts ...grpc.CallOption) (*GetAlarmResultByInstidRespond, error) {
	out := new(GetAlarmResultByInstidRespond)
	err := c.cc.Invoke(ctx, "/alarmpb.Alarm/GetAlarmResultByInstid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alarmClient) GetAlarmResult(ctx context.Context, in *GetAlarmResultRequest, opts ...grpc.CallOption) (*GetAlarmResultRespond, error) {
	out := new(GetAlarmResultRespond)
	err := c.cc.Invoke(ctx, "/alarmpb.Alarm/GetAlarmResult", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AlarmServer is the server API for Alarm service.
type AlarmServer interface {
	Send(context.Context, *SendRequest) (*SendRespond, error)
	GetAlarmStatus(context.Context, *GetAlarmStatusRequest) (*GetAlarmStatusRespond, error)
	GetAlarmItem(context.Context, *GetAlarmItemRequest) (*GetAlarmItemRespond, error)
	GetAlarmResultByInstid(context.Context, *GetAlarmResultByInstidRequest) (*GetAlarmResultByInstidRespond, error)
	GetAlarmResult(context.Context, *GetAlarmResultRequest) (*GetAlarmResultRespond, error)
}

// UnimplementedAlarmServer can be embedded to have forward compatible implementations.
type UnimplementedAlarmServer struct {
}

func (*UnimplementedAlarmServer) Send(ctx context.Context, req *SendRequest) (*SendRespond, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (*UnimplementedAlarmServer) GetAlarmStatus(ctx context.Context, req *GetAlarmStatusRequest) (*GetAlarmStatusRespond, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAlarmStatus not implemented")
}
func (*UnimplementedAlarmServer) GetAlarmItem(ctx context.Context, req *GetAlarmItemRequest) (*GetAlarmItemRespond, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAlarmItem not implemented")
}
func (*UnimplementedAlarmServer) GetAlarmResultByInstid(ctx context.Context, req *GetAlarmResultByInstidRequest) (*GetAlarmResultByInstidRespond, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAlarmResultByInstid not implemented")
}
func (*UnimplementedAlarmServer) GetAlarmResult(ctx context.Context, req *GetAlarmResultRequest) (*GetAlarmResultRespond, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAlarmResult not implemented")
}

func RegisterAlarmServer(s *grpc.Server, srv AlarmServer) {
	s.RegisterService(&_Alarm_serviceDesc, srv)
}

func _Alarm_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlarmServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alarmpb.Alarm/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlarmServer).Send(ctx, req.(*SendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Alarm_GetAlarmStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAlarmStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlarmServer).GetAlarmStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alarmpb.Alarm/GetAlarmStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlarmServer).GetAlarmStatus(ctx, req.(*GetAlarmStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Alarm_GetAlarmItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAlarmItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlarmServer).GetAlarmItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alarmpb.Alarm/GetAlarmItem",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlarmServer).GetAlarmItem(ctx, req.(*GetAlarmItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Alarm_GetAlarmResultByInstid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAlarmResultByInstidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlarmServer).GetAlarmResultByInstid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alarmpb.Alarm/GetAlarmResultByInstid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlarmServer).GetAlarmResultByInstid(ctx, req.(*GetAlarmResultByInstidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Alarm_GetAlarmResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAlarmResultRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlarmServer).GetAlarmResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alarmpb.Alarm/GetAlarmResult",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlarmServer).GetAlarmResult(ctx, req.(*GetAlarmResultRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Alarm_serviceDesc = grpc.ServiceDesc{
	ServiceName: "alarmpb.Alarm",
	HandlerType: (*AlarmServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _Alarm_Send_Handler,
		},
		{
			MethodName: "GetAlarmStatus",
			Handler:    _Alarm_GetAlarmStatus_Handler,
		},
		{
			MethodName: "GetAlarmItem",
			Handler:    _Alarm_GetAlarmItem_Handler,
		},
		{
			MethodName: "GetAlarmResultByInstid",
			Handler:    _Alarm_GetAlarmResultByInstid_Handler,
		},
		{
			MethodName: "GetAlarmResult",
			Handler:    _Alarm_GetAlarmResult_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "alarm.proto",
}
