package baseresponse

import (
	"fmt"
	"net/http"
	"runtu666/common/shared/baseerror"
	"strconv"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
	"github.com/tealeg/xlsx"
)

type response struct {
	Code    int64       `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Time    int64       `json:"time"`
}

const codeServiceUnavailable = 10001

var (
	serviceUnavailable = "服务器竟然开小差，一会儿再试试吧"
	HeaderFrom         = "from"
)

func GetFrom(r *http.Request) (int64, error) {
	from := r.Header.Get(HeaderFrom)
	i, err := strconv.Atoi(from)
	if err != nil {
		return 0, baseerror.NewDefaultError("header no from !!!")
	}
	return int64(i), err
}

func FormatResponseWithRequest(data interface{}, err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		codeErr, ok := baseerror.FromError(err)
		if ok {
			httpBizError(w, codeErr)
		} else {
			httpServerError(w)
		}
		logx.WithContext(r.Context()).Error(err)
	} else {
		HttpOk(w, data)
	}
}

func HttpParamErrorWithRequest(w http.ResponseWriter, r *http.Request, err error) {
	logx.WithContext(r.Context()).Error(err)
	HttpParamError(w, err.Error())
}

func FormatXlsxResponse(file *xlsx.File, name string, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "applicationnd.ms-excel")
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s.xlsx;filename*=utf-8''%s.xlsx", name, name))
	w.Header().Add("Cache-Control", "max-age=0")
	_ = file.Write(w)
}

func FormatImageResponse(png []byte, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(png)))
	_, _ = w.Write(png)
}

func httpBizError(w http.ResponseWriter, err *baseerror.CodeError) {
	HttpError(w, http.StatusOK, err.Code(), err.Desc())
}

func httpServerError(w http.ResponseWriter) {
	HttpError(w, http.StatusOK, codeServiceUnavailable, serviceUnavailable)
}

func HttpOk(w http.ResponseWriter, data interface{}) {
	response := response{
		Code: 200,
		Data: data,
		Time: time.Now().Unix(),
	}
	if data == nil {
		response.Message = "操作成功"
	}
	httpx.WriteJson(w, http.StatusOK, response)

}
func HttpParamError(w http.ResponseWriter, desc string) {
	httpx.WriteJson(w, http.StatusOK, response{
		Code:    http.StatusBadRequest,
		Message: desc,
		Time:    time.Now().Unix(),
	})
}

func HttpError(w http.ResponseWriter, httpCode int, appCode int64, desc string) {
	httpx.WriteJson(w, httpCode, response{
		Code:    appCode,
		Message: desc,
		Time:    time.Now().Unix(),
	})
}
