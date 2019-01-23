package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_encrypt()
func Test_MakeOut(t *testing.T) {
	global := &GlobalInfo{
		TerminalCode:      "0",
		Version:           "2.0",
		AppID:             "ZZS_PT_DZFP",
		UserName:          "111MFWIK",
		PassWord:          "12345678909oyKs7cVo1yYzkuisP9bhA==",
		TaxpayerID:        "310101000000090",
		AuthorizationCode: "3100000090",
		RequestCode:       "111MFWIK",
		RequestTime:       "2016-11-14 14:11:30 301",
		DataExchangeID:    "111MFWIK20180508154134675",
	}
	data := &RequestData{
		Content: &MakeRequestContent{
			FPTXX: FPTXX{
				FPQQLSH:   "111MFWIK201805081541348952",
				DSPTBM:    "111MFWIK",
				NSRSBH:    "310101000000090",
				NSRMC:     "上海爱信诺航天信息有限公司90",
				DKBZ:      "0",
				KPXM:      "A10",
				BMBBBH:    "1.0",
				XHFNSRSBH: "310101000000090",
				XHFMC:     "上海爱信诺航天信息有限公司90",
				XHFDZ:     "天津市-武清区",
				XHFDH:     "13810189561",
				XHFYHZH:   "saleAccountNo",
				GHFMC:     "托管中心",
				GHFQYLX:   "01",
				GHFDZ:     "购买方地址",
				GHFGDDH:   "18810987722",
				GHFSJ:     "18810987721",
				GHFYHZH:   "999",
				KPY:       "一二三四",
				FHR:       "复核人",
				SKY:       "收款人",
				KPRQ:      "2016-11-14 14:12:15",
				KPLX:      "1",
				CZDM:      "10",
				QDBZ:      "0",
				TSCHBZ:    "0",
				KPHJJE:    "20.00",
			},
			XMXXS: XMXXS{
				Items: []XMXX{
					XMXX{
						XMMC:  "A10",
						XMDW:  "双",
						XMSL:  "2",
						HSBZ:  "1",
						SPBM:  "1010101030000000000",
						XMJE:  "20",
						SL:    "0.06",
						FPHXZ: "0",
					},
				},
			},
			DDXX: DDXX{
				DDH: "201805081541348952",
			},
		},
	}
	client := NewBillClient()
	client.Global = global
	client.RequestData = data
	client.Key = "9oyKs7cVo1yYzkuisP9bhA=="
	result, err := client.MakeOut()
	if assert.Nil(t, err) {
		if data, ok := result.(*ReturnStateInfo); ok {
			assert.Equal(t, "0000", data.ReturnCode)
		}
	}
}

func Test_Download(t *testing.T) {
	data := &RequestData{
		Content: &DownloadRequestContent{
			FPQQLSH: "111MFWIK201805081541348952",
			DSPTBM:  "111MFWIK",
			NSRSBH:  "310101000000090",
			DDH:     "201805081541348952",
			PDFXZFS: "2",
		},
	}

	client := NewBillClient()
	client.Global = &GlobalInfo{
		TerminalCode:      "0",
		AppID:             "ZZS_PT_DZFP",
		Version:           "2.0",
		InterfaceCode:     "ECXML.FPXZ.CX.E_INV",
		UserName:          "111MFWIK",
		PassWord:          "12345678909oyKs7cVo1yYzkuisP9bhA==",
		TaxpayerID:        "310101000000090",
		AuthorizationCode: "3100000090",
		RequestCode:       "111MFWIK",
		RequestTime:       time.Now().Format("2006-01-02 15:04:05 700"),
		ResponseCode:      "121",
		DataExchangeID:    "111MFWIK20180508154134634",
	}
	client.RequestData = data
	client.ReturnState = &ReturnStateInfo{}
	client.Key = "9oyKs7cVo1yYzkuisP9bhA=="
	result, err := client.Download()
	if assert.Nil(t, err) {
		if data, ok := result.(*ResponseContent); ok {
			assert.NotNil(t, data.Fpqqlsh)
		}
	}
}
