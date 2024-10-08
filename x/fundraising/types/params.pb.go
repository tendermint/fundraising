// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fundraising/fundraising/v1/params.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
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

// Params defines the parameters for the module.
type Params struct {
	// auction_creation_fee specifies the fee for auction creation.
	// this prevents from spamming attack and it is collected in the community
	// pool
	AuctionCreationFee github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=auction_creation_fee,json=auctionCreationFee,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"auction_creation_fee"`
	// place_bid_fee specifies the fee for placing a bid for an auction.
	// this prevents from spamming attack and it is collected in the community
	// pool
	PlaceBidFee github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,2,rep,name=place_bid_fee,json=placeBidFee,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"place_bid_fee"`
	// extended_period specifies the extended period that determines how long
	// the extended auction round lasts
	ExtendedPeriod uint32 `protobuf:"varint,3,opt,name=extended_period,json=extendedPeriod,proto3" json:"extended_period,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ee333e6a32caa0f, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetAuctionCreationFee() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.AuctionCreationFee
	}
	return nil
}

func (m *Params) GetPlaceBidFee() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.PlaceBidFee
	}
	return nil
}

func (m *Params) GetExtendedPeriod() uint32 {
	if m != nil {
		return m.ExtendedPeriod
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "fundraising.fundraising.v1.Params")
}

func init() {
	proto.RegisterFile("fundraising/fundraising/v1/params.proto", fileDescriptor_3ee333e6a32caa0f)
}

var fileDescriptor_3ee333e6a32caa0f = []byte{
	// 370 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x52, 0x31, 0x4f, 0xc2, 0x40,
	0x18, 0xed, 0x41, 0xc2, 0x50, 0x44, 0x63, 0xc3, 0x00, 0x0c, 0x85, 0xb8, 0x80, 0x24, 0xf6, 0x52,
	0x8d, 0x8b, 0x23, 0x24, 0xac, 0x12, 0x46, 0x97, 0xe6, 0xda, 0x7e, 0xd4, 0x8b, 0xf4, 0xae, 0xe9,
	0x15, 0x02, 0x3f, 0xc0, 0xc5, 0xc9, 0xc4, 0xcd, 0xc9, 0xd1, 0x38, 0xf1, 0x33, 0x18, 0x19, 0x9d,
	0xd4, 0xc0, 0x80, 0xbf, 0xc1, 0xc9, 0xf4, 0x7a, 0x26, 0x75, 0x70, 0x75, 0x69, 0xdf, 0xf7, 0xde,
	0xdd, 0xf7, 0xde, 0xdd, 0x77, 0x7a, 0x7b, 0x3c, 0x65, 0x7e, 0x4c, 0xa8, 0xa0, 0x2c, 0xc0, 0x79,
	0x3c, 0xb3, 0x71, 0x44, 0x62, 0x12, 0x0a, 0x2b, 0x8a, 0x79, 0xc2, 0x8d, 0x46, 0x4e, 0xb4, 0xf2,
	0x78, 0x66, 0x37, 0x0e, 0x49, 0x48, 0x19, 0xc7, 0xf2, 0x9b, 0x2d, 0x6f, 0x98, 0x1e, 0x17, 0x21,
	0x17, 0xd8, 0x25, 0x02, 0xf0, 0xcc, 0x76, 0x21, 0x21, 0x36, 0xf6, 0x38, 0x65, 0x4a, 0xaf, 0x67,
	0xba, 0x23, 0x2b, 0x9c, 0x15, 0x4a, 0xaa, 0x06, 0x3c, 0xe0, 0x19, 0x9f, 0xa2, 0x8c, 0x3d, 0xfa,
	0x2a, 0xe8, 0xa5, 0xa1, 0x0c, 0x64, 0x3c, 0x20, 0xbd, 0x4a, 0xa6, 0x5e, 0x42, 0x39, 0x73, 0xbc,
	0x18, 0x88, 0x04, 0x63, 0x80, 0x1a, 0x6a, 0x15, 0x3b, 0xe5, 0xd3, 0xba, 0xa5, 0xda, 0xa5, 0xde,
	0x96, 0xf2, 0xb6, 0xfa, 0x9c, 0xb2, 0xde, 0x60, 0xf5, 0xd6, 0xd4, 0x5e, 0xde, 0x9b, 0x9d, 0x80,
	0x26, 0xd7, 0x53, 0xd7, 0xf2, 0x78, 0xa8, 0xbc, 0xd5, 0xef, 0x44, 0xf8, 0x37, 0x38, 0x59, 0x44,
	0x20, 0xe4, 0x06, 0xf1, 0xb8, 0x5b, 0x76, 0xf7, 0x26, 0x10, 0x10, 0x6f, 0xe1, 0xa4, 0xe9, 0xc5,
	0xf3, 0x6e, 0xd9, 0x45, 0x23, 0x43, 0xd9, 0xf7, 0x95, 0xfb, 0x00, 0xc0, 0xb8, 0x45, 0x7a, 0x25,
	0x9a, 0x10, 0x0f, 0x1c, 0x97, 0xfa, 0x32, 0x4e, 0xe1, 0xbf, 0xe2, 0x94, 0xa5, 0x6f, 0x8f, 0xfa,
	0x69, 0x8e, 0xb6, 0x7e, 0x00, 0xf3, 0x04, 0x98, 0x0f, 0xbe, 0x13, 0x41, 0x4c, 0xb9, 0x5f, 0x2b,
	0xb6, 0x50, 0xa7, 0x32, 0xda, 0xff, 0xa1, 0x87, 0x92, 0xbd, 0x38, 0xfe, 0x7c, 0x6a, 0xa2, 0xbb,
	0xdd, 0xb2, 0xdb, 0xca, 0xcf, 0x7d, 0xfe, 0xeb, 0x15, 0x64, 0x37, 0xde, 0xbb, 0x5c, 0x6d, 0x4c,
	0xb4, 0xde, 0x98, 0xe8, 0x63, 0x63, 0xa2, 0xfb, 0xad, 0xa9, 0xad, 0xb7, 0xa6, 0xf6, 0xba, 0x35,
	0xb5, 0xab, 0xf3, 0x5c, 0x74, 0xd9, 0x3d, 0x0e, 0x29, 0x4b, 0xf0, 0xdf, 0x1d, 0xe5, 0x69, 0xdc,
	0x92, 0x1c, 0xea, 0xd9, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x1a, 0x41, 0xef, 0x15, 0x7f, 0x02,
	0x00, 0x00,
}

func (this *Params) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Params)
	if !ok {
		that2, ok := that.(Params)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.AuctionCreationFee) != len(that1.AuctionCreationFee) {
		return false
	}
	for i := range this.AuctionCreationFee {
		if !this.AuctionCreationFee[i].Equal(&that1.AuctionCreationFee[i]) {
			return false
		}
	}
	if len(this.PlaceBidFee) != len(that1.PlaceBidFee) {
		return false
	}
	for i := range this.PlaceBidFee {
		if !this.PlaceBidFee[i].Equal(&that1.PlaceBidFee[i]) {
			return false
		}
	}
	if this.ExtendedPeriod != that1.ExtendedPeriod {
		return false
	}
	return true
}
func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ExtendedPeriod != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.ExtendedPeriod))
		i--
		dAtA[i] = 0x18
	}
	if len(m.PlaceBidFee) > 0 {
		for iNdEx := len(m.PlaceBidFee) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PlaceBidFee[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.AuctionCreationFee) > 0 {
		for iNdEx := len(m.AuctionCreationFee) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AuctionCreationFee[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.AuctionCreationFee) > 0 {
		for _, e := range m.AuctionCreationFee {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if len(m.PlaceBidFee) > 0 {
		for _, e := range m.PlaceBidFee {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if m.ExtendedPeriod != 0 {
		n += 1 + sovParams(uint64(m.ExtendedPeriod))
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuctionCreationFee", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AuctionCreationFee = append(m.AuctionCreationFee, types.Coin{})
			if err := m.AuctionCreationFee[len(m.AuctionCreationFee)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PlaceBidFee", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PlaceBidFee = append(m.PlaceBidFee, types.Coin{})
			if err := m.PlaceBidFee[len(m.PlaceBidFee)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExtendedPeriod", wireType)
			}
			m.ExtendedPeriod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExtendedPeriod |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
