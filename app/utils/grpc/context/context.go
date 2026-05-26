package context

import (
	"context"
	"errors"
	"nova-factory-server/app/constant/agent"
	"strconv"
)

func GetGatewayId(ctx context.Context) (int64, error) {
	gatewayId := ctx.Value(agent.GATEWAYID)
	//if !ok {
	//	return 0, errors.New("gateway_id not exist")
	//}

	gatewayIdStr, ok := gatewayId.([]string)
	if !ok {
		return 0, errors.New("gateway_id not exist")
	}

	if len(gatewayIdStr) == 0 {
		return 0, errors.New("gateway_id is empty")
	}

	return strconv.ParseInt(gatewayIdStr[0], 10, 64)
}

func GetUsername(ctx context.Context) (string, error) {
	username := ctx.Value(agent.USERNAME)
	//if !ok {
	//	return 0, errors.New("gateway_id not exist")
	//}

	usernameStr, ok := username.([]string)
	if !ok {
		return "", errors.New("gateway_id not exist")
	}

	if len(usernameStr) == 0 {
		return "", errors.New("gateway_id is empty")
	}

	return usernameStr[0], nil
}

func GetPassword(ctx context.Context) (string, error) {
	password := ctx.Value(agent.PASSWORD)
	//if !ok {
	//	return 0, errors.New("gateway_id not exist")
	//}

	passwordStr, ok := password.([]string)
	if !ok {
		return "", errors.New("gateway_id not exist")
	}

	if len(passwordStr) == 0 {
		return "", errors.New("gateway_id is empty")
	}

	return passwordStr[0], nil
}
