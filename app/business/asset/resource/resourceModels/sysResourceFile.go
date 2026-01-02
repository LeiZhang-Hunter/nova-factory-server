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
	ResourceId        int64  `json:"resourceId,string" db:"resource_id" swaggerignore:"true"` // 主键ID
	Name              string `json:"name" db:"name" binding:"required"`                       // 资源名称
	OriginalName      string `json:"originalName" db:"original_name"`                         // 原始文件名
	Extension         string `json:"extension" db:"extension"`                                // 文件扩展名
	MimeType          string `json:"mimeType" db:"mime_type"`                                 // MIME类型
	Type              string `json:"type" db:"type" binding:"required"`                       // 类型：FILE-文件，FOLDER-文件夹
	Path              string `json:"path" db:"path"`                                          // 完整路径
	IsFolder          int    `json:"isFolder" db:"is_folder"`
	FileSize          int64  `json:"fileSize" db:"file_size"`                           // 文件大小
	Md5Hash           string `json:"md5Hash" db:"md5_hash"`                             // MD5哈希
	StorageKey        string `json:"storageKey" db:"storage_key"`                       // 存储唯一标识
	StorageType       string `json:"storageType" db:"storage_type"`                     // 存储类型
	ParentId          int64  `json:"parentId,string" db:"parent_id"`                    // 父级ID
	Level             int    `json:"level" db:"level"`                                  // 层级深度
	Lineage           string `json:"lineage" db:"lineage"`                              // 层级路径
	OwnerId           int64  `json:"ownerId,string" db:"owner_id"`                      // 所有者ID
	OwnerType         string `json:"ownerType" db:"owner_type"`                         // 所有者类型
	IsPublic          int    `json:"isPublic" db:"is_public"`                           // 是否公开
	PermissionMask    int    `json:"permissionMask" db:"permission_mask"`               // 权限掩码
	Version           int    `json:"version" db:"version"`                              // 版本号
	IsLatest          int    `json:"isLatest" db:"is_latest"`                           // 是否最新版本
	PreviousVersionId int64  `json:"previousVersionId,string" db:"previous_version_id"` // 上一个版本ID
	Category          string `json:"category" db:"category"`                            // 分类
	Tags              string `json:"tags" db:"tags"`                                    // 标签JSON
	Status            bool   `json:"status" db:"status"`                                // 状态
	Description       string `json:"description" db:"description"`                      // 描述
	Metadata          string `json:"metadata" db:"metadata"`                            // 元数据JSON
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
	ResourceId   int64  `json:"resourceId,string" db:"resource_id,string"`
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
	Status       bool   `json:"status" db:"status"`
	Description  string `json:"description" db:"description"`
	Metadata     string `json:"metadata" db:"metadata"`
	DeptId       int64  `json:"deptId,string" db:"dept_id"` // 部门ID
	State        int    `json:"state" db:"state"`           // 操作状态
	baize.BaseEntity
}

func ToSysResourceFile(dml *SysResourceFileDML) *SysResourceFile {
	return &SysResourceFile{
		ResourceId:   dml.ResourceId,
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
		ParentId:     dml.ParentId,
		Level:        dml.Level,
		Lineage:      dml.Lineage,
		OwnerId:      dml.OwnerId,
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
