package model

import "nova-factory-server/app/utils/observer/integration/result"

// BtypeGetData 往来单位数据，实现 result.BtypeGetData。
type BtypeGetData struct {
	BtypeCode      string `json:"btypecode"`
	BtypeName      string `json:"btypename"`
	LinkMan        string `json:"linkman"`
	Tel            string `json:"tel"`
	Address        string `json:"address"`
	Remark         string `json:"remark"`
	DefaultElyCode string `json:"defaultelycode"`
	DefaultElyName string `json:"defaultelyname"`
}

func (d *BtypeGetData) GetBtypeCode() string      { return d.BtypeCode }
func (d *BtypeGetData) GetBtypeName() string      { return d.BtypeName }
func (d *BtypeGetData) GetLinkMan() string        { return d.LinkMan }
func (d *BtypeGetData) GetTel() string            { return d.Tel }
func (d *BtypeGetData) GetAddress() string        { return d.Address }
func (d *BtypeGetData) GetRemark() string         { return d.Remark }
func (d *BtypeGetData) GetDefaultElyCode() string { return d.DefaultElyCode }
func (d *BtypeGetData) GetDefaultElyName() string { return d.DefaultElyName }

// BtypeGetResponse 往来单位查询响应，实现 result.BtypeGetResponse。
type BtypeGetResponse struct {
	Code    int64           `json:"code"`
	Message string          `json:"message"`
	Total   int64           `json:"total"`
	Datas   []*BtypeGetData `json:"datas"`
}

func (r *BtypeGetResponse) GetCode() int64     { return r.Code }
func (r *BtypeGetResponse) GetMessage() string { return r.Message }
func (r *BtypeGetResponse) GetTotal() int64    { return r.Total }
func (r *BtypeGetResponse) GetDatas() []result.BtypeGetData {
	out := make([]result.BtypeGetData, len(r.Datas))
	for i, v := range r.Datas {
		out[i] = v
	}
	return out
}
