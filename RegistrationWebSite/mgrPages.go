package main
/* ************************************************************************
*	Author		TELennon
*	Created		Oct 2014
*
*	Copyright 2014 - Tom Lennon.  All rights reserved.
*	Use of this source code is governed by a MIT-style
*	license that can be found in the LICENSE.md file.
*	
*	mgrPages.go
*		mgrPages is the web site camp managers, coordinators, coaches and
*		billing staff will use to access and manage camp data
*		
************************************************************************* */

import (
	"fmt"
	"net/http"
	"html/template"
	//"time"
	//"github.com/gorilla/mux"

	//"github.com/telennon/dnwTennis/dnwcamp"
)

func managerHandler(w http.ResponseWriter, req *http.Request)  {
	fmt.Println("In the managerHandler")

	tmpl := template.Must(template.ParseFiles("./template/manager.html"))
	
	if (tmpl == nil) {
		return
	}
	
	err := tmpl.Execute(w, "")
	if (err != nil) {
		fmt.Println("Error executing template .. \n\t", err)
		return
	}
}