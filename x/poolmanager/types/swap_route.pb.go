// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: symphony/poolmanager/v1beta1/swap_route.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
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

type SwapAmountInRoute struct {
	PoolId        uint64 `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty" yaml:"pool_id"`
	TokenOutDenom string `protobuf:"bytes,2,opt,name=token_out_denom,json=tokenOutDenom,proto3" json:"token_out_denom,omitempty" yaml:"token_out_denom"`
}

func (m *SwapAmountInRoute) Reset()         { *m = SwapAmountInRoute{} }
func (m *SwapAmountInRoute) String() string { return proto.CompactTextString(m) }
func (*SwapAmountInRoute) ProtoMessage()    {}
func (*SwapAmountInRoute) Descriptor() ([]byte, []int) {
	return fileDescriptor_c0d5ee395675eb84, []int{0}
}
func (m *SwapAmountInRoute) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SwapAmountInRoute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SwapAmountInRoute.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SwapAmountInRoute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwapAmountInRoute.Merge(m, src)
}
func (m *SwapAmountInRoute) XXX_Size() int {
	return m.Size()
}
func (m *SwapAmountInRoute) XXX_DiscardUnknown() {
	xxx_messageInfo_SwapAmountInRoute.DiscardUnknown(m)
}

var xxx_messageInfo_SwapAmountInRoute proto.InternalMessageInfo

func (m *SwapAmountInRoute) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

func (m *SwapAmountInRoute) GetTokenOutDenom() string {
	if m != nil {
		return m.TokenOutDenom
	}
	return ""
}

type SwapAmountOutRoute struct {
	PoolId       uint64 `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty" yaml:"pool_id"`
	TokenInDenom string `protobuf:"bytes,2,opt,name=token_in_denom,json=tokenInDenom,proto3" json:"token_in_denom,omitempty" yaml:"token_in_denom"`
}

func (m *SwapAmountOutRoute) Reset()         { *m = SwapAmountOutRoute{} }
func (m *SwapAmountOutRoute) String() string { return proto.CompactTextString(m) }
func (*SwapAmountOutRoute) ProtoMessage()    {}
func (*SwapAmountOutRoute) Descriptor() ([]byte, []int) {
	return fileDescriptor_c0d5ee395675eb84, []int{1}
}
func (m *SwapAmountOutRoute) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SwapAmountOutRoute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SwapAmountOutRoute.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SwapAmountOutRoute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwapAmountOutRoute.Merge(m, src)
}
func (m *SwapAmountOutRoute) XXX_Size() int {
	return m.Size()
}
func (m *SwapAmountOutRoute) XXX_DiscardUnknown() {
	xxx_messageInfo_SwapAmountOutRoute.DiscardUnknown(m)
}

var xxx_messageInfo_SwapAmountOutRoute proto.InternalMessageInfo

func (m *SwapAmountOutRoute) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

func (m *SwapAmountOutRoute) GetTokenInDenom() string {
	if m != nil {
		return m.TokenInDenom
	}
	return ""
}

type SwapAmountInSplitRoute struct {
	Pools         []SwapAmountInRoute   `protobuf:"bytes,1,rep,name=pools,proto3" json:"pools" yaml:"pools"`
	TokenInAmount cosmossdk_io_math.Int `protobuf:"bytes,2,opt,name=token_in_amount,json=tokenInAmount,proto3,customtype=cosmossdk.io/math.Int" json:"token_in_amount" yaml:"token_in_amount"`
}

func (m *SwapAmountInSplitRoute) Reset()         { *m = SwapAmountInSplitRoute{} }
func (m *SwapAmountInSplitRoute) String() string { return proto.CompactTextString(m) }
func (*SwapAmountInSplitRoute) ProtoMessage()    {}
func (*SwapAmountInSplitRoute) Descriptor() ([]byte, []int) {
	return fileDescriptor_c0d5ee395675eb84, []int{2}
}
func (m *SwapAmountInSplitRoute) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SwapAmountInSplitRoute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SwapAmountInSplitRoute.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SwapAmountInSplitRoute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwapAmountInSplitRoute.Merge(m, src)
}
func (m *SwapAmountInSplitRoute) XXX_Size() int {
	return m.Size()
}
func (m *SwapAmountInSplitRoute) XXX_DiscardUnknown() {
	xxx_messageInfo_SwapAmountInSplitRoute.DiscardUnknown(m)
}

var xxx_messageInfo_SwapAmountInSplitRoute proto.InternalMessageInfo

func (m *SwapAmountInSplitRoute) GetPools() []SwapAmountInRoute {
	if m != nil {
		return m.Pools
	}
	return nil
}

type SwapAmountOutSplitRoute struct {
	Pools          []SwapAmountOutRoute  `protobuf:"bytes,1,rep,name=pools,proto3" json:"pools" yaml:"pools"`
	TokenOutAmount cosmossdk_io_math.Int `protobuf:"bytes,2,opt,name=token_out_amount,json=tokenOutAmount,proto3,customtype=cosmossdk.io/math.Int" json:"token_out_amount" yaml:"token_out_amount"`
}

func (m *SwapAmountOutSplitRoute) Reset()         { *m = SwapAmountOutSplitRoute{} }
func (m *SwapAmountOutSplitRoute) String() string { return proto.CompactTextString(m) }
func (*SwapAmountOutSplitRoute) ProtoMessage()    {}
func (*SwapAmountOutSplitRoute) Descriptor() ([]byte, []int) {
	return fileDescriptor_c0d5ee395675eb84, []int{3}
}
func (m *SwapAmountOutSplitRoute) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SwapAmountOutSplitRoute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SwapAmountOutSplitRoute.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SwapAmountOutSplitRoute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwapAmountOutSplitRoute.Merge(m, src)
}
func (m *SwapAmountOutSplitRoute) XXX_Size() int {
	return m.Size()
}
func (m *SwapAmountOutSplitRoute) XXX_DiscardUnknown() {
	xxx_messageInfo_SwapAmountOutSplitRoute.DiscardUnknown(m)
}

var xxx_messageInfo_SwapAmountOutSplitRoute proto.InternalMessageInfo

func (m *SwapAmountOutSplitRoute) GetPools() []SwapAmountOutRoute {
	if m != nil {
		return m.Pools
	}
	return nil
}

func init() {
	proto.RegisterType((*SwapAmountInRoute)(nil), "symphony.poolmanager.v1beta1.SwapAmountInRoute")
	proto.RegisterType((*SwapAmountOutRoute)(nil), "symphony.poolmanager.v1beta1.SwapAmountOutRoute")
	proto.RegisterType((*SwapAmountInSplitRoute)(nil), "symphony.poolmanager.v1beta1.SwapAmountInSplitRoute")
	proto.RegisterType((*SwapAmountOutSplitRoute)(nil), "symphony.poolmanager.v1beta1.SwapAmountOutSplitRoute")
}

func init() {
	proto.RegisterFile("symphony/poolmanager/v1beta1/swap_route.proto", fileDescriptor_c0d5ee395675eb84)
}

var fileDescriptor_c0d5ee395675eb84 = []byte{
	// 450 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0x4d, 0x6f, 0xd3, 0x30,
	0x18, 0xc7, 0x6b, 0x5e, 0x86, 0x30, 0xa3, 0x40, 0xb4, 0x97, 0x32, 0xa1, 0xa4, 0xca, 0xa9, 0x12,
	0x9a, 0xcd, 0xc6, 0xa1, 0x88, 0x0b, 0x22, 0xe2, 0x92, 0xd3, 0x44, 0x76, 0x03, 0xa4, 0xc8, 0x59,
	0xa3, 0x36, 0x5a, 0xec, 0x27, 0xaa, 0x9d, 0x8d, 0x5c, 0x11, 0x1f, 0x80, 0x8f, 0xb5, 0xe3, 0x4e,
	0x08, 0x81, 0x14, 0xa1, 0xf6, 0x1b, 0xe4, 0x13, 0xa0, 0xbc, 0x98, 0x26, 0x45, 0x2a, 0x88, 0x9b,
	0x5f, 0x9e, 0xbf, 0x9f, 0xdf, 0xff, 0x6f, 0x1b, 0x1f, 0xca, 0x8c, 0x27, 0x33, 0x10, 0x19, 0x4d,
	0x00, 0x62, 0xce, 0x04, 0x9b, 0x86, 0x73, 0x7a, 0x71, 0x14, 0x84, 0x8a, 0x1d, 0x51, 0x79, 0xc9,
	0x12, 0x7f, 0x0e, 0xa9, 0x0a, 0x49, 0x32, 0x07, 0x05, 0xc6, 0x13, 0x5d, 0x4e, 0x5a, 0xe5, 0xa4,
	0x29, 0x3f, 0xd8, 0x99, 0xc2, 0x14, 0xaa, 0x42, 0x5a, 0x8e, 0x6a, 0x8d, 0xfd, 0x19, 0xe1, 0x47,
	0xa7, 0x97, 0x2c, 0x79, 0xcd, 0x21, 0x15, 0xca, 0x15, 0x5e, 0x79, 0x9e, 0xf1, 0x14, 0xdf, 0x29,
	0x8f, 0xf0, 0xa3, 0xc9, 0x00, 0x0d, 0xd1, 0xe8, 0x96, 0x63, 0x14, 0xb9, 0xd5, 0xcf, 0x18, 0x8f,
	0x5f, 0xda, 0xcd, 0x86, 0xed, 0x6d, 0x95, 0x23, 0x77, 0x62, 0x38, 0xf8, 0x81, 0x82, 0xf3, 0x50,
	0xf8, 0x90, 0x2a, 0x7f, 0x12, 0x0a, 0xe0, 0x83, 0x1b, 0x43, 0x34, 0xba, 0xeb, 0x1c, 0x14, 0xb9,
	0xb5, 0x57, 0x8b, 0xd6, 0x0a, 0x6c, 0xef, 0x7e, 0xb5, 0x72, 0x92, 0xaa, 0x37, 0xd5, 0xfc, 0x13,
	0xc2, 0xc6, 0x0a, 0xe3, 0x24, 0x55, 0xff, 0xc1, 0xf1, 0x0a, 0xf7, 0xeb, 0x36, 0x91, 0xe8, 0x60,
	0x3c, 0x2e, 0x72, 0x6b, 0xb7, 0x8d, 0xa1, 0xf7, 0x6d, 0x6f, 0xbb, 0x5a, 0x70, 0x45, 0x0d, 0xf1,
	0x15, 0xe1, 0xbd, 0x76, 0x16, 0xa7, 0x49, 0x1c, 0x35, 0x20, 0xef, 0xf1, 0xed, 0xb2, 0x8b, 0x1c,
	0xa0, 0xe1, 0xcd, 0xd1, 0xbd, 0x63, 0x4a, 0x36, 0x45, 0x4d, 0xfe, 0x08, 0xd4, 0xd9, 0xb9, 0xca,
	0xad, 0x5e, 0x91, 0x5b, 0xdb, 0x2b, 0x76, 0x69, 0x7b, 0xf5, 0x99, 0x86, 0xaf, 0x03, 0x8c, 0x84,
	0xcf, 0x2a, 0x59, 0x43, 0x3e, 0x2e, 0x55, 0xdf, 0x73, 0x6b, 0xf7, 0x0c, 0x24, 0x07, 0x29, 0x27,
	0xe7, 0x24, 0x02, 0xca, 0x99, 0x9a, 0x11, 0x57, 0xa8, 0xf5, 0x74, 0x7f, 0xab, 0x75, 0xba, 0xae,
	0xa8, 0x21, 0xec, 0x1f, 0x08, 0xef, 0x77, 0xd2, 0x6d, 0x39, 0xfb, 0xd0, 0x75, 0xf6, 0xec, 0x5f,
	0x9d, 0xe9, 0x3b, 0xda, 0x6c, 0x2d, 0xc0, 0x0f, 0x57, 0x57, 0xdf, 0xf1, 0xf6, 0xe2, 0x6f, 0xde,
	0xf6, 0xd7, 0x5f, 0x8e, 0x36, 0xd7, 0xd7, 0x4f, 0xa7, 0x06, 0x71, 0xde, 0x5e, 0x2d, 0x4c, 0x74,
	0xbd, 0x30, 0xd1, 0xcf, 0x85, 0x89, 0xbe, 0x2c, 0xcd, 0xde, 0xf5, 0xd2, 0xec, 0x7d, 0x5b, 0x9a,
	0xbd, 0x77, 0xe3, 0x69, 0xa4, 0x66, 0x69, 0x40, 0xce, 0x80, 0xd3, 0xaa, 0x4b, 0x24, 0x0f, 0x63,
	0x16, 0x48, 0x3d, 0xa1, 0x17, 0xc7, 0x63, 0xfa, 0xb1, 0xf3, 0xb9, 0x54, 0x96, 0x84, 0x32, 0xd8,
	0xaa, 0x3e, 0xc7, 0xf3, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x63, 0x3c, 0x24, 0xf6, 0x81, 0x03,
	0x00, 0x00,
}

func (m *SwapAmountInRoute) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SwapAmountInRoute) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SwapAmountInRoute) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TokenOutDenom) > 0 {
		i -= len(m.TokenOutDenom)
		copy(dAtA[i:], m.TokenOutDenom)
		i = encodeVarintSwapRoute(dAtA, i, uint64(len(m.TokenOutDenom)))
		i--
		dAtA[i] = 0x12
	}
	if m.PoolId != 0 {
		i = encodeVarintSwapRoute(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SwapAmountOutRoute) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SwapAmountOutRoute) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SwapAmountOutRoute) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TokenInDenom) > 0 {
		i -= len(m.TokenInDenom)
		copy(dAtA[i:], m.TokenInDenom)
		i = encodeVarintSwapRoute(dAtA, i, uint64(len(m.TokenInDenom)))
		i--
		dAtA[i] = 0x12
	}
	if m.PoolId != 0 {
		i = encodeVarintSwapRoute(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SwapAmountInSplitRoute) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SwapAmountInSplitRoute) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SwapAmountInSplitRoute) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.TokenInAmount.Size()
		i -= size
		if _, err := m.TokenInAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintSwapRoute(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Pools) > 0 {
		for iNdEx := len(m.Pools) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Pools[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSwapRoute(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *SwapAmountOutSplitRoute) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SwapAmountOutSplitRoute) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SwapAmountOutSplitRoute) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.TokenOutAmount.Size()
		i -= size
		if _, err := m.TokenOutAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintSwapRoute(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Pools) > 0 {
		for iNdEx := len(m.Pools) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Pools[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSwapRoute(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintSwapRoute(dAtA []byte, offset int, v uint64) int {
	offset -= sovSwapRoute(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SwapAmountInRoute) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolId != 0 {
		n += 1 + sovSwapRoute(uint64(m.PoolId))
	}
	l = len(m.TokenOutDenom)
	if l > 0 {
		n += 1 + l + sovSwapRoute(uint64(l))
	}
	return n
}

func (m *SwapAmountOutRoute) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolId != 0 {
		n += 1 + sovSwapRoute(uint64(m.PoolId))
	}
	l = len(m.TokenInDenom)
	if l > 0 {
		n += 1 + l + sovSwapRoute(uint64(l))
	}
	return n
}

func (m *SwapAmountInSplitRoute) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Pools) > 0 {
		for _, e := range m.Pools {
			l = e.Size()
			n += 1 + l + sovSwapRoute(uint64(l))
		}
	}
	l = m.TokenInAmount.Size()
	n += 1 + l + sovSwapRoute(uint64(l))
	return n
}

func (m *SwapAmountOutSplitRoute) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Pools) > 0 {
		for _, e := range m.Pools {
			l = e.Size()
			n += 1 + l + sovSwapRoute(uint64(l))
		}
	}
	l = m.TokenOutAmount.Size()
	n += 1 + l + sovSwapRoute(uint64(l))
	return n
}

func sovSwapRoute(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSwapRoute(x uint64) (n int) {
	return sovSwapRoute(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SwapAmountInRoute) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSwapRoute
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
			return fmt.Errorf("proto: SwapAmountInRoute: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SwapAmountInRoute: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwapRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenOutDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwapRoute
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
				return ErrInvalidLengthSwapRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSwapRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenOutDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSwapRoute(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSwapRoute
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
func (m *SwapAmountOutRoute) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSwapRoute
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
			return fmt.Errorf("proto: SwapAmountOutRoute: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SwapAmountOutRoute: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwapRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenInDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwapRoute
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
				return ErrInvalidLengthSwapRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSwapRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenInDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSwapRoute(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSwapRoute
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
func (m *SwapAmountInSplitRoute) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSwapRoute
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
			return fmt.Errorf("proto: SwapAmountInSplitRoute: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SwapAmountInSplitRoute: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pools", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwapRoute
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
				return ErrInvalidLengthSwapRoute
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSwapRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pools = append(m.Pools, SwapAmountInRoute{})
			if err := m.Pools[len(m.Pools)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenInAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwapRoute
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
				return ErrInvalidLengthSwapRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSwapRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TokenInAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSwapRoute(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSwapRoute
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
func (m *SwapAmountOutSplitRoute) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSwapRoute
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
			return fmt.Errorf("proto: SwapAmountOutSplitRoute: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SwapAmountOutSplitRoute: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pools", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwapRoute
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
				return ErrInvalidLengthSwapRoute
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSwapRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pools = append(m.Pools, SwapAmountOutRoute{})
			if err := m.Pools[len(m.Pools)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenOutAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwapRoute
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
				return ErrInvalidLengthSwapRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSwapRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TokenOutAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSwapRoute(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSwapRoute
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
func skipSwapRoute(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSwapRoute
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
					return 0, ErrIntOverflowSwapRoute
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
					return 0, ErrIntOverflowSwapRoute
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
				return 0, ErrInvalidLengthSwapRoute
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSwapRoute
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSwapRoute
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSwapRoute        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSwapRoute          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSwapRoute = fmt.Errorf("proto: unexpected end of group")
)
