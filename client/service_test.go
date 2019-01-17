package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_encrypt()
func Test_encrypt(t *testing.T) {
	data := &Data{
		Content: &RequestContent{
			Fpqqlsh: "fpqqlsh",
			Dsptbm:  "dsptbm",
		},
	}
	client := NewBillClient()
	client.RequestData = data
	client.Key = "123456781234567812345678"
	err := client.MakeOut()
	assert.NotNil(t, err)
}
