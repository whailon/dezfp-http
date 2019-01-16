package tools

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

// Padding 补码
func Padding(plainText []byte, blockSize int) []byte {
	n := blockSize - len(plainText)%blockSize
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

// UnPadding 删除填充
func UnPadding(cipherText []byte) []byte {
	//取出密文最后一个字节end
	end := cipherText[len(cipherText)-1]
	//删除填充
	cipherText = cipherText[:len(cipherText)-int(end)]
	return cipherText
}

// TripleDESCBCEncrypt 三重DES加密
// plainText 明文byte数组
// key 密钥
func TripleDESCBCEncrypt(plainText, key []byte) []byte {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		panic(err)
	}
	plainText = Padding(plainText, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	return cipherText
}

// TripleDESCBCDecrypt 三重解密
// cipherText 密文byte数组
// key 密钥与加密一致
func TripleDESCBCDecrypt(cipherText, key []byte) []byte {
	//指定解密算法，返回一个Block接口对象
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		panic(err)
	}
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	//解密
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	//删除填充
	plainText = UnPadding(plainText)
	//返回明文
	return plainText
}
