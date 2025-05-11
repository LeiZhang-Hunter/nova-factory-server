package file

import (
	"bytes"
	"io"
	"mime/multipart"
)

func CreatFormFiles(b *bytes.Buffer, form *multipart.Form, w *multipart.Writer) error {
	for name, files := range form.File {
		for _, file := range files {
			f, err := file.Open()
			if err != nil {
				return err
			}
			fw, err := w.CreateFormFile(name, file.Filename)
			if err != nil {
				return err
			}
			_, err = io.Copy(fw, f)
			if err != nil {
				return err
			}
		}
	}
	w.Close() //要关闭，会将w.w.boundary刷写到w.writer中
	return nil
}
