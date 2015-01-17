// Copyright 2014 - Tom Lennon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This package is a web server that provides pages for
//	singing up for a camp as well as pages for use by
//	the camp coordinator and pro
//
// Original created Jan 2015

package main 

const (
	CAMPID = 1001							// Identifies the DNWTennis Camp
	TFORMAT = "2006/01/02 15:04:05 MST" 	// Time FOrmat in the DB
)

import (
	"fmt"
	"net/http"
	"html/template"
	"time"

	"github.com/telennon/dnwTennis/dnwcamp"
)

type AllPageData struct {
	Title 				string			// Generic name for camp
	Cost				int64 			// Base cost of camp - Applies to all sections
	RegStart			string			// DateTime signup can begin
	RegEnd				string 			// DateTime signup can end
	RefundDeadline 		string			// Last date to cancel and get a refund
	BillingDate			string			// Date on which charge will be billed -- 1st time use in 2015
	CampOver			string			// Date after which the site should no longer allow registrations and flip to a see you next year page
	Sections			[]Section 		// Array of sections for this camp
	PreSignup			bool			// Not accepting signups yet
	PostSignup			bool			// SIgnup is over buty still accepting signups to be approved by coordinator
	PostCamp			bool			// After camp is over, the page flips to a see you agian next year page			
}

func rootPageHandler(w http.ResponseWriter, r *http.Request) {

	// Get the information for the active Camp
	campData, err := dnwcamp.GetActiveCamp(CAMPID)
	if err != nil {
		fmt.Println("Error attempting to retrieve the active camp record/n", err)
		// Need to display an alternate error HTML Page that the user will see
		// And create some log entry as well so I can find out this happened
	} else {
		pageData := AllPageData{ Title:campData.Title, Cost:campData.Cost, RegStart:campData.RegStart,
			RegEnd:campData.RegEnd, RefundDeadline:campData.RefundDeadline, BillingDate:campData.BillingDate,
			CampOver:campData.CampOver, Sections:campData.Sections, PreSignup:false, PostSignup:false, PostCamp:false }
		currentdate := time.now()
		registrationStart := time.Parse(TFORMAT, campData.RegStart)
		registrationEnd := time.Parse(TFORMAT, campData.RegEnd)
		registrationClosed := time.Parse(TFORMAT, campData.CampOver)

		// Show the see you next year page
		if currentdate.After(registrationClosed) {
			pageData.PostCamp = true
		}

		// Show the waitlist error
		if currentdate.After(registrationEnd) && currentdate.Before(PostCamp) {
			pageData.PostSignup = true
		}

		// Show the too early page
		if currentdate.Before(registrationStart) {
			pageData.PreSignup = true
		}

		//     "RegEnd": "2012/04/28 23:59:59",
		t, err := template.ParseFiles("html/index.html")
		if err != nil {
			fmt.Println(err)
		}
		err = t.Execute(w, pageData)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	http.HandleFunc("/", rootPageHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		//log.Fatal("ListenAndServe: ", err)
		fmt.Println(err)
	}
}