package gmx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	Tag_OK                = 0
	Tag_ParamIllegal      = 1
	Tag_AuthFailed        = 2
	Tag_AuthNeed          = 3
	Tag_AccessDenied      = 4
	Tag_ProcessTimeout    = 5
	Tag_NoFindValue       = 6
	Tag_NoFindField       = 7
	Tag_MapInvalid        = 8
	Tag_NoMoreData        = 9
	Tag_SliceIsNil        = 10
	Tag_NoFindMethod      = 11
	Tag_OpDisable         = 12
	Tag_InternalException = 13
)

type ApiResult struct {
	Tag      int
	Msg      interface{}
	UserData string
}

func apiResponse(rw http.ResponseWriter, tag int, msg interface{}, ud string) bool {
	result := &ApiResult{
		Tag:      tag,
		Msg:      msg,
		UserData: ud,
	}

	data, err := json.MarshalIndent(&result, "", " ")
	if err != nil {
		debugln("apiResponse Marshal error:", err)
		return false
	}

	dataLen := len(data)
	rw.Header().Set("Content-Length", fmt.Sprintf("%v", dataLen))
	rw.WriteHeader(http.StatusOK)
	pos := 0
	for pos < dataLen {
		writeLen, err := rw.Write(data[pos:])
		if err != nil {
			debugln("apiResponse SendData error:", err, " data=", string(data[:]), " pos=", pos, " writelen=", writeLen, " dataLen=", dataLen)
			return false
		}
		pos += writeLen
	}

	return true
}
