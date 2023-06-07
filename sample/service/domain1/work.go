//This is a comment for common purpose
//@MethodGroup:Test2
package domain1

import (
	// "fmt"
	lt2 "rice/model/view"	
)

//This is Dowork2
//It is just a dummy func
//@Method:get,get
//@MethodName:D6/:name
func DoWork6(name,name2,name3 string ,context *lt2.Context) {
	context.Res.Write([]byte(name))
}