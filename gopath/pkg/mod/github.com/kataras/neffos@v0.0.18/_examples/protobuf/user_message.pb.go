// Code generated by protoc-gen-go. DO NOT EDIT.
// source: user_message.proto

package main

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

type UserMessage struct {
	Username             string   `protobuf:"bytes,1,opt,name=Username,proto3" json:"Username,omitempty"`
	Text                 string   `protobuf:"bytes,2,opt,name=Text,proto3" json:"Text,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserMessage) Reset()         { *m = UserMessage{} }
func (m *UserMessage) String() string { return proto.CompactTextString(m) }
func (*UserMessage) ProtoMessage()    {}
func (*UserMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_7965fb6944d08275, []int{0}
}

func (m *UserMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserMessage.Unmarshal(m, b)
}
func (m *UserMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserMessage.Marshal(b, m, deterministic)
}
func (m *UserMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserMessage.Merge(m, src)
}
func (m *UserMessage) XXX_Size() int {
	return xxx_messageInfo_UserMessage.Size(m)
}
func (m *UserMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_UserMessage.DiscardUnknown(m)
}

var xxx_messageInfo_UserMessage proto.InternalMessageInfo

func (m *UserMessage) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *UserMessage) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func init() {
	proto.RegisterType((*UserMessage)(nil), "main.UserMessage")
}

func init() { proto.RegisterFile("user_message.proto", fileDescriptor_7965fb6944d08275) }

var fileDescriptor_7965fb6944d08275 = []byte{
	// 98 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0x2d, 0x4e, 0x2d,
	0x8a, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0xc9, 0x4d, 0xcc, 0xcc, 0x53, 0xb2, 0xe5, 0xe2, 0x0e, 0x2d, 0x4e, 0x2d, 0xf2, 0x85, 0x48, 0x09,
	0x49, 0x71, 0x71, 0x80, 0xb8, 0x79, 0x89, 0xb9, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41,
	0x70, 0xbe, 0x90, 0x10, 0x17, 0x4b, 0x48, 0x6a, 0x45, 0x89, 0x04, 0x13, 0x58, 0x1c, 0xcc, 0x4e,
	0x62, 0x03, 0x9b, 0x65, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x6a, 0x40, 0x7a, 0xd9, 0x61, 0x00,
	0x00, 0x00,
}
