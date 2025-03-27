// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: proto/search/search.proto

package update

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ImgUrl        string                 `protobuf:"bytes,5,opt,name=imgUrl,proto3" json:"imgUrl,omitempty"`
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

func (x *Comics) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
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

var file_proto_search_search_proto_rawDesc = string([]byte{
	0x0a, 0x19, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x73,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x3d, 0x0a, 0x0d, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x68, 0x72, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x70, 0x68, 0x72, 0x61, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d,
	0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22,
	0x42, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x69, 0x63, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x6d, 0x67,
	0x55, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x69, 0x6d, 0x67, 0x55, 0x72,
	0x6c, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x4a, 0x04, 0x08, 0x03, 0x10, 0x04, 0x4a, 0x04, 0x08,
	0x04, 0x10, 0x05, 0x22, 0x3b, 0x0a, 0x11, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64,
	0x65, 0x64, 0x43, 0x6f, 0x6d, 0x69, 0x63, 0x73, 0x12, 0x26, 0x0a, 0x06, 0x63, 0x6f, 0x6d, 0x69,
	0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x69, 0x63, 0x73, 0x52, 0x06, 0x63, 0x6f, 0x6d, 0x69, 0x63, 0x73,
	0x32, 0x80, 0x01, 0x0a, 0x06, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x38, 0x0a, 0x04, 0x50,
	0x69, 0x6e, 0x67, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3c, 0x0a, 0x06, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12,
	0x15, 0x2e, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x43, 0x6f, 0x6d, 0x69, 0x63,
	0x73, 0x22, 0x00, 0x42, 0x1f, 0x5a, 0x1d, 0x79, 0x61, 0x64, 0x72, 0x6f, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x63, 0x6f, 0x75, 0x72, 0x73, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

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
	3, // 3: update.Search.Ping:output_type -> google.protobuf.Empty
	2, // 4: update.Search.Search:output_type -> update.RecommendedComics
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
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
