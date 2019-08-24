// Code generated by protoc-gen-go. DO NOT EDIT.
// source: metadata.proto

package redcloud

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

//
// Declaration of the type of data contained in a table.
type DataUsage int32

const (
	// Missing data usage declaration
	DataUsage_UNKNOWN DataUsage = 0
	// Table contains SPI
	DataUsage_SENSITIVE_PERSONAL_INFORMATION DataUsage = 1
	// Table contains internal business data that is not SPI
	DataUsage_INTERNAL_BUSINESS_DATA DataUsage = 2
)

var DataUsage_name = map[int32]string{
	0: "UNKNOWN",
	1: "SENSITIVE_PERSONAL_INFORMATION",
	2: "INTERNAL_BUSINESS_DATA",
}
var DataUsage_value = map[string]int32{
	"UNKNOWN":                        0,
	"SENSITIVE_PERSONAL_INFORMATION": 1,
	"INTERNAL_BUSINESS_DATA":         2,
}

func (x DataUsage) String() string {
	return proto.EnumName(DataUsage_name, int32(x))
}
func (DataUsage) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

//
// SSTablePathDescription describes all paths holding data for the given
// table/key range.
type SSTablePathDescription struct {
	// Name of the column family which is covered.
	ColumnFamily string `protobuf:"bytes,1,opt,name=column_family,json=columnFamily" json:"column_family,omitempty"`
	//
	// Path to the major sstable holding the current compacted version of the
	// column family in the specified table.
	MajorSstablePath string `protobuf:"bytes,2,opt,name=major_sstable_path,json=majorSstablePath" json:"major_sstable_path,omitempty"`
	//
	// Path to the minor sstable holding the current amended data of the column
	// family in the specified table.
	MinorSstablePath string `protobuf:"bytes,3,opt,name=minor_sstable_path,json=minorSstablePath" json:"minor_sstable_path,omitempty"`
	//
	// Path to all journal files containing further amendments of the minor
	// sstable which have not been (fully) compacted yet.
	RelevantJournalPaths []string `protobuf:"bytes,4,rep,name=relevant_journal_paths,json=relevantJournalPaths" json:"relevant_journal_paths,omitempty"`
	// Size of the major sstable on the most recent compaction.
	MajorSstableSize int64 `protobuf:"varint,5,opt,name=major_sstable_size,json=majorSstableSize" json:"major_sstable_size,omitempty"`
	// Size of the minor sstable on the most recent compaction.
	MinorSstableSize int64 `protobuf:"varint,6,opt,name=minor_sstable_size,json=minorSstableSize" json:"minor_sstable_size,omitempty"`
}

func (m *SSTablePathDescription) Reset()                    { *m = SSTablePathDescription{} }
func (m *SSTablePathDescription) String() string            { return proto.CompactTextString(m) }
func (*SSTablePathDescription) ProtoMessage()               {}
func (*SSTablePathDescription) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *SSTablePathDescription) GetColumnFamily() string {
	if m != nil {
		return m.ColumnFamily
	}
	return ""
}

func (m *SSTablePathDescription) GetMajorSstablePath() string {
	if m != nil {
		return m.MajorSstablePath
	}
	return ""
}

func (m *SSTablePathDescription) GetMinorSstablePath() string {
	if m != nil {
		return m.MinorSstablePath
	}
	return ""
}

func (m *SSTablePathDescription) GetRelevantJournalPaths() []string {
	if m != nil {
		return m.RelevantJournalPaths
	}
	return nil
}

func (m *SSTablePathDescription) GetMajorSstableSize() int64 {
	if m != nil {
		return m.MajorSstableSize
	}
	return 0
}

func (m *SSTablePathDescription) GetMinorSstableSize() int64 {
	if m != nil {
		return m.MinorSstableSize
	}
	return 0
}

//
// ServerTabletMetadata holds metadata for tablets, i.e. individual pieces of
// tables living on specific servers.
type ServerTabletMetadata struct {
	// The first key of the tablet.
	StartKey []byte `protobuf:"bytes,1,opt,name=start_key,json=startKey,proto3" json:"start_key,omitempty"`
	//
	// The first key after the end of the tablet, or an empty byte string
	// if this is the last tablet of the table.
	EndKey []byte `protobuf:"bytes,2,opt,name=end_key,json=endKey,proto3" json:"end_key,omitempty"`
	//
	// host currently holding the tablet. Please note that this may change
	// at any time.
	Host string `protobuf:"bytes,3,opt,name=host" json:"host,omitempty"`
	//
	// Port the host currently holding the tablet is exporting the database
	// service on.
	Port int32 `protobuf:"varint,4,opt,name=port" json:"port,omitempty"`
	//
	// Optional path to the sstable file holding the table data, so future
	// data nodes can pick it up.
	SstablePath []*SSTablePathDescription `protobuf:"bytes,5,rep,name=sstable_path,json=sstablePath" json:"sstable_path,omitempty"`
}

func (m *ServerTabletMetadata) Reset()                    { *m = ServerTabletMetadata{} }
func (m *ServerTabletMetadata) String() string            { return proto.CompactTextString(m) }
func (*ServerTabletMetadata) ProtoMessage()               {}
func (*ServerTabletMetadata) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *ServerTabletMetadata) GetStartKey() []byte {
	if m != nil {
		return m.StartKey
	}
	return nil
}

func (m *ServerTabletMetadata) GetEndKey() []byte {
	if m != nil {
		return m.EndKey
	}
	return nil
}

func (m *ServerTabletMetadata) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *ServerTabletMetadata) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *ServerTabletMetadata) GetSstablePath() []*SSTablePathDescription {
	if m != nil {
		return m.SstablePath
	}
	return nil
}

//
// ColumnFamilyMetadata holds metadata for an individual column family. Only
// column families which have a ColumnFamilyMetadata record will be considered
// as existent.
type ColumnFamilyMetadata struct {
	// Name of the column family registered.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *ColumnFamilyMetadata) Reset()                    { *m = ColumnFamilyMetadata{} }
func (m *ColumnFamilyMetadata) String() string            { return proto.CompactTextString(m) }
func (*ColumnFamilyMetadata) ProtoMessage()               {}
func (*ColumnFamilyMetadata) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

func (m *ColumnFamilyMetadata) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

//
// TableMetadata holds table metadata for redcloud tables. This is the
// user-specified part of the table metadata.
type TableMetadata struct {
	// The name of the table.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	//
	// Desired size after which tablets should be split.
	// If unspecified, defaults to 128 MB.
	SplitSize int64 `protobuf:"varint,2,opt,name=split_size,json=splitSize" json:"split_size,omitempty"`
	//
	// Maximum number of versions to keep of each cell. If 0, cells will never
	// be expired automatically based on the number of them.
	MaxVersions int64 `protobuf:"varint,3,opt,name=max_versions,json=maxVersions" json:"max_versions,omitempty"`
	//
	// Maximum age, in milliseconds, of older versions of cells, i.e. the ones
	// which do not have the highest insertion time stamp. If 0, old versions
	// are not expired based on age.
	MaxVersionAge int64 `protobuf:"varint,4,opt,name=max_version_age,json=maxVersionAge" json:"max_version_age,omitempty"`
	//
	// Path the files will be stored under. Changes to this value will only take
	// effect gradually when tablets are reloaded.
	PathPrefix string `protobuf:"bytes,5,opt,name=path_prefix,json=pathPrefix" json:"path_prefix,omitempty"`
	//
	// List of all column families configured for the table.
	ColumnFamily []*ColumnFamilyMetadata `protobuf:"bytes,6,rep,name=column_family,json=columnFamily" json:"column_family,omitempty"`
	//
	// Declaration of the type of data contained in the table.
	DataUsage DataUsage `protobuf:"varint,7,opt,name=data_usage,json=dataUsage,enum=redcloud.DataUsage" json:"data_usage,omitempty"`
}

func (m *TableMetadata) Reset()                    { *m = TableMetadata{} }
func (m *TableMetadata) String() string            { return proto.CompactTextString(m) }
func (*TableMetadata) ProtoMessage()               {}
func (*TableMetadata) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

func (m *TableMetadata) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TableMetadata) GetSplitSize() int64 {
	if m != nil {
		return m.SplitSize
	}
	return 0
}

func (m *TableMetadata) GetMaxVersions() int64 {
	if m != nil {
		return m.MaxVersions
	}
	return 0
}

func (m *TableMetadata) GetMaxVersionAge() int64 {
	if m != nil {
		return m.MaxVersionAge
	}
	return 0
}

func (m *TableMetadata) GetPathPrefix() string {
	if m != nil {
		return m.PathPrefix
	}
	return ""
}

func (m *TableMetadata) GetColumnFamily() []*ColumnFamilyMetadata {
	if m != nil {
		return m.ColumnFamily
	}
	return nil
}

func (m *TableMetadata) GetDataUsage() DataUsage {
	if m != nil {
		return m.DataUsage
	}
	return DataUsage_UNKNOWN
}

//
// ServerTableMetadata holds Server-side table metadata for redcloud tables.
type ServerTableMetadata struct {
	// The name of the table, for identifying it.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Table metadata associated with the table.
	TableMd *TableMetadata `protobuf:"bytes,2,opt,name=table_md,json=tableMd" json:"table_md,omitempty"`
	//
	// List of all tablets associated with this table and the servers they
	// are loaded on.
	Tablet []*ServerTabletMetadata `protobuf:"bytes,3,rep,name=tablet" json:"tablet,omitempty"`
}

func (m *ServerTableMetadata) Reset()                    { *m = ServerTableMetadata{} }
func (m *ServerTableMetadata) String() string            { return proto.CompactTextString(m) }
func (*ServerTableMetadata) ProtoMessage()               {}
func (*ServerTableMetadata) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{4} }

func (m *ServerTableMetadata) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ServerTableMetadata) GetTableMd() *TableMetadata {
	if m != nil {
		return m.TableMd
	}
	return nil
}

func (m *ServerTableMetadata) GetTablet() []*ServerTabletMetadata {
	if m != nil {
		return m.Tablet
	}
	return nil
}

func init() {
	proto.RegisterType((*SSTablePathDescription)(nil), "redcloud.SSTablePathDescription")
	proto.RegisterType((*ServerTabletMetadata)(nil), "redcloud.ServerTabletMetadata")
	proto.RegisterType((*ColumnFamilyMetadata)(nil), "redcloud.ColumnFamilyMetadata")
	proto.RegisterType((*TableMetadata)(nil), "redcloud.TableMetadata")
	proto.RegisterType((*ServerTableMetadata)(nil), "redcloud.ServerTableMetadata")
	proto.RegisterEnum("redcloud.DataUsage", DataUsage_name, DataUsage_value)
}

func init() { proto.RegisterFile("metadata.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 579 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0x4d, 0x4f, 0xdb, 0x40,
	0x10, 0xad, 0x93, 0x90, 0xe0, 0x49, 0xa0, 0xd1, 0x82, 0x20, 0x6a, 0x55, 0xea, 0xba, 0x52, 0x65,
	0xa1, 0x8a, 0x43, 0x5a, 0xf5, 0x9e, 0x42, 0x90, 0x5c, 0x8a, 0x83, 0xd6, 0x81, 0xf6, 0xb6, 0x5a,
	0xf0, 0x02, 0xa6, 0xfe, 0xd2, 0xee, 0x06, 0x01, 0xff, 0xa3, 0xbf, 0xa1, 0x87, 0xfe, 0x82, 0xfe,
	0xbb, 0x6a, 0xc7, 0x86, 0x84, 0x10, 0x71, 0x1b, 0xbd, 0xf7, 0x76, 0xf4, 0xe6, 0xcd, 0x68, 0x61,
	0x35, 0x15, 0x9a, 0x47, 0x5c, 0xf3, 0x9d, 0x42, 0xe6, 0x3a, 0x27, 0xcb, 0x52, 0x44, 0x67, 0x49,
	0x3e, 0x89, 0xdc, 0xbf, 0x35, 0xd8, 0x08, 0xc3, 0x31, 0x3f, 0x4d, 0xc4, 0x11, 0xd7, 0x97, 0x7b,
	0x42, 0x9d, 0xc9, 0xb8, 0xd0, 0x71, 0x9e, 0x91, 0xf7, 0xb0, 0x72, 0x96, 0x27, 0x93, 0x34, 0x63,
	0xe7, 0x3c, 0x8d, 0x93, 0xdb, 0x9e, 0xe5, 0x58, 0x9e, 0x4d, 0x3b, 0x25, 0xb8, 0x8f, 0x18, 0xf9,
	0x08, 0x24, 0xe5, 0x57, 0xb9, 0x64, 0x4a, 0x69, 0xd3, 0x84, 0x15, 0x5c, 0x5f, 0xf6, 0x6a, 0xa8,
	0xec, 0x22, 0x13, 0x96, 0x84, 0xe9, 0x8e, 0xea, 0x38, 0x9b, 0x57, 0xd7, 0x2b, 0xb5, 0x61, 0x66,
	0xd5, 0x9f, 0x61, 0x43, 0x8a, 0x44, 0x5c, 0xf3, 0x4c, 0xb3, 0xab, 0x7c, 0x22, 0x33, 0x9e, 0xe0,
	0x03, 0xd5, 0x6b, 0x38, 0x75, 0xcf, 0xa6, 0xeb, 0xf7, 0xec, 0xb7, 0x92, 0x34, 0x8f, 0xd4, 0x53,
	0x47, 0x2a, 0xbe, 0x13, 0xbd, 0x25, 0xc7, 0xf2, 0xea, 0x8f, 0x1d, 0x85, 0xf1, 0x9d, 0x78, 0xea,
	0x08, 0xd5, 0xcd, 0x4a, 0x3d, 0xe3, 0xc8, 0xa8, 0xdd, 0x7f, 0x16, 0xac, 0x87, 0x42, 0x5e, 0x0b,
	0x89, 0x89, 0xe9, 0xc3, 0x2a, 0x56, 0xf2, 0x1a, 0x6c, 0xa5, 0xb9, 0xd4, 0xec, 0x97, 0x28, 0x73,
	0xea, 0xd0, 0x65, 0x04, 0x0e, 0xc4, 0x2d, 0xd9, 0x84, 0x96, 0xc8, 0x22, 0xa4, 0x6a, 0x48, 0x35,
	0x45, 0x16, 0x19, 0x82, 0x40, 0xe3, 0x32, 0x57, 0xba, 0x0a, 0x00, 0x6b, 0x83, 0x15, 0xb9, 0xd4,
	0xbd, 0x86, 0x63, 0x79, 0x4b, 0x14, 0x6b, 0xb2, 0x0b, 0x9d, 0x47, 0x81, 0x2d, 0x39, 0x75, 0xaf,
	0xdd, 0x77, 0x76, 0xee, 0xb7, 0xb8, 0xb3, 0x78, 0x83, 0xb4, 0xad, 0xa6, 0x69, 0xba, 0xdb, 0xb0,
	0xbe, 0x3b, 0xb3, 0xb9, 0x07, 0xeb, 0x04, 0x1a, 0x19, 0x4f, 0x45, 0xb5, 0x5d, 0xac, 0xdd, 0x3f,
	0x35, 0x58, 0xc1, 0x8e, 0xcf, 0xa9, 0xc8, 0x1b, 0x00, 0x55, 0x24, 0xb1, 0x2e, 0x33, 0xab, 0x61,
	0x66, 0x36, 0x22, 0x18, 0xed, 0x3b, 0xe8, 0xa4, 0xfc, 0x86, 0x5d, 0x0b, 0xa9, 0xe2, 0x3c, 0x53,
	0x38, 0x65, 0x9d, 0xb6, 0x53, 0x7e, 0x73, 0x52, 0x41, 0xe4, 0x03, 0xbc, 0x9c, 0x91, 0x30, 0x7e,
	0x21, 0x70, 0xee, 0x3a, 0x5d, 0x99, 0xaa, 0x06, 0x17, 0x82, 0xbc, 0x85, 0xb6, 0x19, 0x9c, 0x15,
	0x52, 0x9c, 0xc7, 0x37, 0xb8, 0x4c, 0x9b, 0x82, 0x81, 0x8e, 0x10, 0x21, 0xbb, 0xf3, 0xb7, 0xda,
	0xc4, 0x88, 0xb6, 0xa6, 0x11, 0x2d, 0x9a, 0x7d, 0xee, 0x96, 0xfb, 0x00, 0x06, 0x65, 0x13, 0x65,
	0x8c, 0xb4, 0x1c, 0xcb, 0x5b, 0xed, 0xaf, 0x4d, 0x3b, 0xec, 0x71, 0xcd, 0x8f, 0x0d, 0x45, 0xed,
	0xe8, 0xbe, 0x74, 0x7f, 0x5b, 0xb0, 0x36, 0x73, 0x11, 0xcf, 0xe6, 0xd5, 0x87, 0xe5, 0x72, 0x89,
	0x69, 0x84, 0x69, 0xb5, 0xfb, 0x9b, 0xd3, 0xee, 0x8f, 0x9e, 0xd3, 0x16, 0x0a, 0x0f, 0x23, 0xf2,
	0x05, 0x9a, 0x58, 0x9a, 0x23, 0x99, 0x9b, 0x68, 0xd1, 0x21, 0xd2, 0x4a, 0xbd, 0xfd, 0x13, 0xec,
	0x07, 0xbf, 0xa4, 0x0d, 0xad, 0xe3, 0xe0, 0x20, 0x18, 0xfd, 0x08, 0xba, 0x2f, 0x88, 0x0b, 0x5b,
	0xe1, 0x30, 0x08, 0xfd, 0xb1, 0x7f, 0x32, 0x64, 0x47, 0x43, 0x1a, 0x8e, 0x82, 0xc1, 0x77, 0xe6,
	0x07, 0xfb, 0x23, 0x7a, 0x38, 0x18, 0xfb, 0xa3, 0xa0, 0x6b, 0x91, 0x57, 0xb0, 0xe1, 0x07, 0xe3,
	0x21, 0x35, 0xcc, 0xd7, 0xe3, 0xd0, 0x0f, 0x86, 0x61, 0xc8, 0xf6, 0x06, 0xe3, 0x41, 0xb7, 0x76,
	0xda, 0xc4, 0x2f, 0xe4, 0xd3, 0xff, 0x00, 0x00, 0x00, 0xff, 0xff, 0xeb, 0x0c, 0xcf, 0x91, 0x54,
	0x04, 0x00, 0x00,
}