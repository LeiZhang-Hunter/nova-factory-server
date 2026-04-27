package fileUtils

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// BuildAbsoluteURL converts a relative resource path to an absolute URL using the current request host.
func BuildAbsoluteURL(c *gin.Context, resourcePath string) string {
	resourcePath = strings.TrimSpace(resourcePath)
	if resourcePath == "" || strings.HasPrefix(resourcePath, "http://") || strings.HasPrefix(resourcePath, "https://") {
		return resourcePath
	}

	baseURL := configuredBaseURL(c)
	if baseURL == "" {
		return resourcePath
	}

	if !strings.HasPrefix(resourcePath, "/") {
		resourcePath = "/" + resourcePath
	}

	return strings.TrimRight(baseURL, "/") + resourcePath
}

func configuredBaseURL(c *gin.Context) string {
	for _, key := range []string{"upload_file.domain_name", "host"} {
		if baseURL := strings.TrimSpace(viper.GetString(key)); baseURL != "" {
			return baseURL
		}
	}

	host := requestHost(c)
	if host == "" {
		return ""
	}

	return requestScheme(c) + "://" + host
}

func requestScheme(c *gin.Context) string {
	if c == nil {
		return "http"
	}
	if forwardedProto := strings.TrimSpace(c.GetHeader("X-Forwarded-Proto")); forwardedProto != "" {
		return forwardedProto
	}
	if c.Request != nil && c.Request.TLS != nil {
		return "https"
	}
	return "http"
}

func requestHost(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if forwardedHost := strings.TrimSpace(c.GetHeader("X-Forwarded-Host")); forwardedHost != "" {
		return forwardedHost
	}
	if c.Request == nil {
		return ""
	}
	return strings.TrimSpace(c.Request.Host)
}
