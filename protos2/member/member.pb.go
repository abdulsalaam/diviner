// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protos/member/member.proto

/*
Package member is a generated protocol buffer package.

It is generated from these files:
	protos/member/member.proto

It has these top-level messages:
	Member
	Transfer
*/
package member

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

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
	Address string             `protobuf:"bytes,1,opt,name=address" json:"address,omitempty"`
	Balance float64            `protobuf:"fixed64,2,opt,name=balance" json:"balance,omitempty"`
	Assets  map[string]float64 `protobuf:"bytes,3,rep,name=assets" json:"assets,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"fixed64,2,opt,name=value"`
	Subsidy float64            `protobuf:"fixed64,4,opt,name=subsidy" json:"subsidy,omitempty"`
	Blocked bool               `protobuf:"varint,5,opt,name=blocked" json:"blocked,omitempty"`
}

func (m *Member) Reset()                    { *m = Member{} }
func (m *Member) String() string            { return proto.CompactTextString(m) }
func (*Member) ProtoMessage()               {}
func (*Member) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

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

func (m *Member) GetSubsidy() float64 {
	if m != nil {
		return m.Subsidy
	}
	return 0
}

func (m *Member) GetBlocked() bool {
	if m != nil {
		return m.Blocked
	}
	return false
}

type Transfer struct {
	Target string                     `protobuf:"bytes,1,opt,name=target" json:"target,omitempty"`
	Amount float64                    `protobuf:"fixed64,2,opt,name=amount" json:"amount,omitempty"`
	Time   *google_protobuf.Timestamp `protobuf:"bytes,3,opt,name=time" json:"time,omitempty"`
}

func (m *Transfer) Reset()                    { *m = Transfer{} }
func (m *Transfer) String() string            { return proto.CompactTextString(m) }
func (*Transfer) ProtoMessage()               {}
func (*Transfer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Transfer) GetTarget() string {
	if m != nil {
		return m.Target
	}
	return ""
}

func (m *Transfer) GetAmount() float64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *Transfer) GetTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

func init() {
	proto.RegisterType((*Member)(nil), "member.Member")
	proto.RegisterType((*Transfer)(nil), "member.Transfer")
}

func init() { proto.RegisterFile("protos/member/member.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 269 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xbf, 0x4e, 0xc3, 0x30,
	0x10, 0xc6, 0xe5, 0xa6, 0x0d, 0xe5, 0xba, 0x20, 0x0b, 0x21, 0x2b, 0x0b, 0x51, 0xa7, 0x4c, 0x8e,
	0x54, 0x16, 0x60, 0x63, 0x60, 0x64, 0x89, 0xfa, 0x02, 0x4e, 0x73, 0x8d, 0x4a, 0x93, 0xb8, 0xf2,
	0x39, 0x48, 0x79, 0x5a, 0x5e, 0x05, 0xf9, 0x9f, 0xc4, 0x64, 0xff, 0xee, 0xbe, 0xfb, 0x7c, 0xfe,
	0xa0, 0xb8, 0x19, 0x6d, 0x35, 0xd5, 0x23, 0x8e, 0x2d, 0x9a, 0x78, 0x48, 0x5f, 0xe4, 0x79, 0xa0,
	0xe2, 0xb9, 0xd7, 0xba, 0x1f, 0xb0, 0xf6, 0xd5, 0x76, 0x3e, 0xd7, 0xf6, 0x32, 0x22, 0x59, 0x35,
	0xde, 0x82, 0x70, 0xff, 0xcb, 0x20, 0xff, 0xf2, 0x5a, 0x2e, 0xe0, 0x4e, 0x75, 0x9d, 0x41, 0x22,
	0xc1, 0x4a, 0x56, 0xdd, 0x37, 0x09, 0x5d, 0xa7, 0x55, 0x83, 0x9a, 0x4e, 0x28, 0x56, 0x25, 0xab,
	0x58, 0x93, 0x90, 0x1f, 0x20, 0x57, 0x44, 0x68, 0x49, 0x64, 0x65, 0x56, 0xed, 0x0e, 0x85, 0x8c,
	0x6b, 0x04, 0x4f, 0xf9, 0xe1, 0x9b, 0x9f, 0x93, 0x35, 0x4b, 0x13, 0x95, 0xce, 0x8d, 0xe6, 0x96,
	0x2e, 0xdd, 0x22, 0xd6, 0xc1, 0x2d, 0xa2, 0x7f, 0x67, 0xd0, 0xa7, 0x2b, 0x76, 0x62, 0x53, 0xb2,
	0x6a, 0xdb, 0x24, 0x2c, 0xde, 0x60, 0xf7, 0xcf, 0x8a, 0x3f, 0x40, 0x76, 0xc5, 0x25, 0xae, 0xe9,
	0xae, 0xfc, 0x11, 0x36, 0x3f, 0x6a, 0x98, 0xd3, 0x82, 0x01, 0xde, 0x57, 0xaf, 0x6c, 0xff, 0x0d,
	0xdb, 0xa3, 0x51, 0x13, 0x9d, 0xd1, 0xf0, 0x27, 0xc8, 0xad, 0x32, 0x3d, 0xda, 0x38, 0x1a, 0xc9,
	0xd5, 0xd5, 0xa8, 0xe7, 0xc9, 0xc6, 0xf1, 0x48, 0x5c, 0xc2, 0xda, 0x05, 0x26, 0xb2, 0x92, 0xf9,
	0xcf, 0x85, 0x34, 0x65, 0x4a, 0x53, 0x1e, 0x53, 0x9a, 0x8d, 0xd7, 0xb5, 0xb9, 0xef, 0xbc, 0xfc,
	0x05, 0x00, 0x00, 0xff, 0xff, 0xaf, 0x68, 0x4c, 0x03, 0x9b, 0x01, 0x00, 0x00,
}
