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
	ParentId     *int64 `form:"parent_id" db:"parent_id"`        // 父级ID
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
	ResourceID     int64  `gorm:"column:resource_id;primaryKey;comment:主键ID" json:"resource_id,string"`        // 主键ID
	Name           string `gorm:"column:name;not null;comment:资源名称" json:"name"`                               // 资源名称
	OriginalName   string `gorm:"column:original_name;comment:原始文件名（仅文件有）" json:"original_name"`               // 原始文件名（仅文件有）
	Extension      string `gorm:"column:extension;comment:文件扩展名（仅文件有）" json:"extension"`                       // 文件扩展名（仅文件有）
	MimeType       string `gorm:"column:mime_type;comment:MIME类型（仅文件有）" json:"mime_type"`                      // MIME类型（仅文件有）
	Type           string `gorm:"column:type;not null;default:FILE;comment:类型：FILE-文件，FOLDER-文件夹" json:"type"` // 类型：FILE-文件，FOLDER-文件夹
	Path           string `gorm:"column:path;not null;comment:完整路径（不含域名）" json:"path"`                         // 完整路径（不含域名）
	FileSize       int64  `gorm:"column:file_size;comment:文件大小（字节）" json:"file_size"`                          // 文件大小（字节）
	Md5Hash        string `gorm:"column:md5_hash;comment:文件MD5哈希值" json:"md5_hash"`                            // 文件MD5哈希值
	StorageKey     string `gorm:"column:storage_key;comment:存储服务中的唯一标识" json:"storage_key"`                    // 存储服务中的唯一标识
	StorageType    string `gorm:"column:storage_type;default:LOCAL;comment:存储类型" json:"storage_type"`          // 存储类型
	ParentID       int64  `gorm:"column:parent_id;comment:父级ID（0表示根目录）" json:"parent_id,string"`               // 父级ID（0表示根目录）
	Level          int32  `gorm:"column:level;comment:层级深度" json:"level"`                                      // 层级深度
	Lineage        string `gorm:"column:lineage;comment:层级路径，格式：/父ID/祖父ID/..." json:"lineage"`                 // 层级路径，格式：/父ID/祖父ID/...
	PermissionMask int32  `gorm:"column:permission_mask;comment:权限掩码" json:"permission_mask,string"`           // 权限掩码
	Category       string `gorm:"column:category;comment:分类：IMAGE, DOCUMENT, VIDEO, AUDIO等" json:"category"`   // 分类：IMAGE, DOCUMENT, VIDEO, AUDIO等
	Status         int32  `gorm:"column:status;default:1;comment:状态：0-删除，1-正常，2-隐藏" json:"status"`             // 状态：0-删除，1-正常，2-隐藏
	Description    string `gorm:"column:description;comment:描述" json:"description"`                            // 描述
}

// SysResourceFileVo ResourceFile VO struct
type SysResourceFileVo struct {
	ResourceID     int64  `gorm:"column:resource_id;primaryKey;comment:主键ID" json:"resource_id,string"`        // 主键ID
	Name           string `gorm:"column:name;not null;comment:资源名称" json:"name"`                               // 资源名称
	OriginalName   string `gorm:"column:original_name;comment:原始文件名（仅文件有）" json:"original_name"`               // 原始文件名（仅文件有）
	Extension      string `gorm:"column:extension;comment:文件扩展名（仅文件有）" json:"extension"`                       // 文件扩展名（仅文件有）
	MimeType       string `gorm:"column:mime_type;comment:MIME类型（仅文件有）" json:"mime_type"`                      // MIME类型（仅文件有）
	Type           string `gorm:"column:type;not null;default:FILE;comment:类型：FILE-文件，FOLDER-文件夹" json:"type"` // 类型：FILE-文件，FOLDER-文件夹
	Path           string `gorm:"column:path;not null;comment:完整路径（不含域名）" json:"path"`                         // 完整路径（不含域名）
	FileSize       int64  `gorm:"column:file_size;comment:文件大小（字节）" json:"file_size"`                          // 文件大小（字节）
	Md5Hash        string `gorm:"column:md5_hash;comment:文件MD5哈希值" json:"md5_hash"`                            // 文件MD5哈希值
	StorageKey     string `gorm:"column:storage_key;comment:存储服务中的唯一标识" json:"storage_key"`                    // 存储服务中的唯一标识
	StorageType    string `gorm:"column:storage_type;default:LOCAL;comment:存储类型" json:"storage_type"`          // 存储类型
	ParentID       int64  `gorm:"column:parent_id;comment:父级ID（0表示根目录）" json:"parent_id,string"`               // 父级ID（0表示根目录）
	Level          int32  `gorm:"column:level;comment:层级深度" json:"level"`                                      // 层级深度
	Lineage        string `gorm:"column:lineage;comment:层级路径，格式：/父ID/祖父ID/..." json:"lineage"`                 // 层级路径，格式：/父ID/祖父ID/...
	PermissionMask int32  `gorm:"column:permission_mask;comment:权限掩码" json:"permission_mask,string"`           // 权限掩码
	Category       string `gorm:"column:category;comment:分类：IMAGE, DOCUMENT, VIDEO, AUDIO等" json:"category"`   // 分类：IMAGE, DOCUMENT, VIDEO, AUDIO等
	Status         bool   `gorm:"column:status;default:1;comment:状态：0-删除，1-正常，2-隐藏" json:"status"`             // 状态：0-删除，1-正常，2-隐藏
	Description    string `gorm:"column:description;comment:描述" json:"description"`                            // 描述
	baize.BaseEntity
}

// SysResourceFile 资源文件
type SysResourceFile struct {
	ResourceID     int64  `gorm:"column:resource_id;primaryKey;comment:主键ID" json:"resource_id,string"`        // 主键ID
	Name           string `gorm:"column:name;not null;comment:资源名称" json:"name"`                               // 资源名称
	OriginalName   string `gorm:"column:original_name;comment:原始文件名（仅文件有）" json:"original_name"`               // 原始文件名（仅文件有）
	Extension      string `gorm:"column:extension;comment:文件扩展名（仅文件有）" json:"extension"`                       // 文件扩展名（仅文件有）
	MimeType       string `gorm:"column:mime_type;comment:MIME类型（仅文件有）" json:"mime_type"`                      // MIME类型（仅文件有）
	Type           string `gorm:"column:type;not null;default:FILE;comment:类型：FILE-文件，FOLDER-文件夹" json:"type"` // 类型：FILE-文件，FOLDER-文件夹
	Path           string `gorm:"column:path;not null;comment:完整路径（不含域名）" json:"path"`                         // 完整路径（不含域名）
	FileSize       int64  `gorm:"column:file_size;comment:文件大小（字节）" json:"file_size"`                          // 文件大小（字节）
	Md5Hash        string `gorm:"column:md5_hash;comment:文件MD5哈希值" json:"md5_hash"`                            // 文件MD5哈希值
	StorageKey     string `gorm:"column:storage_key;comment:存储服务中的唯一标识" json:"storage_key"`                    // 存储服务中的唯一标识
	StorageType    string `gorm:"column:storage_type;default:LOCAL;comment:存储类型" json:"storage_type"`          // 存储类型
	ParentID       int64  `gorm:"column:parent_id;comment:父级ID（0表示根目录）" json:"parent_id"`                      // 父级ID（0表示根目录）
	Level          int32  `gorm:"column:level;comment:层级深度" json:"level"`                                      // 层级深度
	Lineage        string `gorm:"column:lineage;comment:层级路径，格式：/父ID/祖父ID/..." json:"lineage"`                 // 层级路径，格式：/父ID/祖父ID/...
	PermissionMask int32  `gorm:"column:permission_mask;comment:权限掩码" json:"permission_mask,string"`           // 权限掩码
	Category       string `gorm:"column:category;comment:分类：IMAGE, DOCUMENT, VIDEO, AUDIO等" json:"category"`   // 分类：IMAGE, DOCUMENT, VIDEO, AUDIO等
	Status         int32  `gorm:"column:status;default:1;comment:状态：0-删除，1-正常，2-隐藏" json:"status"`             // 状态：0-删除，1-正常，2-隐藏
	Description    string `gorm:"column:description;comment:描述" json:"description"`                            // 描述
	DeptId         int64  `json:"deptId,string" db:"dept_id"`                                                  // 部门ID
	State          int    `json:"state" db:"state"`                                                            // 操作状态
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
		Path:         dml.Path,
		FileSize:     dml.FileSize,
		Md5Hash:      dml.Md5Hash,
		StorageType:  dml.StorageType,
		ParentID:     dml.ParentID,
		Level:        dml.Level,
		Lineage:      dml.Lineage,
		Category:     dml.Category,
		Status:       dml.Status,
		Description:  dml.Description,
	}
}

// SysResourceFileList 资源列表
type SysResourceFileList struct {
	List  []*SysResourceFileVo `json:"list"`
	Total int64                `json:"total"`
}
