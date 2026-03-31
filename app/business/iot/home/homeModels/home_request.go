package homeModels

type HomeRequest struct {
	Platform string `form:"platform" json:"platform"`
	Version  string `form:"version" json:"version"`
}
