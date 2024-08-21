// Code generated by protoc-gen-go-pulsar. DO NOT EDIT.
package fundraisingv1

import (
	fmt "fmt"
	io "io"
	reflect "reflect"
	sync "sync"

	_ "cosmossdk.io/api/amino"
	v1beta1 "cosmossdk.io/api/cosmos/base/v1beta1"
	_ "github.com/cosmos/cosmos-proto"
	runtime "github.com/cosmos/cosmos-proto/runtime"
	_ "github.com/cosmos/gogoproto/gogoproto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoiface "google.golang.org/protobuf/runtime/protoiface"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

var _ protoreflect.List = (*_Params_1_list)(nil)

type _Params_1_list struct {
	list *[]*v1beta1.Coin
}

func (x *_Params_1_list) Len() int {
	if x.list == nil {
		return 0
	}
	return len(*x.list)
}

func (x *_Params_1_list) Get(i int) protoreflect.Value {
	return protoreflect.ValueOfMessage((*x.list)[i].ProtoReflect())
}

func (x *_Params_1_list) Set(i int, value protoreflect.Value) {
	valueUnwrapped := value.Message()
	concreteValue := valueUnwrapped.Interface().(*v1beta1.Coin)
	(*x.list)[i] = concreteValue
}

func (x *_Params_1_list) Append(value protoreflect.Value) {
	valueUnwrapped := value.Message()
	concreteValue := valueUnwrapped.Interface().(*v1beta1.Coin)
	*x.list = append(*x.list, concreteValue)
}

func (x *_Params_1_list) AppendMutable() protoreflect.Value {
	v := new(v1beta1.Coin)
	*x.list = append(*x.list, v)
	return protoreflect.ValueOfMessage(v.ProtoReflect())
}

func (x *_Params_1_list) Truncate(n int) {
	for i := n; i < len(*x.list); i++ {
		(*x.list)[i] = nil
	}
	*x.list = (*x.list)[:n]
}

func (x *_Params_1_list) NewElement() protoreflect.Value {
	v := new(v1beta1.Coin)
	return protoreflect.ValueOfMessage(v.ProtoReflect())
}

func (x *_Params_1_list) IsValid() bool {
	return x.list != nil
}

var _ protoreflect.List = (*_Params_2_list)(nil)

type _Params_2_list struct {
	list *[]*v1beta1.Coin
}

func (x *_Params_2_list) Len() int {
	if x.list == nil {
		return 0
	}
	return len(*x.list)
}

func (x *_Params_2_list) Get(i int) protoreflect.Value {
	return protoreflect.ValueOfMessage((*x.list)[i].ProtoReflect())
}

func (x *_Params_2_list) Set(i int, value protoreflect.Value) {
	valueUnwrapped := value.Message()
	concreteValue := valueUnwrapped.Interface().(*v1beta1.Coin)
	(*x.list)[i] = concreteValue
}

func (x *_Params_2_list) Append(value protoreflect.Value) {
	valueUnwrapped := value.Message()
	concreteValue := valueUnwrapped.Interface().(*v1beta1.Coin)
	*x.list = append(*x.list, concreteValue)
}

func (x *_Params_2_list) AppendMutable() protoreflect.Value {
	v := new(v1beta1.Coin)
	*x.list = append(*x.list, v)
	return protoreflect.ValueOfMessage(v.ProtoReflect())
}

func (x *_Params_2_list) Truncate(n int) {
	for i := n; i < len(*x.list); i++ {
		(*x.list)[i] = nil
	}
	*x.list = (*x.list)[:n]
}

func (x *_Params_2_list) NewElement() protoreflect.Value {
	v := new(v1beta1.Coin)
	return protoreflect.ValueOfMessage(v.ProtoReflect())
}

func (x *_Params_2_list) IsValid() bool {
	return x.list != nil
}

var (
	md_Params                      protoreflect.MessageDescriptor
	fd_Params_auction_creation_fee protoreflect.FieldDescriptor
	fd_Params_place_bid_fee        protoreflect.FieldDescriptor
	fd_Params_extended_period      protoreflect.FieldDescriptor
)

func init() {
	file_fundraising_fundraising_v1_params_proto_init()
	md_Params = File_fundraising_fundraising_v1_params_proto.Messages().ByName("Params")
	fd_Params_auction_creation_fee = md_Params.Fields().ByName("auction_creation_fee")
	fd_Params_place_bid_fee = md_Params.Fields().ByName("place_bid_fee")
	fd_Params_extended_period = md_Params.Fields().ByName("extended_period")
}

var _ protoreflect.Message = (*fastReflection_Params)(nil)

type fastReflection_Params Params

func (x *Params) ProtoReflect() protoreflect.Message {
	return (*fastReflection_Params)(x)
}

func (x *Params) slowProtoReflect() protoreflect.Message {
	mi := &file_fundraising_fundraising_v1_params_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

var (
	_fastReflection_Params_messageType fastReflection_Params_messageType
	_                                  protoreflect.MessageType = fastReflection_Params_messageType{}
)

type fastReflection_Params_messageType struct{}

func (x fastReflection_Params_messageType) Zero() protoreflect.Message {
	return (*fastReflection_Params)(nil)
}

func (x fastReflection_Params_messageType) New() protoreflect.Message {
	return new(fastReflection_Params)
}

func (x fastReflection_Params_messageType) Descriptor() protoreflect.MessageDescriptor {
	return md_Params
}

// Descriptor returns message descriptor, which contains only the protobuf
// type information for the message.
func (x *fastReflection_Params) Descriptor() protoreflect.MessageDescriptor {
	return md_Params
}

// Type returns the message type, which encapsulates both Go and protobuf
// type information. If the Go type information is not needed,
// it is recommended that the message descriptor be used instead.
func (x *fastReflection_Params) Type() protoreflect.MessageType {
	return _fastReflection_Params_messageType
}

// New returns a newly allocated and mutable empty message.
func (x *fastReflection_Params) New() protoreflect.Message {
	return new(fastReflection_Params)
}

// Interface unwraps the message reflection interface and
// returns the underlying ProtoMessage interface.
func (x *fastReflection_Params) Interface() protoreflect.ProtoMessage {
	return (*Params)(x)
}

// Range iterates over every populated field in an undefined order,
// calling f for each field descriptor and value encountered.
// Range returns immediately if f returns false.
// While iterating, mutating operations may only be performed
// on the current field descriptor.
func (x *fastReflection_Params) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	if len(x.AuctionCreationFee) != 0 {
		value := protoreflect.ValueOfList(&_Params_1_list{list: &x.AuctionCreationFee})
		if !f(fd_Params_auction_creation_fee, value) {
			return
		}
	}
	if len(x.PlaceBidFee) != 0 {
		value := protoreflect.ValueOfList(&_Params_2_list{list: &x.PlaceBidFee})
		if !f(fd_Params_place_bid_fee, value) {
			return
		}
	}
	if x.ExtendedPeriod != uint32(0) {
		value := protoreflect.ValueOfUint32(x.ExtendedPeriod)
		if !f(fd_Params_extended_period, value) {
			return
		}
	}
}

// Has reports whether a field is populated.
//
// Some fields have the property of nullability where it is possible to
// distinguish between the default value of a field and whether the field
// was explicitly populated with the default value. Singular message fields,
// member fields of a oneof, and proto2 scalar fields are nullable. Such
// fields are populated only if explicitly set.
//
// In other cases (aside from the nullable cases above),
// a proto3 scalar field is populated if it contains a non-zero value, and
// a repeated field is populated if it is non-empty.
func (x *fastReflection_Params) Has(fd protoreflect.FieldDescriptor) bool {
	switch fd.FullName() {
	case "fundraising.fundraising.v1.Params.auction_creation_fee":
		return len(x.AuctionCreationFee) != 0
	case "fundraising.fundraising.v1.Params.place_bid_fee":
		return len(x.PlaceBidFee) != 0
	case "fundraising.fundraising.v1.Params.extended_period":
		return x.ExtendedPeriod != uint32(0)
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: fundraising.fundraising.v1.Params"))
		}
		panic(fmt.Errorf("message fundraising.fundraising.v1.Params does not contain field %s", fd.FullName()))
	}
}

// Clear clears the field such that a subsequent Has call reports false.
//
// Clearing an extension field clears both the extension type and value
// associated with the given field number.
//
// Clear is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Params) Clear(fd protoreflect.FieldDescriptor) {
	switch fd.FullName() {
	case "fundraising.fundraising.v1.Params.auction_creation_fee":
		x.AuctionCreationFee = nil
	case "fundraising.fundraising.v1.Params.place_bid_fee":
		x.PlaceBidFee = nil
	case "fundraising.fundraising.v1.Params.extended_period":
		x.ExtendedPeriod = uint32(0)
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: fundraising.fundraising.v1.Params"))
		}
		panic(fmt.Errorf("message fundraising.fundraising.v1.Params does not contain field %s", fd.FullName()))
	}
}

// Get retrieves the value for a field.
//
// For unpopulated scalars, it returns the default value, where
// the default value of a bytes scalar is guaranteed to be a copy.
// For unpopulated composite types, it returns an empty, read-only view
// of the value; to obtain a mutable reference, use Mutable.
func (x *fastReflection_Params) Get(descriptor protoreflect.FieldDescriptor) protoreflect.Value {
	switch descriptor.FullName() {
	case "fundraising.fundraising.v1.Params.auction_creation_fee":
		if len(x.AuctionCreationFee) == 0 {
			return protoreflect.ValueOfList(&_Params_1_list{})
		}
		listValue := &_Params_1_list{list: &x.AuctionCreationFee}
		return protoreflect.ValueOfList(listValue)
	case "fundraising.fundraising.v1.Params.place_bid_fee":
		if len(x.PlaceBidFee) == 0 {
			return protoreflect.ValueOfList(&_Params_2_list{})
		}
		listValue := &_Params_2_list{list: &x.PlaceBidFee}
		return protoreflect.ValueOfList(listValue)
	case "fundraising.fundraising.v1.Params.extended_period":
		value := x.ExtendedPeriod
		return protoreflect.ValueOfUint32(value)
	default:
		if descriptor.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: fundraising.fundraising.v1.Params"))
		}
		panic(fmt.Errorf("message fundraising.fundraising.v1.Params does not contain field %s", descriptor.FullName()))
	}
}

// Set stores the value for a field.
//
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType.
// When setting a composite type, it is unspecified whether the stored value
// aliases the source's memory in any way. If the composite value is an
// empty, read-only value, then it panics.
//
// Set is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Params) Set(fd protoreflect.FieldDescriptor, value protoreflect.Value) {
	switch fd.FullName() {
	case "fundraising.fundraising.v1.Params.auction_creation_fee":
		lv := value.List()
		clv := lv.(*_Params_1_list)
		x.AuctionCreationFee = *clv.list
	case "fundraising.fundraising.v1.Params.place_bid_fee":
		lv := value.List()
		clv := lv.(*_Params_2_list)
		x.PlaceBidFee = *clv.list
	case "fundraising.fundraising.v1.Params.extended_period":
		x.ExtendedPeriod = uint32(value.Uint())
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: fundraising.fundraising.v1.Params"))
		}
		panic(fmt.Errorf("message fundraising.fundraising.v1.Params does not contain field %s", fd.FullName()))
	}
}

// Mutable returns a mutable reference to a composite type.
//
// If the field is unpopulated, it may allocate a composite value.
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType
// if not already stored.
// It panics if the field does not contain a composite type.
//
// Mutable is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Params) Mutable(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "fundraising.fundraising.v1.Params.auction_creation_fee":
		if x.AuctionCreationFee == nil {
			x.AuctionCreationFee = []*v1beta1.Coin{}
		}
		value := &_Params_1_list{list: &x.AuctionCreationFee}
		return protoreflect.ValueOfList(value)
	case "fundraising.fundraising.v1.Params.place_bid_fee":
		if x.PlaceBidFee == nil {
			x.PlaceBidFee = []*v1beta1.Coin{}
		}
		value := &_Params_2_list{list: &x.PlaceBidFee}
		return protoreflect.ValueOfList(value)
	case "fundraising.fundraising.v1.Params.extended_period":
		panic(fmt.Errorf("field extended_period of message fundraising.fundraising.v1.Params is not mutable"))
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: fundraising.fundraising.v1.Params"))
		}
		panic(fmt.Errorf("message fundraising.fundraising.v1.Params does not contain field %s", fd.FullName()))
	}
}

// NewField returns a new value that is assignable to the field
// for the given descriptor. For scalars, this returns the default value.
// For lists, maps, and messages, this returns a new, empty, mutable value.
func (x *fastReflection_Params) NewField(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "fundraising.fundraising.v1.Params.auction_creation_fee":
		list := []*v1beta1.Coin{}
		return protoreflect.ValueOfList(&_Params_1_list{list: &list})
	case "fundraising.fundraising.v1.Params.place_bid_fee":
		list := []*v1beta1.Coin{}
		return protoreflect.ValueOfList(&_Params_2_list{list: &list})
	case "fundraising.fundraising.v1.Params.extended_period":
		return protoreflect.ValueOfUint32(uint32(0))
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: fundraising.fundraising.v1.Params"))
		}
		panic(fmt.Errorf("message fundraising.fundraising.v1.Params does not contain field %s", fd.FullName()))
	}
}

// WhichOneof reports which field within the oneof is populated,
// returning nil if none are populated.
// It panics if the oneof descriptor does not belong to this message.
func (x *fastReflection_Params) WhichOneof(d protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	switch d.FullName() {
	default:
		panic(fmt.Errorf("%s is not a oneof field in fundraising.fundraising.v1.Params", d.FullName()))
	}
	panic("unreachable")
}

// GetUnknown retrieves the entire list of unknown fields.
// The caller may only mutate the contents of the RawFields
// if the mutated bytes are stored back into the message with SetUnknown.
func (x *fastReflection_Params) GetUnknown() protoreflect.RawFields {
	return x.unknownFields
}

// SetUnknown stores an entire list of unknown fields.
// The raw fields must be syntactically valid according to the wire format.
// An implementation may panic if this is not the case.
// Once stored, the caller must not mutate the content of the RawFields.
// An empty RawFields may be passed to clear the fields.
//
// SetUnknown is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Params) SetUnknown(fields protoreflect.RawFields) {
	x.unknownFields = fields
}

// IsValid reports whether the message is valid.
//
// An invalid message is an empty, read-only value.
//
// An invalid message often corresponds to a nil pointer of the concrete
// message type, but the details are implementation dependent.
// Validity is not part of the protobuf data model, and may not
// be preserved in marshaling or other operations.
func (x *fastReflection_Params) IsValid() bool {
	return x != nil
}

// ProtoMethods returns optional fastReflectionFeature-path implementations of various operations.
// This method may return nil.
//
// The returned methods type is identical to
// "google.golang.org/protobuf/runtime/protoiface".Methods.
// Consult the protoiface package documentation for details.
func (x *fastReflection_Params) ProtoMethods() *protoiface.Methods {
	size := func(input protoiface.SizeInput) protoiface.SizeOutput {
		x := input.Message.Interface().(*Params)
		if x == nil {
			return protoiface.SizeOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Size:              0,
			}
		}
		options := runtime.SizeInputToOptions(input)
		_ = options
		var n int
		var l int
		_ = l
		if len(x.AuctionCreationFee) > 0 {
			for _, e := range x.AuctionCreationFee {
				l = options.Size(e)
				n += 1 + l + runtime.Sov(uint64(l))
			}
		}
		if len(x.PlaceBidFee) > 0 {
			for _, e := range x.PlaceBidFee {
				l = options.Size(e)
				n += 1 + l + runtime.Sov(uint64(l))
			}
		}
		if x.ExtendedPeriod != 0 {
			n += 1 + runtime.Sov(uint64(x.ExtendedPeriod))
		}
		if x.unknownFields != nil {
			n += len(x.unknownFields)
		}
		return protoiface.SizeOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Size:              n,
		}
	}

	marshal := func(input protoiface.MarshalInput) (protoiface.MarshalOutput, error) {
		x := input.Message.Interface().(*Params)
		if x == nil {
			return protoiface.MarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Buf:               input.Buf,
			}, nil
		}
		options := runtime.MarshalInputToOptions(input)
		_ = options
		size := options.Size(x)
		dAtA := make([]byte, size)
		i := len(dAtA)
		_ = i
		var l int
		_ = l
		if x.unknownFields != nil {
			i -= len(x.unknownFields)
			copy(dAtA[i:], x.unknownFields)
		}
		if x.ExtendedPeriod != 0 {
			i = runtime.EncodeVarint(dAtA, i, uint64(x.ExtendedPeriod))
			i--
			dAtA[i] = 0x18
		}
		if len(x.PlaceBidFee) > 0 {
			for iNdEx := len(x.PlaceBidFee) - 1; iNdEx >= 0; iNdEx-- {
				encoded, err := options.Marshal(x.PlaceBidFee[iNdEx])
				if err != nil {
					return protoiface.MarshalOutput{
						NoUnkeyedLiterals: input.NoUnkeyedLiterals,
						Buf:               input.Buf,
					}, err
				}
				i -= len(encoded)
				copy(dAtA[i:], encoded)
				i = runtime.EncodeVarint(dAtA, i, uint64(len(encoded)))
				i--
				dAtA[i] = 0x12
			}
		}
		if len(x.AuctionCreationFee) > 0 {
			for iNdEx := len(x.AuctionCreationFee) - 1; iNdEx >= 0; iNdEx-- {
				encoded, err := options.Marshal(x.AuctionCreationFee[iNdEx])
				if err != nil {
					return protoiface.MarshalOutput{
						NoUnkeyedLiterals: input.NoUnkeyedLiterals,
						Buf:               input.Buf,
					}, err
				}
				i -= len(encoded)
				copy(dAtA[i:], encoded)
				i = runtime.EncodeVarint(dAtA, i, uint64(len(encoded)))
				i--
				dAtA[i] = 0xa
			}
		}
		if input.Buf != nil {
			input.Buf = append(input.Buf, dAtA...)
		} else {
			input.Buf = dAtA
		}
		return protoiface.MarshalOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Buf:               input.Buf,
		}, nil
	}
	unmarshal := func(input protoiface.UnmarshalInput) (protoiface.UnmarshalOutput, error) {
		x := input.Message.Interface().(*Params)
		if x == nil {
			return protoiface.UnmarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Flags:             input.Flags,
			}, nil
		}
		options := runtime.UnmarshalInputToOptions(input)
		_ = options
		dAtA := input.Buf
		l := len(dAtA)
		iNdEx := 0
		for iNdEx < l {
			preIndex := iNdEx
			var wire uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
				}
				if iNdEx >= l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: Params: wiretype end group for non-group")
			}
			if fieldNum <= 0 {
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
			}
			switch fieldNum {
			case 1:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field AuctionCreationFee", wireType)
				}
				var msglen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					msglen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if msglen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + msglen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.AuctionCreationFee = append(x.AuctionCreationFee, &v1beta1.Coin{})
				if err := options.Unmarshal(dAtA[iNdEx:postIndex], x.AuctionCreationFee[len(x.AuctionCreationFee)-1]); err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				iNdEx = postIndex
			case 2:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field PlaceBidFee", wireType)
				}
				var msglen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					msglen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if msglen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + msglen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.PlaceBidFee = append(x.PlaceBidFee, &v1beta1.Coin{})
				if err := options.Unmarshal(dAtA[iNdEx:postIndex], x.PlaceBidFee[len(x.PlaceBidFee)-1]); err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				iNdEx = postIndex
			case 3:
				if wireType != 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field ExtendedPeriod", wireType)
				}
				x.ExtendedPeriod = 0
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					x.ExtendedPeriod |= uint32(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
			default:
				iNdEx = preIndex
				skippy, err := runtime.Skip(dAtA[iNdEx:])
				if err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				if (skippy < 0) || (iNdEx+skippy) < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if (iNdEx + skippy) > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				if !options.DiscardUnknown {
					x.unknownFields = append(x.unknownFields, dAtA[iNdEx:iNdEx+skippy]...)
				}
				iNdEx += skippy
			}
		}

		if iNdEx > l {
			return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
		}
		return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, nil
	}
	return &protoiface.Methods{
		NoUnkeyedLiterals: struct{}{},
		Flags:             protoiface.SupportMarshalDeterministic | protoiface.SupportUnmarshalDiscardUnknown,
		Size:              size,
		Marshal:           marshal,
		Unmarshal:         unmarshal,
		Merge:             nil,
		CheckInitialized:  nil,
	}
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.0
// 	protoc        (unknown)
// source: fundraising/fundraising/v1/params.proto

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Params defines the parameters for the module.
type Params struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// auction_creation_fee specifies the fee for auction creation.
	// this prevents from spamming attack and it is collected in the community
	// pool
	AuctionCreationFee []*v1beta1.Coin `protobuf:"bytes,1,rep,name=auction_creation_fee,json=auctionCreationFee,proto3" json:"auction_creation_fee,omitempty"`
	// place_bid_fee specifies the fee for placing a bid for an auction.
	// this prevents from spamming attack and it is collected in the community
	// pool
	PlaceBidFee []*v1beta1.Coin `protobuf:"bytes,2,rep,name=place_bid_fee,json=placeBidFee,proto3" json:"place_bid_fee,omitempty"`
	// extended_period specifies the extended period that determines how long
	// the extended auction round lasts
	ExtendedPeriod uint32 `protobuf:"varint,3,opt,name=extended_period,json=extendedPeriod,proto3" json:"extended_period,omitempty"`
}

func (x *Params) Reset() {
	*x = Params{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fundraising_fundraising_v1_params_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Params) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Params) ProtoMessage() {}

// Deprecated: Use Params.ProtoReflect.Descriptor instead.
func (*Params) Descriptor() ([]byte, []int) {
	return file_fundraising_fundraising_v1_params_proto_rawDescGZIP(), []int{0}
}

func (x *Params) GetAuctionCreationFee() []*v1beta1.Coin {
	if x != nil {
		return x.AuctionCreationFee
	}
	return nil
}

func (x *Params) GetPlaceBidFee() []*v1beta1.Coin {
	if x != nil {
		return x.PlaceBidFee
	}
	return nil
}

func (x *Params) GetExtendedPeriod() uint32 {
	if x != nil {
		return x.ExtendedPeriod
	}
	return 0
}

var File_fundraising_fundraising_v1_params_proto protoreflect.FileDescriptor

var file_fundraising_fundraising_v1_params_proto_rawDesc = []byte{
	0x0a, 0x27, 0x66, 0x75, 0x6e, 0x64, 0x72, 0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x2f, 0x66, 0x75,
	0x6e, 0x64, 0x72, 0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x66, 0x75, 0x6e, 0x64, 0x72,
	0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x2e, 0x66, 0x75, 0x6e, 0x64, 0x72, 0x61, 0x69, 0x73, 0x69,
	0x6e, 0x67, 0x2e, 0x76, 0x31, 0x1a, 0x11, 0x61, 0x6d, 0x69, 0x6e, 0x6f, 0x2f, 0x61, 0x6d, 0x69,
	0x6e, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73,
	0x2f, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x63, 0x6f,
	0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73,
	0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x67, 0x6f, 0x67, 0x6f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67,
	0x6f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfa, 0x02, 0x0a, 0x06, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x12, 0x93, 0x01, 0x0a, 0x14, 0x61, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x66, 0x65, 0x65, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x62, 0x61, 0x73,
	0x65, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x43, 0x6f, 0x69, 0x6e, 0x42, 0x46,
	0xc8, 0xde, 0x1f, 0x00, 0xaa, 0xdf, 0x1f, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73,
	0x2d, 0x73, 0x64, 0x6b, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x43, 0x6f, 0x69, 0x6e, 0x73,
	0x9a, 0xe7, 0xb0, 0x2a, 0x0c, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x5f, 0x63, 0x6f, 0x69, 0x6e,
	0x73, 0xa8, 0xe7, 0xb0, 0x2a, 0x01, 0x52, 0x12, 0x61, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x65, 0x65, 0x12, 0x85, 0x01, 0x0a, 0x0d, 0x70,
	0x6c, 0x61, 0x63, 0x65, 0x5f, 0x62, 0x69, 0x64, 0x5f, 0x66, 0x65, 0x65, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x62, 0x61, 0x73, 0x65,
	0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x43, 0x6f, 0x69, 0x6e, 0x42, 0x46, 0xc8,
	0xde, 0x1f, 0x00, 0xaa, 0xdf, 0x1f, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2d,
	0x73, 0x64, 0x6b, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x43, 0x6f, 0x69, 0x6e, 0x73, 0x9a,
	0xe7, 0xb0, 0x2a, 0x0c, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x5f, 0x63, 0x6f, 0x69, 0x6e, 0x73,
	0xa8, 0xe7, 0xb0, 0x2a, 0x01, 0x52, 0x0b, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x42, 0x69, 0x64, 0x46,
	0x65, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x5f, 0x70,
	0x65, 0x72, 0x69, 0x6f, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0e, 0x65, 0x78, 0x74,
	0x65, 0x6e, 0x64, 0x65, 0x64, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x3a, 0x29, 0xe8, 0xa0, 0x1f,
	0x01, 0x8a, 0xe7, 0xb0, 0x2a, 0x20, 0x66, 0x75, 0x6e, 0x64, 0x72, 0x61, 0x69, 0x73, 0x69, 0x6e,
	0x67, 0x2f, 0x78, 0x2f, 0x66, 0x75, 0x6e, 0x64, 0x72, 0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x2f,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x42, 0x87, 0x02, 0x0a, 0x1e, 0x63, 0x6f, 0x6d, 0x2e, 0x66,
	0x75, 0x6e, 0x64, 0x72, 0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x2e, 0x66, 0x75, 0x6e, 0x64, 0x72,
	0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x42, 0x0b, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x74, 0x2f,
	0x66, 0x75, 0x6e, 0x64, 0x72, 0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x66, 0x75, 0x6e, 0x64, 0x72, 0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x2f, 0x66, 0x75, 0x6e, 0x64,
	0x72, 0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x31, 0x3b, 0x66, 0x75, 0x6e, 0x64, 0x72,
	0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x46, 0x46, 0x58, 0xaa, 0x02,
	0x1a, 0x46, 0x75, 0x6e, 0x64, 0x72, 0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x2e, 0x46, 0x75, 0x6e,
	0x64, 0x72, 0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x1a, 0x46, 0x75,
	0x6e, 0x64, 0x72, 0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x5c, 0x46, 0x75, 0x6e, 0x64, 0x72, 0x61,
	0x69, 0x73, 0x69, 0x6e, 0x67, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x26, 0x46, 0x75, 0x6e, 0x64, 0x72,
	0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x5c, 0x46, 0x75, 0x6e, 0x64, 0x72, 0x61, 0x69, 0x73, 0x69,
	0x6e, 0x67, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x1c, 0x46, 0x75, 0x6e, 0x64, 0x72, 0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x3a,
	0x3a, 0x46, 0x75, 0x6e, 0x64, 0x72, 0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x3a, 0x3a, 0x56, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_fundraising_fundraising_v1_params_proto_rawDescOnce sync.Once
	file_fundraising_fundraising_v1_params_proto_rawDescData = file_fundraising_fundraising_v1_params_proto_rawDesc
)

func file_fundraising_fundraising_v1_params_proto_rawDescGZIP() []byte {
	file_fundraising_fundraising_v1_params_proto_rawDescOnce.Do(func() {
		file_fundraising_fundraising_v1_params_proto_rawDescData = protoimpl.X.CompressGZIP(file_fundraising_fundraising_v1_params_proto_rawDescData)
	})
	return file_fundraising_fundraising_v1_params_proto_rawDescData
}

var (
	file_fundraising_fundraising_v1_params_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
	file_fundraising_fundraising_v1_params_proto_goTypes  = []interface{}{
		(*Params)(nil),       // 0: fundraising.fundraising.v1.Params
		(*v1beta1.Coin)(nil), // 1: cosmos.base.v1beta1.Coin
	}
)

var file_fundraising_fundraising_v1_params_proto_depIdxs = []int32{
	1, // 0: fundraising.fundraising.v1.Params.auction_creation_fee:type_name -> cosmos.base.v1beta1.Coin
	1, // 1: fundraising.fundraising.v1.Params.place_bid_fee:type_name -> cosmos.base.v1beta1.Coin
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_fundraising_fundraising_v1_params_proto_init() }
func file_fundraising_fundraising_v1_params_proto_init() {
	if File_fundraising_fundraising_v1_params_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fundraising_fundraising_v1_params_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Params); i {
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
			RawDescriptor: file_fundraising_fundraising_v1_params_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_fundraising_fundraising_v1_params_proto_goTypes,
		DependencyIndexes: file_fundraising_fundraising_v1_params_proto_depIdxs,
		MessageInfos:      file_fundraising_fundraising_v1_params_proto_msgTypes,
	}.Build()
	File_fundraising_fundraising_v1_params_proto = out.File
	file_fundraising_fundraising_v1_params_proto_rawDesc = nil
	file_fundraising_fundraising_v1_params_proto_goTypes = nil
	file_fundraising_fundraising_v1_params_proto_depIdxs = nil
}