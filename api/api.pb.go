// https://developers.google.com/protocol-buffers/docs/proto3#nested
// https://developers.google.com/protocol-buffers/docs/overview

//
//
// protoc --go_out=. --go_opt=paths=source_relative \
//--go-grpc_out=. --go-grpc_opt=paths=source_relative \
//api/api.proto
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: api/api.proto

package collective

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

type DataUpdates struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DataUpdates) Reset() {
	*x = DataUpdates{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataUpdates) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataUpdates) ProtoMessage() {}

func (x *DataUpdates) ProtoReflect() protoreflect.Message {
	mi := &file_api_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataUpdates.ProtoReflect.Descriptor instead.
func (*DataUpdates) Descriptor() ([]byte, []int) {
	return file_api_api_proto_rawDescGZIP(), []int{0}
}

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key      string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Database string `protobuf:"bytes,2,opt,name=database,proto3" json:"database,omitempty"`
	Data     []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_api_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_api_api_proto_rawDescGZIP(), []int{1}
}

func (x *Data) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Data) GetDatabase() string {
	if x != nil {
		return x.Database
	}
	return ""
}

func (x *Data) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type Updated struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UpdatedSuccessfully bool `protobuf:"varint,1,opt,name=updatedSuccessfully,proto3" json:"updatedSuccessfully,omitempty"`
}

func (x *Updated) Reset() {
	*x = Updated{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Updated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Updated) ProtoMessage() {}

func (x *Updated) ProtoReflect() protoreflect.Message {
	mi := &file_api_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Updated.ProtoReflect.Descriptor instead.
func (*Updated) Descriptor() ([]byte, []int) {
	return file_api_api_proto_rawDescGZIP(), []int{2}
}

func (x *Updated) GetUpdatedSuccessfully() bool {
	if x != nil {
		return x.UpdatedSuccessfully
	}
	return false
}

type SyncIp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IpAddress string `protobuf:"bytes,1,opt,name=ipAddress,proto3" json:"ipAddress,omitempty"`
}

func (x *SyncIp) Reset() {
	*x = SyncIp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncIp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncIp) ProtoMessage() {}

func (x *SyncIp) ProtoReflect() protoreflect.Message {
	mi := &file_api_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncIp.ProtoReflect.Descriptor instead.
func (*SyncIp) Descriptor() ([]byte, []int) {
	return file_api_api_proto_rawDescGZIP(), []int{3}
}

func (x *SyncIp) GetIpAddress() string {
	if x != nil {
		return x.IpAddress
	}
	return ""
}

type DataUpdates_CollectiveDataUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Update     bool  `protobuf:"varint,1,opt,name=update,proto3" json:"update,omitempty"`
	UpdateType int32 `protobuf:"varint,2,opt,name=updateType,proto3" json:"updateType,omitempty"`
}

func (x *DataUpdates_CollectiveDataUpdate) Reset() {
	*x = DataUpdates_CollectiveDataUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataUpdates_CollectiveDataUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataUpdates_CollectiveDataUpdate) ProtoMessage() {}

func (x *DataUpdates_CollectiveDataUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_api_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataUpdates_CollectiveDataUpdate.ProtoReflect.Descriptor instead.
func (*DataUpdates_CollectiveDataUpdate) Descriptor() ([]byte, []int) {
	return file_api_api_proto_rawDescGZIP(), []int{0, 0}
}

func (x *DataUpdates_CollectiveDataUpdate) GetUpdate() bool {
	if x != nil {
		return x.Update
	}
	return false
}

func (x *DataUpdates_CollectiveDataUpdate) GetUpdateType() int32 {
	if x != nil {
		return x.UpdateType
	}
	return 0
}

type DataUpdates_CollectiveReplicaUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Update     bool  `protobuf:"varint,1,opt,name=update,proto3" json:"update,omitempty"`
	UpdateType int32 `protobuf:"varint,2,opt,name=updateType,proto3" json:"updateType,omitempty"`
}

func (x *DataUpdates_CollectiveReplicaUpdate) Reset() {
	*x = DataUpdates_CollectiveReplicaUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataUpdates_CollectiveReplicaUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataUpdates_CollectiveReplicaUpdate) ProtoMessage() {}

func (x *DataUpdates_CollectiveReplicaUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_api_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataUpdates_CollectiveReplicaUpdate.ProtoReflect.Descriptor instead.
func (*DataUpdates_CollectiveReplicaUpdate) Descriptor() ([]byte, []int) {
	return file_api_api_proto_rawDescGZIP(), []int{0, 1}
}

func (x *DataUpdates_CollectiveReplicaUpdate) GetUpdate() bool {
	if x != nil {
		return x.Update
	}
	return false
}

func (x *DataUpdates_CollectiveReplicaUpdate) GetUpdateType() int32 {
	if x != nil {
		return x.UpdateType
	}
	return 0
}

type DataUpdates_CollectiveDataUpdate_Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReplicaNodeGroup  int32    `protobuf:"varint,1,opt,name=replicaNodeGroup,proto3" json:"replicaNodeGroup,omitempty"`
	DataKey           string   `protobuf:"bytes,2,opt,name=dataKey,proto3" json:"dataKey,omitempty"`
	Database          string   `protobuf:"bytes,3,opt,name=database,proto3" json:"database,omitempty"`
	ReplicatedNodeIds []string `protobuf:"bytes,4,rep,name=replicatedNodeIds,proto3" json:"replicatedNodeIds,omitempty"`
}

func (x *DataUpdates_CollectiveDataUpdate_Data) Reset() {
	*x = DataUpdates_CollectiveDataUpdate_Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataUpdates_CollectiveDataUpdate_Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataUpdates_CollectiveDataUpdate_Data) ProtoMessage() {}

func (x *DataUpdates_CollectiveDataUpdate_Data) ProtoReflect() protoreflect.Message {
	mi := &file_api_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataUpdates_CollectiveDataUpdate_Data.ProtoReflect.Descriptor instead.
func (*DataUpdates_CollectiveDataUpdate_Data) Descriptor() ([]byte, []int) {
	return file_api_api_proto_rawDescGZIP(), []int{0, 0, 0}
}

func (x *DataUpdates_CollectiveDataUpdate_Data) GetReplicaNodeGroup() int32 {
	if x != nil {
		return x.ReplicaNodeGroup
	}
	return 0
}

func (x *DataUpdates_CollectiveDataUpdate_Data) GetDataKey() string {
	if x != nil {
		return x.DataKey
	}
	return ""
}

func (x *DataUpdates_CollectiveDataUpdate_Data) GetDatabase() string {
	if x != nil {
		return x.Database
	}
	return ""
}

func (x *DataUpdates_CollectiveDataUpdate_Data) GetReplicatedNodeIds() []string {
	if x != nil {
		return x.ReplicatedNodeIds
	}
	return nil
}

type DataUpdates_CollectiveReplicaUpdate_UpdateReplica struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReplicaNodeGroup int32 `protobuf:"varint,1,opt,name=replicaNodeGroup,proto3" json:"replicaNodeGroup,omitempty"`
	FullGroup        bool  `protobuf:"varint,2,opt,name=fullGroup,proto3" json:"fullGroup,omitempty"`
}

func (x *DataUpdates_CollectiveReplicaUpdate_UpdateReplica) Reset() {
	*x = DataUpdates_CollectiveReplicaUpdate_UpdateReplica{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_api_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataUpdates_CollectiveReplicaUpdate_UpdateReplica) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataUpdates_CollectiveReplicaUpdate_UpdateReplica) ProtoMessage() {}

func (x *DataUpdates_CollectiveReplicaUpdate_UpdateReplica) ProtoReflect() protoreflect.Message {
	mi := &file_api_api_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataUpdates_CollectiveReplicaUpdate_UpdateReplica.ProtoReflect.Descriptor instead.
func (*DataUpdates_CollectiveReplicaUpdate_UpdateReplica) Descriptor() ([]byte, []int) {
	return file_api_api_proto_rawDescGZIP(), []int{0, 1, 0}
}

func (x *DataUpdates_CollectiveReplicaUpdate_UpdateReplica) GetReplicaNodeGroup() int32 {
	if x != nil {
		return x.ReplicaNodeGroup
	}
	return 0
}

func (x *DataUpdates_CollectiveReplicaUpdate_UpdateReplica) GetFullGroup() bool {
	if x != nil {
		return x.FullGroup
	}
	return false
}

type DataUpdates_CollectiveReplicaUpdate_UpdateReplicaReplicaNodes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RodeId    string `protobuf:"bytes,1,opt,name=rodeId,proto3" json:"rodeId,omitempty"`
	IpAddress string `protobuf:"bytes,2,opt,name=ipAddress,proto3" json:"ipAddress,omitempty"`
}

func (x *DataUpdates_CollectiveReplicaUpdate_UpdateReplicaReplicaNodes) Reset() {
	*x = DataUpdates_CollectiveReplicaUpdate_UpdateReplicaReplicaNodes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_api_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataUpdates_CollectiveReplicaUpdate_UpdateReplicaReplicaNodes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataUpdates_CollectiveReplicaUpdate_UpdateReplicaReplicaNodes) ProtoMessage() {}

func (x *DataUpdates_CollectiveReplicaUpdate_UpdateReplicaReplicaNodes) ProtoReflect() protoreflect.Message {
	mi := &file_api_api_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataUpdates_CollectiveReplicaUpdate_UpdateReplicaReplicaNodes.ProtoReflect.Descriptor instead.
func (*DataUpdates_CollectiveReplicaUpdate_UpdateReplicaReplicaNodes) Descriptor() ([]byte, []int) {
	return file_api_api_proto_rawDescGZIP(), []int{0, 1, 0, 0}
}

func (x *DataUpdates_CollectiveReplicaUpdate_UpdateReplicaReplicaNodes) GetRodeId() string {
	if x != nil {
		return x.RodeId
	}
	return ""
}

func (x *DataUpdates_CollectiveReplicaUpdate_UpdateReplicaReplicaNodes) GetIpAddress() string {
	if x != nil {
		return x.IpAddress
	}
	return ""
}

var File_api_api_proto protoreflect.FileDescriptor

var file_api_api_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x04, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0xed, 0x03, 0x0a, 0x0b, 0x44, 0x61, 0x74, 0x61, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x73, 0x1a, 0xe7, 0x01, 0x0a, 0x14, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63,
	0x74, 0x69, 0x76, 0x65, 0x44, 0x61, 0x74, 0x61, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x1a, 0x96, 0x01, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x2a, 0x0a, 0x10, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x4e, 0x6f, 0x64, 0x65, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x72, 0x65, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x4e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x64,
	0x61, 0x74, 0x61, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x61,
	0x74, 0x61, 0x4b, 0x65, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73,
	0x65, 0x12, 0x2c, 0x0a, 0x11, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0x4e,
	0x6f, 0x64, 0x65, 0x49, 0x64, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x11, 0x72, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x73, 0x1a,
	0xf3, 0x01, 0x0a, 0x17, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x1a, 0x9f, 0x01, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x12, 0x2a, 0x0a, 0x10, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x4e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x10, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x4e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x75, 0x6c, 0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x66, 0x75, 0x6c, 0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x1a,
	0x44, 0x0a, 0x0c, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12,
	0x16, 0x0a, 0x06, 0x72, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x72, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x70, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x70, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x48, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x3b, 0x0a, 0x07, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x12, 0x30, 0x0a, 0x13, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x66, 0x75, 0x6c, 0x6c,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x13, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x66, 0x75, 0x6c, 0x6c, 0x79, 0x22, 0x26, 0x0a, 0x06,
	0x53, 0x79, 0x6e, 0x63, 0x49, 0x70, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x70, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x70, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x32, 0xe7, 0x02, 0x0a, 0x0a, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x47, 0x75,
	0x69, 0x64, 0x65, 0x12, 0x37, 0x0a, 0x0d, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x44, 0x61, 0x74, 0x61,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x1a, 0x0d, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x2f, 0x0a, 0x0f,
	0x53, 0x79, 0x6e, 0x63, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0c, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x49, 0x70, 0x1a, 0x0a, 0x2e,
	0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x22, 0x00, 0x30, 0x01, 0x12, 0x3a, 0x0a,
	0x10, 0x44, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x72, 0x79, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x12, 0x11, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x73, 0x1a, 0x0d, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x2d, 0x0a, 0x0a, 0x44, 0x61, 0x74,
	0x61, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x0a, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x44,
	0x61, 0x74, 0x61, 0x1a, 0x0d, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x34, 0x0a, 0x11, 0x52, 0x65, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x44, 0x61, 0x74, 0x61, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x0a, 0x2e,
	0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x0d, 0x2e, 0x6d, 0x61, 0x69, 0x6e,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x23,
	0x0a, 0x07, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x0a, 0x2e, 0x6d, 0x61, 0x69, 0x6e,
	0x2e, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x0a, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x44, 0x61, 0x74,
	0x61, 0x22, 0x00, 0x12, 0x29, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x0a, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x0d, 0x2e,
	0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x22, 0x00, 0x42, 0x0e,
	0x5a, 0x0c, 0x2e, 0x2f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_api_proto_rawDescOnce sync.Once
	file_api_api_proto_rawDescData = file_api_api_proto_rawDesc
)

func file_api_api_proto_rawDescGZIP() []byte {
	file_api_api_proto_rawDescOnce.Do(func() {
		file_api_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_api_proto_rawDescData)
	})
	return file_api_api_proto_rawDescData
}

var file_api_api_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_api_api_proto_goTypes = []interface{}{
	(*DataUpdates)(nil),                                       // 0: main.DataUpdates
	(*Data)(nil),                                              // 1: main.Data
	(*Updated)(nil),                                           // 2: main.Updated
	(*SyncIp)(nil),                                            // 3: main.SyncIp
	(*DataUpdates_CollectiveDataUpdate)(nil),                  // 4: main.DataUpdates.CollectiveDataUpdate
	(*DataUpdates_CollectiveReplicaUpdate)(nil),               // 5: main.DataUpdates.CollectiveReplicaUpdate
	(*DataUpdates_CollectiveDataUpdate_Data)(nil),             // 6: main.DataUpdates.CollectiveDataUpdate.Data
	(*DataUpdates_CollectiveReplicaUpdate_UpdateReplica)(nil), // 7: main.DataUpdates.CollectiveReplicaUpdate.UpdateReplica
	(*DataUpdates_CollectiveReplicaUpdate_UpdateReplicaReplicaNodes)(nil), // 8: main.DataUpdates.CollectiveReplicaUpdate.UpdateReplica.replicaNodes
}
var file_api_api_proto_depIdxs = []int32{
	0, // 0: main.RouteGuide.ReplicaUpdate:input_type -> main.DataUpdates
	3, // 1: main.RouteGuide.SyncDataRequest:input_type -> main.SyncIp
	0, // 2: main.RouteGuide.DictionaryUpdate:input_type -> main.DataUpdates
	1, // 3: main.RouteGuide.DataUpdate:input_type -> main.Data
	1, // 4: main.RouteGuide.ReplicaDataUpdate:input_type -> main.Data
	1, // 5: main.RouteGuide.GetData:input_type -> main.Data
	1, // 6: main.RouteGuide.DeleteData:input_type -> main.Data
	2, // 7: main.RouteGuide.ReplicaUpdate:output_type -> main.Updated
	1, // 8: main.RouteGuide.SyncDataRequest:output_type -> main.Data
	2, // 9: main.RouteGuide.DictionaryUpdate:output_type -> main.Updated
	2, // 10: main.RouteGuide.DataUpdate:output_type -> main.Updated
	2, // 11: main.RouteGuide.ReplicaDataUpdate:output_type -> main.Updated
	1, // 12: main.RouteGuide.GetData:output_type -> main.Data
	2, // 13: main.RouteGuide.DeleteData:output_type -> main.Updated
	7, // [7:14] is the sub-list for method output_type
	0, // [0:7] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_api_proto_init() }
func file_api_api_proto_init() {
	if File_api_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataUpdates); i {
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
		file_api_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data); i {
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
		file_api_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Updated); i {
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
		file_api_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncIp); i {
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
		file_api_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataUpdates_CollectiveDataUpdate); i {
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
		file_api_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataUpdates_CollectiveReplicaUpdate); i {
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
		file_api_api_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataUpdates_CollectiveDataUpdate_Data); i {
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
		file_api_api_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataUpdates_CollectiveReplicaUpdate_UpdateReplica); i {
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
		file_api_api_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataUpdates_CollectiveReplicaUpdate_UpdateReplicaReplicaNodes); i {
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
			RawDescriptor: file_api_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_api_proto_goTypes,
		DependencyIndexes: file_api_api_proto_depIdxs,
		MessageInfos:      file_api_api_proto_msgTypes,
	}.Build()
	File_api_api_proto = out.File
	file_api_api_proto_rawDesc = nil
	file_api_api_proto_goTypes = nil
	file_api_api_proto_depIdxs = nil
}
