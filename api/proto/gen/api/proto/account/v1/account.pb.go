// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        (unknown)
// source: api/proto/account/v1/account.proto

package accountv1

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

type AccountType int32

const (
	AccountType_AccountTypeUnknown AccountType = 0
	// 个人赞赏账号
	AccountType_AccountTypeReward AccountType = 1
	// 平台分成账号
	AccountType_AccountTypeSystem AccountType = 2
)

// Enum value maps for AccountType.
var (
	AccountType_name = map[int32]string{
		0: "AccountTypeUnknown",
		1: "AccountTypeReward",
		2: "AccountTypeSystem",
	}
	AccountType_value = map[string]int32{
		"AccountTypeUnknown": 0,
		"AccountTypeReward":  1,
		"AccountTypeSystem":  2,
	}
)

func (x AccountType) Enum() *AccountType {
	p := new(AccountType)
	*p = x
	return p
}

func (x AccountType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AccountType) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_account_v1_account_proto_enumTypes[0].Descriptor()
}

func (AccountType) Type() protoreflect.EnumType {
	return &file_api_proto_account_v1_account_proto_enumTypes[0]
}

func (x AccountType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AccountType.Descriptor instead.
func (AccountType) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_account_v1_account_proto_rawDescGZIP(), []int{0}
}

type CreditRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	//  唯一标识业务的
	Biz   string `protobuf:"bytes,1,opt,name=biz,proto3" json:"biz,omitempty"`
	BizId int64  `protobuf:"varint,2,opt,name=biz_id,json=bizId,proto3" json:"biz_id,omitempty"`
	// 后续如果   还有退款，部分退款，平台垫资等需求，在这里加字段
	// 注意一点，就是账号服务一般来说会和很多服务的数据关联在一起
	// 后续对账、统计、报表，账号都是一个核心
	// 不同的账号金额变动,这里可能是一加一减的
	Items         []*CreditItem `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreditRequest) Reset() {
	*x = CreditRequest{}
	mi := &file_api_proto_account_v1_account_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreditRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreditRequest) ProtoMessage() {}

func (x *CreditRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_account_v1_account_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreditRequest.ProtoReflect.Descriptor instead.
func (*CreditRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_account_v1_account_proto_rawDescGZIP(), []int{0}
}

func (x *CreditRequest) GetBiz() string {
	if x != nil {
		return x.Biz
	}
	return ""
}

func (x *CreditRequest) GetBizId() int64 {
	if x != nil {
		return x.BizId
	}
	return 0
}

func (x *CreditRequest) GetItems() []*CreditItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type CreditItem struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// 在一些复杂的系统里面，用户可能有多个账号，还有虚拟账号，退款账号等乱七八糟的划分
	Account int64 `protobuf:"varint,1,opt,name=account,proto3" json:"account,omitempty"`
	// 账号类型
	AccountType AccountType `protobuf:"varint,2,opt,name=account_type,json=accountType,proto3,enum=account.v1.AccountType" json:"account_type,omitempty"`
	// 金额
	Amt int64 `protobuf:"varint,3,opt,name=amt,proto3" json:"amt,omitempty"`
	// 货币，正常来说它类似于支付，最开始就尽量把货币的问题纳入考虑范围
	Currency      string `protobuf:"bytes,4,opt,name=currency,proto3" json:"currency,omitempty"`
	Uid           int64  `protobuf:"varint,5,opt,name=uid,proto3" json:"uid,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreditItem) Reset() {
	*x = CreditItem{}
	mi := &file_api_proto_account_v1_account_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreditItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreditItem) ProtoMessage() {}

func (x *CreditItem) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_account_v1_account_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreditItem.ProtoReflect.Descriptor instead.
func (*CreditItem) Descriptor() ([]byte, []int) {
	return file_api_proto_account_v1_account_proto_rawDescGZIP(), []int{1}
}

func (x *CreditItem) GetAccount() int64 {
	if x != nil {
		return x.Account
	}
	return 0
}

func (x *CreditItem) GetAccountType() AccountType {
	if x != nil {
		return x.AccountType
	}
	return AccountType_AccountTypeUnknown
}

func (x *CreditItem) GetAmt() int64 {
	if x != nil {
		return x.Amt
	}
	return 0
}

func (x *CreditItem) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

func (x *CreditItem) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

type CreditResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreditResponse) Reset() {
	*x = CreditResponse{}
	mi := &file_api_proto_account_v1_account_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreditResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreditResponse) ProtoMessage() {}

func (x *CreditResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_account_v1_account_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreditResponse.ProtoReflect.Descriptor instead.
func (*CreditResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_account_v1_account_proto_rawDescGZIP(), []int{2}
}

var File_api_proto_account_v1_account_proto protoreflect.FileDescriptor

var file_api_proto_account_v1_account_proto_rawDesc = []byte{
	0x0a, 0x22, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x76, 0x31,
	0x22, 0x66, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x64, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x10, 0x0a, 0x03, 0x62, 0x69, 0x7a, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x62, 0x69, 0x7a, 0x12, 0x15, 0x0a, 0x06, 0x62, 0x69, 0x7a, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x62, 0x69, 0x7a, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x05, 0x69, 0x74,
	0x65, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x64, 0x69, 0x74, 0x49, 0x74, 0x65,
	0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0xa2, 0x01, 0x0a, 0x0a, 0x43, 0x72, 0x65,
	0x64, 0x69, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x3a, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x0b, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x61, 0x6d, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x61, 0x6d, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x69, 0x64, 0x22, 0x10, 0x0a,
	0x0e, 0x43, 0x72, 0x65, 0x64, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2a,
	0x53, 0x0a, 0x0b, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16,
	0x0a, 0x12, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x55, 0x6e, 0x6b,
	0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x10, 0x01, 0x12, 0x15, 0x0a,
	0x11, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x53, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x10, 0x02, 0x32, 0x51, 0x0a, 0x0e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3f, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x64, 0x69, 0x74,
	0x12, 0x19, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72,
	0x65, 0x64, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x64, 0x69, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x9a, 0x01, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x2e,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x31, 0x65, 0x63, 0x6f, 0x64,
	0x65, 0x70, 0x75, 0x62, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x2f, 0x76, 0x31, 0x3b, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x76, 0x31, 0xa2, 0x02, 0x03,
	0x41, 0x58, 0x58, 0xaa, 0x02, 0x0a, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x56, 0x31,
	0xca, 0x02, 0x0a, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x16,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0b, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_account_v1_account_proto_rawDescOnce sync.Once
	file_api_proto_account_v1_account_proto_rawDescData = file_api_proto_account_v1_account_proto_rawDesc
)

func file_api_proto_account_v1_account_proto_rawDescGZIP() []byte {
	file_api_proto_account_v1_account_proto_rawDescOnce.Do(func() {
		file_api_proto_account_v1_account_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_account_v1_account_proto_rawDescData)
	})
	return file_api_proto_account_v1_account_proto_rawDescData
}

var file_api_proto_account_v1_account_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_proto_account_v1_account_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_proto_account_v1_account_proto_goTypes = []any{
	(AccountType)(0),       // 0: account.v1.AccountType
	(*CreditRequest)(nil),  // 1: account.v1.CreditRequest
	(*CreditItem)(nil),     // 2: account.v1.CreditItem
	(*CreditResponse)(nil), // 3: account.v1.CreditResponse
}
var file_api_proto_account_v1_account_proto_depIdxs = []int32{
	2, // 0: account.v1.CreditRequest.items:type_name -> account.v1.CreditItem
	0, // 1: account.v1.CreditItem.account_type:type_name -> account.v1.AccountType
	1, // 2: account.v1.AccountService.Credit:input_type -> account.v1.CreditRequest
	3, // 3: account.v1.AccountService.Credit:output_type -> account.v1.CreditResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_proto_account_v1_account_proto_init() }
func file_api_proto_account_v1_account_proto_init() {
	if File_api_proto_account_v1_account_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proto_account_v1_account_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_account_v1_account_proto_goTypes,
		DependencyIndexes: file_api_proto_account_v1_account_proto_depIdxs,
		EnumInfos:         file_api_proto_account_v1_account_proto_enumTypes,
		MessageInfos:      file_api_proto_account_v1_account_proto_msgTypes,
	}.Build()
	File_api_proto_account_v1_account_proto = out.File
	file_api_proto_account_v1_account_proto_rawDesc = nil
	file_api_proto_account_v1_account_proto_goTypes = nil
	file_api_proto_account_v1_account_proto_depIdxs = nil
}
