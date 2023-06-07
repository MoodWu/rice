package router

import (
	"net/http"	
	"rice/model/view"
	"time"
	"encoding/json"
	"fmt"
	"rice/common/consts"
	"strconv"
)

/*Get请求的参数获取*/
func GetBoolValue(r *http.Request,name string) bool {
	s := r.FormValue(name)
	v,err := strconv.ParseBool(s)
	if err !=nil {
		fmt.Println("转换参数值为Bool出错，参数名：",name,",参数值:",s)
	}
	return v
}

func GetByteValue(r *http.Request,name string) byte {
	return 0
}

func GetInt8Value(r *http.Request,name string) int8 {
	s := r.FormValue(name)
	v,err:=strconv.ParseInt(s,10,8)
	if err !=nil {
		fmt.Println("转换参数值为int8出错，参数名：",name,",参数值:",s)
	}
	return int8(v)
}

func GetInt16Value(r *http.Request,name string) int16 {
	s := r.FormValue(name)
	v,err:=strconv.ParseInt(s,10,16)
	if err !=nil {
		fmt.Println("转换参数值为int16出错，参数名：",name,",参数值:",s)
	}
	return int16(v)
}

func GetInt32Value(r *http.Request,name string) int32 {
	s := r.FormValue(name)
	v,err:=strconv.ParseInt(s,10,32)
	if err !=nil {
		fmt.Println("转换参数值为int32出错，参数名：",name,",参数值:",s)
	}
	return int32(v)
}

func GetInt64Value(r *http.Request,name string) int64 {
	s := r.FormValue(name)
	v,err:=strconv.ParseInt(s,10,64)
	if err !=nil {
		fmt.Println("转换参数值为int64出错，参数名：",name,",参数值:",s)
	}
	return v
}

func GetIntValue(r *http.Request,name string) int {
	s := r.FormValue(name)
	v,err:=strconv.ParseInt(s,10,0)
	if err !=nil {
		fmt.Println("转换参数值为int出错，参数名：",name,",参数值:",s)
	}
	return int(v)
}

func GetUintValue(r *http.Request,name string) uint {
	s := r.FormValue(name)
	v,err:=strconv.ParseUint(s,10,0)
	if err !=nil {
		fmt.Println("转换参数值为uint出错，参数名：",name,",参数值:",s)
	}
	return uint(v)
}

func GetUintptrValue(r *http.Request,name string) uintptr {
	return 0
}

func GetFloat32Value(r *http.Request,name string) float32 {
	s := r.FormValue(name)
	v,err:=strconv.ParseFloat(s,32)
	if err !=nil {
		fmt.Println("转换参数值为float32出错，参数名：",name,",参数值:",s)
	}
	return float32(v)
}

func GetFloat64Value(r *http.Request,name string) float64 {
	s := r.FormValue(name)
	v,err:=strconv.ParseFloat(s,64)
	if err !=nil {
		fmt.Println("转换参数值为float64出错，参数名：",name,",参数值:",s)
	}
	return v
}

func GetComplex64Value(r *http.Request,name string) complex64 {
	return 0
}

func GetComplex128Value(r *http.Request,name string) complex128 {
	return 0
}

func GetStringValue(r *http.Request,name string) string {
	return r.FormValue(name)
}

func GetTimeValue(r *http.Request,name string) time.Time {
	s := r.FormValue(name)
	v,err := time.Parse(consts.TimeFormat,s)
	if err !=nil {
		fmt.Println("转换参数值为time出错，参数名：",name,",参数值:",s)
	}
	return v
}

func GetObjectValue(r *http.Request,name string,outObj interface{})  {
	item := r.FormValue(name)

	if item != "" {
		err := json.Unmarshal([]byte(item), &outObj)
		if err != nil {
			fmt.Println("转换参数错误：", err)
		}
	}

}


/*Post请求的参数处理*/
func PostBoolValue(r *http.Request,name string) bool {
	s := r.PostFormValue(name)
	v,err := strconv.ParseBool(s)
	if err !=nil {
		fmt.Println("转换参数值为Bool出错，参数名：",name,",参数值:",s)
	}
	return v
}

func PostByteValue(r *http.Request,name string) byte {
	return 0
}

func PostInt8Value(r *http.Request,name string) int8 {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseInt(s,10,8)
	if err !=nil {
		fmt.Println("转换参数值为int8出错，参数名：",name,",参数值:",s)
	}
	return int8(v)
}

func PostInt16Value(r *http.Request,name string) int16 {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseInt(s,10,16)
	if err !=nil {
		fmt.Println("转换参数值为int16出错，参数名：",name,",参数值:",s)
	}
	return int16(v)
}

func PostInt32Value(r *http.Request,name string) int32 {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseInt(s,10,32)
	if err !=nil {
		fmt.Println("转换参数值为int32出错，参数名：",name,",参数值:",s)
	}
	return int32(v)
}

func PostInt64Value(r *http.Request,name string) int64 {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseInt(s,10,64)
	if err !=nil {
		fmt.Println("转换参数值为int64出错，参数名：",name,",参数值:",s)
	}
	return v
}

func PostIntValue(r *http.Request,name string) int {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseInt(s,10,0)
	if err !=nil {
		fmt.Println("转换参数值为int出错，参数名：",name,",参数值:",s)
	}
	return int(v)
}

func PostUintValue(r *http.Request,name string) uint {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseUint(s,10,0)
	if err !=nil {
		fmt.Println("转换参数值为uint出错，参数名：",name,",参数值:",s)
	}
	return uint(v)
}

func PostUintptrValue(r *http.Request,name string) uintptr {
	return 0
}

func PostFloat32Value(r *http.Request,name string) float32 {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseFloat(s,32)
	if err !=nil {
		fmt.Println("转换参数值为float32出错，参数名：",name,",参数值:",s)
	}
	return float32(v)
}

func PostFloat64Value(r *http.Request,name string) float64 {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseFloat(s,64)
	if err !=nil {
		fmt.Println("转换参数值为float64出错，参数名：",name,",参数值:",s)
	}
	return v
}

func PostComplex64Value(r *http.Request,name string) complex64 {
	return 0
}

func PostComplex128Value(r *http.Request,name string) complex128 {
	return 0
}

func PostStringValue(r *http.Request,name string) string {
	return r.PostFormValue(name)
}

func PostTimeValue(r *http.Request,name string) time.Time {
	s := r.PostFormValue(name)
	v,err := time.Parse(consts.TimeFormat,s)
	if err !=nil {
		fmt.Println("转换参数值为time出错，参数名：",name,",参数值:",s)
	}
	return v
}

func PostObjectValue(r *http.Request,name string,outObj interface{})  {
	//目前项目的约定，Post请求的参数都封装在一个大的对象中，参数名就为item,定义在consts.POST_PARAM_NAME
	// item := r.PostFormValue(name)
	item := r.PostFormValue(consts.POST_PARAM_NAME)

	if item != "" {
		err := json.Unmarshal([]byte(item), &outObj)
		if err != nil {
			fmt.Println("转换参数错误：", err)
		}
	}

}


/*创建调用上下文*/
func CreateContext(w http.ResponseWriter,r *http.Request) *view.Context {
	user := GetUserInfo(r)
	return &view.Context {
		Res:w,
		Req:r,
		User:user,
	}
}

func GetUserInfo(r *http.Request) view.User {
	staffId := r.Header.Get(consts.RIOStaffID)
	staffName := r.Header.Get(consts.RIOStaffID)

	return view.User{
		UserID :staffId,
		UserName:staffName,
	}
}


func RenderResult(w http.ResponseWriter,data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}