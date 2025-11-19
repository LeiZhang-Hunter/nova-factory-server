package alertController

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"nova-factory-server/app/business/ai/aiDataSetService"
	"nova-factory-server/app/business/alert/alertDao"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/alert/alertService"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/utils/bufferpool"
	"nova-factory-server/app/utils/template"
	"sync"
)

// Runner  运行异步状态分析告警数据
type Runner struct {
	taskNumber    uint32
	done          chan struct{}
	alertChan     chan *alertModels.AlertLogInfo
	wait          sync.WaitGroup
	service       alertService.AlertRuleService
	runnerService alertService.RunnerService
	chatService   aiDataSetService.IChartService
	bp            *bufferpool.BufferPool
	deviceDao     deviceDao.IDeviceDao
	logDao        alertDao.AlertLogDao
}

var runner *Runner

func NewRunner(service alertService.AlertRuleService, runnerService alertService.RunnerService,
	chatService aiDataSetService.IChartService, deviceDao deviceDao.IDeviceDao, logDao alertDao.AlertLogDao) *Runner {
	taskNumber := viper.GetUint32("task_number")
	if taskNumber == 0 {
		taskNumber = 2
	}
	taskChannelNumber := viper.GetUint32("task_channel_number")
	if taskChannelNumber == 0 {
		taskChannelNumber = 32
	}
	r := &Runner{
		taskNumber:    taskNumber,
		done:          make(chan struct{}),
		alertChan:     make(chan *alertModels.AlertLogInfo, taskChannelNumber),
		service:       service,
		runnerService: runnerService,
		chatService:   chatService,
		bp:            bufferpool.NewBufferPool(1024),
		deviceDao:     deviceDao,
		logDao:        logDao,
	}
	runner = r
	return r
}

func GetAlertRunner() *Runner {
	return runner
}

func (r *Runner) Push(data *alertModels.AlertLogInfo) {
	select {
	case r.alertChan <- data:
		return
	}
}

func (r *Runner) Run() {
	for i := 0; i < int(r.taskNumber); i++ {
		r.wait.Add(1)
		go func() {
			defer r.wait.Done()
			r.run()
		}()
	}
	return
}

func (r *Runner) run() {
	for {
		select {
		case <-r.done:
			return
		case alertData := <-r.alertChan:
			{
				r.handle(alertData)
			}
		}
	}
}

func (r *Runner) Stop() {
	close(r.done)
	r.wait.Wait()
}

func (r *Runner) handle(data *alertModels.AlertLogInfo) {
	if data == nil {
		return
	}

	if data.GatewayId <= 0 {
		return
	}
	ctx := gin.Context{}
	// 读取网关下面配置的告警策略
	reason, err := r.service.GetReasonByGatewayId(&ctx, data.GatewayId)
	if err != nil {
		zap.L().Error("GetReasonByGatewayId", zap.Error(err))
		return
	}

	if reason == nil {
		return
	}

	makeTemplate, err := template.MakeTemplate(reason.Message)
	if err != nil {
		zap.L().Error("MakeTemplate failed", zap.Error(err))
		return
	}

	var param map[string]interface{} = make(map[string]interface{})
	if data.Alert.Labels.DeviceId != "" {
		deviceInfo, err := r.deviceDao.GetByIdString(&ctx, data.Alert.Labels.DeviceId)
		if err != nil {
			zap.L().Error("get by device id failed", zap.Error(err))
		}
		if deviceInfo != nil {
			param["Device"] = deviceInfo.Name
		} else {
			param["Device"] = ""
		}
	}
	param["Message"] = data.Alert.Labels.Message

	buffer := r.bp.Get()
	defer r.bp.Put(buffer)
	err = makeTemplate.Execute(buffer, param)
	if err != nil {
		return
	}

	return
}
