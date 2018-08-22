package controllers

import (
	"encoding/json"
	"fmt"
	"godsp/lib"
	"net"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	xj "github.com/basgys/goxml2json"
)

// Operations about dsp
type DspController struct {
	beego.Controller
}

var connectionPool = lib.NewConnectionPool()

func init() {
	//open socket DSP
	//go test.StartServer()
	c := make(chan int)
	_total := beego.AppConfig.String("dspPool")
	total, err := strconv.Atoi(_total)
	if err != nil {
		total = 1
	}

	for i := 0; i < total; i++ {
		go startClient(connectionPool, i, c)
		clientID := <-c
		beego.Info("connectionID", clientID)
	}
	//end open socket
}

// @Title Dsp Select
// @Description query from config
// @Param	oprid		query 	string	true		"Operation from config"
// @Success 200 {string} success
// @Failure 403 data not found
// @router /select [get]
func (q *DspController) Select() {

	_oprid := q.GetString("oprid")
	beego.Info(_oprid)
	//oprid := beego.AppConfig.String("oprid." + _oprid)
	query := "test"

	q.Data["xml"] = query
	q.ServeXML()

	/*mymap := q.Ctx.Request.URL.Query()
	keys := reflect.ValueOf(mymap).MapKeys()
	strkeys := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
		query = strings.Replace(query, "["+strkeys[i]+"]", mymap[strkeys[i]][0], -1)

	}*/
	query = strings.Replace(query, "[msisdn]", "628118003585", 1)
	payload := getDataRequest(query)
	cr := make(chan string)
	go submitRequest(connectionPool, payload, cr)
	msg := <-cr
	response, err := parseResultData(msg)
	if err != nil {
		q.Data["json"] = msg
		beego.Info(response)
		q.ServeJSON()
	}
	var resp map[string]interface{}
	if resp == nil {
		resp = make(map[string]interface{})
	}
	q.Data["json"] = response
	beego.Info(response)
	q.ServeJSON()

}

func getResponse(connectionPool *lib.ConnectionPool, connectionID int, c net.Conn, resp chan string) {
	for {
		message := make([]byte, 4096)
		length, err := c.Read(message)
		if err != nil {
			c.Close()
			cr := make(chan int)
			go startClient(connectionPool, connectionID, cr)
			connectionID := <-cr
			beego.Info("parse-restart:", connectionID)
			break
		}
		if length > 0 {
			data := message[:length]
			resp <- string(data)
			//data contains invalid should retry

			break
		}
	}

}

func submitLogin(connectionPool *lib.ConnectionPool, connectionID int, c net.Conn, data string, response chan string) {
	resp := make(chan string)
	go getResponse(connectionPool, connectionID, c, resp)
	fmt.Fprintf(c, data)
	msg := <-resp
	beego.Info("LOGIN", msg)
	response <- msg

}
func submitRequest(connectionPool *lib.ConnectionPool, data string, response chan string) {
	c, connectionID, session := connectionPool.GetWithId()
	beego.Info("SESSION:", session)
	beego.Info("CONNECTION_ID:", connectionID)
	data = strings.Replace(data, "[session]", session, 1)
	beego.Info(data)
	resp := make(chan string)
	go getResponse(connectionPool, connectionID, c, resp)
	fmt.Fprintf(c, data)
	msg := <-resp
	response <- msg
}

func startClient(connectionPool *lib.ConnectionPool, connectionID int, c chan int) {
	socket, err := net.Dial("tcp", beego.AppConfig.String("dspIP")+":"+beego.AppConfig.String("dspPort"))
	if err != nil {
		beego.Info("error")
		panic(err)
	}
	session := "1234567"
	clogin := make(chan string)
	beego.Info(getDataLogin())

	go submitLogin(connectionPool, connectionID, socket, getDataLogin(), clogin)
	loginResp := <-clogin
	msg := strings.Split(string(loginResp), "\r\n")
	session = strings.Split(string(msg[1]), "/")[3]
	beego.Info("SESSION: " + session)
	connectionID = connectionPool.AddWithId(socket, connectionID, session)
	beego.Info("START", connectionID)
	c <- connectionID

}

func getDataLogin() string {
	data := "<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\"><soapenv:Body><LGI><HLRSN>1</HLRSN><OPNAME>[user]</OPNAME><PWD>[pwd]</PWD></LGI></soapenv:Body></soapenv:Envelope>"
	data = strings.Replace(data, "[user]", beego.AppConfig.String("dspUser"), 1)
	data = strings.Replace(data, "[pwd]", beego.AppConfig.String("dspPwd"), 1)
	header := "POST " + "/" + " HTTP/1.1\r\n"
	header += "HOST: " + beego.AppConfig.String("localIP") + ":" + beego.AppConfig.String("localPort") + "\r\n"
	header += "Content-Length: " + strconv.Itoa(len(data)) + "\r\n"
	header += "Content-Type: text/xml;charset=UTF-8\r\n"
	header += "\r\n"
	return header + data
}
func getDataRequest(command string) string {
	data := "<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\"><soapenv:Body>" + command + "</soapenv:Body></soapenv:Envelope>"
	header := "POST " + "/" + "[session]" + " HTTP/1.1\r\n"
	header += "HOST: " + beego.AppConfig.String("localIP") + ":" + beego.AppConfig.String("localPort") + "\r\n"
	header += "Content-Length: " + strconv.Itoa(len(data)) + "\r\n"
	header += "Content-Type: text/xml;charset=UTF-8\r\n"
	header += "\r\n"
	return header + data
}

func parseResult(msg string) map[string]interface{} {
	_msg := strings.Split(msg, "\r\n")
	body := _msg[5]
	beego.Info("BODY:", body)
	r1 := strings.Split(body, "<Result>")
	r2 := strings.Split(r1[1], "</Result>")
	_result := r2[0]
	beego.Info("RESULT:", _result)

	xml := strings.NewReader(_result)
	_json, err := xj.Convert(xml)
	if err != nil {
		beego.Error("ERROR_PARSE_RESULT")
	}
	body = _json.String()
	beego.Info("RESULT:", body)

	var data map[string]interface{}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		beego.Error("ERROR_PARSE_RESULT")
	}
	//beego.Info(data["ResultData"])

	return data
}
func parseResultData(msg string) (map[string]interface{}, error) {
	_msg := strings.Split(msg, "\r\n")
	body := _msg[5]
	beego.Info("BODY:", body)
	r1 := strings.Split(body, "<ResultData>")
	r2 := strings.Split(r1[1], "</ResultData>")
	_result := r2[0]
	beego.Info("RESULT:", _result)

	xml := strings.NewReader(_result)
	_json, err := xj.Convert(xml)
	if err != nil {
		beego.Error("ERROR_PARSE_RESULT")
		//errors.New("ERROR_PARSE_RESULT")
	}
	body = _json.String()
	beego.Info("RESULT:", body)

	var data map[string]interface{}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		beego.Error("ERROR_PARSE_RESULT")
		//errors.New("ERROR_PARSE_RESULT")
	}
	//beego.Info(data["ResultData"])

	return data, err
}
