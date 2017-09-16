/*
	router authority defined
*/

package router

import (
	"github.com/valyala/fasthttp"
)

// Authority struct
type Authority struct {
	Admin bool
}

//GetAuthority return request's authority
func GetAuthority(ctx *fasthttp.RequestCtx) Authority {
	a := Authority{}
	a.Admin = false
	return a
}
