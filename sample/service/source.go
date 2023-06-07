//This is a comment for common purpose
//@MethodGroup:Test
package service

import (
	"fmt"	
	"time"	
	lt "rice/model/view"
)

type MyParam struct {	
	ID   string
	Name string
}



//type int string

/* This is a Cooment */
//@Method:get
func DoWork(p1 int, p2 string, p3 time.Time)(code int,k *int,ret string,p *MyParam) {
	fmt.Println("This is a test")
}

//This is Dowork2
//It is just a dummy func
//@Method:get,get
//@MethodName:MyDoWoirk
//@RightCheck:OpCode
func DoWork2(p1 int, p2 string, p3 lt.MyType) string {
	//Inner one
	type Test struct {
		ID int
	}
	fmt.Println("DoWork2")
	return ""
}


func DoWork3(p MyParam, p2 []int, p3 map[string]string, p4 [4]int, extra ...int) (MyRet string) {
	fmt.Println("DoWork3")
	return ""
}

//Reciver 1
//@Method:get
func DowWork4(p1 int, p2 string) (MyRet string) {
	return ""
}

//Reciver 2
//@Method:get
func DoWork5(p1 int) string {
	return ""
}

//Do My Thing
//@Method:get
func DoWork6(context *lt.Context) string {
	context.Res.Write([]byte("Got you"))
	return ""
}

//Do Dummy Thing
//@Method:get
func DoWork7() {
	fmt.Println("Done")	
}

