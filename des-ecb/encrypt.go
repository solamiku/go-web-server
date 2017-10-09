package ldes

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"io"
)

// 用zlib压缩
func ZlibEncode(src []byte) []byte {
	var bf bytes.Buffer
	//	w := zlib.NewWriter(&bf)
	w, err := zlib.NewWriterLevel(&bf, zlib.BestCompression)
	if err != nil {
		return []byte("")
	}
	w.Write(src)
	w.Close()
	return bf.Bytes()
}

// 用zlib解压缩
func ZlibDecode(src []byte) ([]byte, error) {
	br := bytes.NewReader(src)
	var out bytes.Buffer
	r, err := zlib.NewReader(br)
	if err != nil {
		return nil, err
	}
	io.Copy(&out, r)
	r.Close()
	return out.Bytes(), nil
}

// zlib压缩->base64编码
//	ZlibBase64Encode(src) use stander base64-encoding
//	or ZlibBase64Encode(src, base64.URLEncoding)  use url-base64-encoding
func ZlibBase64Encode(src []byte, custom_b64enc ...*base64.Encoding) []byte {
	b64enc := base64.StdEncoding
	if len(custom_b64enc) > 0 {
		b64enc = custom_b64enc[0]
	}
	var bf bytes.Buffer
	w0 := base64.NewEncoder(b64enc, &bf)
	w, err := zlib.NewWriterLevel(w0, zlib.BestCompression)
	if err != nil {
		return []byte("")
	}
	w.Write(src)
	w.Close()
	w0.Close()
	return bf.Bytes()
}

// base64解码->zlib解压
func Base64ZlibDecode(src []byte, custom_b64enc ...*base64.Encoding) ([]byte, error) {
	b64enc := base64.StdEncoding
	if len(custom_b64enc) > 0 {
		b64enc = custom_b64enc[0]
	}
	br := bytes.NewReader(src)
	b64 := base64.NewDecoder(b64enc, br)
	var out bytes.Buffer
	r, err := zlib.NewReader(b64)
	if err != nil {
		return nil, err
	}
	io.Copy(&out, r)
	r.Close()
	return out.Bytes(), nil
}

func Base64Encode(src []byte, custom_b64enc ...*base64.Encoding) []byte {
	b64enc := base64.StdEncoding
	if len(custom_b64enc) > 0 {
		b64enc = custom_b64enc[0]
	}
	var bf bytes.Buffer
	w0 := base64.NewEncoder(b64enc, &bf)
	w0.Write(src)
	w0.Close()
	return bf.Bytes()
}
func Base64Decode(src []byte, custom_b64enc ...*base64.Encoding) []byte {
	b64enc := base64.StdEncoding
	if len(custom_b64enc) > 0 {
		b64enc = custom_b64enc[0]
	}
	br := bytes.NewReader(src)
	b64 := base64.NewDecoder(b64enc, br)
	var out bytes.Buffer
	io.Copy(&out, b64)
	return out.Bytes()
}

// 包0协议 打包
// 压缩 > 加密？ > base64
// pwd为8字节启用加密
func Pack0Encode(src []byte, pwd []byte) []byte {
	// 压缩
	s1 := ZlibEncode(src)
	//	fmt.Println("s0 ", len(s0))
	// 是否加密
	if len(pwd) == 8 {
		//		var err error
		s1, _ = DesEcbEncryptPkcs5(s1, pwd)
	}
	//	fmt.Println("des ", len(s1))
	// base64
	return Base64Encode(s1)
}

// 包0协议 解包
// base64 > 解密？ > 解压
// pwd为8字节启用解密
func Pack0Decode(src []byte, pwd []byte) []byte {
	s1 := Base64Decode(src)
	if len(s1) == 0 || len(s1)%8 != 0 {
		return []byte{}
	}
	// 是否解密
	if len(pwd) == 8 {
		s1, _ = DesEcbDecryptPkcs5(s1, pwd)
	}
	// 解压
	s1, _ = ZlibDecode(s1)
	return s1
}
