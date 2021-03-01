package gmx

import (
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

func op(rw http.ResponseWriter, req *http.Request) {

	var body *string
	var tag int
	var msg interface{}
	start := time.Now()
	defer func() {
		var bd string
		if body != nil {
			bd = *body
		}
		logf("op : RemoteAddr:%s (query=%s) body=%s tag=%d msg=%v take=%d\n", req.RemoteAddr, req.URL.String(), bd, tag, msg, time.Now().Sub(start))
	}()

	if !Cfg.EnableOp {
		apiResponse(rw, Tag_OpDisable, nil, "")
		return
	}

	strs := strings.Split(req.RemoteAddr, ":")
	if len(strs) < 2 {
		apiResponse(rw, Tag_AccessDenied, nil, "")
		return
	}

	ip := strs[0]
	if !checkWhiteList(ip) {
		apiResponse(rw, Tag_AccessDenied, nil, "")
		return
	}

	m := req.URL.Query()
	token := m.Get("token")
	if !checkAuth(token, ip) {
		apiResponse(rw, Tag_AuthNeed, nil, "")
		return
	}

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}

	param := string(data)
	body = &param
	hn := m.Get("handler")
	handler := GetOpHandler(hn)
	if handler != nil {
		res := make(chan *ApiResult, 1)
		oe := &OpEvent{h: handler, param: param, res: res, token: token}
		if threadSafeSandBox != nil {
			threadSafeSandBox(oe)
		} else {
			oe.Done()
		}

		//用户透传数据
		ud := m.Get("ud")
		select {
		case scr := <-res:
			tag = scr.Tag
			msg = scr.Msg
			apiResponse(rw, scr.Tag, scr.Msg, ud)
		case <-time.After(time.Second * time.Duration(Cfg.Timeout)):
			tag = Tag_ProcessTimeout
			apiResponse(rw, Tag_ProcessTimeout, "process timeout", ud)
		}
		return
	}
}

var handlerMgr sync.Map

type Handler interface {
	Serve(*OpEvent) *ApiResult
}

type HandlerWrapper func(*OpEvent) *ApiResult

func (hw HandlerWrapper) Serve(event *OpEvent) *ApiResult {
	return hw(event)
}

func RegisteOpHandler(op string, h Handler) {
	handlerMgr.Store(op, h)
}

func GetOpHandler(op string) Handler {
	if hi, ok := handlerMgr.Load(op); ok && hi != nil {
		return hi.(Handler)
	}
	return nil
}

type OpEvent struct {
	h     Handler
	param string
	token string
	res   chan *ApiResult
}

func (oe *OpEvent) Done() error {
	defer func() {
		err := recover()
		if err != nil {
			oe.res <- &ApiResult{Tag: Tag_InternalException, Msg: err}
		}
	}()
	ret := oe.h.Serve(oe)
	oe.res <- ret
	return nil
}

func init() {
	http.HandleFunc("/op", op)
}
