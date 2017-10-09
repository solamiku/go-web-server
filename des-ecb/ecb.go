// des_ecb project main.go
package ldes

import (
	"bytes"
	"crypto/cipher"
)

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

//返回ECB方式的加密器
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("dec_ecb: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("dec_ecb: output smaller than input")
	}

	for len(src) > 0 {
		// Write to the dst
		x.b.Encrypt(dst[:x.blockSize], src[:x.blockSize])

		// Move to the next block
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

//返回ECB方式的解密器
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("dec_ecb: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("dec_ecb: output smaller than input")
	}

	if len(src) == 0 {
		return
	}

	for len(src) > 0 {
		// Write to the dst
		x.b.Decrypt(dst[:x.blockSize], src[:x.blockSize])

		// Move to the next block
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

//末尾填充0，bytes长度保持在blockSize的倍数
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

//末尾去除0
func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

//用末尾需要填充的个数的值来填充末尾，当刚好时填充blockSize个blockSize
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去除PKCS5填充的值
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	if length < unpadding || unpadding > 8 || unpadding <= 0 {
		return []byte{}
	}
	return origData[:(length - unpadding)]
}

//根据给定的长度截取bytes
func UnPaddingByLength(origData []byte, length int) []byte {
	olen := len(origData)
	if length > olen {
		panic("dec_ecb : unpadding length is longger than origin data")
	}
	return origData[:(olen - length)]
}
