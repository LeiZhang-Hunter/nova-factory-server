package provider

import (
	"errors"
	"fmt"
	"sync"
)

var (
	registryMu sync.RWMutex
	providers  = make(map[int]PaymentMethod)
)

// RegisterPaymentMethod 注册支付方法。
func RegisterPaymentMethod(p PaymentMethod) {
	registryMu.Lock()
	defer registryMu.Unlock()
	providers[p.Channel()] = p
}

// GetPaymentMethod 按支付通道获取支付方法。
func GetPaymentMethod(payChannel int) (PaymentMethod, error) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	p, ok := providers[payChannel]
	if !ok {
		return nil, errors.New(fmt.Sprintf("不支持的支付通道: %d", payChannel))
	}
	return p, nil
}
