// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: mediafile/mediafile.proto

package mediafile

import (
	common "github.com/Zzocker/book-labs/protos/common"
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

type MediaFileType int32

const (
	MediaFileType_UNKNOWN MediaFileType = 0
	MediaFileType_BOOK    MediaFileType = 1
	MediaFileType_USER    MediaFileType = 2
	MediaFileType_COMMENT MediaFileType = 3
)

// Enum value maps for MediaFileType.
var (
	MediaFileType_name = map[int32]string{
		0: "UNKNOWN",
		1: "BOOK",
		2: "USER",
		3: "COMMENT",
	}
	MediaFileType_value = map[string]int32{
		"UNKNOWN": 0,
		"BOOK":    1,
		"USER":    2,
		"COMMENT": 3,
	}
)

func (x MediaFileType) Enum() *MediaFileType {
	p := new(MediaFileType)
	*p = x
	return p
}

func (x MediaFileType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MediaFileType) Descriptor() protoreflect.EnumDescriptor {
	return file_mediafile_mediafile_proto_enumTypes[0].Descriptor()
}

func (MediaFileType) Type() protoreflect.EnumType {
	return &file_mediafile_mediafile_proto_enumTypes[0]
}

func (x MediaFileType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MediaFileType.Descriptor instead.
func (MediaFileType) EnumDescriptor() ([]byte, []int) {
	return file_mediafile_mediafile_proto_rawDescGZIP(), []int{0}
}

type MediaFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID        string        `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Data      []byte        `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data,omitempty"`
	Extension string        `protobuf:"bytes,3,opt,name=Extension,proto3" json:"Extension,omitempty"`
	Type      MediaFileType `protobuf:"varint,4,opt,name=Type,proto3,enum=mediafile.MediaFileType" json:"Type,omitempty"`
}

func (x *MediaFile) Reset() {
	*x = MediaFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediafile_mediafile_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaFile) ProtoMessage() {}

func (x *MediaFile) ProtoReflect() protoreflect.Message {
	mi := &file_mediafile_mediafile_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaFile.ProtoReflect.Descriptor instead.
func (*MediaFile) Descriptor() ([]byte, []int) {
	return file_mediafile_mediafile_proto_rawDescGZIP(), []int{0}
}

func (x *MediaFile) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *MediaFile) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *MediaFile) GetExtension() string {
	if x != nil {
		return x.Extension
	}
	return ""
}

func (x *MediaFile) GetType() MediaFileType {
	if x != nil {
		return x.Type
	}
	return MediaFileType_UNKNOWN
}

type MediaFileID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *MediaFileID) Reset() {
	*x = MediaFileID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediafile_mediafile_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaFileID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaFileID) ProtoMessage() {}

func (x *MediaFileID) ProtoReflect() protoreflect.Message {
	mi := &file_mediafile_mediafile_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaFileID.ProtoReflect.Descriptor instead.
func (*MediaFileID) Descriptor() ([]byte, []int) {
	return file_mediafile_mediafile_proto_rawDescGZIP(), []int{1}
}

func (x *MediaFileID) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

var File_mediafile_mediafile_proto protoreflect.FileDescriptor

var file_mediafile_mediafile_proto_rawDesc = []byte{
	0x0a, 0x19, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x66, 0x69, 0x6c, 0x65, 0x2f, 0x6d, 0x65, 0x64, 0x69,
	0x61, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x6d, 0x65, 0x64,
	0x69, 0x61, 0x66, 0x69, 0x6c, 0x65, 0x1a, 0x13, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7b, 0x0a, 0x09, 0x4d,
	0x65, 0x64, 0x69, 0x61, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1c, 0x0a, 0x09,
	0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2c, 0x0a, 0x04, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61,
	0x66, 0x69, 0x6c, 0x65, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x22, 0x1d, 0x0a, 0x0b, 0x4d, 0x65, 0x64, 0x69,
	0x61, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x2a, 0x3d, 0x0a, 0x0d, 0x4d, 0x65, 0x64, 0x69, 0x61,
	0x46, 0x69, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e,
	0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x42, 0x4f, 0x4f, 0x4b, 0x10, 0x01, 0x12,
	0x08, 0x0a, 0x04, 0x55, 0x53, 0x45, 0x52, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x4d,
	0x4d, 0x45, 0x4e, 0x54, 0x10, 0x03, 0x32, 0xb0, 0x01, 0x0a, 0x10, 0x4d, 0x65, 0x64, 0x69, 0x61,
	0x46, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x36, 0x0a, 0x06, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x14, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x66, 0x69, 0x6c,
	0x65, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x46, 0x69, 0x6c, 0x65, 0x1a, 0x16, 0x2e, 0x6d, 0x65,
	0x64, 0x69, 0x61, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x46, 0x69, 0x6c,
	0x65, 0x49, 0x44, 0x12, 0x33, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x16, 0x2e, 0x6d, 0x65, 0x64,
	0x69, 0x61, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x46, 0x69, 0x6c, 0x65,
	0x49, 0x44, 0x1a, 0x14, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x4d,
	0x65, 0x64, 0x69, 0x61, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x2f, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x12, 0x16, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x4d,
	0x65, 0x64, 0x69, 0x61, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x1a, 0x0d, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x5a, 0x7a, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2f,
	0x62, 0x6f, 0x6f, 0x6b, 0x2d, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x2f, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x66, 0x69, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_mediafile_mediafile_proto_rawDescOnce sync.Once
	file_mediafile_mediafile_proto_rawDescData = file_mediafile_mediafile_proto_rawDesc
)

func file_mediafile_mediafile_proto_rawDescGZIP() []byte {
	file_mediafile_mediafile_proto_rawDescOnce.Do(func() {
		file_mediafile_mediafile_proto_rawDescData = protoimpl.X.CompressGZIP(file_mediafile_mediafile_proto_rawDescData)
	})
	return file_mediafile_mediafile_proto_rawDescData
}

var file_mediafile_mediafile_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_mediafile_mediafile_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_mediafile_mediafile_proto_goTypes = []interface{}{
	(MediaFileType)(0),   // 0: mediafile.MediaFileType
	(*MediaFile)(nil),    // 1: mediafile.MediaFile
	(*MediaFileID)(nil),  // 2: mediafile.MediaFileID
	(*common.Empty)(nil), // 3: common.Empty
}
var file_mediafile_mediafile_proto_depIdxs = []int32{
	0, // 0: mediafile.MediaFile.Type:type_name -> mediafile.MediaFileType
	1, // 1: mediafile.MediaFileService.Upload:input_type -> mediafile.MediaFile
	2, // 2: mediafile.MediaFileService.Get:input_type -> mediafile.MediaFileID
	2, // 3: mediafile.MediaFileService.Delete:input_type -> mediafile.MediaFileID
	2, // 4: mediafile.MediaFileService.Upload:output_type -> mediafile.MediaFileID
	1, // 5: mediafile.MediaFileService.Get:output_type -> mediafile.MediaFile
	3, // 6: mediafile.MediaFileService.Delete:output_type -> common.Empty
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_mediafile_mediafile_proto_init() }
func file_mediafile_mediafile_proto_init() {
	if File_mediafile_mediafile_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mediafile_mediafile_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MediaFile); i {
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
		file_mediafile_mediafile_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MediaFileID); i {
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
			RawDescriptor: file_mediafile_mediafile_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mediafile_mediafile_proto_goTypes,
		DependencyIndexes: file_mediafile_mediafile_proto_depIdxs,
		EnumInfos:         file_mediafile_mediafile_proto_enumTypes,
		MessageInfos:      file_mediafile_mediafile_proto_msgTypes,
	}.Build()
	File_mediafile_mediafile_proto = out.File
	file_mediafile_mediafile_proto_rawDesc = nil
	file_mediafile_mediafile_proto_goTypes = nil
	file_mediafile_mediafile_proto_depIdxs = nil
}
