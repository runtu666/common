package httpresp

import (
	"net/http"
	"strconv"

	"gitee.com/xunmeng_1/common-go/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	codeServerError    = 1
	serviceUnavailable = "服务器竟然开小差，一会儿再试试吧"
)

type errResp struct {
	Code int    `json:"code"`
	Desc string `json:"desc,omitempty"`
}

func (e *errResp) Error() string {
	return strconv.Itoa(e.Code) + ":" + e.Desc
}

func Http(w http.ResponseWriter, r *http.Request, data interface{}, err error) {
	if err != nil {
		HttpErr(w, r, err)
	} else {
		HttpOkJson(w, r, data)
	}
}

func HttpErr(w http.ResponseWriter, r *http.Request, err error) {
	codeErr, ok := errorx.FromError(err)
	if ok {
		httpx.WriteJson(w, codeErr.Status(), errResp{
			Code: codeErr.Code(),
			Desc: codeErr.Error(),
		})
	} else {
		httpx.WriteJson(w, http.StatusInternalServerError, errResp{
			Code: codeServerError,
			Desc: serviceUnavailable,
		})
		logx.WithContext(r.Context()).Error(err)
	}
}

func HttpOk(w http.ResponseWriter, r *http.Request) {
	httpx.Ok(w)
}

func HttpOkJson(w http.ResponseWriter, r *http.Request, data interface{}) {
	if data == nil {
		HttpOk(w, r)
	} else {
		httpx.OkJson(w, data)
	}
}
