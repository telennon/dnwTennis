package main
/* ************************************************************************
*	Author		TELennon
*	Created		Oct 2014
*
*	Copyright 2014 - Tom Lennon.  All rights reserved.
*	Use of this source code is governed by a MIT-style
*	license that can be found in the LICENSE.md file.
*	
*	regPages.go
*		regPages is the web site users will access to register for 
*		a camp
*		
************************************************************************* */

import (
	"fmt"
	"net/http"
	"html/template"
	"time"
	//"io/ioutil"
	"github.com/gorilla/mux"

	"github.com/telennon/dnwTennis/dnwcamp"
)

type CampDisplay struct {
	CampRecord	 		*dnwcamp.Camp
	NotTimeYet			bool
	PastTime			bool
	CampOver 			bool
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", homeHandler)
    r.HandleFunc("/manager", managerHandler)
    r.HandleFunc("/Manager", managerHandler)
    r.HandleFunc("/Mgr", managerHandler)
    r.HandleFunc("/mgr", managerHandler)
    r.HandleFunc("/reg/{id:[0-9]+}", existingRegHandler)
    http.Handle("/", r)

    http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("In the homeHandler")
	var cd CampDisplay
	tennisCamp, err := dnwcamp.ListCurrentCamp()
	if err != nil {
		fmt.Println(err)
		return
	}
	cd.CampRecord = tennisCamp


	// Manipulate some of the dates to make them display better
	const dbShortForm = "Jan 2"
	const dbLongForm = "2006/01/02 15:04:05"
	
	t_now := time.Now()
	t_RegStart, _ := time.Parse(dbLongForm, cd.CampRecord.RegStart)
	t_RegEnd, _ := time.Parse(dbLongForm, cd.CampRecord.RegEnd)
	t_RefundDeadline, _ := time.Parse(dbLongForm, cd.CampRecord.RefundDeadline)
	t_CampOver, _ := time.Parse(dbLongForm, cd.CampRecord.CampOver)

	cd.NotTimeYet = t_now.Before(t_RegStart)
	cd.PastTime = t_RegEnd.Before(t_now) && t_now.Before(t_CampOver)
	cd.CampOver = t_CampOver.Before(t_now)

	cd.CampRecord.RegStart = t_RegStart.Format(dbShortForm)
	cd.CampRecord.RegEnd = t_RegEnd.Format(dbShortForm)
	cd.CampRecord.RefundDeadline = t_RefundDeadline.Format(dbShortForm)
	
	for s := range cd.CampRecord.Sections {
		t1, _ := time.Parse(dbLongForm, cd.CampRecord.Sections[s].Start)
		cd.CampRecord.Sections[s].Start = t1.Format(dbShortForm)
		t2, _ := time.Parse(dbLongForm, cd.CampRecord.Sections[s].End)
		cd.CampRecord.Sections[s].End = t2.Format(dbShortForm)
	}

	//tmpl, err := template.New("welcome").ParseFiles("./template/test.html")
	tmpl := template.Must(template.ParseFiles("./template/welcome.html"))
	if (err != nil) {
		fmt.Println("Template failed to parse")
		return
	}
	if (tmpl == nil) {
		return
	}
	err = tmpl.Execute(w, cd)
	if (err != nil) {
		fmt.Println("Error executing template .. \n\t", err)
		return
	}
	//filename := "./template/test.html"
	//body, err := ioutil.ReadFile(filename)
	//if err != nil {
	//	fmt.Println("Error Reading .. \n\t", err)
	//}
	//fmt.Println(body)
	fmt.Println("Leaving the homeHandler without errors")
	return
	
}

func managerHandler(w http.ResponseWriter, req *http.Request)  {
	fmt.Println("In the managerHandler")
}

func existingRegHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("In the existingRegHandler")
}
