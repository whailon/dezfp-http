package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_encrypt()
func Test_MakeOut(t *testing.T) {
	data := &Data{
		Content: &RequestContent{
			Fpqqlsh: "fpqqlsh",
			Dsptbm:  "dsptbm",
		},
	}
	client := NewBillClient()
	client.RequestData = data
	client.Key = "123456781234567812345678"
	_, err := client.MakeOut()
	assert.Nil(t, err)
}

func Test_Download(t *testing.T) {
	data := &Data{
		Content: &RequestContent{
			Fpqqlsh: "fpqqlsh",
			Dsptbm:  "dsptbm",
		},
	}
	client := NewBillClient()
	client.RequestData = data
	client.Key = "123456781234567812345678"
	_, err := client.Download()
	assert.Nil(t, err)
}
