// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: symphony/incentives/group.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	_ "github.com/osmosis-labs/osmosis/v27/x/lockup/types"
	_ "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
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

// SplittingPolicy determines the way we want to split incentives in groupGauges
type SplittingPolicy int32

const (
	ByVolume SplittingPolicy = 0
)

var SplittingPolicy_name = map[int32]string{
	0: "ByVolume",
}

var SplittingPolicy_value = map[string]int32{
	"ByVolume": 0,
}

func (x SplittingPolicy) String() string {
	return proto.EnumName(SplittingPolicy_name, int32(x))
}

func (SplittingPolicy) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2ec2db04d4afc72a, []int{0}
}

// Note that while both InternalGaugeInfo and InternalGaugeRecord could
// technically be replaced by DistrInfo and DistrRecord from the pool-incentives
// module, we create separate types here to keep our abstractions clean and
// readable (pool-incentives distribution abstractions are used in a very
// specific way that does not directly relate to gauge logic). This also helps
// us sidestep a refactor to avoid an import cycle.
type InternalGaugeInfo struct {
	TotalWeight  cosmossdk_io_math.Int `protobuf:"bytes,1,opt,name=total_weight,json=totalWeight,proto3,customtype=cosmossdk.io/math.Int" json:"total_weight" yaml:"total_weight"`
	GaugeRecords []InternalGaugeRecord `protobuf:"bytes,2,rep,name=gauge_records,json=gaugeRecords,proto3" json:"gauge_records"`
}

func (m *InternalGaugeInfo) Reset()         { *m = InternalGaugeInfo{} }
func (m *InternalGaugeInfo) String() string { return proto.CompactTextString(m) }
func (*InternalGaugeInfo) ProtoMessage()    {}
func (*InternalGaugeInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_2ec2db04d4afc72a, []int{0}
}
func (m *InternalGaugeInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InternalGaugeInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InternalGaugeInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InternalGaugeInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InternalGaugeInfo.Merge(m, src)
}
func (m *InternalGaugeInfo) XXX_Size() int {
	return m.Size()
}
func (m *InternalGaugeInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_InternalGaugeInfo.DiscardUnknown(m)
}

var xxx_messageInfo_InternalGaugeInfo proto.InternalMessageInfo

func (m *InternalGaugeInfo) GetGaugeRecords() []InternalGaugeRecord {
	if m != nil {
		return m.GaugeRecords
	}
	return nil
}

type InternalGaugeRecord struct {
	GaugeId uint64 `protobuf:"varint,1,opt,name=gauge_id,json=gaugeId,proto3" json:"gauge_id,omitempty" yaml:"gauge_id"`
	// CurrentWeight is the current weight of this gauge being distributed to for
	// this epoch. For instance, for volume splitting policy, this stores the
	// volume generated in the last epoch of the linked pool.
	CurrentWeight cosmossdk_io_math.Int `protobuf:"bytes,2,opt,name=current_weight,json=currentWeight,proto3,customtype=cosmossdk.io/math.Int" json:"current_weight"`
	// CumulativeWeight serves as a snapshot of the accumulator being tracked
	// based on splitting policy. For instance, for volume splitting policy, this
	// stores the cumulative volume for the linked pool at time of last update.
	CumulativeWeight cosmossdk_io_math.Int `protobuf:"bytes,3,opt,name=cumulative_weight,json=cumulativeWeight,proto3,customtype=cosmossdk.io/math.Int" json:"cumulative_weight"`
}

func (m *InternalGaugeRecord) Reset()         { *m = InternalGaugeRecord{} }
func (m *InternalGaugeRecord) String() string { return proto.CompactTextString(m) }
func (*InternalGaugeRecord) ProtoMessage()    {}
func (*InternalGaugeRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_2ec2db04d4afc72a, []int{1}
}
func (m *InternalGaugeRecord) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InternalGaugeRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InternalGaugeRecord.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InternalGaugeRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InternalGaugeRecord.Merge(m, src)
}
func (m *InternalGaugeRecord) XXX_Size() int {
	return m.Size()
}
func (m *InternalGaugeRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_InternalGaugeRecord.DiscardUnknown(m)
}

var xxx_messageInfo_InternalGaugeRecord proto.InternalMessageInfo

func (m *InternalGaugeRecord) GetGaugeId() uint64 {
	if m != nil {
		return m.GaugeId
	}
	return 0
}

// Group is an object that stores a 1:1 mapped gauge ID, a list of pool gauge
// info, and a splitting policy. These are grouped into a single abstraction to
// allow for distribution of group incentives to internal gauges according to
// the specified splitting policy.
type Group struct {
	GroupGaugeId      uint64            `protobuf:"varint,1,opt,name=group_gauge_id,json=groupGaugeId,proto3" json:"group_gauge_id,omitempty"`
	InternalGaugeInfo InternalGaugeInfo `protobuf:"bytes,2,opt,name=internal_gauge_info,json=internalGaugeInfo,proto3" json:"internal_gauge_info"`
	SplittingPolicy   SplittingPolicy   `protobuf:"varint,3,opt,name=splitting_policy,json=splittingPolicy,proto3,enum=symphony.incentives.SplittingPolicy" json:"splitting_policy,omitempty"`
}

func (m *Group) Reset()         { *m = Group{} }
func (m *Group) String() string { return proto.CompactTextString(m) }
func (*Group) ProtoMessage()    {}
func (*Group) Descriptor() ([]byte, []int) {
	return fileDescriptor_2ec2db04d4afc72a, []int{2}
}
func (m *Group) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Group) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Group.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Group) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Group.Merge(m, src)
}
func (m *Group) XXX_Size() int {
	return m.Size()
}
func (m *Group) XXX_DiscardUnknown() {
	xxx_messageInfo_Group.DiscardUnknown(m)
}

var xxx_messageInfo_Group proto.InternalMessageInfo

func (m *Group) GetGroupGaugeId() uint64 {
	if m != nil {
		return m.GroupGaugeId
	}
	return 0
}

func (m *Group) GetInternalGaugeInfo() InternalGaugeInfo {
	if m != nil {
		return m.InternalGaugeInfo
	}
	return InternalGaugeInfo{}
}

func (m *Group) GetSplittingPolicy() SplittingPolicy {
	if m != nil {
		return m.SplittingPolicy
	}
	return ByVolume
}

// CreateGroup is called via governance to create a new group.
// It takes an array of pool IDs to split the incentives across.
type CreateGroup struct {
	PoolIds []uint64 `protobuf:"varint,1,rep,packed,name=pool_ids,json=poolIds,proto3" json:"pool_ids,omitempty"`
}

func (m *CreateGroup) Reset()         { *m = CreateGroup{} }
func (m *CreateGroup) String() string { return proto.CompactTextString(m) }
func (*CreateGroup) ProtoMessage()    {}
func (*CreateGroup) Descriptor() ([]byte, []int) {
	return fileDescriptor_2ec2db04d4afc72a, []int{3}
}
func (m *CreateGroup) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CreateGroup) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CreateGroup.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CreateGroup) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateGroup.Merge(m, src)
}
func (m *CreateGroup) XXX_Size() int {
	return m.Size()
}
func (m *CreateGroup) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateGroup.DiscardUnknown(m)
}

var xxx_messageInfo_CreateGroup proto.InternalMessageInfo

func (m *CreateGroup) GetPoolIds() []uint64 {
	if m != nil {
		return m.PoolIds
	}
	return nil
}

// GroupsWithGauge is a helper struct that stores a group and its
// associated gauge.
type GroupsWithGauge struct {
	Group Group `protobuf:"bytes,1,opt,name=group,proto3" json:"group"`
	Gauge Gauge `protobuf:"bytes,2,opt,name=gauge,proto3" json:"gauge"`
}

func (m *GroupsWithGauge) Reset()         { *m = GroupsWithGauge{} }
func (m *GroupsWithGauge) String() string { return proto.CompactTextString(m) }
func (*GroupsWithGauge) ProtoMessage()    {}
func (*GroupsWithGauge) Descriptor() ([]byte, []int) {
	return fileDescriptor_2ec2db04d4afc72a, []int{4}
}
func (m *GroupsWithGauge) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GroupsWithGauge) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GroupsWithGauge.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GroupsWithGauge) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GroupsWithGauge.Merge(m, src)
}
func (m *GroupsWithGauge) XXX_Size() int {
	return m.Size()
}
func (m *GroupsWithGauge) XXX_DiscardUnknown() {
	xxx_messageInfo_GroupsWithGauge.DiscardUnknown(m)
}

var xxx_messageInfo_GroupsWithGauge proto.InternalMessageInfo

func (m *GroupsWithGauge) GetGroup() Group {
	if m != nil {
		return m.Group
	}
	return Group{}
}

func (m *GroupsWithGauge) GetGauge() Gauge {
	if m != nil {
		return m.Gauge
	}
	return Gauge{}
}

func init() {
	proto.RegisterEnum("symphony.incentives.SplittingPolicy", SplittingPolicy_name, SplittingPolicy_value)
	proto.RegisterType((*InternalGaugeInfo)(nil), "symphony.incentives.InternalGaugeInfo")
	proto.RegisterType((*InternalGaugeRecord)(nil), "symphony.incentives.InternalGaugeRecord")
	proto.RegisterType((*Group)(nil), "symphony.incentives.Group")
	proto.RegisterType((*CreateGroup)(nil), "symphony.incentives.CreateGroup")
	proto.RegisterType((*GroupsWithGauge)(nil), "symphony.incentives.GroupsWithGauge")
}

func init() { proto.RegisterFile("symphony/incentives/group.proto", fileDescriptor_2ec2db04d4afc72a) }

var fileDescriptor_2ec2db04d4afc72a = []byte{
	// 619 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0x4f, 0x8b, 0xd3, 0x40,
	0x18, 0xc6, 0x93, 0xdd, 0xae, 0xbb, 0x4e, 0xbb, 0xdb, 0x36, 0x55, 0xa8, 0x05, 0x93, 0x12, 0x56,
	0x29, 0x82, 0x19, 0xb6, 0xca, 0x0a, 0x7b, 0xac, 0x42, 0xa9, 0x17, 0x97, 0x2c, 0x58, 0x10, 0xa1,
	0x4c, 0x92, 0x69, 0x3a, 0x6c, 0x92, 0x09, 0x99, 0x49, 0x35, 0x47, 0x6f, 0x1e, 0xfd, 0x08, 0x82,
	0x9f, 0xc4, 0xdb, 0x1e, 0xf7, 0x28, 0xa2, 0x45, 0xda, 0x8b, 0xe7, 0xfd, 0x04, 0x92, 0x49, 0xb2,
	0xdd, 0xae, 0xf5, 0xcf, 0x29, 0x79, 0x67, 0x9e, 0xe7, 0xc9, 0xfb, 0x7b, 0x79, 0x03, 0x34, 0x96,
	0xf8, 0xe1, 0x84, 0x06, 0x09, 0x24, 0x81, 0x8d, 0x03, 0x4e, 0xa6, 0x98, 0x41, 0x37, 0xa2, 0x71,
	0x68, 0x84, 0x11, 0xe5, 0x54, 0x69, 0x14, 0x02, 0x63, 0x29, 0x68, 0xdd, 0x72, 0xa9, 0x4b, 0xc5,
	0x3d, 0x4c, 0xdf, 0x32, 0x69, 0x4b, 0x75, 0x29, 0x75, 0x3d, 0x0c, 0x45, 0x65, 0xc5, 0x63, 0xe8,
	0xc4, 0x11, 0xe2, 0x84, 0x06, 0xf9, 0xbd, 0x76, 0xfd, 0x9e, 0x13, 0x1f, 0x33, 0x8e, 0xfc, 0xb0,
	0x08, 0xb0, 0x29, 0xf3, 0x29, 0x83, 0x16, 0x62, 0x18, 0x4e, 0x0f, 0x2c, 0xcc, 0xd1, 0x01, 0xb4,
	0x29, 0x29, 0x02, 0x5a, 0x97, 0xcd, 0x7a, 0xd4, 0x3e, 0x8d, 0x43, 0xf1, 0x28, 0xc2, 0xd7, 0x82,
	0xa0, 0xd8, 0xc5, 0x99, 0x40, 0xff, 0x2c, 0x83, 0xfa, 0x20, 0xe0, 0x38, 0x0a, 0x90, 0xd7, 0x4f,
	0xcf, 0x07, 0xc1, 0x98, 0x2a, 0x43, 0x50, 0xe1, 0x94, 0x23, 0x6f, 0xf4, 0x06, 0x13, 0x77, 0xc2,
	0x9b, 0x72, 0x5b, 0xee, 0xdc, 0xec, 0x3d, 0x3e, 0x9b, 0x69, 0xd2, 0xd7, 0x99, 0x76, 0x3b, 0x6b,
	0x88, 0x39, 0xa7, 0x06, 0xa1, 0xd0, 0x47, 0x7c, 0x62, 0x0c, 0x02, 0x7e, 0x31, 0xd3, 0x1a, 0x09,
	0xf2, 0xbd, 0x23, 0xfd, 0xaa, 0x55, 0x37, 0xcb, 0xa2, 0x1c, 0x8a, 0x4a, 0x39, 0x01, 0xbb, 0xe2,
	0xeb, 0xa3, 0x08, 0xdb, 0x34, 0x72, 0x58, 0x73, 0xa3, 0xbd, 0xd9, 0x29, 0x77, 0x3b, 0xc6, 0x9a,
	0x79, 0x1a, 0x2b, 0x7d, 0x99, 0xc2, 0xd0, 0x2b, 0xa5, 0x3d, 0x98, 0x15, 0x77, 0x79, 0xc4, 0xf4,
	0xef, 0x32, 0x68, 0xac, 0xd1, 0x2a, 0x06, 0xd8, 0xc9, 0x3e, 0x46, 0x1c, 0x41, 0x50, 0xea, 0x35,
	0x2e, 0x66, 0x5a, 0x35, 0x6b, 0xb2, 0xb8, 0xd1, 0xcd, 0x6d, 0xf1, 0x3a, 0x70, 0x94, 0x67, 0x60,
	0xcf, 0x8e, 0xa3, 0x08, 0x07, 0xbc, 0xe0, 0xde, 0x10, 0xdc, 0x77, 0xff, 0xca, 0x6d, 0xee, 0xe6,
	0xa6, 0x1c, 0xf1, 0x39, 0xa8, 0xdb, 0xb1, 0x1f, 0x7b, 0x28, 0x85, 0x28, 0x82, 0x36, 0xff, 0x27,
	0xa8, 0xb6, 0xf4, 0x65, 0x59, 0x47, 0xa5, 0x9f, 0x1f, 0x35, 0x59, 0xff, 0x26, 0x83, 0xad, 0x7e,
	0xba, 0x7c, 0xca, 0x3e, 0xd8, 0x13, 0x5b, 0x38, 0x5a, 0xe5, 0x32, 0x2b, 0xe2, 0xb4, 0x9f, 0x73,
	0xbc, 0x06, 0x0d, 0x92, 0x8f, 0xa3, 0x10, 0x06, 0x63, 0x2a, 0x60, 0xca, 0xdd, 0xfb, 0xff, 0x1e,
	0x75, 0xba, 0x02, 0xf9, 0xa0, 0xeb, 0xe4, 0xb7, 0xdd, 0x78, 0x01, 0x6a, 0x2c, 0xf4, 0x08, 0xe7,
	0x24, 0x70, 0x47, 0x21, 0xf5, 0x88, 0x9d, 0x08, 0xbc, 0xbd, 0xee, 0xfe, 0xda, 0xe8, 0x93, 0x42,
	0x7c, 0x2c, 0xb4, 0x66, 0x95, 0xad, 0x1e, 0xe8, 0x1d, 0x50, 0x7e, 0x1a, 0x61, 0xc4, 0x71, 0xc6,
	0x78, 0x07, 0xec, 0x84, 0x94, 0x7a, 0x23, 0xe2, 0xb0, 0xa6, 0xdc, 0xde, 0xec, 0x94, 0xcc, 0xed,
	0xb4, 0x1e, 0x38, 0x4c, 0x7f, 0x27, 0x83, 0xaa, 0x10, 0xb1, 0x21, 0xe1, 0x13, 0xd1, 0x92, 0x72,
	0x08, 0xb6, 0x04, 0xbc, 0x98, 0x44, 0xb9, 0xdb, 0x5a, 0xdb, 0x83, 0x30, 0xe5, 0x48, 0x99, 0x5c,
	0xf8, 0xd2, 0x80, 0x7c, 0x2c, 0x7f, 0xf0, 0xa5, 0x8a, 0x4b, 0x5f, 0x5a, 0x3c, 0xb8, 0x07, 0xaa,
	0xd7, 0x88, 0x94, 0x0a, 0xd8, 0xe9, 0x25, 0x2f, 0xa9, 0x17, 0xfb, 0xb8, 0x26, 0xb5, 0x4a, 0xef,
	0x3f, 0xa9, 0x52, 0xef, 0xf8, 0x6c, 0xae, 0xca, 0xe7, 0x73, 0x55, 0xfe, 0x31, 0x57, 0xe5, 0x0f,
	0x0b, 0x55, 0x3a, 0x5f, 0xa8, 0xd2, 0x97, 0x85, 0x2a, 0xbd, 0x3a, 0x74, 0x09, 0x9f, 0xc4, 0x96,
	0x61, 0x53, 0x1f, 0x8a, 0x35, 0x20, 0xec, 0xa1, 0x87, 0x2c, 0x56, 0x14, 0x70, 0xda, 0x7d, 0x02,
	0xdf, 0x5e, 0xfd, 0x5f, 0x79, 0x12, 0x62, 0x66, 0xdd, 0x10, 0x3f, 0xec, 0xa3, 0x5f, 0x01, 0x00,
	0x00, 0xff, 0xff, 0xf8, 0xf3, 0x64, 0x09, 0x9c, 0x04, 0x00, 0x00,
}

func (this *InternalGaugeRecord) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*InternalGaugeRecord)
	if !ok {
		that2, ok := that.(InternalGaugeRecord)
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
	if this.GaugeId != that1.GaugeId {
		return false
	}
	if !this.CurrentWeight.Equal(that1.CurrentWeight) {
		return false
	}
	if !this.CumulativeWeight.Equal(that1.CumulativeWeight) {
		return false
	}
	return true
}
func (m *InternalGaugeInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InternalGaugeInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *InternalGaugeInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.GaugeRecords) > 0 {
		for iNdEx := len(m.GaugeRecords) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.GaugeRecords[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGroup(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size := m.TotalWeight.Size()
		i -= size
		if _, err := m.TotalWeight.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGroup(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *InternalGaugeRecord) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InternalGaugeRecord) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *InternalGaugeRecord) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.CumulativeWeight.Size()
		i -= size
		if _, err := m.CumulativeWeight.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGroup(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.CurrentWeight.Size()
		i -= size
		if _, err := m.CurrentWeight.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGroup(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.GaugeId != 0 {
		i = encodeVarintGroup(dAtA, i, uint64(m.GaugeId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Group) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Group) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Group) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SplittingPolicy != 0 {
		i = encodeVarintGroup(dAtA, i, uint64(m.SplittingPolicy))
		i--
		dAtA[i] = 0x18
	}
	{
		size, err := m.InternalGaugeInfo.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGroup(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.GroupGaugeId != 0 {
		i = encodeVarintGroup(dAtA, i, uint64(m.GroupGaugeId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *CreateGroup) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CreateGroup) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CreateGroup) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PoolIds) > 0 {
		dAtA3 := make([]byte, len(m.PoolIds)*10)
		var j2 int
		for _, num := range m.PoolIds {
			for num >= 1<<7 {
				dAtA3[j2] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j2++
			}
			dAtA3[j2] = uint8(num)
			j2++
		}
		i -= j2
		copy(dAtA[i:], dAtA3[:j2])
		i = encodeVarintGroup(dAtA, i, uint64(j2))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *GroupsWithGauge) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GroupsWithGauge) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GroupsWithGauge) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Gauge.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGroup(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.Group.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGroup(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGroup(dAtA []byte, offset int, v uint64) int {
	offset -= sovGroup(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *InternalGaugeInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.TotalWeight.Size()
	n += 1 + l + sovGroup(uint64(l))
	if len(m.GaugeRecords) > 0 {
		for _, e := range m.GaugeRecords {
			l = e.Size()
			n += 1 + l + sovGroup(uint64(l))
		}
	}
	return n
}

func (m *InternalGaugeRecord) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.GaugeId != 0 {
		n += 1 + sovGroup(uint64(m.GaugeId))
	}
	l = m.CurrentWeight.Size()
	n += 1 + l + sovGroup(uint64(l))
	l = m.CumulativeWeight.Size()
	n += 1 + l + sovGroup(uint64(l))
	return n
}

func (m *Group) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.GroupGaugeId != 0 {
		n += 1 + sovGroup(uint64(m.GroupGaugeId))
	}
	l = m.InternalGaugeInfo.Size()
	n += 1 + l + sovGroup(uint64(l))
	if m.SplittingPolicy != 0 {
		n += 1 + sovGroup(uint64(m.SplittingPolicy))
	}
	return n
}

func (m *CreateGroup) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.PoolIds) > 0 {
		l = 0
		for _, e := range m.PoolIds {
			l += sovGroup(uint64(e))
		}
		n += 1 + sovGroup(uint64(l)) + l
	}
	return n
}

func (m *GroupsWithGauge) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Group.Size()
	n += 1 + l + sovGroup(uint64(l))
	l = m.Gauge.Size()
	n += 1 + l + sovGroup(uint64(l))
	return n
}

func sovGroup(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGroup(x uint64) (n int) {
	return sovGroup(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *InternalGaugeInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGroup
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
			return fmt.Errorf("proto: InternalGaugeInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InternalGaugeInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalWeight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGroup
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
				return ErrInvalidLengthGroup
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGroup
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TotalWeight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GaugeRecords", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGroup
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
				return ErrInvalidLengthGroup
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGroup
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GaugeRecords = append(m.GaugeRecords, InternalGaugeRecord{})
			if err := m.GaugeRecords[len(m.GaugeRecords)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGroup(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGroup
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
func (m *InternalGaugeRecord) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGroup
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
			return fmt.Errorf("proto: InternalGaugeRecord: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InternalGaugeRecord: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GaugeId", wireType)
			}
			m.GaugeId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGroup
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GaugeId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentWeight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGroup
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
				return ErrInvalidLengthGroup
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGroup
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CurrentWeight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CumulativeWeight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGroup
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
				return ErrInvalidLengthGroup
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGroup
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CumulativeWeight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGroup(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGroup
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
func (m *Group) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGroup
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
			return fmt.Errorf("proto: Group: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Group: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupGaugeId", wireType)
			}
			m.GroupGaugeId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGroup
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GroupGaugeId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InternalGaugeInfo", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGroup
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
				return ErrInvalidLengthGroup
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGroup
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.InternalGaugeInfo.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SplittingPolicy", wireType)
			}
			m.SplittingPolicy = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGroup
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SplittingPolicy |= SplittingPolicy(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGroup(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGroup
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
func (m *CreateGroup) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGroup
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
			return fmt.Errorf("proto: CreateGroup: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CreateGroup: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGroup
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.PoolIds = append(m.PoolIds, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGroup
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthGroup
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthGroup
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.PoolIds) == 0 {
					m.PoolIds = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGroup
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.PoolIds = append(m.PoolIds, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolIds", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGroup(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGroup
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
func (m *GroupsWithGauge) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGroup
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
			return fmt.Errorf("proto: GroupsWithGauge: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GroupsWithGauge: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Group", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGroup
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
				return ErrInvalidLengthGroup
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGroup
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Group.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Gauge", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGroup
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
				return ErrInvalidLengthGroup
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGroup
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Gauge.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGroup(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGroup
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
func skipGroup(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGroup
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
					return 0, ErrIntOverflowGroup
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
					return 0, ErrIntOverflowGroup
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
				return 0, ErrInvalidLengthGroup
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGroup
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGroup
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGroup        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGroup          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGroup = fmt.Errorf("proto: unexpected end of group")
)
