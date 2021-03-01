package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/idealeak/goserver.v3/core/admin"
	"github.com/idealeak/goserver.v3/core/logger"
	"github.com/idealeak/goserver.v3/core/module"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	STATE_OK  = 1
	STATE_ERR = 0
)

const (
	RESPONSE_STATE     = "State"
	RESPONSE_ERRMSG    = "ErrMes"
	RESPONSE_PAGECOUNT = "PageCount"
	RESPONSE_PAGENO    = "PageNo"
	RESPONSE_TOTAL     = "Total"
	RESPONSE_DATA      = "Data"
)

func (clm *APIHttpData) Init() {
	clm.RouteMap["/api/test/47"] = "http://127.0.0.1:9898/api/Game/AddCoinById?ts=1575889767&sign=123"
	clm.RouteMap["/api/test/240"] = "http://192.168.1.240:9898/api/Game/AddCoinById?ts=1575889767&sign=123"
}
func init() {
	admin.MyAdminApp.Route("/api/test/47", WorldSrvApi)
	admin.MyAdminApp.Route("/api/test/240", WorldSrvApi)
	module.RegisteModule(APIHttpSington, time.Second, 1)
}

var APIHttpSington = &APIHttpData{
	RouteMap: make(map[string]string),
}

type APIHttpData struct {
	RouteMap map[string]string
}

func (clm *APIHttpData) ModuleName() string {
	return "APIHttp"
}
func (clm *APIHttpData) Update() {
}

func (clm *APIHttpData) Shutdown() {
	module.UnregisteModule(clm)
}
func WorldSrvApi(rw http.ResponseWriter, req *http.Request) {
	Path := req.URL.Path
	Method := req.Method
	fmt.Println(Path, Method)
	if Method == "POST" {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.Logger.Info("Body err.", err)
			WebApiResponse(rw, map[string]interface{}{RESPONSE_STATE: STATE_ERR, RESPONSE_ERRMSG: "Post data is null!"})
			return
		}
		fmt.Println(data)
	} else {
		m := req.URL.Query()
		id := m.Get("id")
		snid, err := strconv.Atoi(id)
		if err != nil {
			WebApiResponse(rw, map[string]interface{}{RESPONSE_STATE: STATE_ERR, RESPONSE_ERRMSG: err})
			return
		}
		//fmt.Println(snid)
		if url, ok := APIHttpSington.RouteMap[Path]; ok {
			WebApiResponseByte(rw, httpPost(url, snid))
		} else {
			WebApiResponse(rw, map[string]interface{}{RESPONSE_STATE: STATE_ERR, RESPONSE_ERRMSG: "route is nil." + Path})
		}
	}
	return
}
func WebApiResponse(rw http.ResponseWriter, params map[string]interface{}) bool {
	data, err := json.Marshal(params)
	if err != nil {
		logger.Logger.Info("webApiResponse Marshal error:", err)
		return false
	}

	dataLen := len(data)
	rw.Header().Set("Content-Length", fmt.Sprintf("%v", dataLen))
	rw.WriteHeader(http.StatusOK)
	pos := 0
	for pos < dataLen {
		writeLen, err := rw.Write(data[pos:])
		if err != nil {
			logger.Logger.Info("webApiResponse SendData error:", err, " data=", string(data[:]), " pos=", pos, " writelen=", writeLen, " dataLen=", dataLen)
			return false
		}
		pos += writeLen
	}
	return true
}
func WebApiResponseByte(rw http.ResponseWriter, data []byte) bool {
	dataLen := len(data)
	rw.Header().Set("Content-Length", fmt.Sprintf("%v", dataLen))
	rw.WriteHeader(http.StatusOK)
	pos := 0
	for pos < dataLen {
		writeLen, err := rw.Write(data[pos:])
		if err != nil {
			logger.Logger.Info("webApiResponse SendData error:", err, " data=", string(data[:]), " pos=", pos, " writelen=", writeLen, " dataLen=", dataLen)
			return false
		}
		pos += writeLen
	}
	return true
}
func httpPost(url string, id int) []byte {
	fmt.Println("url:>", url)
	val := make(map[string]interface{})
	val["ID"] = id
	val["Gold"] = 50000000
	val["GoldExt"] = 0
	val["Platform"] = "1"
	val["BillNo"] = rand.Intn(1000000)
	val["IsAccTodayRecharge"] = true
	nval := make(map[string]interface{})
	nval["Param"] = val
	c, _ := json.Marshal(nval)
	post := string(c)
	//fmt.Println(url, "post", post)

	var jsonStr = []byte(post)
	//fmt.Println("jsonStr", jsonStr)
	new_str := bytes.NewBuffer(jsonStr)
	fmt.Println("new_str", new_str)

	req, err := http.NewRequest("POST", url, new_str)
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("status", resp.Status)
	//fmt.Println("response:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return body
}
