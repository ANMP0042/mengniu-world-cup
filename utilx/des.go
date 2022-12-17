/**
 * @Author: YMBoom
 * @Description:
 * @File:  des
 * @Version: 1.0.0
 * @Date: 2022/12/05 18:39
 */
package utilx

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
	"errors"
	"fmt"
)

func DesRK(rk, desKey string) string {
	data, err := base64.StdEncoding.DecodeString(rk)
	if err != nil {
		return ""
	}

	key := []byte(desKey)[0:8]
	res := DesECBDecrypter(data, key)
	return string(res)
}

func DesECBDecrypter(data, key []byte) []byte {
	//NewCipher创建一个新的加密块
	block, err := des.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return nil
	}

	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		//Encrypt加密第一个块，将其结果保存到dst
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}

	// pkcs5填充
	out = pkcs5UnPadding(out)

	return out
}

//pkcs5减码算法
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func DesECBEncrypt(data, key []byte) ([]byte, error) {
	//NewCipher创建一个新的加密块
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	bs := block.BlockSize()
	data = Pkcs5Padding(data, bs)
	if len(data)%bs != 0 {
		return nil, errors.New("need a multiple of the blocksize")
	}

	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		//Encrypt加密第一个块，将其结果保存到dst
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

func Pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
