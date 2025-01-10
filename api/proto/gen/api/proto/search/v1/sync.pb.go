// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        (unknown)
// source: api/proto/search/v1/sync.proto

package searchv1

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

type InputAnyRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	IndexName     string                 `protobuf:"bytes,1,opt,name=index_name,json=indexName,proto3" json:"index_name,omitempty"`
	DocId         string                 `protobuf:"bytes,2,opt,name=doc_id,json=docId,proto3" json:"doc_id,omitempty"`
	Data          string                 `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputAnyRequest) Reset() {
	*x = InputAnyRequest{}
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputAnyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputAnyRequest) ProtoMessage() {}

func (x *InputAnyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputAnyRequest.ProtoReflect.Descriptor instead.
func (*InputAnyRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_search_v1_sync_proto_rawDescGZIP(), []int{0}
}

func (x *InputAnyRequest) GetIndexName() string {
	if x != nil {
		return x.IndexName
	}
	return ""
}

func (x *InputAnyRequest) GetDocId() string {
	if x != nil {
		return x.DocId
	}
	return ""
}

func (x *InputAnyRequest) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type InputAnyResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputAnyResponse) Reset() {
	*x = InputAnyResponse{}
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputAnyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputAnyResponse) ProtoMessage() {}

func (x *InputAnyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputAnyResponse.ProtoReflect.Descriptor instead.
func (*InputAnyResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_search_v1_sync_proto_rawDescGZIP(), []int{1}
}

type InputUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *User                  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputUserRequest) Reset() {
	*x = InputUserRequest{}
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputUserRequest) ProtoMessage() {}

func (x *InputUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputUserRequest.ProtoReflect.Descriptor instead.
func (*InputUserRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_search_v1_sync_proto_rawDescGZIP(), []int{2}
}

func (x *InputUserRequest) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type InputUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputUserResponse) Reset() {
	*x = InputUserResponse{}
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputUserResponse) ProtoMessage() {}

func (x *InputUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputUserResponse.ProtoReflect.Descriptor instead.
func (*InputUserResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_search_v1_sync_proto_rawDescGZIP(), []int{3}
}

type InputArticleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Article       *Article               `protobuf:"bytes,1,opt,name=article,proto3" json:"article,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputArticleRequest) Reset() {
	*x = InputArticleRequest{}
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputArticleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputArticleRequest) ProtoMessage() {}

func (x *InputArticleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputArticleRequest.ProtoReflect.Descriptor instead.
func (*InputArticleRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_search_v1_sync_proto_rawDescGZIP(), []int{4}
}

func (x *InputArticleRequest) GetArticle() *Article {
	if x != nil {
		return x.Article
	}
	return nil
}

type InputArticleResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputArticleResponse) Reset() {
	*x = InputArticleResponse{}
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputArticleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputArticleResponse) ProtoMessage() {}

func (x *InputArticleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputArticleResponse.ProtoReflect.Descriptor instead.
func (*InputArticleResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_search_v1_sync_proto_rawDescGZIP(), []int{5}
}

type Article struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Status        int32                  `protobuf:"varint,3,opt,name=status,proto3" json:"status,omitempty"`
	Content       string                 `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Tags          []string               `protobuf:"bytes,5,rep,name=tags,proto3" json:"tags,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Article) Reset() {
	*x = Article{}
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Article) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Article) ProtoMessage() {}

func (x *Article) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Article.ProtoReflect.Descriptor instead.
func (*Article) Descriptor() ([]byte, []int) {
	return file_api_proto_search_v1_sync_proto_rawDescGZIP(), []int{6}
}

func (x *Article) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Article) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Article) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *Article) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Article) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

type User struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Email         string                 `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Nickname      string                 `protobuf:"bytes,3,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Phone         string                 `protobuf:"bytes,4,opt,name=phone,proto3" json:"phone,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_search_v1_sync_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_api_proto_search_v1_sync_proto_rawDescGZIP(), []int{7}
}

func (x *User) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *User) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

var File_api_proto_search_v1_sync_proto protoreflect.FileDescriptor

var file_api_proto_search_v1_sync_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x09, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x22, 0x5b, 0x0a, 0x0f, 0x49,
	0x6e, 0x70, 0x75, 0x74, 0x41, 0x6e, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d,
	0x0a, 0x0a, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x15, 0x0a,
	0x06, 0x64, 0x6f, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64,
	0x6f, 0x63, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x12, 0x0a, 0x10, 0x49, 0x6e, 0x70, 0x75,
	0x74, 0x41, 0x6e, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x37, 0x0a, 0x10,
	0x49, 0x6e, 0x70, 0x75, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x23, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x04, 0x75, 0x73, 0x65, 0x72, 0x22, 0x13, 0x0a, 0x11, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x43, 0x0a, 0x13, 0x49, 0x6e,
	0x70, 0x75, 0x74, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x2c, 0x0a, 0x07, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x07, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x22,
	0x16, 0x0a, 0x14, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x75, 0x0a, 0x07, 0x41, 0x72, 0x74, 0x69, 0x63,
	0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61,
	0x67, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x22, 0x5e,
	0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08,
	0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x32, 0xeb,
	0x01, 0x0a, 0x0b, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x46,
	0x0a, 0x09, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1b, 0x2e, 0x73, 0x65,
	0x61, 0x72, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4f, 0x0a, 0x0c, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x41,
	0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x12, 0x1e, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e,
	0x76, 0x31, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e,
	0x76, 0x31, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x08, 0x49, 0x6e, 0x70, 0x75, 0x74,
	0x41, 0x6e, 0x79, 0x12, 0x1a, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e,
	0x49, 0x6e, 0x70, 0x75, 0x74, 0x41, 0x6e, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1b, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x70, 0x75,
	0x74, 0x41, 0x6e, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x90, 0x01, 0x0a,
	0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x42, 0x09,
	0x53, 0x79, 0x6e, 0x63, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2f, 0x65, 0x63, 0x6f,
	0x64, 0x65, 0x70, 0x75, 0x62, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x2f, 0x76, 0x31, 0x3b, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x53,
	0x58, 0x58, 0xaa, 0x02, 0x09, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x56, 0x31, 0xca, 0x02,
	0x09, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x15, 0x53, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0xea, 0x02, 0x0a, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x3a, 0x3a, 0x56, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_search_v1_sync_proto_rawDescOnce sync.Once
	file_api_proto_search_v1_sync_proto_rawDescData = file_api_proto_search_v1_sync_proto_rawDesc
)

func file_api_proto_search_v1_sync_proto_rawDescGZIP() []byte {
	file_api_proto_search_v1_sync_proto_rawDescOnce.Do(func() {
		file_api_proto_search_v1_sync_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_search_v1_sync_proto_rawDescData)
	})
	return file_api_proto_search_v1_sync_proto_rawDescData
}

var file_api_proto_search_v1_sync_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_api_proto_search_v1_sync_proto_goTypes = []any{
	(*InputAnyRequest)(nil),      // 0: search.v1.InputAnyRequest
	(*InputAnyResponse)(nil),     // 1: search.v1.InputAnyResponse
	(*InputUserRequest)(nil),     // 2: search.v1.InputUserRequest
	(*InputUserResponse)(nil),    // 3: search.v1.InputUserResponse
	(*InputArticleRequest)(nil),  // 4: search.v1.InputArticleRequest
	(*InputArticleResponse)(nil), // 5: search.v1.InputArticleResponse
	(*Article)(nil),              // 6: search.v1.Article
	(*User)(nil),                 // 7: search.v1.User
}
var file_api_proto_search_v1_sync_proto_depIdxs = []int32{
	7, // 0: search.v1.InputUserRequest.user:type_name -> search.v1.User
	6, // 1: search.v1.InputArticleRequest.article:type_name -> search.v1.Article
	2, // 2: search.v1.SyncService.InputUser:input_type -> search.v1.InputUserRequest
	4, // 3: search.v1.SyncService.InputArticle:input_type -> search.v1.InputArticleRequest
	0, // 4: search.v1.SyncService.InputAny:input_type -> search.v1.InputAnyRequest
	3, // 5: search.v1.SyncService.InputUser:output_type -> search.v1.InputUserResponse
	5, // 6: search.v1.SyncService.InputArticle:output_type -> search.v1.InputArticleResponse
	1, // 7: search.v1.SyncService.InputAny:output_type -> search.v1.InputAnyResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_proto_search_v1_sync_proto_init() }
func file_api_proto_search_v1_sync_proto_init() {
	if File_api_proto_search_v1_sync_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proto_search_v1_sync_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_search_v1_sync_proto_goTypes,
		DependencyIndexes: file_api_proto_search_v1_sync_proto_depIdxs,
		MessageInfos:      file_api_proto_search_v1_sync_proto_msgTypes,
	}.Build()
	File_api_proto_search_v1_sync_proto = out.File
	file_api_proto_search_v1_sync_proto_rawDesc = nil
	file_api_proto_search_v1_sync_proto_goTypes = nil
	file_api_proto_search_v1_sync_proto_depIdxs = nil
}
