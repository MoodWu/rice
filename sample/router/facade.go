/*
此文件为自动生成，生成时间:2020-05-19 15:09:39
*/
package router

import (
	"net/http"
	lt2 "rice/model/view"
	service "rice/service"
	domain1 "rice/service/domain1"
	time "time"
)

//This is Dowork2
//It is just a dummy func
func domain1DoWork6FacadeGet(w http.ResponseWriter, r *http.Request) {
	_ticker_id, _ticker := StartTimer("domain1DoWork6FacadeGet")
	name := GetStringValue(r, "name")
	name2 := GetStringValue(r, "name2")
	name3 := GetStringValue(r, "name3")
	var context *lt2.Context
	context = CreateContext(w, r)
	domain1.DoWork6(name, name2, name3, context)
	EndTimer(_ticker_id, _ticker)
}

/* This is a Cooment */
func serviceDoWorkFacadeGet(w http.ResponseWriter, r *http.Request) {
	_ticker_id, _ticker := StartTimer("serviceDoWorkFacadeGet")
	p1, err := GetIntValue(r, "p1")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	p2 := GetStringValue(r, "p2")
	var p3 time.Time
	p3, err = GetTimeValue(r, "p3")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	service.DoWork(p1, p2, p3)
	EndTimer(_ticker_id, _ticker)
}

//This is Dowork2
//It is just a dummy func
//@RightCheck:OpCode
func serviceDoWork2FacadeGet(w http.ResponseWriter, r *http.Request) {
	_ticker_id, _ticker := StartTimer("serviceDoWork2FacadeGet")
	p1, err := GetIntValue(r, "p1")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	p2 := GetStringValue(r, "p2")
	tmp_p3 := lt2.MyType{}
	if GetObjectValue(r, "p3", &tmp_p3) != nil {
		w.WriteHeader(500)
		w.Write([]byte(GetObjectValue(r, "p3", &tmp_p3).Error()))
		return
	}
	p3 := tmp_p3
	ret := service.DoWork2(p1, p2, p3)
	RenderResult(w, ret)
	EndTimer(_ticker_id, _ticker)
}

//Reciver 1
func serviceDowWork4FacadeGet(w http.ResponseWriter, r *http.Request) {
	_ticker_id, _ticker := StartTimer("serviceDowWork4FacadeGet")
	p1, err := GetIntValue(r, "p1")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	p2 := GetStringValue(r, "p2")
	ret := service.DowWork4(p1, p2)
	RenderResult(w, ret)
	EndTimer(_ticker_id, _ticker)
}

//Reciver 2
func serviceDoWork5FacadeGet(w http.ResponseWriter, r *http.Request) {
	_ticker_id, _ticker := StartTimer("serviceDoWork5FacadeGet")
	p1, err := GetIntValue(r, "p1")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	ret := service.DoWork5(p1)
	RenderResult(w, ret)
	EndTimer(_ticker_id, _ticker)
}

//Do My Thing
func serviceDoWork6FacadeGet(w http.ResponseWriter, r *http.Request) {
	_ticker_id, _ticker := StartTimer("serviceDoWork6FacadeGet")
	var context *lt2.Context
	context = CreateContext(w, r)
	ret := service.DoWork6(context)
	RenderResult(w, ret)
	EndTimer(_ticker_id, _ticker)
}

//Do Dummy Thing
func serviceDoWork7FacadeGet(w http.ResponseWriter, r *http.Request) {
	_ticker_id, _ticker := StartTimer("serviceDoWork7FacadeGet")
	service.DoWork7()
	EndTimer(_ticker_id, _ticker)
}
