package ocr

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ParseFileData struct {
	FileType string         `json:"file_type"`
	Text     string         `json:"txt"`
	Extra    map[string]any `json:"-"`
}

func (p *ParseFileData) UnmarshalJSON(data []byte) error {
	type alias struct {
		FileType string `json:"file_type"`
		Text     string `json:"txt"`
	}
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	p.FileType = a.FileType
	p.Text = a.Text
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	delete(m, "file_type")
	delete(m, "txt")
	p.Extra = m
	return nil
}

type ParseFileResponse struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data *ParseFileData `json:"data"`
}

type ParseFileRequest struct {
	UploadFile  *multipart.FileHeader
	StartPage   int
	NumberWords int
}

type FileOCRClient struct {
	baseURL    string
	parsePath  string
	httpClient *http.Client
}

func NewFileOCRClient(baseURL string, httpClient *http.Client) *FileOCRClient {
	client := httpClient
	if client == nil {
		client = &http.Client{Timeout: 2 * time.Minute}
	}
	return &FileOCRClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		parsePath:  "/",
		httpClient: client,
	}
}

func (f *FileOCRClient) ParseFile(ctx context.Context, req *ParseFileRequest) (*ParseFileResponse, error) {
	if req == nil {
		return nil, errors.New("request不能为空")
	}
	startPage := req.StartPage
	numberWords := req.NumberWords
	if numberWords <= 0 {
		numberWords = 100
	}
	return f.parseByUploadHeader(ctx, req.UploadFile, startPage, numberWords)

}

func (f *FileOCRClient) parseByUploadHeader(ctx context.Context, header *multipart.FileHeader, startPage int, numberWords int) (*ParseFileResponse, error) {
	if header == nil {
		return nil, errors.New("upload file不能为空")
	}
	file, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	formFile, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(formFile, file); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}
	u, err := url.Parse(f.baseURL + f.parsePath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("start_page", toIntString(startPage))
	query.Set("number_words", toIntString(numberWords))
	u.RawQuery = query.Encode()
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := f.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return decodeParseResponse(resp.Body)
}

func decodeParseResponse(body io.Reader) (*ParseFileResponse, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var ret ParseFileResponse
	if err = json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}
	if ret.Code != 200 {
		if strings.TrimSpace(ret.Msg) == "" {
			ret.Msg = "解析失败"
		}
		return nil, errors.New(ret.Msg)
	}
	return &ret, nil
}

func isRemoteFile(path string) bool {
	lower := strings.ToLower(strings.TrimSpace(path))
	return strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")
}

func toIntString(v int) string {
	return strconv.Itoa(v)
}
