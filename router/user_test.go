package router

import (
	"crypto/md5"
	"fmt"
	"testing"

	"github.com/solamiku/go-utility/crypto"
)

func Test_crypto(t *testing.T) {
	testmd5 := func(pwd string) string {
		m := md5.Sum([]byte(fmt.Sprintf("%s_%s", "ttest", pwd)))
		return string(m[:])
	}
	key := "10KncU7_"

	ori := "admin" + ";" + testmd5("adminnimda")
	fmt.Println("ori ", ori, "byte:", []byte(ori))
	en, _ := crypto.DesECB([]byte(ori), []byte(key), true)
	fmt.Println("after crypto", en)

	bs, _ := crypto.DesECB([]byte(en), []byte(key), false)
	fmt.Println("after uncrypto", bs)

	// fmt.Println(deStr)
}
