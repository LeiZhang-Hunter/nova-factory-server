package fileUtils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

func GetFileContentType(file multipart.File) string {

	data := make([]byte, 512)
	_, err := file.Read(data)
	if err != nil {
		panic(err)
	}
	return http.DetectContentType(data)
}

// ValidateAllowedDocumentFiles validates uploaded files by filename extension.
func ValidateAllowedDocumentFiles(files []*multipart.FileHeader) error {
	allowedExt := map[string]struct{}{
		".md":   {},
		".mdx":  {},
		".docx": {},
		".xlsx": {},
		".xls":  {},
		".pptx": {},
		".pdf":  {},
		".txt":  {},
		".jpeg": {},
		".jpg":  {},
		".png":  {},
		".tif":  {},
		".gif":  {},
		".csv":  {},
		".json": {},
		".eml":  {},
		".html": {},
	}
	for _, f := range files {
		if f == nil {
			continue
		}
		ext := strings.ToLower(filepath.Ext(strings.TrimSpace(f.Filename)))
		if _, ok := allowedExt[ext]; !ok {
			return errors.New(fmt.Sprintf("%s %s", f.Filename,
				"不支持的文件格式，仅支持：MD、MDX、DOCX、XLSX、XLS、PPTX、PDF、TXT、JPEG、JPG、PNG、TIF、GIF、CSV、JSON、EML、HTML"))
		}
	}
	return nil
}
