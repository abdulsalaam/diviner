// Code generated by protoc-gen-go. DO NOT EDIT.
// source: member/member.proto

/*
Package member is a generated protocol buffer package.

It is generated from these files:
	member/member.proto

It has these top-level messages:
	Member
*/
package member

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Member struct {
	Id      string             `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Address string             `protobuf:"bytes,2,opt,name=address" json:"address,omitempty"`
	Balance float64            `protobuf:"fixed64,3,opt,name=balance" json:"balance,omitempty"`
	Assets  map[string]float64 `protobuf:"bytes,4,rep,name=assets" json:"assets,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"fixed64,2,opt,name=value"`
}

func (m *Member) Reset()                    { *m = Member{} }
func (m *Member) String() string            { return proto.CompactTextString(m) }
func (*Member) ProtoMessage()               {}
func (*Member) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Member) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Member) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Member) GetBalance() float64 {
	if m != nil {
		return m.Balance
	}
	return 0
}

func (m *Member) GetAssets() map[string]float64 {
	if m != nil {
		return m.Assets
	}
	return nil
}

func init() {
	proto.RegisterType((*Member)(nil), "member.Member")
}

func init() { proto.RegisterFile("member/member.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 179 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xce, 0x4d, 0xcd, 0x4d,
	0x4a, 0x2d, 0xd2, 0x87, 0x50, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0x6c, 0x10, 0x9e, 0xd2,
	0x6e, 0x46, 0x2e, 0x36, 0x5f, 0x30, 0x53, 0x88, 0x8f, 0x8b, 0x29, 0x33, 0x45, 0x82, 0x51, 0x81,
	0x51, 0x83, 0x33, 0x88, 0x29, 0x33, 0x45, 0x48, 0x82, 0x8b, 0x3d, 0x31, 0x25, 0xa5, 0x28, 0xb5,
	0xb8, 0x58, 0x82, 0x09, 0x2c, 0x08, 0xe3, 0x82, 0x64, 0x92, 0x12, 0x73, 0x12, 0xf3, 0x92, 0x53,
	0x25, 0x98, 0x15, 0x18, 0x35, 0x18, 0x83, 0x60, 0x5c, 0x21, 0x23, 0x2e, 0xb6, 0xc4, 0xe2, 0xe2,
	0xd4, 0x92, 0x62, 0x09, 0x16, 0x05, 0x66, 0x0d, 0x6e, 0x23, 0x29, 0x3d, 0xa8, 0xad, 0x10, 0x3b,
	0xf4, 0x1c, 0xc1, 0x92, 0xae, 0x79, 0x25, 0x45, 0x95, 0x41, 0x50, 0x95, 0x52, 0x96, 0x5c, 0xdc,
	0x48, 0xc2, 0x42, 0x02, 0x5c, 0xcc, 0xd9, 0xa9, 0x95, 0x50, 0x77, 0x80, 0x98, 0x42, 0x22, 0x5c,
	0xac, 0x65, 0x89, 0x39, 0xa5, 0xa9, 0x60, 0x67, 0x30, 0x06, 0x41, 0x38, 0x56, 0x4c, 0x16, 0x8c,
	0x49, 0x6c, 0x60, 0xcf, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x15, 0xcb, 0xe2, 0xd0, 0xe3,
	0x00, 0x00, 0x00,
}
