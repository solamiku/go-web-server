package ldes

import (
	"crypto/des"
	"fmt"
)

func DesEcbEncryptPkcs5(originData, key []byte) ([]byte, error) {
	if len(originData) == 0 {
		return []byte{}, nil
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	originData = PKCS5Padding(originData, block.BlockSize())
	blockMode := NewECBEncrypter(block)
	crypted := make([]byte, len(originData))
	blockMode.CryptBlocks(crypted, originData)
	return crypted, nil
}

func DesEcbDecryptPkcs5(crypted, key []byte) ([]byte, error) {
	if len(crypted) == 0 {
		return []byte{}, nil
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(crypted)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("error size")
	}
	blockMode := NewECBDecrypter(block)
	originData := make([]byte, len(crypted))
	blockMode.CryptBlocks(originData, crypted)
	originData = PKCS5UnPadding(originData)
	return originData, nil
}

func DesEcbEncryptZf(originData, key []byte) ([]byte, error) {
	if len(originData) == 0 {
		return []byte{}, nil
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	originData = ZeroPadding(originData, block.BlockSize())
	blockMode := NewECBEncrypter(block)
	crypted := make([]byte, len(originData))
	blockMode.CryptBlocks(crypted, originData)
	return crypted, nil
}

func DesEcbDecryptZf(crypted, key []byte) ([]byte, error) {
	if len(crypted) == 0 {
		return []byte{}, nil
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(crypted)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("error size")
	}
	blockMode := NewECBDecrypter(block)
	originData := make([]byte, len(crypted))
	blockMode.CryptBlocks(originData, crypted)
	originData = ZeroUnPadding(originData)
	return originData, nil
}
