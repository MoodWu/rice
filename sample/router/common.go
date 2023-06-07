package router

import (
	"net/http"	
	"rice/model/view"
	"time"
	"encoding/json"
	"errors"
	"rice/common/util"
	"rice/common/consts"
	"strconv"
	"io/ioutil"
	"math/rand"
	"fmt"
	"strings"
)

/*Get请求的参数获取*/
func GetBoolValue(r *http.Request,name string) (bool,error) {
	s := r.FormValue(name)
	v,err := strconv.ParseBool(s)
	if err !=nil {
		util.LogError("转换参数值为Bool出错，参数名：",name,",参数值:",s)
	}
	return v,err
}

func GetByteValue(r *http.Request,name string) (byte,error) {
	return 0,nil
}

func GetInt8Value(r *http.Request,name string) (int8,error) {
	s := r.FormValue(name)
	v,err:=strconv.ParseInt(s,10,8)
	if err !=nil {
		util.LogError("转换参数值为int8出错，参数名：",name,",参数值:",s)
	}
	return int8(v),err
}

func GetInt16Value(r *http.Request,name string) (int16,error) {
	s := r.FormValue(name)
	v,err:=strconv.ParseInt(s,10,16)
	if err !=nil {
		util.LogError("转换参数值为int16出错，参数名：",name,",参数值:",s)
	}
	return int16(v),err
}

func GetInt32Value(r *http.Request,name string) (int32,error) {
	s := r.FormValue(name)
	v,err:=strconv.ParseInt(s,10,32)
	if err !=nil {
		util.LogError("转换参数值为int32出错，参数名：",name,",参数值:",s)
	}
	return int32(v),err
}

func GetInt64Value(r *http.Request,name string) (int64,error) {
	s := r.FormValue(name)
	v,err:=strconv.ParseInt(s,10,64)
	if err !=nil {
		util.LogError("转换参数值为int64出错，参数名：",name,",参数值:",s)
	}
	return v,err
}

func GetIntValue(r *http.Request,name string) (int,error) {
	s := r.FormValue(name)
	v,err:=strconv.ParseInt(s,10,0)
	if err !=nil {
		util.LogError("转换参数值为int出错，参数名：",name,",参数值:",s)
	}
	return int(v),err
}

func GetUintValue(r *http.Request,name string) (uint,error) {
	s := r.FormValue(name)
	v,err:=strconv.ParseUint(s,10,0)
	if err !=nil {
		util.LogError("转换参数值为uint出错，参数名：",name,",参数值:",s)
	}
	return uint(v),err
}

func GetUintptrValue(r *http.Request,name string) (uintptr,error) {
	return 0,nil 
}

func GetFloat32Value(r *http.Request,name string) (float32,error) {
	s := r.FormValue(name)
	v,err:=strconv.ParseFloat(s,32)
	if err !=nil {
		util.LogError("转换参数值为float32出错，参数名：",name,",参数值:",s)
	}
	return float32(v),err
}

func GetFloat64Value(r *http.Request,name string) (float64,error) {
	s := r.FormValue(name)
	v,err:=strconv.ParseFloat(s,64)
	if err !=nil {
		util.LogError("转换参数值为float64出错，参数名：",name,",参数值:",s)
	}
	return v,err
}

func GetComplex64Value(r *http.Request,name string) (complex64,error) {
	return 0,nil
}

func GetComplex128Value(r *http.Request,name string) (complex128,error) {
	return 0,nil 
}

func GetStringValue(r *http.Request,name string) string {
	return r.FormValue(name)
}

func GetTimeValue(r *http.Request,name string) (time.Time,error) {
	s := r.FormValue(name)
	v,err := time.Parse(consts.TimeFormat,s)
	if err !=nil {
		util.LogError("转换参数值为time出错，参数名：",name,",参数值:",s)
	}
	return v,err
}

func GetObjectValue(r *http.Request,name string,outObj interface{}) error {
	item := r.FormValue(name)

	if item != "" {
		err := json.Unmarshal([]byte(item), &outObj)
		if err != nil {
			util.LogError("转换参数错误：", err)
		}
		return err
	} else {
		return errors.New("Unmarshal，值不能为空") 
	}

}


/*Post请求的参数处理*/
func PostBoolValue(r *http.Request,name string) (bool,error) {
	s := r.PostFormValue(name)
	v,err := strconv.ParseBool(s)
	if err !=nil {
		util.LogError("转换参数值为Bool出错，参数名：",name,",参数值:",s)
	}
	return v,err
}

func PostByteValue(r *http.Request,name string) (byte,error) {
	return 0,nil
}

func PostInt8Value(r *http.Request,name string) (int8,error) {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseInt(s,10,8)
	if err !=nil {
		util.LogError("转换参数值为int8出错，参数名：",name,",参数值:",s)
	}
	return int8(v),err
}

func PostInt16Value(r *http.Request,name string) (int16,error) {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseInt(s,10,16)
	if err !=nil {
		util.LogError("转换参数值为int16出错，参数名：",name,",参数值:",s)
	}
	return int16(v),err
}

func PostInt32Value(r *http.Request,name string) (int32,error) {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseInt(s,10,32)
	if err !=nil {
		util.LogError("转换参数值为int32出错，参数名：",name,",参数值:",s)
	}
	return int32(v),err
}

func PostInt64Value(r *http.Request,name string) (int64,error) {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseInt(s,10,64)
	if err !=nil {
		util.LogError("转换参数值为int64出错，参数名：",name,",参数值:",s)
	}
	return v,err
}

func PostIntValue(r *http.Request,name string) (int,error) {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseInt(s,10,0)
	if err !=nil {
		util.LogError("转换参数值为int出错，参数名：",name,",参数值:",s)
	}
	return int(v),err
}

func PostUintValue(r *http.Request,name string) (uint,error) {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseUint(s,10,0)
	if err !=nil {
		util.LogError("转换参数值为uint出错，参数名：",name,",参数值:",s)
	}
	return uint(v),err
}

func PostUintptrValue(r *http.Request,name string) (uintptr,error) {
	return 0,nil 
}

func PostFloat32Value(r *http.Request,name string) (float32,error) {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseFloat(s,32)
	if err !=nil {
		util.LogError("转换参数值为float32出错，参数名：",name,",参数值:",s)
	}
	return float32(v),err
}

func PostFloat64Value(r *http.Request,name string) (float64,error) {
	s := r.PostFormValue(name)
	v,err:=strconv.ParseFloat(s,64)
	if err !=nil {
		util.LogError("转换参数值为float64出错，参数名：",name,",参数值:",s)
	}
	return v,err
}

func PostComplex64Value(r *http.Request,name string) (complex64,error) {
	return 0,nil
}

func PostComplex128Value(r *http.Request,name string) (complex128,error) {
	return 0,nil
}

func PostStringValue(r *http.Request,name string) string {
	return r.PostFormValue(name)
}

func PostTimeValue(r *http.Request,name string) (time.Time,error) {
	s := r.PostFormValue(name)
	v,err := time.Parse(consts.TimeFormat,s)
	if err !=nil {
		util.LogError("转换参数值为time出错，参数名：",name,",参数值:",s)
	}
	return v,err
}

func PostObjectValue(r *http.Request,name string,outObj interface{}) error {
	//目前项目的约定，Post请求的参数都封装在一个大的对象中，参数名就为item,定义在consts.POST_PARAM_NAME
	// item := r.PostFormValue(name)
	header := r.Header["Content-Type"]
	
	bFlag := false
	for _,ct := range header {
		ct = strings.ToLower(ct)
		//判断是否用formdata传参
		if strings.HasPrefix(ct,"multipart/form-data") || ct== "application/x-www-form-urlencoded" {
			bFlag = true
			break
		}
	}

	if bFlag {
		item := r.PostFormValue(consts.POST_PARAM_NAME)

		if item != "" {
			err := json.Unmarshal([]byte(item), &outObj)
			if err != nil {
				util.LogError("转换Form参数错误：", err)
			}
			return err
		} else {
			return errors.New("Unmarshal，值不能为空") 
		}
	} else {
		item, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(item, &outObj)
		if err != nil {
				util.LogError("转换Body参数错误：", err)
		}
		return err	
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
	staffId,_ := strconv.Atoi(r.Header.Get(consts.RIO_STAFFID))
	staffName := r.Header.Get(consts.RIO_STAFFNAME)

	return view.User{
		UserId :staffId,
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


func StartTimer(name string) (string,int64) {
	//生成一个随机数
	now := time.Now().UnixNano()	
	rand.Seed(now)
	rnd := rand.Intn(1000000)	
	return fmt.Sprintf("%v-%s-%v",rnd,name,now),now
}

func EndTimer(name string,startTick int64) int64{
	now := time.Now().UnixNano()
	differ := now - startTick 

	util.LogPerf("name:",name,",nanoseconds:",differ,",micro:",differ/1000000)
	return differ

}