// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protos/market/data.proto

/*
Package market is a generated protocol buffer package.

It is generated from these files:
	protos/market/data.proto

It has these top-level messages:
	Outcome
	Event
	Events
	Market
	Markets
	Asset
	Assets
*/
package market

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

// Event
type Outcome struct {
	Id    string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Title string `protobuf:"bytes,2,opt,name=title" json:"title,omitempty"`
}

func (m *Outcome) Reset()                    { *m = Outcome{} }
func (m *Outcome) String() string            { return proto.CompactTextString(m) }
func (*Outcome) ProtoMessage()               {}
func (*Outcome) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Outcome) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Outcome) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

type Event struct {
	Id       string                     `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	User     string                     `protobuf:"bytes,2,opt,name=user" json:"user,omitempty"`
	Title    string                     `protobuf:"bytes,3,opt,name=title" json:"title,omitempty"`
	Outcomes []*Outcome                 `protobuf:"bytes,4,rep,name=outcomes" json:"outcomes,omitempty"`
	Result   string                     `protobuf:"bytes,5,opt,name=result" json:"result,omitempty"`
	Approved bool                       `protobuf:"varint,6,opt,name=approved" json:"approved,omitempty"`
	End      *google_protobuf.Timestamp `protobuf:"bytes,7,opt,name=end" json:"end,omitempty"`
	Allowed  bool                       `protobuf:"varint,8,opt,name=allowed" json:"allowed,omitempty"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Event) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Event) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *Event) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Event) GetOutcomes() []*Outcome {
	if m != nil {
		return m.Outcomes
	}
	return nil
}

func (m *Event) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

func (m *Event) GetApproved() bool {
	if m != nil {
		return m.Approved
	}
	return false
}

func (m *Event) GetEnd() *google_protobuf.Timestamp {
	if m != nil {
		return m.End
	}
	return nil
}

func (m *Event) GetAllowed() bool {
	if m != nil {
		return m.Allowed
	}
	return false
}

type Events struct {
	List []*Event `protobuf:"bytes,1,rep,name=list" json:"list,omitempty"`
}

func (m *Events) Reset()                    { *m = Events{} }
func (m *Events) String() string            { return proto.CompactTextString(m) }
func (*Events) ProtoMessage()               {}
func (*Events) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Events) GetList() []*Event {
	if m != nil {
		return m.List
	}
	return nil
}

// Market
type Market struct {
	Id        string                     `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	User      string                     `protobuf:"bytes,2,opt,name=user" json:"user,omitempty"`
	Event     string                     `protobuf:"bytes,3,opt,name=event" json:"event,omitempty"`
	Liquidity float64                    `protobuf:"fixed64,4,opt,name=liquidity" json:"liquidity,omitempty"`
	Fund      float64                    `protobuf:"fixed64,5,opt,name=fund" json:"fund,omitempty"`
	Cost      float64                    `protobuf:"fixed64,6,opt,name=cost" json:"cost,omitempty"`
	Shares    map[string]float64         `protobuf:"bytes,7,rep,name=shares" json:"shares,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"fixed64,2,opt,name=value"`
	Prices    map[string]float64         `protobuf:"bytes,8,rep,name=prices" json:"prices,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"fixed64,2,opt,name=value"`
	Settled   bool                       `protobuf:"varint,9,opt,name=settled" json:"settled,omitempty"`
	Start     *google_protobuf.Timestamp `protobuf:"bytes,10,opt,name=start" json:"start,omitempty"`
	End       *google_protobuf.Timestamp `protobuf:"bytes,11,opt,name=end" json:"end,omitempty"`
	Allowed   bool                       `protobuf:"varint,12,opt,name=allowed" json:"allowed,omitempty"`
}

func (m *Market) Reset()                    { *m = Market{} }
func (m *Market) String() string            { return proto.CompactTextString(m) }
func (*Market) ProtoMessage()               {}
func (*Market) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Market) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Market) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *Market) GetEvent() string {
	if m != nil {
		return m.Event
	}
	return ""
}

func (m *Market) GetLiquidity() float64 {
	if m != nil {
		return m.Liquidity
	}
	return 0
}

func (m *Market) GetFund() float64 {
	if m != nil {
		return m.Fund
	}
	return 0
}

func (m *Market) GetCost() float64 {
	if m != nil {
		return m.Cost
	}
	return 0
}

func (m *Market) GetShares() map[string]float64 {
	if m != nil {
		return m.Shares
	}
	return nil
}

func (m *Market) GetPrices() map[string]float64 {
	if m != nil {
		return m.Prices
	}
	return nil
}

func (m *Market) GetSettled() bool {
	if m != nil {
		return m.Settled
	}
	return false
}

func (m *Market) GetStart() *google_protobuf.Timestamp {
	if m != nil {
		return m.Start
	}
	return nil
}

func (m *Market) GetEnd() *google_protobuf.Timestamp {
	if m != nil {
		return m.End
	}
	return nil
}

func (m *Market) GetAllowed() bool {
	if m != nil {
		return m.Allowed
	}
	return false
}

type Markets struct {
	List []*Market `protobuf:"bytes,1,rep,name=list" json:"list,omitempty"`
}

func (m *Markets) Reset()                    { *m = Markets{} }
func (m *Markets) String() string            { return proto.CompactTextString(m) }
func (*Markets) ProtoMessage()               {}
func (*Markets) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Markets) GetList() []*Market {
	if m != nil {
		return m.List
	}
	return nil
}

type Asset struct {
	Id     string  `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Volume float64 `protobuf:"fixed64,2,opt,name=volume" json:"volume,omitempty"`
}

func (m *Asset) Reset()                    { *m = Asset{} }
func (m *Asset) String() string            { return proto.CompactTextString(m) }
func (*Asset) ProtoMessage()               {}
func (*Asset) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Asset) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Asset) GetVolume() float64 {
	if m != nil {
		return m.Volume
	}
	return 0
}

type Assets struct {
	List []*Asset `protobuf:"bytes,1,rep,name=list" json:"list,omitempty"`
}

func (m *Assets) Reset()                    { *m = Assets{} }
func (m *Assets) String() string            { return proto.CompactTextString(m) }
func (*Assets) ProtoMessage()               {}
func (*Assets) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Assets) GetList() []*Asset {
	if m != nil {
		return m.List
	}
	return nil
}

func init() {
	proto.RegisterType((*Outcome)(nil), "market.Outcome")
	proto.RegisterType((*Event)(nil), "market.Event")
	proto.RegisterType((*Events)(nil), "market.Events")
	proto.RegisterType((*Market)(nil), "market.Market")
	proto.RegisterType((*Markets)(nil), "market.Markets")
	proto.RegisterType((*Asset)(nil), "market.Asset")
	proto.RegisterType((*Assets)(nil), "market.Assets")
}

func init() { proto.RegisterFile("protos/market/data.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 479 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xd5, 0x26, 0xf1, 0x47, 0x26, 0x50, 0xd0, 0x0a, 0x55, 0x2b, 0x0b, 0x09, 0xe3, 0x93, 0xa5,
	0x82, 0x8d, 0xc2, 0x05, 0xb8, 0x71, 0xe8, 0x11, 0x81, 0x16, 0xfe, 0x80, 0x1b, 0x6f, 0x8b, 0xd5,
	0x75, 0xd6, 0x78, 0xc7, 0x41, 0xf9, 0xcb, 0xfc, 0x04, 0x4e, 0x68, 0x67, 0xed, 0x36, 0x11, 0x45,
	0xb4, 0xb7, 0x79, 0x33, 0x6f, 0x66, 0xdf, 0xbe, 0x19, 0x10, 0x5d, 0x6f, 0xd0, 0xd8, 0xb2, 0xad,
	0xfa, 0x6b, 0x85, 0x65, 0x5d, 0x61, 0x55, 0x50, 0x8a, 0x87, 0x3e, 0x95, 0xbc, 0xb8, 0x32, 0xe6,
	0x4a, 0xab, 0x92, 0xb2, 0x17, 0xc3, 0x65, 0x89, 0x4d, 0xab, 0x2c, 0x56, 0x6d, 0xe7, 0x89, 0x59,
	0x09, 0xd1, 0xe7, 0x01, 0x37, 0xa6, 0x55, 0xfc, 0x04, 0x66, 0x4d, 0x2d, 0x58, 0xca, 0xf2, 0xa5,
	0x9c, 0x35, 0x35, 0x7f, 0x06, 0x01, 0x36, 0xa8, 0x95, 0x98, 0x51, 0xca, 0x83, 0xec, 0x17, 0x83,
	0xe0, 0x7c, 0xa7, 0xb6, 0xf8, 0x17, 0x9f, 0xc3, 0x62, 0xb0, 0xaa, 0x1f, 0xe9, 0x14, 0xdf, 0xce,
	0x98, 0x1f, 0xcc, 0xe0, 0x67, 0x10, 0x1b, 0xff, 0xa8, 0x15, 0x8b, 0x74, 0x9e, 0xaf, 0xd6, 0x4f,
	0x0a, 0x2f, 0xb8, 0x18, 0xc5, 0xc8, 0x1b, 0x02, 0x3f, 0x85, 0xb0, 0x57, 0x76, 0xd0, 0x28, 0x02,
	0x9a, 0x31, 0x22, 0x9e, 0x40, 0x5c, 0x75, 0x5d, 0x6f, 0x76, 0xaa, 0x16, 0x61, 0xca, 0xf2, 0x58,
	0xde, 0x60, 0xfe, 0x0a, 0xe6, 0x6a, 0x5b, 0x8b, 0x28, 0x65, 0xf9, 0x6a, 0x9d, 0x14, 0xde, 0x84,
	0x62, 0x32, 0xa1, 0xf8, 0x36, 0x99, 0x20, 0x1d, 0x8d, 0x0b, 0x88, 0x2a, 0xad, 0xcd, 0x4f, 0x55,
	0x8b, 0x98, 0x06, 0x4d, 0x30, 0x3b, 0x83, 0x90, 0xfe, 0x6a, 0xf9, 0x4b, 0x58, 0xe8, 0xc6, 0xa2,
	0x60, 0x24, 0xf7, 0xf1, 0x24, 0x97, 0xaa, 0x92, 0x4a, 0xd9, 0xef, 0x39, 0x84, 0x9f, 0x28, 0x7d,
	0x5f, 0x6b, 0x94, 0xeb, 0x9e, 0xac, 0x21, 0xc0, 0x9f, 0xc3, 0x52, 0x37, 0x3f, 0x86, 0xa6, 0x6e,
	0x70, 0x2f, 0x16, 0x29, 0xcb, 0x99, 0xbc, 0x4d, 0xb8, 0x39, 0x97, 0xc3, 0xb6, 0x26, 0x27, 0x98,
	0xa4, 0xd8, 0xe5, 0x36, 0xc6, 0x22, 0x79, 0xc0, 0x24, 0xc5, 0x7c, 0x0d, 0xa1, 0xfd, 0x5e, 0xf5,
	0xca, 0x8a, 0x88, 0xf4, 0x26, 0x93, 0x5e, 0xaf, 0xaf, 0xf8, 0x4a, 0xc5, 0xf3, 0x2d, 0xf6, 0x7b,
	0x39, 0x32, 0x5d, 0x4f, 0xd7, 0x37, 0x1b, 0x65, 0x45, 0x7c, 0x67, 0xcf, 0x17, 0x2a, 0x8e, 0x3d,
	0x9e, 0xe9, 0x9c, 0xb3, 0x0a, 0x51, 0xab, 0x5a, 0x2c, 0xbd, 0x73, 0x23, 0xe4, 0x6f, 0x20, 0xb0,
	0x58, 0xf5, 0x28, 0xe0, 0xbf, 0x3b, 0xf0, 0xc4, 0x69, 0x67, 0xab, 0x07, 0xef, 0xec, 0xd1, 0xd1,
	0xce, 0x92, 0xf7, 0xb0, 0x3a, 0xf8, 0x1e, 0x7f, 0x0a, 0xf3, 0x6b, 0xb5, 0x1f, 0x77, 0xe1, 0x42,
	0x67, 0xfc, 0xae, 0xd2, 0x83, 0xbf, 0x6b, 0x26, 0x3d, 0xf8, 0x30, 0x7b, 0xc7, 0x5c, 0xeb, 0xc1,
	0x2f, 0x1f, 0xd2, 0x9a, 0xbd, 0x86, 0xc8, 0xfb, 0x64, 0x79, 0x76, 0x74, 0x2a, 0x27, 0xc7, 0x36,
	0x8e, 0xb7, 0x52, 0x42, 0xf0, 0xd1, 0xda, 0x3b, 0x2e, 0xe5, 0x14, 0xc2, 0x9d, 0xd1, 0x43, 0x3b,
	0x3d, 0x31, 0x22, 0x77, 0x89, 0xd4, 0xf0, 0xcf, 0x4b, 0xa4, 0xaa, 0x9f, 0x7e, 0x11, 0x92, 0x6b,
	0x6f, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0xec, 0x49, 0x28, 0x64, 0x20, 0x04, 0x00, 0x00,
}
