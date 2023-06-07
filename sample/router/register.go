/*
此文件为自动生成，生成时间:2020-05-19 15:09:39
*/
package router

import (
	"rice/lib/router"
)

func Register() (*router.Router, error) {
	r := router.New()
	r.BasePath("/api")

	r.Group("/Test2", func() {
		//This is Dowork2
		//It is just a dummy func
		r.Get("/D6/:name", domain1DoWork6FacadeGet)
	})

	r.Group("/Test", func() {
		/* This is a Cooment */
		r.Get("/DoWork", serviceDoWorkFacadeGet)
		//This is Dowork2
		//It is just a dummy func
		//@RightCheck:OpCode
		r.Get("/MyDoWoirk", serviceDoWork2FacadeGet)
		//Reciver 1
		r.Get("/DowWork4", serviceDowWork4FacadeGet)
		//Reciver 2
		r.Get("/DoWork5", serviceDoWork5FacadeGet)
		//Do My Thing
		r.Get("/DoWork6", serviceDoWork6FacadeGet)
		//Do Dummy Thing
		r.Get("/DoWork7", serviceDoWork7FacadeGet)
	})

	return r, nil
}
