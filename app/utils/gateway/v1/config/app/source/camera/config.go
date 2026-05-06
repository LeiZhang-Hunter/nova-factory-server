package camera

import "nova-factory-server/app/utils/gateway/v1/config/app/source/camera/ipc"

type Config struct {
	Devices    []ipc.Device     `yaml:"devices"`
	GrpcClient GrpcClientConfig `yaml:"grpc_client"`
}

type GrpcClientConfig struct {
	Enabled         bool     `yaml:"enabled"`
	Address         string   `yaml:"address"`
	Node            string   `yaml:"node"`
	SubscribeMethod string   `yaml:"subscribe_method"`
	AckMethod       string   `yaml:"ack_method"`
	Go2RTCBase      string   `yaml:"go2rtc_base"`
	ICEServers      []string `yaml:"ice_servers"`
	ICEUsername     string   `yaml:"ice_username"`
	ICECredential   string   `yaml:"ice_credential"`
	WebRTCPortMin   uint16   `yaml:"webrtc_port_min"`
	WebRTCPortMax   uint16   `yaml:"webrtc_port_max"`
}
