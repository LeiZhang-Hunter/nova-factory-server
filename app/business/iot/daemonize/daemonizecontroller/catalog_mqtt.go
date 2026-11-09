package daemonizecontroller

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"nova-factory-server/app/business/iot/daemonize/catalogwire"
)

const (
	catalogQueueSize       = 64
	catalogWorkerCount     = 4
	catalogMaxRequestBytes = 4096
)

// StartCatalogMQTT 启动网关配置目录响应服务，返回的清理函数应在服务关闭时调用。
func (c *Config) StartCatalogMQTT(ctx context.Context) (func(), error) {
	if !viper.GetBool("gateway_mqtt.enabled") {
		return func() {}, nil
	}

	opts, err := catalogMQTTOptions()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	jobs := make(chan catalogwire.Request, catalogQueueSize)
	subscribed := make(chan error, 1)
	opts.SetOnConnectHandler(catalogOnConnect(ctx, jobs, subscribed))
	client := mqtt.NewClient(opts)

	var workers sync.WaitGroup
	for i := 0; i < catalogWorkerCount; i++ {
		workers.Add(1)
		go c.runCatalogWorker(ctx, client, jobs, &workers)
	}

	cleanup := func() {
		cancel()
		client.Disconnect(250)
		workers.Wait()
	}

	token := client.Connect()
	if !token.WaitTimeout(10*time.Second) || token.Error() != nil {
		cleanup()
		return nil, errors.New("gateway MQTT connection failed")
	}

	select {
	case <-ctx.Done():
		cleanup()
		return nil, ctx.Err()
	case err := <-subscribed:
		if err != nil {
			cleanup()
			return nil, err
		}
	case <-time.After(6 * time.Second):
		cleanup()
		return nil, errors.New("gateway MQTT subscription timed out")
	}

	return cleanup, nil
}

// catalogMQTTOptions 读取并校验 MQTT 配置，构造客户端选项。
func catalogMQTTOptions() (*mqtt.ClientOptions, error) {
	broker := viper.GetString("gateway_mqtt.broker")
	u, err := url.Parse(broker)
	if err != nil || u.Hostname() == "" || u.Port() == "" || u.User != nil || u.Path != "" ||
		u.RawQuery != "" || u.Fragment != "" || (u.Scheme != "tcp" && u.Scheme != "ssl") {
		return nil, errors.New("invalid gateway_mqtt broker")
	}

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(viper.GetString("gateway_mqtt.client_id")).
		SetUsername(viper.GetString("gateway_mqtt.username")).
		SetPassword(viper.GetString("gateway_mqtt.password")).
		SetAutoReconnect(true).
		SetConnectRetry(false).
		SetCleanSession(true).
		SetConnectTimeout(5 * time.Second).
		SetWriteTimeout(5 * time.Second)

	if opts.ClientID == "" {
		return nil, errors.New("gateway_mqtt.client_id is required")
	}

	opts.SetTLSConfig(&tls.Config{
		MinVersion: tls.VersionTLS12,
	})

	return opts, nil
}

// catalogOnConnect 创建连接回调，订阅目录请求并将有效消息放入队列。
func catalogOnConnect(ctx context.Context, jobs chan<- catalogwire.Request, subscribed chan<- error) mqtt.OnConnectHandler {
	handler := func(_ mqtt.Client, m mqtt.Message) {
		if m.Retained() || len(m.Payload()) > catalogMaxRequestBytes {
			return
		}

		var req catalogwire.Request
		if json.Unmarshal(m.Payload(), &req) != nil || catalogwire.TopicGateway(m.Topic()) != req.GatewayID ||
			!catalogwire.ValidID(req.GatewayID) ||
			!catalogwire.ValidRequestID(req.RequestID) {
			return
		}

		select {
		case jobs <- req:
		case <-ctx.Done():
		default:
			zap.L().Warn("gateway MQTT catalog queue full")
		}
	}

	return func(client mqtt.Client) {
		token := client.Subscribe(catalogwire.Prefix+"+/catalog/request", 1, handler)
		var err error
		if !token.WaitTimeout(5*time.Second) || token.Error() != nil {
			err = errors.New("gateway MQTT catalog subscription failed")
			zap.L().Error("gateway MQTT catalog subscription failed")
		}

		select {
		case subscribed <- err:
		default:
		}
	}
}

// runCatalogWorker 消费目录请求，并为每次处理设置超时时间。
func (c *Config) runCatalogWorker(ctx context.Context, client mqtt.Client, jobs <-chan catalogwire.Request, workers *sync.WaitGroup) {
	defer workers.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case req := <-jobs:
			op, done := context.WithTimeout(ctx, 20*time.Second)
			c.answerCatalog(op, client, req)
			done()
		}
	}
}

// answerCatalog 验证网关请求，生成并发布签名后的目录响应。
func (c *Config) answerCatalog(ctx context.Context, client mqtt.Client, req catalogwire.Request) {
	id, _ := strconv.ParseUint(req.GatewayID, 10, 64)
	agent, err := c.agentService.GetByObjectId(ctx, id)
	if err != nil || agent == nil || agent.Username != req.Username || !req.Verify(agent.Password, time.Now()) {
		return
	}

	reply := catalogwire.Response{
		RequestID: req.RequestID,
		GatewayID: req.GatewayID,
	}

	catalog, err := c.service.Catalog(ctx, int64(id))
	if err == nil {
		reply.Catalog, err = json.Marshal(catalog)
	}

	if err != nil || len(reply.Catalog) > catalogwire.MaxBytes-4096 {
		reply.Error = "catalog unavailable"
		reply.Catalog = nil
		zap.L().Warn("gateway MQTT catalog generation failed")
	}

	reply.Signature = reply.Sign(agent.Password)
	body, err := json.Marshal(reply)
	if err != nil {
		return
	}

	token := client.Publish(catalogwire.ResponseTopic(req.GatewayID, req.RequestID), 1, false, body)
	select {
	case <-ctx.Done():
	case <-token.Done():
		if token.Error() != nil {
			zap.L().Warn("gateway MQTT catalog response failed")
		}
	}
}
