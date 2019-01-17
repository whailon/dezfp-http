package client

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"

	"dezfp-http/tools"
)

// URL 电子发票路由
const (
	// URL 请求地址
	URL = "http://fw1test.shdzfp.com:9000/sajt-shdzfp-sl-http/SvrServlet"

	// 接口编码
	// MakeOutBill 开具发票
	MakeOutBill = "ECXML.FPKJ.BC.E_INV"
	// DownloadBill 下载发票
	DownloadBill = "ECXML.FPXZ.CX.E_INV"
)

// GlobalInfo 全局数据项
type GlobalInfo struct {
	TerminalCode  string `xml:"terminalCode"`
	AppID         string `xml:"appId"`
	Version       string `xml:"version"`
	InterfaceCode string `xml:"interfaceCode"`

	UserName          string `xml:"userName"`
	PassWord          string `xml:"passWord"`
	TaxpayerID        string `xml:"taxpayerId"`
	AuthorizationCode string `xml:"authorizationCode"`
	RequestCode       string `xml:"requestCode"`
	RequestTime       string `xml:"requestTime"`
	ResponseCode      string `xml:"responseCode"`
	DataExchangeID    string `xml:"dataExchangeId"`
}

// ReturnStateInfo 数据交换请求返回状态信息
type ReturnStateInfo struct {
	ReturnCode    string `xml:"returnCode"`
	ReturnMessage string `xml:"returnMessage"`
}

// DataDescription 交换数据属性描述
type DataDescription struct {
	ZipCode     string `xml:"zipCode"`
	EncryptCode string `xml:"encryptCode"`
	CodeType    string `xml:"codeType"`
}

// RequestContent 交换数据内容
type RequestContent struct {
	XMLName xml.Name `xml:"REQUEST_FPXXXZ_NEW"`
	Fpqqlsh string   `xml:"FPQQLSH"`
	Dsptbm  string   `xml:"DSPTBM"`
	Nsrsbh  string   `xml:"NSRSBH"`
	Ddh     string   `xml:"DDH"`
	PdfXzfs string   `xml:"PDF_XZFS"`
}

// ResponseContent 交换数据内容
type ResponseContent struct {
	XMLName       xml.Name `xml:"REQUEST_FPKJXX_FPJGXX_NEW"`
	Fpqqlsh       string   `xml:"FPQQLSH"`
	Ddh           string   `xml:"DDH"`
	Kplsh         string   `xml:"KPLSH"`
	Fwm           string   `xml:"FWM"`
	Ewn           string   `xml:"EWM"`
	FpzlDm        string   `xml:"FPZL_DM"`
	FpDm          string   `xml:"FP_DM"`
	FpHm          string   `xml:"FP_HM"`
	KPRQ          string   `xml:"KPRQ"`
	KPLX          string   `xml:"KPLX"`
	HJBHSJE       string   `xml:"HJBHSJE"`
	KPHJSE        string   `xml:"KPHJSE"`
	PdfFile       string   `xml:"PDF_FILE"`
	PdfURL        string   `xml:"PDF_URL"`
	CZDM          string   `xml:"CZDM"`
	RETURNCODE    string   `xml:"RETURNCODE"`
	RETURNMESSAGE string   `xml:"RETURNMESSAGE"`
}

// Data 交换数据
type Data struct {
	Description *DataDescription `xml:"dataDescription"`
	// EncryptContent 根据Content加密生成
	EncryptContent string `xml:"content"`
	// Content 交换数据内容明文，必需
	Content *RequestContent `xml:"-"`
}

type ResponseData struct {
	Content *ResponseContent `xml:"content"`
}

// BillClient xml请求
type BillClient struct {
	XMLName     xml.Name         `xml:"interface"`
	Global      *GlobalInfo      `xml:"globalInfo"`
	ReturnState *ReturnStateInfo `xml:"returnStateInfo"`
	RequestData *Data            `xml:"data"`
	Key         string           `xml:"-"`
}

type SecBillClient struct {
	XMLName      xml.Name      `xml:"interface"`
	ResponseData *ResponseData `xml:"data"`
}

// ToString 转成字符串
func (c *Data) encrypt(key []byte) {
	code, _ := xml.Marshal(c.Content)
	requestType := `<REQUEST_FPXXXZ_NEW class='REQUEST_FPXXXZ_NEW'>`
	str := strings.Replace(string(code), "<REQUEST_FPXXXZ_NEW>", requestType, -1)
	res := util.TripleDESCBCEncrypt([]byte(str), key)
	c.EncryptContent = string(res)
}

func (c *Data) defaultDescription() {
	c.Description = &DataDescription{
		ZipCode:     "0",
		EncryptCode: "1",
		CodeType:    "3DES",
	}
}

// ToString 转成字符串
func (s *BillClient) toString() string {
	code, _ := xml.Marshal(s)
	interfactType := `<?xml version='1.0' encoding='utf-8'?>
	<interface xmlns='' xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance'
	xsi:schemaLocation='http://www.chinatax.gov.cn/tirip/dataspec/interfaces.xsd'
	version='DZFP1.0'>`
	return strings.Replace(string(code), "<interface>", interfactType, -1)
}

// NewBillClient 实例化发票对象
func NewBillClient() *BillClient {
	return &BillClient{
		Global: &GlobalInfo{
			Version: "2.0",
		},
		RequestData: &Data{},
	}
}

// MakeOut 开具发票
func (s *BillClient) MakeOut() (interface{}, error) {
	// 发起请求
	return s.doAction(MakeOutBill)
}

// Download 开具发票
func (s *BillClient) Download() (interface{}, error) {
	// 发起请求
	return s.doAction(DownloadBill)
}

func (s *BillClient) init(interfaceCode string) []byte {
	s.setInterfaceCode(interfaceCode)
	s.RequestData.encrypt([]byte(s.Key))
	if s.RequestData.Description == nil {
		s.RequestData.defaultDescription()
	}

	return []byte(s.toString())
}

func (s *BillClient) getInterfaceCode() string {
	return s.Global.InterfaceCode
}

func (s *BillClient) setInterfaceCode(code string) {
	s.Global.InterfaceCode = code
}

func (s *BillClient) doAction(interfaceCode string) (interface{}, error) {
	// BillClient初始化
	xmlStr := s.init(interfaceCode)
	//发送请求.
	req, err := http.NewRequest("POST", URL, bytes.NewReader(xmlStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/xml")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if s.getInterfaceCode() == MakeOutBill {
		var xmlRe BillClient
		err = xml.Unmarshal(body, &xmlRe)
		return xmlRe.ReturnState, err
	}
	var xmlRe SecBillClient
	err = xml.Unmarshal(body, &xmlRe)
	return xmlRe.ResponseData, err
}
