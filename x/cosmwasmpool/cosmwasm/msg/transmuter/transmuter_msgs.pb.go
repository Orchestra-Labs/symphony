// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: symphony/cosmwasmpool/v1beta1/model/transmuter_msgs.proto

package transmuter

import (
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

// ===================== JoinPoolExecuteMsg
type EmptyRequest struct {
}

func (m *EmptyRequest) Reset()         { *m = EmptyRequest{} }
func (m *EmptyRequest) String() string { return proto.CompactTextString(m) }
func (*EmptyRequest) ProtoMessage()    {}
func (*EmptyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_803a225850f3aa24, []int{0}
}
func (m *EmptyRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EmptyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EmptyRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EmptyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptyRequest.Merge(m, src)
}
func (m *EmptyRequest) XXX_Size() int {
	return m.Size()
}
func (m *EmptyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EmptyRequest proto.InternalMessageInfo

type JoinPoolExecuteMsgRequest struct {
	// join_pool is the structure containing all request fields of the join pool
	// execute message.
	JoinPool EmptyRequest `protobuf:"bytes,1,opt,name=join_pool,json=joinPool,proto3" json:"join_pool"`
}

func (m *JoinPoolExecuteMsgRequest) Reset()         { *m = JoinPoolExecuteMsgRequest{} }
func (m *JoinPoolExecuteMsgRequest) String() string { return proto.CompactTextString(m) }
func (*JoinPoolExecuteMsgRequest) ProtoMessage()    {}
func (*JoinPoolExecuteMsgRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_803a225850f3aa24, []int{1}
}
func (m *JoinPoolExecuteMsgRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *JoinPoolExecuteMsgRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_JoinPoolExecuteMsgRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *JoinPoolExecuteMsgRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinPoolExecuteMsgRequest.Merge(m, src)
}
func (m *JoinPoolExecuteMsgRequest) XXX_Size() int {
	return m.Size()
}
func (m *JoinPoolExecuteMsgRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinPoolExecuteMsgRequest.DiscardUnknown(m)
}

var xxx_messageInfo_JoinPoolExecuteMsgRequest proto.InternalMessageInfo

func (m *JoinPoolExecuteMsgRequest) GetJoinPool() EmptyRequest {
	if m != nil {
		return m.JoinPool
	}
	return EmptyRequest{}
}

type JoinPoolExecuteMsgResponse struct {
}

func (m *JoinPoolExecuteMsgResponse) Reset()         { *m = JoinPoolExecuteMsgResponse{} }
func (m *JoinPoolExecuteMsgResponse) String() string { return proto.CompactTextString(m) }
func (*JoinPoolExecuteMsgResponse) ProtoMessage()    {}
func (*JoinPoolExecuteMsgResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_803a225850f3aa24, []int{2}
}
func (m *JoinPoolExecuteMsgResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *JoinPoolExecuteMsgResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_JoinPoolExecuteMsgResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *JoinPoolExecuteMsgResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinPoolExecuteMsgResponse.Merge(m, src)
}
func (m *JoinPoolExecuteMsgResponse) XXX_Size() int {
	return m.Size()
}
func (m *JoinPoolExecuteMsgResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinPoolExecuteMsgResponse.DiscardUnknown(m)
}

var xxx_messageInfo_JoinPoolExecuteMsgResponse proto.InternalMessageInfo

// ===================== ExitPoolExecuteMsg
type ExitPoolExecuteMsgRequest struct {
	// exit_pool is the structure containing all request fields of the exit pool
	// execute message.
	ExitPool EmptyRequest `protobuf:"bytes,1,opt,name=exit_pool,json=exitPool,proto3" json:"exit_pool"`
}

func (m *ExitPoolExecuteMsgRequest) Reset()         { *m = ExitPoolExecuteMsgRequest{} }
func (m *ExitPoolExecuteMsgRequest) String() string { return proto.CompactTextString(m) }
func (*ExitPoolExecuteMsgRequest) ProtoMessage()    {}
func (*ExitPoolExecuteMsgRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_803a225850f3aa24, []int{3}
}
func (m *ExitPoolExecuteMsgRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExitPoolExecuteMsgRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExitPoolExecuteMsgRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExitPoolExecuteMsgRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExitPoolExecuteMsgRequest.Merge(m, src)
}
func (m *ExitPoolExecuteMsgRequest) XXX_Size() int {
	return m.Size()
}
func (m *ExitPoolExecuteMsgRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ExitPoolExecuteMsgRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ExitPoolExecuteMsgRequest proto.InternalMessageInfo

func (m *ExitPoolExecuteMsgRequest) GetExitPool() EmptyRequest {
	if m != nil {
		return m.ExitPool
	}
	return EmptyRequest{}
}

type ExitPoolExecuteMsgResponse struct {
}

func (m *ExitPoolExecuteMsgResponse) Reset()         { *m = ExitPoolExecuteMsgResponse{} }
func (m *ExitPoolExecuteMsgResponse) String() string { return proto.CompactTextString(m) }
func (*ExitPoolExecuteMsgResponse) ProtoMessage()    {}
func (*ExitPoolExecuteMsgResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_803a225850f3aa24, []int{4}
}
func (m *ExitPoolExecuteMsgResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExitPoolExecuteMsgResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExitPoolExecuteMsgResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExitPoolExecuteMsgResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExitPoolExecuteMsgResponse.Merge(m, src)
}
func (m *ExitPoolExecuteMsgResponse) XXX_Size() int {
	return m.Size()
}
func (m *ExitPoolExecuteMsgResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ExitPoolExecuteMsgResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ExitPoolExecuteMsgResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*EmptyRequest)(nil), "symphony.cosmwasmpool.v1beta1.EmptyRequest")
	proto.RegisterType((*JoinPoolExecuteMsgRequest)(nil), "symphony.cosmwasmpool.v1beta1.JoinPoolExecuteMsgRequest")
	proto.RegisterType((*JoinPoolExecuteMsgResponse)(nil), "symphony.cosmwasmpool.v1beta1.JoinPoolExecuteMsgResponse")
	proto.RegisterType((*ExitPoolExecuteMsgRequest)(nil), "symphony.cosmwasmpool.v1beta1.ExitPoolExecuteMsgRequest")
	proto.RegisterType((*ExitPoolExecuteMsgResponse)(nil), "symphony.cosmwasmpool.v1beta1.ExitPoolExecuteMsgResponse")
}

func init() {
	proto.RegisterFile("symphony/cosmwasmpool/v1beta1/model/transmuter_msgs.proto", fileDescriptor_803a225850f3aa24)
}

var fileDescriptor_803a225850f3aa24 = []byte{
	// 293 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x91, 0x31, 0x4b, 0xc3, 0x40,
	0x14, 0xc7, 0x13, 0x10, 0xd1, 0x28, 0x0e, 0xc5, 0x41, 0x8b, 0x9e, 0xd2, 0x49, 0x10, 0xef, 0x68,
	0x1d, 0xc4, 0xb5, 0x90, 0xa5, 0xa0, 0x48, 0x47, 0x97, 0x92, 0xa4, 0xc7, 0xf5, 0x6a, 0x2e, 0x2f,
	0xe6, 0xbd, 0xd4, 0xe4, 0x5b, 0xf8, 0xb1, 0x3a, 0x76, 0x74, 0x12, 0x49, 0xbe, 0x88, 0xa4, 0x4d,
	0x20, 0x42, 0x75, 0x71, 0xbb, 0xc7, 0xbd, 0xff, 0xfb, 0xfd, 0xe0, 0xef, 0xdc, 0x63, 0x6e, 0xe2,
	0x19, 0x44, 0xb9, 0x08, 0x00, 0xcd, 0x9b, 0x87, 0x26, 0x06, 0x08, 0xc5, 0xa2, 0xef, 0x4b, 0xf2,
	0xfa, 0xc2, 0xc0, 0x54, 0x86, 0x82, 0x12, 0x2f, 0x42, 0x93, 0x92, 0x4c, 0x26, 0x06, 0x15, 0xf2,
	0x38, 0x01, 0x82, 0xce, 0x79, 0x13, 0xe5, 0xed, 0x28, 0xaf, 0xa3, 0xdd, 0x63, 0x05, 0x0a, 0xd6,
	0x9b, 0xa2, 0x7a, 0x6d, 0x42, 0xbd, 0x23, 0xe7, 0xd0, 0x35, 0x31, 0xe5, 0x63, 0xf9, 0x9a, 0x4a,
	0xa4, 0xde, 0x8b, 0x73, 0x3a, 0x02, 0x1d, 0x3d, 0x01, 0x84, 0x6e, 0x26, 0x83, 0x94, 0xe4, 0x03,
	0xaa, 0xfa, 0xb3, 0xf3, 0xe8, 0xec, 0xcf, 0x41, 0x47, 0x93, 0xea, 0xee, 0x89, 0x7d, 0x69, 0x5f,
	0x1d, 0x0c, 0xae, 0xf9, 0x9f, 0x54, 0xde, 0x3e, 0x3e, 0xdc, 0x59, 0x7e, 0x5e, 0x58, 0xe3, 0xbd,
	0x79, 0x0d, 0xe8, 0x9d, 0x39, 0xdd, 0x6d, 0x30, 0x8c, 0x21, 0x42, 0x59, 0xa9, 0xb8, 0x99, 0xa6,
	0x5f, 0x55, 0x64, 0xa6, 0xe9, 0xbf, 0x2a, 0xb2, 0x06, 0x54, 0x2a, 0xdb, 0x60, 0x1b, 0x95, 0xe1,
	0x74, 0x59, 0x30, 0x7b, 0x55, 0x30, 0xfb, 0xab, 0x60, 0xf6, 0x7b, 0xc9, 0xac, 0x55, 0xc9, 0xac,
	0x8f, 0x92, 0x59, 0xcf, 0x23, 0xa5, 0x69, 0x96, 0xfa, 0x3c, 0x00, 0x23, 0x00, 0x0d, 0xa0, 0xc6,
	0x9b, 0xd0, 0xf3, 0xb1, 0x19, 0xc4, 0x62, 0x70, 0x27, 0xb2, 0x9f, 0x65, 0x36, 0x83, 0x30, 0xa8,
	0x5a, 0x5d, 0xfa, 0xbb, 0xeb, 0x4a, 0x6e, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0x2f, 0x5c, 0xf0,
	0xd5, 0x04, 0x02, 0x00, 0x00,
}

func (m *EmptyRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EmptyRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EmptyRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *JoinPoolExecuteMsgRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *JoinPoolExecuteMsgRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *JoinPoolExecuteMsgRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.JoinPool.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTransmuterMsgs(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *JoinPoolExecuteMsgResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *JoinPoolExecuteMsgResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *JoinPoolExecuteMsgResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *ExitPoolExecuteMsgRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExitPoolExecuteMsgRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExitPoolExecuteMsgRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.ExitPool.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTransmuterMsgs(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *ExitPoolExecuteMsgResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExitPoolExecuteMsgResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExitPoolExecuteMsgResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTransmuterMsgs(dAtA []byte, offset int, v uint64) int {
	offset -= sovTransmuterMsgs(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *EmptyRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *JoinPoolExecuteMsgRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.JoinPool.Size()
	n += 1 + l + sovTransmuterMsgs(uint64(l))
	return n
}

func (m *JoinPoolExecuteMsgResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *ExitPoolExecuteMsgRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.ExitPool.Size()
	n += 1 + l + sovTransmuterMsgs(uint64(l))
	return n
}

func (m *ExitPoolExecuteMsgResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTransmuterMsgs(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTransmuterMsgs(x uint64) (n int) {
	return sovTransmuterMsgs(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EmptyRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTransmuterMsgs
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
			return fmt.Errorf("proto: EmptyRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EmptyRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTransmuterMsgs(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTransmuterMsgs
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
func (m *JoinPoolExecuteMsgRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTransmuterMsgs
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
			return fmt.Errorf("proto: JoinPoolExecuteMsgRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: JoinPoolExecuteMsgRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field JoinPool", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransmuterMsgs
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
				return ErrInvalidLengthTransmuterMsgs
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTransmuterMsgs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.JoinPool.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTransmuterMsgs(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTransmuterMsgs
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
func (m *JoinPoolExecuteMsgResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTransmuterMsgs
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
			return fmt.Errorf("proto: JoinPoolExecuteMsgResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: JoinPoolExecuteMsgResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTransmuterMsgs(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTransmuterMsgs
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
func (m *ExitPoolExecuteMsgRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTransmuterMsgs
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
			return fmt.Errorf("proto: ExitPoolExecuteMsgRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExitPoolExecuteMsgRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExitPool", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransmuterMsgs
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
				return ErrInvalidLengthTransmuterMsgs
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTransmuterMsgs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ExitPool.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTransmuterMsgs(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTransmuterMsgs
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
func (m *ExitPoolExecuteMsgResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTransmuterMsgs
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
			return fmt.Errorf("proto: ExitPoolExecuteMsgResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExitPoolExecuteMsgResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTransmuterMsgs(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTransmuterMsgs
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
func skipTransmuterMsgs(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTransmuterMsgs
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
					return 0, ErrIntOverflowTransmuterMsgs
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
					return 0, ErrIntOverflowTransmuterMsgs
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
				return 0, ErrInvalidLengthTransmuterMsgs
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTransmuterMsgs
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTransmuterMsgs
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTransmuterMsgs        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTransmuterMsgs          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTransmuterMsgs = fmt.Errorf("proto: unexpected end of group")
)
