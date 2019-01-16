package client

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// URL 电子发票路由
const URL = "http://fw1test.shdzfp.com:9000/sajt-shdzfp-sl-http/SvrServlet"

// GlobalInfo 全局数据项
type GlobalInfo struct {
	TerminalCode      string `xml:"terminalCode"`
	AppID             string `xml:"appId"`
	Version           string `xml:"version"`
	InterfaceCode     string `xml:"interfaceCode"`
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

// Data 交换数据
type Data struct {
	Description DataDescription `xml:"dataDescription"`
	Content     string          `xml:"content"`
}

// XMLString xml请求
type XMLString struct {
	XMLName     xml.Name        `xml:"interface"`
	Global      GlobalInfo      `xml:"globalInfo"`
	ReturnState ReturnStateInfo `xml:"returnStateInfo"`
	RequestData Data            `xml:"data"`
}

// ToString 转成字符串
func (c *RequestContent) ToString() string {
	code, _ := xml.Marshal(c)
	requestType := `<REQUEST_FPXXXZ_NEW class='REQUEST_FPXXXZ_NEW'>`
	return strings.Replace(string(code), "<REQUEST_FPXXXZ_NEW>", requestType, -1)
}

// ToString 转成字符串
func (s *XMLString) ToString() string {
	code, _ := xml.Marshal(s)
	interfactType := `<?xml version='1.0' encoding='utf-8'?>
	<interface xmlns='' xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance'
	xsi:schemaLocation='http://www.chinatax.gov.cn/tirip/dataspec/interfaces.xsd'
	version='DZFP1.0'>`
	return strings.Replace(string(code), "<interface>", interfactType, -1)
}

// MakeOut 开具发票
func MakeOut(xmlStr []byte) error {
	//发送请求.
	req, err := http.NewRequest("POST", URL, bytes.NewReader(xmlStr))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/xml")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	var xmlRe ReturnStateInfo
	respBytes, err := ioutil.ReadAll(resp.Body)
	err = xml.Unmarshal(respBytes, &xmlRe)
	if err != nil {
		return err
	}

	if xmlRe.ReturnCode != "0000" {
		return errors.New(xmlRe.ReturnMessage)
	}

	return nil
}
