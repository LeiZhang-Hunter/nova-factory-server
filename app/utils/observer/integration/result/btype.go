package result

// BtypeGetResponse 往来单位查询响应接口，对应 emall.btype.get 返回。
type BtypeGetResponse interface {
	GetCode() int64
	GetMessage() string
	GetTotal() int64
	GetDatas() []BtypeGetData
}

// BtypeGetData 往来单位数据接口，对应 datas[{btypecode, ..., defaultelyname}]。
type BtypeGetData interface {
	GetBtypeCode() string
	GetBtypeName() string
	GetLinkMan() string
	GetTel() string
	GetAddress() string
	GetRemark() string
	GetDefaultElyCode() string
	GetDefaultElyName() string
}
