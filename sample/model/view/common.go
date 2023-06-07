package view

import (
	"net/http"	
)

type User struct {
	UserID        string
	UserName      string
}

type Context struct{
	Res http.ResponseWriter
	Req *http.Request
	User
}

type MyType struct {
	ID string
	Name string
}


