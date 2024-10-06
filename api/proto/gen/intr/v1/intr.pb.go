// interactive.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: intr/v1/intr.proto

package intrv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Like struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID        uint64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	LikeCount int64  `protobuf:"varint,2,opt,name=LikeCount,proto3" json:"LikeCount,omitempty"`
	Ctime     int64  `protobuf:"varint,3,opt,name=Ctime,proto3" json:"Ctime,omitempty"`
}

func (x *Like) Reset() {
	*x = Like{}
	if protoimpl.UnsafeEnabled {
		mi := &file_intr_v1_intr_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Like) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Like) ProtoMessage() {}

func (x *Like) ProtoReflect() protoreflect.Message {
	mi := &file_intr_v1_intr_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Like.ProtoReflect.Descriptor instead.
func (*Like) Descriptor() ([]byte, []int) {
	return file_intr_v1_intr_proto_rawDescGZIP(), []int{0}
}

func (x *Like) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *Like) GetLikeCount() int64 {
	if x != nil {
		return x.LikeCount
	}
	return 0
}

func (x *Like) GetCtime() int64 {
	if x != nil {
		return x.Ctime
	}
	return 0
}

type IncrReadCountRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ObjType string `protobuf:"bytes,1,opt,name=ObjType,proto3" json:"ObjType,omitempty"`
	ObjID   uint64 `protobuf:"varint,2,opt,name=ObjID,proto3" json:"ObjID,omitempty"`
}

func (x *IncrReadCountRequest) Reset() {
	*x = IncrReadCountRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_intr_v1_intr_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IncrReadCountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IncrReadCountRequest) ProtoMessage() {}

func (x *IncrReadCountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_intr_v1_intr_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IncrReadCountRequest.ProtoReflect.Descriptor instead.
func (*IncrReadCountRequest) Descriptor() ([]byte, []int) {
	return file_intr_v1_intr_proto_rawDescGZIP(), []int{1}
}

func (x *IncrReadCountRequest) GetObjType() string {
	if x != nil {
		return x.ObjType
	}
	return ""
}

func (x *IncrReadCountRequest) GetObjID() uint64 {
	if x != nil {
		return x.ObjID
	}
	return 0
}

type IncrReadCountResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *IncrReadCountResponse) Reset() {
	*x = IncrReadCountResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_intr_v1_intr_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IncrReadCountResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IncrReadCountResponse) ProtoMessage() {}

func (x *IncrReadCountResponse) ProtoReflect() protoreflect.Message {
	mi := &file_intr_v1_intr_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IncrReadCountResponse.ProtoReflect.Descriptor instead.
func (*IncrReadCountResponse) Descriptor() ([]byte, []int) {
	return file_intr_v1_intr_proto_rawDescGZIP(), []int{2}
}

type LikeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ObjType string `protobuf:"bytes,1,opt,name=ObjType,proto3" json:"ObjType,omitempty"`
	ObjID   uint64 `protobuf:"varint,2,opt,name=ObjID,proto3" json:"ObjID,omitempty"`
	Uid     uint64 `protobuf:"varint,3,opt,name=uid,proto3" json:"uid,omitempty"`
}

func (x *LikeRequest) Reset() {
	*x = LikeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_intr_v1_intr_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LikeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LikeRequest) ProtoMessage() {}

func (x *LikeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_intr_v1_intr_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LikeRequest.ProtoReflect.Descriptor instead.
func (*LikeRequest) Descriptor() ([]byte, []int) {
	return file_intr_v1_intr_proto_rawDescGZIP(), []int{3}
}

func (x *LikeRequest) GetObjType() string {
	if x != nil {
		return x.ObjType
	}
	return ""
}

func (x *LikeRequest) GetObjID() uint64 {
	if x != nil {
		return x.ObjID
	}
	return 0
}

func (x *LikeRequest) GetUid() uint64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

type LikeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LikeResponse) Reset() {
	*x = LikeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_intr_v1_intr_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LikeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LikeResponse) ProtoMessage() {}

func (x *LikeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_intr_v1_intr_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LikeResponse.ProtoReflect.Descriptor instead.
func (*LikeResponse) Descriptor() ([]byte, []int) {
	return file_intr_v1_intr_proto_rawDescGZIP(), []int{4}
}

type UnLikeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ObjType string `protobuf:"bytes,1,opt,name=ObjType,proto3" json:"ObjType,omitempty"`
	ObjID   uint64 `protobuf:"varint,2,opt,name=ObjID,proto3" json:"ObjID,omitempty"`
	Uid     uint64 `protobuf:"varint,3,opt,name=uid,proto3" json:"uid,omitempty"`
}

func (x *UnLikeRequest) Reset() {
	*x = UnLikeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_intr_v1_intr_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnLikeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnLikeRequest) ProtoMessage() {}

func (x *UnLikeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_intr_v1_intr_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnLikeRequest.ProtoReflect.Descriptor instead.
func (*UnLikeRequest) Descriptor() ([]byte, []int) {
	return file_intr_v1_intr_proto_rawDescGZIP(), []int{5}
}

func (x *UnLikeRequest) GetObjType() string {
	if x != nil {
		return x.ObjType
	}
	return ""
}

func (x *UnLikeRequest) GetObjID() uint64 {
	if x != nil {
		return x.ObjID
	}
	return 0
}

func (x *UnLikeRequest) GetUid() uint64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

type UnLikeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UnLikeResponse) Reset() {
	*x = UnLikeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_intr_v1_intr_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnLikeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnLikeResponse) ProtoMessage() {}

func (x *UnLikeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_intr_v1_intr_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnLikeResponse.ProtoReflect.Descriptor instead.
func (*UnLikeResponse) Descriptor() ([]byte, []int) {
	return file_intr_v1_intr_proto_rawDescGZIP(), []int{6}
}

type CollectRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ObjType string `protobuf:"bytes,1,opt,name=ObjType,proto3" json:"ObjType,omitempty"`
	ObjID   uint64 `protobuf:"varint,2,opt,name=ObjID,proto3" json:"ObjID,omitempty"`
	Uid     uint64 `protobuf:"varint,3,opt,name=uid,proto3" json:"uid,omitempty"`
}

func (x *CollectRequest) Reset() {
	*x = CollectRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_intr_v1_intr_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CollectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CollectRequest) ProtoMessage() {}

func (x *CollectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_intr_v1_intr_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CollectRequest.ProtoReflect.Descriptor instead.
func (*CollectRequest) Descriptor() ([]byte, []int) {
	return file_intr_v1_intr_proto_rawDescGZIP(), []int{7}
}

func (x *CollectRequest) GetObjType() string {
	if x != nil {
		return x.ObjType
	}
	return ""
}

func (x *CollectRequest) GetObjID() uint64 {
	if x != nil {
		return x.ObjID
	}
	return 0
}

func (x *CollectRequest) GetUid() uint64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

type CollectResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CollectResponse) Reset() {
	*x = CollectResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_intr_v1_intr_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CollectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CollectResponse) ProtoMessage() {}

func (x *CollectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_intr_v1_intr_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CollectResponse.ProtoReflect.Descriptor instead.
func (*CollectResponse) Descriptor() ([]byte, []int) {
	return file_intr_v1_intr_proto_rawDescGZIP(), []int{8}
}

type GetListBatchOfLikesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Limit     int32  `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
	OrderBy   string `protobuf:"bytes,2,opt,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	Desc      bool   `protobuf:"varint,3,opt,name=desc,proto3" json:"desc,omitempty"`
	LastValue int64  `protobuf:"varint,4,opt,name=lastValue,proto3" json:"lastValue,omitempty"`
	ObjType   string `protobuf:"bytes,5,opt,name=ObjType,proto3" json:"ObjType,omitempty"`
}

func (x *GetListBatchOfLikesRequest) Reset() {
	*x = GetListBatchOfLikesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_intr_v1_intr_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetListBatchOfLikesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetListBatchOfLikesRequest) ProtoMessage() {}

func (x *GetListBatchOfLikesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_intr_v1_intr_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetListBatchOfLikesRequest.ProtoReflect.Descriptor instead.
func (*GetListBatchOfLikesRequest) Descriptor() ([]byte, []int) {
	return file_intr_v1_intr_proto_rawDescGZIP(), []int{9}
}

func (x *GetListBatchOfLikesRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *GetListBatchOfLikesRequest) GetOrderBy() string {
	if x != nil {
		return x.OrderBy
	}
	return ""
}

func (x *GetListBatchOfLikesRequest) GetDesc() bool {
	if x != nil {
		return x.Desc
	}
	return false
}

func (x *GetListBatchOfLikesRequest) GetLastValue() int64 {
	if x != nil {
		return x.LastValue
	}
	return 0
}

func (x *GetListBatchOfLikesRequest) GetObjType() string {
	if x != nil {
		return x.ObjType
	}
	return ""
}

type GetListBatchOfLikesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*Like `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *GetListBatchOfLikesResponse) Reset() {
	*x = GetListBatchOfLikesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_intr_v1_intr_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetListBatchOfLikesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetListBatchOfLikesResponse) ProtoMessage() {}

func (x *GetListBatchOfLikesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_intr_v1_intr_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetListBatchOfLikesResponse.ProtoReflect.Descriptor instead.
func (*GetListBatchOfLikesResponse) Descriptor() ([]byte, []int) {
	return file_intr_v1_intr_proto_rawDescGZIP(), []int{10}
}

func (x *GetListBatchOfLikesResponse) GetData() []*Like {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_intr_v1_intr_proto protoreflect.FileDescriptor

var file_intr_v1_intr_proto_rawDesc = []byte{
	0x0a, 0x12, 0x69, 0x6e, 0x74, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x74, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x69, 0x6e, 0x74, 0x72, 0x2e, 0x76, 0x31, 0x22, 0x4a, 0x0a,
	0x04, 0x4c, 0x69, 0x6b, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x4c, 0x69, 0x6b, 0x65, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x4c, 0x69, 0x6b, 0x65, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x43, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x43, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x46, 0x0a, 0x14, 0x49, 0x6e, 0x63,
	0x72, 0x52, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x18, 0x0a, 0x07, 0x4f, 0x62, 0x6a, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x4f, 0x62, 0x6a, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x4f,
	0x62, 0x6a, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x4f, 0x62, 0x6a, 0x49,
	0x44, 0x22, 0x17, 0x0a, 0x15, 0x49, 0x6e, 0x63, 0x72, 0x52, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x4f, 0x0a, 0x0b, 0x4c, 0x69,
	0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x4f, 0x62, 0x6a,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4f, 0x62, 0x6a, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x4f, 0x62, 0x6a, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x05, 0x4f, 0x62, 0x6a, 0x49, 0x44, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x75, 0x69, 0x64, 0x22, 0x0e, 0x0a, 0x0c, 0x4c,
	0x69, 0x6b, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x51, 0x0a, 0x0d, 0x55,
	0x6e, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x4f, 0x62, 0x6a, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4f,
	0x62, 0x6a, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x4f, 0x62, 0x6a, 0x49, 0x44, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x4f, 0x62, 0x6a, 0x49, 0x44, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x75, 0x69, 0x64, 0x22, 0x10,
	0x0a, 0x0e, 0x55, 0x6e, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x52, 0x0a, 0x0e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x4f, 0x62, 0x6a, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x4f, 0x62, 0x6a, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x4f, 0x62, 0x6a, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x4f, 0x62, 0x6a,
	0x49, 0x44, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x03, 0x75, 0x69, 0x64, 0x22, 0x11, 0x0a, 0x0f, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x99, 0x01, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x4c,
	0x69, 0x73, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x4f, 0x66, 0x4c, 0x69, 0x6b, 0x65, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x19, 0x0a, 0x08,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x62, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x63, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63, 0x12, 0x1c, 0x0a, 0x09, 0x6c,
	0x61, 0x73, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x6c, 0x61, 0x73, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x4f, 0x62, 0x6a,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4f, 0x62, 0x6a, 0x54,
	0x79, 0x70, 0x65, 0x22, 0x40, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x4f, 0x66, 0x4c, 0x69, 0x6b, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x21, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0d, 0x2e, 0x69, 0x6e, 0x74, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x6b, 0x65, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0xed, 0x02, 0x0a, 0x0b, 0x4c, 0x69, 0x6b, 0x65, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4e, 0x0a, 0x0d, 0x49, 0x6e, 0x63, 0x72, 0x52, 0x65, 0x61,
	0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1d, 0x2e, 0x69, 0x6e, 0x74, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x49, 0x6e, 0x63, 0x72, 0x52, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x69, 0x6e, 0x74, 0x72, 0x2e, 0x76, 0x31, 0x2e,
	0x49, 0x6e, 0x63, 0x72, 0x52, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x33, 0x0a, 0x04, 0x4c, 0x69, 0x6b, 0x65, 0x12, 0x14, 0x2e,
	0x69, 0x6e, 0x74, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x69, 0x6e, 0x74, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69,
	0x6b, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x06, 0x55, 0x6e,
	0x4c, 0x69, 0x6b, 0x65, 0x12, 0x16, 0x2e, 0x69, 0x6e, 0x74, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55,
	0x6e, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x69,
	0x6e, 0x74, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x6e, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x07, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x12, 0x17, 0x2e, 0x69, 0x6e, 0x74, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x69, 0x6e, 0x74, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x60, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x4f, 0x66, 0x4c, 0x69, 0x6b, 0x65, 0x73, 0x12, 0x23, 0x2e, 0x69, 0x6e, 0x74,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x4f, 0x66, 0x4c, 0x69, 0x6b, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x24, 0x2e, 0x69, 0x6e, 0x74, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x73,
	0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x4f, 0x66, 0x4c, 0x69, 0x6b, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x78, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x6e, 0x74,
	0x72, 0x2e, 0x76, 0x31, 0x42, 0x09, 0x49, 0x6e, 0x74, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x21, 0x70, 0x6f, 0x73, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x69, 0x6e, 0x74, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x69, 0x6e,
	0x74, 0x72, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x58, 0x58, 0xaa, 0x02, 0x07, 0x49, 0x6e, 0x74,
	0x72, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x07, 0x49, 0x6e, 0x74, 0x72, 0x5c, 0x56, 0x31, 0xe2, 0x02,
	0x13, 0x49, 0x6e, 0x74, 0x72, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x08, 0x49, 0x6e, 0x74, 0x72, 0x3a, 0x3a, 0x56, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_intr_v1_intr_proto_rawDescOnce sync.Once
	file_intr_v1_intr_proto_rawDescData = file_intr_v1_intr_proto_rawDesc
)

func file_intr_v1_intr_proto_rawDescGZIP() []byte {
	file_intr_v1_intr_proto_rawDescOnce.Do(func() {
		file_intr_v1_intr_proto_rawDescData = protoimpl.X.CompressGZIP(file_intr_v1_intr_proto_rawDescData)
	})
	return file_intr_v1_intr_proto_rawDescData
}

var file_intr_v1_intr_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_intr_v1_intr_proto_goTypes = []any{
	(*Like)(nil),                        // 0: intr.v1.Like
	(*IncrReadCountRequest)(nil),        // 1: intr.v1.IncrReadCountRequest
	(*IncrReadCountResponse)(nil),       // 2: intr.v1.IncrReadCountResponse
	(*LikeRequest)(nil),                 // 3: intr.v1.LikeRequest
	(*LikeResponse)(nil),                // 4: intr.v1.LikeResponse
	(*UnLikeRequest)(nil),               // 5: intr.v1.UnLikeRequest
	(*UnLikeResponse)(nil),              // 6: intr.v1.UnLikeResponse
	(*CollectRequest)(nil),              // 7: intr.v1.CollectRequest
	(*CollectResponse)(nil),             // 8: intr.v1.CollectResponse
	(*GetListBatchOfLikesRequest)(nil),  // 9: intr.v1.GetListBatchOfLikesRequest
	(*GetListBatchOfLikesResponse)(nil), // 10: intr.v1.GetListBatchOfLikesResponse
}
var file_intr_v1_intr_proto_depIdxs = []int32{
	0,  // 0: intr.v1.GetListBatchOfLikesResponse.data:type_name -> intr.v1.Like
	1,  // 1: intr.v1.LikeService.IncrReadCount:input_type -> intr.v1.IncrReadCountRequest
	3,  // 2: intr.v1.LikeService.Like:input_type -> intr.v1.LikeRequest
	5,  // 3: intr.v1.LikeService.UnLike:input_type -> intr.v1.UnLikeRequest
	7,  // 4: intr.v1.LikeService.Collect:input_type -> intr.v1.CollectRequest
	9,  // 5: intr.v1.LikeService.GetListBatchOfLikes:input_type -> intr.v1.GetListBatchOfLikesRequest
	2,  // 6: intr.v1.LikeService.IncrReadCount:output_type -> intr.v1.IncrReadCountResponse
	4,  // 7: intr.v1.LikeService.Like:output_type -> intr.v1.LikeResponse
	6,  // 8: intr.v1.LikeService.UnLike:output_type -> intr.v1.UnLikeResponse
	8,  // 9: intr.v1.LikeService.Collect:output_type -> intr.v1.CollectResponse
	10, // 10: intr.v1.LikeService.GetListBatchOfLikes:output_type -> intr.v1.GetListBatchOfLikesResponse
	6,  // [6:11] is the sub-list for method output_type
	1,  // [1:6] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_intr_v1_intr_proto_init() }
func file_intr_v1_intr_proto_init() {
	if File_intr_v1_intr_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_intr_v1_intr_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Like); i {
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
		file_intr_v1_intr_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*IncrReadCountRequest); i {
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
		file_intr_v1_intr_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*IncrReadCountResponse); i {
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
		file_intr_v1_intr_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*LikeRequest); i {
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
		file_intr_v1_intr_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*LikeResponse); i {
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
		file_intr_v1_intr_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*UnLikeRequest); i {
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
		file_intr_v1_intr_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*UnLikeResponse); i {
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
		file_intr_v1_intr_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*CollectRequest); i {
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
		file_intr_v1_intr_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*CollectResponse); i {
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
		file_intr_v1_intr_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*GetListBatchOfLikesRequest); i {
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
		file_intr_v1_intr_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*GetListBatchOfLikesResponse); i {
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
			RawDescriptor: file_intr_v1_intr_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_intr_v1_intr_proto_goTypes,
		DependencyIndexes: file_intr_v1_intr_proto_depIdxs,
		MessageInfos:      file_intr_v1_intr_proto_msgTypes,
	}.Build()
	File_intr_v1_intr_proto = out.File
	file_intr_v1_intr_proto_rawDesc = nil
	file_intr_v1_intr_proto_goTypes = nil
	file_intr_v1_intr_proto_depIdxs = nil
}
