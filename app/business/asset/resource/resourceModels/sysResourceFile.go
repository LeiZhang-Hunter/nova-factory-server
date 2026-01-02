package resourceModels

import (
	"nova-factory-server/app/baize"
)

// SysResourceFileDQL ResourceFile DQL struct
type SysResourceFileDQL struct {
	Name         string `form:"name" db:"name"`                  // 资源名称
	OriginalName string `form:"originalName" db:"original_name"` // 原始文件名
	Extension    string `form:"extension" db:"extension"`        // 文件扩展名
	Type         string `form:"type" db:"type"`                  // 类型：FILE-文件，FOLDER-文件夹
	Category     string `form:"category" db:"category"`          // 分类
	ParentId     int64  `form:"parentId" db:"parent_id"`         // 父级ID
	OwnerId      int64  `form:"ownerId" db:"owner_id"`           // 所有者ID
	OwnerType    string `form:"ownerType" db:"owner_type"`       // 所有者类型
	IsPublic     *int   `form:"isPublic" db:"is_public"`         // 是否公开
	Status       *int   `form:"status" db:"status"`              // 状态
	BeginTime    string `form:"beginTime" db:"begin_time"`       // 开始时间
	EndTime      string `form:"endTime" db:"end_time"`           // 结束时间
	baize.BaseEntityDQL
}

// SysResourceFileDML ResourceFile DML struct
type SysResourceFileDML struct {
	ResourceID        int64  `gorm:"column:resource_id;primaryKey;comment:主键ID" json:"resource_id"`               // 主键ID
	Name              string `gorm:"column:name;not null;comment:资源名称" json:"name"`                               // 资源名称
	OriginalName      string `gorm:"column:original_name;comment:原始文件名（仅文件有）" json:"original_name"`               // 原始文件名（仅文件有）
	Extension         string `gorm:"column:extension;comment:文件扩展名（仅文件有）" json:"extension"`                       // 文件扩展名（仅文件有）
	MimeType          string `gorm:"column:mime_type;comment:MIME类型（仅文件有）" json:"mime_type"`                      // MIME类型（仅文件有）
	Type              string `gorm:"column:type;not null;default:FILE;comment:类型：FILE-文件，FOLDER-文件夹" json:"type"` // 类型：FILE-文件，FOLDER-文件夹
	IsFolder          bool   `gorm:"column:is_folder;comment:是否文件夹（虚拟列，方便查询）" json:"is_folder"`                   // 是否文件夹（虚拟列，方便查询）
	Path              string `gorm:"column:path;not null;comment:完整路径（不含域名）" json:"path"`                         // 完整路径（不含域名）
	FileSize          int64  `gorm:"column:file_size;comment:文件大小（字节）" json:"file_size"`                          // 文件大小（字节）
	Md5Hash           string `gorm:"column:md5_hash;comment:文件MD5哈希值" json:"md5_hash"`                            // 文件MD5哈希值
	StorageKey        string `gorm:"column:storage_key;comment:存储服务中的唯一标识" json:"storage_key"`                    // 存储服务中的唯一标识
	StorageType       string `gorm:"column:storage_type;default:LOCAL;comment:存储类型" json:"storage_type"`          // 存储类型
	ParentID          int64  `gorm:"column:parent_id;comment:父级ID（0表示根目录）" json:"parent_id"`                      // 父级ID（0表示根目录）
	Level             int32  `gorm:"column:level;comment:层级深度" json:"level"`                                      // 层级深度
	Lineage           string `gorm:"column:lineage;comment:层级路径，格式：/父ID/祖父ID/..." json:"lineage"`                 // 层级路径，格式：/父ID/祖父ID/...
	OwnerID           int64  `gorm:"column:owner_id;comment:所有者ID" json:"owner_id"`                               // 所有者ID
	OwnerType         string `gorm:"column:owner_type;comment:所有者类型：USER-用户，DEPARTMENT-部门等" json:"owner_type"`    // 所有者类型：USER-用户，DEPARTMENT-部门等
	IsPublic          bool   `gorm:"column:is_public;comment:是否公开" json:"is_public"`                              // 是否公开
	PermissionMask    int32  `gorm:"column:permission_mask;comment:权限掩码" json:"permission_mask"`                  // 权限掩码
	Version           int32  `gorm:"column:version;default:1;comment:版本号" json:"version"`                         // 版本号
	IsLatest          bool   `gorm:"column:is_latest;default:1;comment:是否最新版本" json:"is_latest"`                  // 是否最新版本
	PreviousVersionID int64  `gorm:"column:previous_version_id;comment:上一个版本ID" json:"previous_version_id"`       // 上一个版本ID
	Category          string `gorm:"column:category;comment:分类：IMAGE, DOCUMENT, VIDEO, AUDIO等" json:"category"`   // 分类：IMAGE, DOCUMENT, VIDEO, AUDIO等
	Tags              string `gorm:"column:tags;comment:标签数组" json:"tags"`                                        // 标签数组
	Status            int32  `gorm:"column:status;default:1;comment:状态：0-删除，1-正常，2-隐藏" json:"status"`             // 状态：0-删除，1-正常，2-隐藏
	Description       string `gorm:"column:description;comment:描述" json:"description"`                            // 描述
	Metadata          string `gorm:"column:metadata;comment:元数据（如：图片宽高、视频时长等）" json:"metadata"`                   // 元数据（如：图片宽高、视频时长等）
}

// SysResourceFileVo ResourceFile VO struct
type SysResourceFileVo struct {
	ResourceId   int64  `json:"resourceId,string" db:"resource_id"`
	Name         string `json:"name" db:"name"`
	OriginalName string `json:"originalName" db:"original_name"`
	Extension    string `json:"extension" db:"extension"`
	MimeType     string `json:"mimeType" db:"mime_type"`
	Type         string `json:"type" db:"type"`
	IsFolder     int    `json:"isFolder" db:"is_folder"`
	Path         string `json:"path" db:"path"`
	FileSize     int64  `json:"fileSize" db:"file_size"`
	Md5Hash      string `json:"md5Hash" db:"md5_hash"`
	StorageType  string `json:"storageType" db:"storage_type"`
	ParentId     int64  `json:"parentId,string" db:"parent_id"`
	Level        int    `json:"level" db:"level"`
	Lineage      string `json:"lineage" db:"lineage"`
	OwnerId      int64  `json:"ownerId,string" db:"owner_id"`
	OwnerType    string `json:"ownerType" db:"owner_type"`
	IsPublic     int    `json:"isPublic" db:"is_public"`
	Category     string `json:"category" db:"category"`
	Tags         string `json:"tags" db:"tags"`
	Status       string `json:"status" db:"status"`
	Description  string `json:"description" db:"description"`
	Metadata     string `json:"metadata" db:"metadata"`
	CreateBy     int64  `json:"createBy,string" db:"create_by"`
	CreateTime   string `json:"createTime" db:"create_time"`
	UpdateBy     int64  `json:"updateBy,string" db:"update_by"`
	UpdateTime   string `json:"updateTime" db:"update_time"`
	baize.BaseEntity
}

// SysResourceFile 资源文件
type SysResourceFile struct {
	ResourceID        int64  `gorm:"column:resource_id;primaryKey;comment:主键ID" json:"resource_id"`               // 主键ID
	Name              string `gorm:"column:name;not null;comment:资源名称" json:"name"`                               // 资源名称
	OriginalName      string `gorm:"column:original_name;comment:原始文件名（仅文件有）" json:"original_name"`               // 原始文件名（仅文件有）
	Extension         string `gorm:"column:extension;comment:文件扩展名（仅文件有）" json:"extension"`                       // 文件扩展名（仅文件有）
	MimeType          string `gorm:"column:mime_type;comment:MIME类型（仅文件有）" json:"mime_type"`                      // MIME类型（仅文件有）
	Type              string `gorm:"column:type;not null;default:FILE;comment:类型：FILE-文件，FOLDER-文件夹" json:"type"` // 类型：FILE-文件，FOLDER-文件夹
	IsFolder          bool   `gorm:"column:is_folder;comment:是否文件夹（虚拟列，方便查询）" json:"is_folder"`                   // 是否文件夹（虚拟列，方便查询）
	Path              string `gorm:"column:path;not null;comment:完整路径（不含域名）" json:"path"`                         // 完整路径（不含域名）
	FileSize          int64  `gorm:"column:file_size;comment:文件大小（字节）" json:"file_size"`                          // 文件大小（字节）
	Md5Hash           string `gorm:"column:md5_hash;comment:文件MD5哈希值" json:"md5_hash"`                            // 文件MD5哈希值
	StorageKey        string `gorm:"column:storage_key;comment:存储服务中的唯一标识" json:"storage_key"`                    // 存储服务中的唯一标识
	StorageType       string `gorm:"column:storage_type;default:LOCAL;comment:存储类型" json:"storage_type"`          // 存储类型
	ParentID          int64  `gorm:"column:parent_id;comment:父级ID（0表示根目录）" json:"parent_id"`                      // 父级ID（0表示根目录）
	Level             int32  `gorm:"column:level;comment:层级深度" json:"level"`                                      // 层级深度
	Lineage           string `gorm:"column:lineage;comment:层级路径，格式：/父ID/祖父ID/..." json:"lineage"`                 // 层级路径，格式：/父ID/祖父ID/...
	OwnerID           int64  `gorm:"column:owner_id;comment:所有者ID" json:"owner_id"`                               // 所有者ID
	OwnerType         string `gorm:"column:owner_type;comment:所有者类型：USER-用户，DEPARTMENT-部门等" json:"owner_type"`    // 所有者类型：USER-用户，DEPARTMENT-部门等
	IsPublic          bool   `gorm:"column:is_public;comment:是否公开" json:"is_public"`                              // 是否公开
	PermissionMask    int32  `gorm:"column:permission_mask;comment:权限掩码" json:"permission_mask"`                  // 权限掩码
	Version           int32  `gorm:"column:version;default:1;comment:版本号" json:"version"`                         // 版本号
	IsLatest          bool   `gorm:"column:is_latest;default:1;comment:是否最新版本" json:"is_latest"`                  // 是否最新版本
	PreviousVersionID int64  `gorm:"column:previous_version_id;comment:上一个版本ID" json:"previous_version_id"`       // 上一个版本ID
	Category          string `gorm:"column:category;comment:分类：IMAGE, DOCUMENT, VIDEO, AUDIO等" json:"category"`   // 分类：IMAGE, DOCUMENT, VIDEO, AUDIO等
	Tags              string `gorm:"column:tags;comment:标签数组" json:"tags"`                                        // 标签数组
	Status            int32  `gorm:"column:status;default:1;comment:状态：0-删除，1-正常，2-隐藏" json:"status"`             // 状态：0-删除，1-正常，2-隐藏
	Description       string `gorm:"column:description;comment:描述" json:"description"`                            // 描述
	Metadata          string `gorm:"column:metadata;comment:元数据（如：图片宽高、视频时长等）" json:"metadata"`                   // 元数据（如：图片宽高、视频时长等）
	DeptId            int64  `json:"deptId,string" db:"dept_id"`                                                  // 部门ID
	State             int    `json:"state" db:"state"`                                                            // 操作状态
	baize.BaseEntity
}

func ToSysResourceFile(dml *SysResourceFileDML) *SysResourceFile {
	return &SysResourceFile{
		ResourceID:   dml.ResourceID,
		Name:         dml.Name,
		OriginalName: dml.OriginalName,
		Extension:    dml.Extension,
		MimeType:     dml.MimeType,
		Type:         dml.Type,
		IsFolder:     dml.IsFolder,
		Path:         dml.Path,
		FileSize:     dml.FileSize,
		Md5Hash:      dml.Md5Hash,
		StorageType:  dml.StorageType,
		ParentID:     dml.ParentID,
		Level:        dml.Level,
		Lineage:      dml.Lineage,
		OwnerID:      dml.OwnerID,
		OwnerType:    dml.OwnerType,
		IsPublic:     dml.IsPublic,
		Category:     dml.Category,
		Tags:         dml.Tags,
		Status:       dml.Status,
		Description:  dml.Description,
		Metadata:     dml.Metadata,
	}
}

// SysResourceFileList 资源列表
type SysResourceFileList struct {
	List  []*SysResourceFileVo `json:"list"`
	Total int64                `json:"total"`
}
