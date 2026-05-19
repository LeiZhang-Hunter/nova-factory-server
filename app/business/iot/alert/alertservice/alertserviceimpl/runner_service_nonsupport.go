//go:build !ai
// +build !ai

package alertserviceimpl

import (
	"errors"
	"nova-factory-server/app/business/iot/alert/alertservice"

	"github.com/gin-gonic/gin"
)

type RunnerServiceImpl struct{}

func NewRunnerServiceImpl() alertservice.RunnerService {
	return &RunnerServiceImpl{}
}

func (r *RunnerServiceImpl) Load(_ *gin.Context, _ string) (string, error) {
	return "", errors.New("ai runner is disabled in this build")
}
