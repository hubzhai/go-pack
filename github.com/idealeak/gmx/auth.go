package gmx

import (
	"net/http"
	"strings"
)

func auth(rw http.ResponseWriter, req *http.Request) {
	logf("auth : RemoteAddr:%s (query=%s)...\n", req.RemoteAddr, req.URL.String())

	m := req.URL.Query()
	username := m.Get("username")
	if username == "" {
		apiResponse(rw, Tag_ParamIllegal, "username not allow null.", "")
		return
	}

	password := m.Get("password")
	if password == "" {
		apiResponse(rw, Tag_ParamIllegal, "password not allow null.", "")
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

	if username != Cfg.Username || password != Cfg.Password {
		apiResponse(rw, Tag_AuthFailed, nil, "")
		return
	}

	sid := RandomString(32)
	apiResponse(rw, Tag_OK, sid, "")
	//创建一个会话
	s := NewSession(sid)
	if s != nil {
		s.Set(SESSION_AUTH_IP, ip)
		SetSession(s)
	}
}

func init() {
	http.HandleFunc("/auth", auth)
}
