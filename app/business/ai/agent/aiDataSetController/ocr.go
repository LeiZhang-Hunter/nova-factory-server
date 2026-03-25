package aiDataSetController

import (
	//"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/ocr"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type OCR struct {
	client *ocr.FileOCRClient
}

func NewOCR() *OCR {
	baseURL := strings.TrimSpace(viper.GetString("ocr.host"))
	if baseURL == "" {
		baseURL = strings.TrimSpace(viper.GetString("ocr.base_url"))
	}
	if baseURL == "" {
		baseURL = "http://127.0.0.1:5000"
	}
	return &OCR{
		client: ocr.NewFileOCRClient(baseURL, nil),
	}
}

func (o *OCR) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/ocr")
	//group.POST("/predict", middlewares.HasPermission("ai:ocr:predict"), o.ExtractOCR)
	group.POST("/predict", o.ExtractOCR)
}

// ExtractOCR OCR识别
// @Summary OCR识别
// @Description 上传文件并调用OCR服务进行文字提取
// @Tags 工业智能体/OCR
// @Security BearerAuth
// @Accept multipart/form-data
// @Param file formData file true "上传文件"
// @Param start_page query int false "起始页码"
// @Param number_words query int false "最多返回词数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "识别成功"
// @Router /ai/ocr/predict [post]
func (o *OCR) ExtractOCR(c *gin.Context) {
	file, err := c.FormFile("file")
	startPage, _ := strconv.Atoi(c.DefaultQuery("start_page", "0"))
	numberWords, _ := strconv.Atoi(c.DefaultQuery("number_words", "100"))
	if err != nil || file == nil {
		baizeContext.ParameterError(c)
		return
	}
	resp, err := o.client.ParseFile(c, &ocr.ParseFileRequest{
		UploadFile:  file,
		StartPage:   startPage,
		NumberWords: numberWords,
	})
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, resp.Data)
}
