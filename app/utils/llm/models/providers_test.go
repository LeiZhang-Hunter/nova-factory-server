package models

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"charm.land/fantasy"
)

func TestCreateHTTPClientWithTLSConfig(t *testing.T) {
	tests := []struct {
		name         string
		skipVerify   bool
		wantInsecure bool
	}{
		{
			name:         "skip verify disabled",
			skipVerify:   false,
			wantInsecure: false,
		},
		{
			name:         "skip verify enabled",
			skipVerify:   true,
			wantInsecure: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := createHTTPClientWithTLSConfig(tt.skipVerify)

			if client == nil {
				t.Fatal("expected non-nil client")
			}

			// Check if the client has a custom transport when skipVerify is true
			if tt.skipVerify {
				transport, ok := client.Transport.(*http.Transport)
				if !ok {
					t.Fatal("expected *http.Transport when skipVerify is true")
				}

				if transport.TLSClientConfig == nil {
					t.Fatal("expected non-nil TLSClientConfig when skipVerify is true")
				}

				if transport.TLSClientConfig.InsecureSkipVerify != tt.wantInsecure {
					t.Errorf("InsecureSkipVerify = %v, want %v",
						transport.TLSClientConfig.InsecureSkipVerify, tt.wantInsecure)
				}
			}
		})
	}
}

func TestCreateOAuthHTTPClient(t *testing.T) {
	tests := []struct {
		name         string
		accessToken  string
		skipVerify   bool
		wantInsecure bool
	}{
		{
			name:         "oauth with skip verify disabled",
			accessToken:  "test-token",
			skipVerify:   false,
			wantInsecure: false,
		},
		{
			name:         "oauth with skip verify enabled",
			accessToken:  "test-token",
			skipVerify:   true,
			wantInsecure: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := createOAuthHTTPClient(tt.accessToken, tt.skipVerify)

			if client == nil {
				t.Fatal("expected non-nil client")
			}

			// Check that the transport is an oauthTransport
			oauthTransport, ok := client.Transport.(*oauthTransport)
			if !ok {
				t.Fatal("expected *oauthTransport")
			}

			if oauthTransport.accessToken != tt.accessToken {
				t.Errorf("accessToken = %v, want %v", oauthTransport.accessToken, tt.accessToken)
			}

			// Check the base transport when skipVerify is true
			if tt.skipVerify {
				baseTransport, ok := oauthTransport.base.(*http.Transport)
				if !ok {
					t.Fatal("expected base transport to be *http.Transport when skipVerify is true")
				}

				if baseTransport.TLSClientConfig == nil {
					t.Fatal("expected non-nil TLSClientConfig when skipVerify is true")
				}

				if baseTransport.TLSClientConfig.InsecureSkipVerify != tt.wantInsecure {
					t.Errorf("InsecureSkipVerify = %v, want %v",
						baseTransport.TLSClientConfig.InsecureSkipVerify, tt.wantInsecure)
				}
			}
		})
	}
}

func TestProviderConfigTLSSkipVerify(t *testing.T) {
	// Test that ProviderConfig properly stores TLSSkipVerify
	config := &ProviderConfig{
		ModelString:   "test:model",
		TLSSkipVerify: true,
	}

	if !config.TLSSkipVerify {
		t.Error("expected TLSSkipVerify to be true")
	}
}

func TestCreateProviderDeepSeekAPIExample(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	result, err := CreateProvider(ctx, &ProviderConfig{
		ModelString:    "deepseek/deepseek-chat",
		ProviderAPIKey: "sk-885afc3090024bf8939c81b1a3429775",
		ProviderURL:    "https://api.deepseek.com/v1",
		MaxTokens:      120,
		SystemPrompt:   "\n      <role>\n       你是一个销售订单数据格式转换专家，可以把用户输入的销售订单数据转化为管家婆全渠道的数据格式。\n      </role>\n      <instructions>\n        指示：\n\n仅输出管家婆全渠道JSON格式数据，不输出任何解释、注释或额外文本\n\n输出结构包含以下字段：\n\norders数组（每个订单对象包含）：\n\ntid：网店订单编号，字符串，必填\n\nweight：重量，数字\n\nsize：尺寸，数字\n\nbuyernick：买家账号，字符串\n\nbuyermessage：卖家留言，字符串\n\nsellermemo：卖家备注，字符串\n\ntotal：订单总金额，浮点数，包含运费\n\nprivilege：优惠金额，浮点数，必填\n\npostfee：运费，浮点数，必填\n\nreceivername：收货人名称，字符串，必填\n\nreceiverstate：收货省，字符串，必填\n\nreceivercity：收货市，字符串，必填\n\nreceiverdistrict：收货区，字符串，必填\n\nreceiveraddress：收货地址，字符串，必填\n\nreceiverphone：收货电话，字符串\n\nreceivermobile：收货人手机号，字符串，必填\n\ncreated：订单创建时间，字符串，格式yyyy-MM-dd HH:mm:ss，必填\n\nstatus：订单状态，字符串，必填，转换为NoPay、Payed、Sended、TradeSuccess、TradeClosed、PartSend\n\ntype：订单类型，字符串，必填，货到付款转换为Cod，非货到付款转换为NoCod\n\ninvoicename：发票抬头，字符串\n\nsellerflag：卖家旗帜，字符串\n\npaytime：付款时间，字符串，格式yyyy-MM-dd HH:mm:ss\n\nlogistbtypecode：物流公司编码，字符串\n\nlogistbillcode：物流单号，字符串\n\nbtypecode：往来单位编码，字符串\n\ndetails：订单商品信息数组，必填\n\naccounts：订单账户信息数组\n\ndetails数组（每个商品对象包含）：\n\noid：网店订单明细编号，字符串，必填\n\nbarcode：商品条码，字符串\n\neshopgoodsid：网店商品ID，字符串，方案A必填\n\nouteriid：网店商家编码，字符串\n\neshopgoodsname：网店商品名称，字符串，必填，仅填入商品主体名称不含规格\n\neshopskuid：网店商品SKUID，字符串，方案A且商品有SKU则必填\n\neshopskuname：网店商品SKU名称，字符串，仅填入规格属性名称\n\nnumiid：商品ID，长整数，方案B必填\n\nskuid：规格ID，长整数，方案B必填\n\nnum：基本单位数量，数字，必填，不能为0\n\npayment：商品总额，浮点数，必填\n\npicpath：商品图片路径，字符串\n\nweight：重量，数字\n\nsize：尺寸，数字\n\nunitid：销售单位ID，长整数\n\nunitqty：销售单位数量，数字\n\naccounts数组（每个账户对象包含）：\n\nfinanceCode：账户编码，字符串，必填\n\ntotal：收款金额，数字，必填\n\n方案判断规则：\n\n若输入包含numiid字段或sku_id字段，按方案B处理：将值填入numiid和skuid，eshopgoodsid和eshopskuid设为null\n\n若输入包含eshopgoodsid字段或eshopskuid字段，按方案A处理：将值填入eshopgoodsid和eshopskuid，numiid和skuid设为null\n\n商品名称拆分规则：\n\neshopgoodsname仅填入商品主体名称，不含规格属性\n\neshopskuname仅填入规格属性名称\n\n若输入只有一个商品全称字段，按常见格式将主体名称和规格拆分后分别填入\n\n若输入已分别提供商品名和规格名，直接对应填入\n\n输入中不存在的字段输出null，但必须保留字段位置\n      </instructions>",
	})
	if err != nil {
		t.Fatalf("CreateProvider failed: %v", err)
	}
	if result == nil || result.Model == nil {
		t.Fatal("provider result or model is nil")
	}
	messageV := fantasy.NewUserMessage("请根据下面的信息生成销售订单：\n数量\t含税单价\t货期\t起订量\nZHIDE/质德 聚氨酯IDU密封圈 IDU天蓝色-50*58*10 1个 销售单位：个\t质德\tIDU天蓝色-50*58*10\t1000\t2.7\t7-10天\t100\nZHIDE/质德 骨架油封 氟橡胶 25X45X10 TC型 1个 销售单位：个\t质德\t25X45X10TC型\t1000\t6.24\t7-10天\t100\nZHIDE/质德 菱形八角胶 GR75(81×161)mm 1个 销售单位：个\t质德\tGR75(81×161)mm\t1000\t56.25\t7-10天\t100\nZHIDE/质德 菱形八角胶 GR48(51×104)mm 1个 销售单位：个\t质德\tGR48(51×104)mm\t1000\t14.25\t7-10天\t100\nZHIDE/质德 菱形八角胶 GR38(38×80)mm 1个 销售单位：个\t质德\tGR38(38×80)mm\t1000\t7.5\t7-10天\t100\nZHIDE/质德 聚氨酯弹性垫 18×35mm 1个 销售单位：个\t质德\t18×35mm\t1000\t0.8\t7-10天\t100\nZHIDE/质德 TG丁腈橡胶骨架油封 100×130×12mm NBR 硬度70° 黑色 1个 销售单位：个\t质德\t100×130×12mmNBR 硬度70° 黑色\t1000\t8.19\t7-10天\t100\nZHIDE/质德 聚四氟乙烯垫片 10×18×2mm PTFE 白色 温度250℃ 1个 销售单位：个\t质德\t10×18×2mmPTFE 白色 温度250℃\t1000\t0.9\t7-10天\t100\nZHIDE/质德 丁腈橡胶骨架油封 70×110×12mm TC型 1个 销售单位：个\t质德\t70×110×12mmTC型\t1000\t5\t7-10天\t100\nZHIDE/质德 氟胶骨架油封 35×52×7mm 1个 销售单位：个\t质德\t35×52×7mm\t1000\t7.91\t7-10天\t100\nZHIDE/质德 氟胶骨架油封 30×47×7mm TC FKM 硬度75° 棕色 1个 销售单位：个\t质德\t30×47×7mmTC FKM 硬度75° 棕色\t1000\t6.95\t7-10天\t100\nZHIDE/质德 氟胶O形圈 内径φ8×1.8mm FKM 1个 销售单位：个\t质德\t内径φ8×1.8mmFKM\t1000\t0.24\t7-10天\t100\nZHIDE/质德 氟胶O形圈 外径φ8×2.4mm FKM 1个 销售单位：个\t质德\t外径φ8×2.4mmFKM\t1000\t0.21\t7-10天\t100")

	resp, err := result.Model.Generate(ctx, fantasy.Call{
		Prompt: fantasy.Prompt{
			messageV,
		},
	})
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if resp.Content.Text() == "" {
		t.Fatal("empty response content")
	}
	fmt.Println(resp.Content.Text())
}
