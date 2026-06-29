package api

type ExpressTraces interface {
	AcceptTime() string    // 轨迹发生时间，示例：2021-01-01 09:00:00
	AcceptStation() string // 轨迹描述
	Location() string      // 当前城市
	Action() string        // 增值物流状态： 0-暂无轨迹信息 1-已揽收 2-在途中 201-到达派件城市 202-派件中 211-已放入快递柜或驿站 3-已签收 301-正常签收 302-派件异常后最终签收 304-代收签收 311-快递柜或驿站签收 4-问题件 401-发货无信息 402-超时未签收 403-超时未更新 404-拒收(退件) 405-派件异常 406-退货签收 407-退货未签收 412-快递柜或驿站超时未取 10-待揽件
	Remark() string        // 备注
}

type ExpressQueryResult interface {
	OrderCode() string
	ShipperCode() string
	LogisticCode() string // 快递单号
	Callback() string     // 用户自定义回传字段
	Success() bool        // 成功与否(true/false)
	Reason() string       // 失败原因
	State() string        // 普通物流状态：0-暂无轨迹信息 1-已揽收 2-在途中 3-签收 4-问题件
	StateEx() string      // 增值物流状态： 0-暂无轨迹信息 1-已揽收 2-在途中 201-到达派件城市 202-派件中 211-已放入快递柜或驿站 3-已签收 301-正常签收 302-派件异常后最终签收 304-代收签收 311-快递柜或驿站签收 4-问题件 401-发货无信息 402-超时未签收 403-超时未更新 404-拒收(退件) 405-派件异常 406-退货签收 407-退货未签收 412-快递柜或驿站超时未取 10-待揽件
	Location() string     // 所在城市
	Traces() []ExpressTraces
	Station() string     // 派件网点的名称
	StationTel() string  // 派件网点的电话
	StationAdd() string  // 派件网点的地址
	DeliveryMan() string // 派件快递员
	DeliveryManTel() string
	NextCity() string
	GetStateName() string
	OriginInfo() string
	SignedTime() string
}
