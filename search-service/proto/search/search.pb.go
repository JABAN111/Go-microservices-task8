// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: proto/search/search.proto

package update

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SearchRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Phrase        string                 `protobuf:"bytes,1,opt,name=phrase,proto3" json:"phrase,omitempty"`
	Limit         int64                  `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchRequest) Reset() {
	*x = SearchRequest{}
	mi := &file_proto_search_search_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchRequest) ProtoMessage() {}

func (x *SearchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_search_search_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchRequest.ProtoReflect.Descriptor instead.
func (*SearchRequest) Descriptor() ([]byte, []int) {
	return file_proto_search_search_proto_rawDescGZIP(), []int{0}
}

func (x *SearchRequest) GetPhrase() string {
	if x != nil {
		return x.Phrase
	}
	return ""
}

func (x *SearchRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type Comics struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ImgUrl        string                 `protobuf:"bytes,5,opt,name=img_url,json=imgUrl,proto3" json:"img_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Comics) Reset() {
	*x = Comics{}
	mi := &file_proto_search_search_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Comics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Comics) ProtoMessage() {}

func (x *Comics) ProtoReflect() protoreflect.Message {
	mi := &file_proto_search_search_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Comics.ProtoReflect.Descriptor instead.
func (*Comics) Descriptor() ([]byte, []int) {
	return file_proto_search_search_proto_rawDescGZIP(), []int{1}
}

func (x *Comics) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Comics) GetImgUrl() string {
	if x != nil {
		return x.ImgUrl
	}
	return ""
}

type RecommendedComics struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Comics        []*Comics              `protobuf:"bytes,1,rep,name=comics,proto3" json:"comics,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RecommendedComics) Reset() {
	*x = RecommendedComics{}
	mi := &file_proto_search_search_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RecommendedComics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecommendedComics) ProtoMessage() {}

func (x *RecommendedComics) ProtoReflect() protoreflect.Message {
	mi := &file_proto_search_search_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecommendedComics.ProtoReflect.Descriptor instead.
func (*RecommendedComics) Descriptor() ([]byte, []int) {
	return file_proto_search_search_proto_rawDescGZIP(), []int{2}
}

func (x *RecommendedComics) GetComics() []*Comics {
	if x != nil {
		return x.Comics
	}
	return nil
}

var File_proto_search_search_proto protoreflect.FileDescriptor

const file_proto_search_search_proto_rawDesc = "" +
	"\n" +
	"\x19proto/search/search.proto\x12\x06update\x1a\x1bgoogle/protobuf/empty.proto\"=\n" +
	"\rSearchRequest\x12\x16\n" +
	"\x06phrase\x18\x01 \x01(\tR\x06phrase\x12\x14\n" +
	"\x05limit\x18\x02 \x01(\x03R\x05limit\"C\n" +
	"\x06Comics\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x17\n" +
	"\aimg_url\x18\x05 \x01(\tR\x06imgUrlJ\x04\b\x02\x10\x03J\x04\b\x03\x10\x04J\x04\b\x04\x10\x05\";\n" +
	"\x11RecommendedComics\x12&\n" +
	"\x06comics\x18\x01 \x03(\v2\x0e.update.ComicsR\x06comics2\xbf\x01\n" +
	"\x06Search\x128\n" +
	"\x04Ping\x12\x16.google.protobuf.Empty\x1a\x16.google.protobuf.Empty\"\x00\x12<\n" +
	"\x06Search\x12\x15.update.SearchRequest\x1a\x19.update.RecommendedComics\"\x00\x12=\n" +
	"\aISearch\x12\x15.update.SearchRequest\x1a\x19.update.RecommendedComics\"\x00B\x1fZ\x1dyadro.com/course/proto/updateb\x06proto3"

var (
	file_proto_search_search_proto_rawDescOnce sync.Once
	file_proto_search_search_proto_rawDescData []byte
)

func file_proto_search_search_proto_rawDescGZIP() []byte {
	file_proto_search_search_proto_rawDescOnce.Do(func() {
		file_proto_search_search_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_search_search_proto_rawDesc), len(file_proto_search_search_proto_rawDesc)))
	})
	return file_proto_search_search_proto_rawDescData
}

var file_proto_search_search_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_search_search_proto_goTypes = []any{
	(*SearchRequest)(nil),     // 0: update.SearchRequest
	(*Comics)(nil),            // 1: update.Comics
	(*RecommendedComics)(nil), // 2: update.RecommendedComics
	(*emptypb.Empty)(nil),     // 3: google.protobuf.Empty
}
var file_proto_search_search_proto_depIdxs = []int32{
	1, // 0: update.RecommendedComics.comics:type_name -> update.Comics
	3, // 1: update.Search.Ping:input_type -> google.protobuf.Empty
	0, // 2: update.Search.Search:input_type -> update.SearchRequest
	0, // 3: update.Search.ISearch:input_type -> update.SearchRequest
	3, // 4: update.Search.Ping:output_type -> google.protobuf.Empty
	2, // 5: update.Search.Search:output_type -> update.RecommendedComics
	2, // 6: update.Search.ISearch:output_type -> update.RecommendedComics
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_search_search_proto_init() }
func file_proto_search_search_proto_init() {
	if File_proto_search_search_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_search_search_proto_rawDesc), len(file_proto_search_search_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_search_search_proto_goTypes,
		DependencyIndexes: file_proto_search_search_proto_depIdxs,
		MessageInfos:      file_proto_search_search_proto_msgTypes,
	}.Build()
	File_proto_search_search_proto = out.File
	file_proto_search_search_proto_goTypes = nil
	file_proto_search_search_proto_depIdxs = nil
}
