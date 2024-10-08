// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fundraising/fundraising/v1/bid.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// BidType enumerates the valid types of a bid.
type BidType int32

const (
	// BID_TYPE_UNSPECIFIED defines the default bid type
	BidTypeNil BidType = 0
	// BID_TYPE_FIXED_PRICE defines a bid type for a fixed price auction type
	BidTypeFixedPrice BidType = 1
	// BID_TYPE_BATCH_WORTH defines a bid type for How-Much-Worth-to-Buy of a
	// batch auction
	BidTypeBatchWorth BidType = 2
	// BID_TYPE_BATCH_MANY defines a bid type for How-Many-Coins-to-Buy of a batch
	// auction
	BidTypeBatchMany BidType = 3
)

var BidType_name = map[int32]string{
	0: "BID_TYPE_UNSPECIFIED",
	1: "BID_TYPE_FIXED_PRICE",
	2: "BID_TYPE_BATCH_WORTH",
	3: "BID_TYPE_BATCH_MANY",
}

var BidType_value = map[string]int32{
	"BID_TYPE_UNSPECIFIED": 0,
	"BID_TYPE_FIXED_PRICE": 1,
	"BID_TYPE_BATCH_WORTH": 2,
	"BID_TYPE_BATCH_MANY":  3,
}

func (x BidType) String() string {
	return proto.EnumName(BidType_name, int32(x))
}

func (BidType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_527f99309e242c82, []int{0}
}

// AddressType enumerates the available types of a address.
type AddressType int32

const (
	// the 32 bytes length address type of ADR 028.
	AddressType32Bytes AddressType = 0
	// the default 20 bytes length address type.
	AddressType20Bytes AddressType = 1
)

var AddressType_name = map[int32]string{
	0: "ADDRESS_TYPE_32_BYTES",
	1: "ADDRESS_TYPE_20_BYTES",
}

var AddressType_value = map[string]int32{
	"ADDRESS_TYPE_32_BYTES": 0,
	"ADDRESS_TYPE_20_BYTES": 1,
}

func (x AddressType) String() string {
	return proto.EnumName(AddressType_name, int32(x))
}

func (AddressType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_527f99309e242c82, []int{1}
}

// Bid defines a standard bid for an auction.
type Bid struct {
	// auction_id specifies the id of the auction
	AuctionId uint64 `protobuf:"varint,1,opt,name=auction_id,json=auctionId,proto3" json:"auction_id,omitempty"`
	// bidder specifies the bech32-encoded address that bids for the auction
	Bidder string `protobuf:"bytes,2,opt,name=bidder,proto3" json:"bidder,omitempty"`
	// id specifies an index of a bid for the bidder
	Id uint64 `protobuf:"varint,3,opt,name=id,proto3" json:"id,omitempty"`
	// type specifies the bid type; type 1 is fixed price, 2 is how-much-worth, 3
	// is how-many-coins
	Type BidType `protobuf:"varint,4,opt,name=type,proto3,enum=fundraising.fundraising.v1.BidType" json:"type,omitempty"`
	// price specifies the bid price in which price the bidder places the bid
	Price cosmossdk_io_math.LegacyDec `protobuf:"bytes,5,opt,name=price,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"price"`
	// coin specifies the amount of coin that the bidder bids
	// for a fixed price auction, the denom is of the paying coin.
	// for a batch auction of how-much-worth, the denom is of the paying coin.
	// for a batch auction of how-many-coins, the denom is of the selling coin.
	Coin types.Coin `protobuf:"bytes,6,opt,name=coin,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coin" json:"coin"`
	// is_matched specifies the bid that is a winning bid and enables the bidder
	// to purchase the selling coin
	IsMatched bool `protobuf:"varint,7,opt,name=is_matched,json=isMatched,proto3" json:"is_matched,omitempty"`
}

func (m *Bid) Reset()         { *m = Bid{} }
func (m *Bid) String() string { return proto.CompactTextString(m) }
func (*Bid) ProtoMessage()    {}
func (*Bid) Descriptor() ([]byte, []int) {
	return fileDescriptor_527f99309e242c82, []int{0}
}
func (m *Bid) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Bid) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Bid.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Bid) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bid.Merge(m, src)
}
func (m *Bid) XXX_Size() int {
	return m.Size()
}
func (m *Bid) XXX_DiscardUnknown() {
	xxx_messageInfo_Bid.DiscardUnknown(m)
}

var xxx_messageInfo_Bid proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("fundraising.fundraising.v1.BidType", BidType_name, BidType_value)
	proto.RegisterEnum("fundraising.fundraising.v1.AddressType", AddressType_name, AddressType_value)
	proto.RegisterType((*Bid)(nil), "fundraising.fundraising.v1.Bid")
}

func init() {
	proto.RegisterFile("fundraising/fundraising/v1/bid.proto", fileDescriptor_527f99309e242c82)
}

var fileDescriptor_527f99309e242c82 = []byte{
	// 622 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x93, 0xc1, 0x4b, 0xdc, 0x4e,
	0x14, 0xc7, 0x33, 0xeb, 0xaa, 0x3f, 0xc7, 0x1f, 0xb2, 0xa6, 0x2a, 0x31, 0xa5, 0xd9, 0xd0, 0x16,
	0xba, 0x2c, 0x98, 0xb8, 0x2b, 0xa5, 0xd0, 0xdb, 0x66, 0x13, 0xeb, 0x42, 0xd5, 0x25, 0xbb, 0xc5,
	0xda, 0x4b, 0xc8, 0x66, 0xa6, 0xbb, 0x83, 0x6e, 0x66, 0xc9, 0x44, 0x31, 0xff, 0x81, 0xec, 0xa9,
	0xe7, 0xc2, 0x42, 0xa1, 0x97, 0xd2, 0x93, 0x87, 0xfe, 0x11, 0xd2, 0x93, 0xf4, 0x54, 0x7a, 0xb0,
	0x45, 0x0f, 0xd2, 0xff, 0xa2, 0x64, 0x32, 0x87, 0x58, 0xe9, 0x25, 0x99, 0xf7, 0xde, 0xf7, 0x33,
	0x6f, 0xde, 0x77, 0x18, 0xf8, 0xf8, 0xed, 0x51, 0x88, 0x22, 0x9f, 0x30, 0x12, 0xf6, 0xcd, 0xfc,
	0xfa, 0xb8, 0x66, 0xf6, 0x08, 0x32, 0x46, 0x11, 0x8d, 0xa9, 0xac, 0xe6, 0x2a, 0x46, 0x7e, 0x7d,
	0x5c, 0x53, 0x17, 0xfd, 0x21, 0x09, 0xa9, 0xc9, 0xbf, 0x99, 0x5c, 0xd5, 0x02, 0xca, 0x86, 0x94,
	0x99, 0x3d, 0x9f, 0x61, 0xf3, 0xb8, 0xd6, 0xc3, 0xb1, 0x5f, 0x33, 0x03, 0x4a, 0x42, 0x51, 0x5f,
	0xcd, 0xea, 0x1e, 0x8f, 0xcc, 0x2c, 0x10, 0xa5, 0xa5, 0x3e, 0xed, 0xd3, 0x2c, 0x9f, 0xae, 0xb2,
	0xec, 0xc3, 0xdf, 0x05, 0x38, 0x65, 0x11, 0x24, 0x3f, 0x80, 0xd0, 0x3f, 0x0a, 0x62, 0x42, 0x43,
	0x8f, 0x20, 0x05, 0xe8, 0xa0, 0x52, 0x74, 0xe7, 0x44, 0xa6, 0x85, 0xe4, 0x15, 0x38, 0xd3, 0x23,
	0x08, 0xe1, 0x48, 0x29, 0xe8, 0xa0, 0x32, 0xe7, 0x8a, 0x48, 0x5e, 0x80, 0x05, 0x82, 0x94, 0x29,
	0x2e, 0x2f, 0x10, 0x24, 0x3f, 0x83, 0xc5, 0x38, 0x19, 0x61, 0xa5, 0xa8, 0x83, 0xca, 0x42, 0xfd,
	0x91, 0xf1, 0xef, 0xe9, 0x0c, 0x8b, 0xa0, 0x6e, 0x32, 0xc2, 0x2e, 0x07, 0xe4, 0x17, 0x70, 0x7a,
	0x14, 0x91, 0x00, 0x2b, 0xd3, 0xe9, 0xfe, 0x56, 0xed, 0xfc, 0xb2, 0x2c, 0xfd, 0xb8, 0x2c, 0xdf,
	0xcf, 0x46, 0x60, 0xe8, 0xc0, 0x20, 0xd4, 0x1c, 0xfa, 0xf1, 0xc0, 0x78, 0x89, 0xfb, 0x7e, 0x90,
	0xd8, 0x38, 0xf8, 0xf6, 0x65, 0x0d, 0x8a, 0x09, 0x6d, 0x1c, 0xb8, 0x19, 0x2f, 0xc7, 0xb0, 0x98,
	0xfa, 0xa1, 0xcc, 0xe8, 0xa0, 0x32, 0x5f, 0x5f, 0x35, 0x84, 0x22, 0x35, 0xcc, 0x10, 0x86, 0x19,
	0x4d, 0x4a, 0x42, 0xcb, 0x49, 0x5b, 0x7c, 0xfe, 0x59, 0x7e, 0xd2, 0x27, 0xf1, 0xe0, 0xa8, 0x67,
	0x04, 0x74, 0x28, 0x0c, 0x13, 0xbf, 0x35, 0x86, 0x0e, 0xcc, 0xf4, 0x70, 0x8c, 0x03, 0xef, 0x6f,
	0xce, 0xaa, 0xff, 0x1f, 0xf2, 0xe6, 0x5e, 0xda, 0x81, 0x7d, 0xba, 0x39, 0xab, 0x02, 0x97, 0x77,
	0x4b, 0xed, 0x23, 0xcc, 0x1b, 0xfa, 0x71, 0x30, 0xc0, 0x48, 0x99, 0xd5, 0x41, 0xe5, 0x3f, 0x77,
	0x8e, 0xb0, 0xed, 0x2c, 0xf1, 0xbc, 0x78, 0xfa, 0xa1, 0x2c, 0x55, 0xbf, 0x02, 0x38, 0x2b, 0xa6,
	0x96, 0x2b, 0x70, 0xc9, 0x6a, 0xd9, 0x5e, 0x77, 0xbf, 0xed, 0x78, 0xaf, 0x76, 0x3a, 0x6d, 0xa7,
	0xd9, 0xda, 0x6c, 0x39, 0x76, 0x49, 0x52, 0x17, 0xc6, 0x13, 0x1d, 0x0a, 0xd9, 0x0e, 0x39, 0x94,
	0xcd, 0x9c, 0x72, 0xb3, 0xf5, 0xda, 0xb1, 0xbd, 0xb6, 0xdb, 0x6a, 0x3a, 0x25, 0xa0, 0x2e, 0x8f,
	0x27, 0xfa, 0xa2, 0x50, 0x6e, 0x92, 0x13, 0x8c, 0xda, 0xdc, 0x81, 0x3c, 0x60, 0x35, 0xba, 0xcd,
	0x2d, 0x6f, 0x6f, 0xd7, 0xed, 0x6e, 0x95, 0x0a, 0xb7, 0x00, 0x2b, 0x3d, 0xda, 0x1e, 0x8d, 0xe2,
	0x81, 0xbc, 0x06, 0xef, 0xfd, 0x05, 0x6c, 0x37, 0x76, 0xf6, 0x4b, 0x53, 0xea, 0xd2, 0x78, 0xa2,
	0x97, 0xf2, 0xfa, 0x6d, 0x3f, 0x4c, 0xd4, 0xe2, 0xe9, 0x47, 0x4d, 0xaa, 0x26, 0x70, 0xbe, 0x81,
	0x50, 0x84, 0x19, 0xe3, 0xf3, 0xd4, 0xe0, 0x72, 0xc3, 0xb6, 0x5d, 0xa7, 0xd3, 0xc9, 0xf6, 0xd9,
	0xa8, 0x7b, 0xd6, 0x7e, 0xd7, 0xe9, 0x94, 0x24, 0x75, 0x65, 0x3c, 0xd1, 0xe5, 0x9c, 0x76, 0xa3,
	0x6e, 0x25, 0x31, 0x66, 0x77, 0x90, 0xfa, 0xba, 0x40, 0xc0, 0x1d, 0xa4, 0xbe, 0xce, 0x91, 0xac,
	0xb5, 0xb5, 0x7b, 0x7e, 0xa5, 0x81, 0x8b, 0x2b, 0x0d, 0xfc, 0xba, 0xd2, 0xc0, 0xbb, 0x6b, 0x4d,
	0xba, 0xb8, 0xd6, 0xa4, 0xef, 0xd7, 0x9a, 0xf4, 0xe6, 0x69, 0xee, 0x2e, 0x63, 0x1c, 0x22, 0x1c,
	0x0d, 0x49, 0x18, 0xdf, 0x7a, 0x7d, 0x27, 0xb7, 0x22, 0x7e, 0xbd, 0xbd, 0x19, 0xfe, 0x16, 0x36,
	0xfe, 0x04, 0x00, 0x00, 0xff, 0xff, 0xa2, 0x56, 0x83, 0x81, 0xb3, 0x03, 0x00, 0x00,
}

func (m *Bid) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Bid) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Bid) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.IsMatched {
		i--
		if m.IsMatched {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x38
	}
	{
		size, err := m.Coin.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintBid(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.Price.Size()
		i -= size
		if _, err := m.Price.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintBid(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if m.Type != 0 {
		i = encodeVarintBid(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x20
	}
	if m.Id != 0 {
		i = encodeVarintBid(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Bidder) > 0 {
		i -= len(m.Bidder)
		copy(dAtA[i:], m.Bidder)
		i = encodeVarintBid(dAtA, i, uint64(len(m.Bidder)))
		i--
		dAtA[i] = 0x12
	}
	if m.AuctionId != 0 {
		i = encodeVarintBid(dAtA, i, uint64(m.AuctionId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintBid(dAtA []byte, offset int, v uint64) int {
	offset -= sovBid(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Bid) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AuctionId != 0 {
		n += 1 + sovBid(uint64(m.AuctionId))
	}
	l = len(m.Bidder)
	if l > 0 {
		n += 1 + l + sovBid(uint64(l))
	}
	if m.Id != 0 {
		n += 1 + sovBid(uint64(m.Id))
	}
	if m.Type != 0 {
		n += 1 + sovBid(uint64(m.Type))
	}
	l = m.Price.Size()
	n += 1 + l + sovBid(uint64(l))
	l = m.Coin.Size()
	n += 1 + l + sovBid(uint64(l))
	if m.IsMatched {
		n += 2
	}
	return n
}

func sovBid(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBid(x uint64) (n int) {
	return sovBid(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Bid) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBid
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Bid: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Bid: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuctionId", wireType)
			}
			m.AuctionId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBid
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AuctionId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bidder", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBid
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBid
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBid
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Bidder = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBid
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBid
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= BidType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBid
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBid
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBid
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Price.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Coin", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBid
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthBid
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBid
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Coin.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsMatched", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBid
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsMatched = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipBid(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBid
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipBid(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBid
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBid
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBid
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthBid
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBid
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBid
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBid        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBid          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBid = fmt.Errorf("proto: unexpected end of group")
)
