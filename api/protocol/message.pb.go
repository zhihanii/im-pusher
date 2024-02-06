// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.0--rc2
// source: api/protocol/message.proto

package protocol

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

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version   uint32 `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	MagicNum  uint32 `protobuf:"varint,2,opt,name=magic_num,json=magicNum,proto3" json:"magic_num,omitempty"`
	Operation uint32 `protobuf:"varint,3,opt,name=operation,proto3" json:"operation,omitempty"` //区分消息类型
	Sequence  uint32 `protobuf:"varint,4,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Data      []byte `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"` //消息内容
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protocol_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_api_protocol_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_api_protocol_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Message) GetMagicNum() uint32 {
	if x != nil {
		return x.MagicNum
	}
	return 0
}

func (x *Message) GetOperation() uint32 {
	if x != nil {
		return x.Operation
	}
	return 0
}

func (x *Message) GetSequence() uint32 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *Message) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type TransMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      uint32   `protobuf:"varint,1,opt,name=type,proto3" json:"type,omitempty"` //消息的发送类型: Channel Room Broadcast
	Priority  uint32   `protobuf:"varint,2,opt,name=priority,proto3" json:"priority,omitempty"`
	Server    string   `protobuf:"bytes,3,opt,name=server,proto3" json:"server,omitempty"`
	Room      string   `protobuf:"bytes,4,opt,name=room,proto3" json:"room,omitempty"`
	Keys      []string `protobuf:"bytes,5,rep,name=keys,proto3" json:"keys,omitempty"`
	Operation uint32   `protobuf:"varint,6,opt,name=operation,proto3" json:"operation,omitempty"`
	Sequence  uint32   `protobuf:"varint,7,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Data      []byte   `protobuf:"bytes,8,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *TransMessage) Reset() {
	*x = TransMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protocol_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransMessage) ProtoMessage() {}

func (x *TransMessage) ProtoReflect() protoreflect.Message {
	mi := &file_api_protocol_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransMessage.ProtoReflect.Descriptor instead.
func (*TransMessage) Descriptor() ([]byte, []int) {
	return file_api_protocol_message_proto_rawDescGZIP(), []int{1}
}

func (x *TransMessage) GetType() uint32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *TransMessage) GetPriority() uint32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

func (x *TransMessage) GetServer() string {
	if x != nil {
		return x.Server
	}
	return ""
}

func (x *TransMessage) GetRoom() string {
	if x != nil {
		return x.Room
	}
	return ""
}

func (x *TransMessage) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *TransMessage) GetOperation() uint32 {
	if x != nil {
		return x.Operation
	}
	return 0
}

func (x *TransMessage) GetSequence() uint32 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *TransMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type DispatcherMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Receivers []uint64 `protobuf:"varint,1,rep,packed,name=receivers,proto3" json:"receivers,omitempty"`
	Operation uint32   `protobuf:"varint,2,opt,name=operation,proto3" json:"operation,omitempty"` //区分消息类型
	Sequence  uint32   `protobuf:"varint,3,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Data      []byte   `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"` //消息内容
}

func (x *DispatcherMessage) Reset() {
	*x = DispatcherMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protocol_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DispatcherMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DispatcherMessage) ProtoMessage() {}

func (x *DispatcherMessage) ProtoReflect() protoreflect.Message {
	mi := &file_api_protocol_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DispatcherMessage.ProtoReflect.Descriptor instead.
func (*DispatcherMessage) Descriptor() ([]byte, []int) {
	return file_api_protocol_message_proto_rawDescGZIP(), []int{2}
}

func (x *DispatcherMessage) GetReceivers() []uint64 {
	if x != nil {
		return x.Receivers
	}
	return nil
}

func (x *DispatcherMessage) GetOperation() uint32 {
	if x != nil {
		return x.Operation
	}
	return 0
}

func (x *DispatcherMessage) GetSequence() uint32 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *DispatcherMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type ChatSendMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConversationId uint64 `protobuf:"varint,1,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
	From           uint64 `protobuf:"varint,2,opt,name=from,proto3" json:"from,omitempty"`
	To             uint64 `protobuf:"varint,3,opt,name=to,proto3" json:"to,omitempty"`
	Content        []byte `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *ChatSendMessage) Reset() {
	*x = ChatSendMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protocol_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatSendMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatSendMessage) ProtoMessage() {}

func (x *ChatSendMessage) ProtoReflect() protoreflect.Message {
	mi := &file_api_protocol_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatSendMessage.ProtoReflect.Descriptor instead.
func (*ChatSendMessage) Descriptor() ([]byte, []int) {
	return file_api_protocol_message_proto_rawDescGZIP(), []int{3}
}

func (x *ChatSendMessage) GetConversationId() uint64 {
	if x != nil {
		return x.ConversationId
	}
	return 0
}

func (x *ChatSendMessage) GetFrom() uint64 {
	if x != nil {
		return x.From
	}
	return 0
}

func (x *ChatSendMessage) GetTo() uint64 {
	if x != nil {
		return x.To
	}
	return 0
}

func (x *ChatSendMessage) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

type ChatReceiveMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConversationId uint64 `protobuf:"varint,1,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
	From           uint64 `protobuf:"varint,2,opt,name=from,proto3" json:"from,omitempty"`
	Content        []byte `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *ChatReceiveMessage) Reset() {
	*x = ChatReceiveMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protocol_message_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatReceiveMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatReceiveMessage) ProtoMessage() {}

func (x *ChatReceiveMessage) ProtoReflect() protoreflect.Message {
	mi := &file_api_protocol_message_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatReceiveMessage.ProtoReflect.Descriptor instead.
func (*ChatReceiveMessage) Descriptor() ([]byte, []int) {
	return file_api_protocol_message_proto_rawDescGZIP(), []int{4}
}

func (x *ChatReceiveMessage) GetConversationId() uint64 {
	if x != nil {
		return x.ConversationId
	}
	return 0
}

func (x *ChatReceiveMessage) GetFrom() uint64 {
	if x != nil {
		return x.From
	}
	return 0
}

func (x *ChatReceiveMessage) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

type GroupChatMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupId  uint64 `protobuf:"varint,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	SenderId uint64 `protobuf:"varint,2,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	Content  []byte `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *GroupChatMessage) Reset() {
	*x = GroupChatMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protocol_message_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupChatMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupChatMessage) ProtoMessage() {}

func (x *GroupChatMessage) ProtoReflect() protoreflect.Message {
	mi := &file_api_protocol_message_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupChatMessage.ProtoReflect.Descriptor instead.
func (*GroupChatMessage) Descriptor() ([]byte, []int) {
	return file_api_protocol_message_proto_rawDescGZIP(), []int{5}
}

func (x *GroupChatMessage) GetGroupId() uint64 {
	if x != nil {
		return x.GroupId
	}
	return 0
}

func (x *GroupChatMessage) GetSenderId() uint64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *GroupChatMessage) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

type ServiceMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content []byte `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *ServiceMessage) Reset() {
	*x = ServiceMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protocol_message_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceMessage) ProtoMessage() {}

func (x *ServiceMessage) ProtoReflect() protoreflect.Message {
	mi := &file_api_protocol_message_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceMessage.ProtoReflect.Descriptor instead.
func (*ServiceMessage) Descriptor() ([]byte, []int) {
	return file_api_protocol_message_proto_rawDescGZIP(), []int{6}
}

func (x *ServiceMessage) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

var File_api_protocol_message_proto protoreflect.FileDescriptor

var file_api_protocol_message_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x61, 0x70,
	0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x22, 0x8e, 0x01, 0x0a, 0x07, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x61, 0x67, 0x69, 0x63, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x08, 0x6d, 0x61, 0x67, 0x69, 0x63, 0x4e, 0x75, 0x6d, 0x12, 0x1c, 0x0a,
	0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x73,
	0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x73,
	0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xcc, 0x01, 0x0a, 0x0c,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6f, 0x6d, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x72, 0x6f, 0x6f, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12, 0x1c, 0x0a, 0x09,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65,
	0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x73, 0x65,
	0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x7f, 0x0a, 0x11, 0x44, 0x69,
	0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x04, 0x52, 0x09, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x73, 0x12, 0x1c, 0x0a,
	0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x73,
	0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x73,
	0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x78, 0x0a, 0x0f, 0x43,
	0x68, 0x61, 0x74, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x27,
	0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x74,
	0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x6b, 0x0a, 0x12, 0x43, 0x68, 0x61, 0x74, 0x52, 0x65, 0x63,
	0x65, 0x69, 0x76, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x63,
	0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x22, 0x64, 0x0a, 0x10, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x43, 0x68, 0x61, 0x74, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49,
	0x64, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x2a, 0x0a, 0x0e, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x7a, 0x68, 0x69, 0x68, 0x61, 0x6e, 0x69, 0x69, 0x2f, 0x69, 0x6d, 0x2d, 0x70,
	0x75, 0x73, 0x68, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_api_protocol_message_proto_rawDescOnce sync.Once
	file_api_protocol_message_proto_rawDescData = file_api_protocol_message_proto_rawDesc
)

func file_api_protocol_message_proto_rawDescGZIP() []byte {
	file_api_protocol_message_proto_rawDescOnce.Do(func() {
		file_api_protocol_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_protocol_message_proto_rawDescData)
	})
	return file_api_protocol_message_proto_rawDescData
}

var file_api_protocol_message_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_api_protocol_message_proto_goTypes = []interface{}{
	(*Message)(nil),            // 0: api.protocol.Message
	(*TransMessage)(nil),       // 1: api.protocol.TransMessage
	(*DispatcherMessage)(nil),  // 2: api.protocol.DispatcherMessage
	(*ChatSendMessage)(nil),    // 3: api.protocol.ChatSendMessage
	(*ChatReceiveMessage)(nil), // 4: api.protocol.ChatReceiveMessage
	(*GroupChatMessage)(nil),   // 5: api.protocol.GroupChatMessage
	(*ServiceMessage)(nil),     // 6: api.protocol.ServiceMessage
}
var file_api_protocol_message_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_protocol_message_proto_init() }
func file_api_protocol_message_proto_init() {
	if File_api_protocol_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_protocol_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_api_protocol_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransMessage); i {
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
		file_api_protocol_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DispatcherMessage); i {
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
		file_api_protocol_message_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatSendMessage); i {
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
		file_api_protocol_message_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatReceiveMessage); i {
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
		file_api_protocol_message_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GroupChatMessage); i {
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
		file_api_protocol_message_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceMessage); i {
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
			RawDescriptor: file_api_protocol_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_protocol_message_proto_goTypes,
		DependencyIndexes: file_api_protocol_message_proto_depIdxs,
		MessageInfos:      file_api_protocol_message_proto_msgTypes,
	}.Build()
	File_api_protocol_message_proto = out.File
	file_api_protocol_message_proto_rawDesc = nil
	file_api_protocol_message_proto_goTypes = nil
	file_api_protocol_message_proto_depIdxs = nil
}
