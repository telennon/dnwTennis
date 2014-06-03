package dnwcamp
// Copyright 2014 - Tom Lennon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// registrations.go
// Data structures and CRUD functions for a camp registration

import (
	"fmt"
	"strconv"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)


type Request struct {
	ReqSectionId			bson.ObjectId 	`bson:"reqSectionId"`	// bson ObjectId of the camp section being requested
	ReqPriority				int64			"reqPriority"			// 1st through nth choice - using a priority > 1 time is not legal			// Human Readable camp name
}


type Camper struct {
	Name 					string 			"camperName"
	Age 					string 			"camperAge"
	ShirtSize				string 			"camperShirtSize"		// Other items can be added
																	// Could do this as a uniform where we having clothing and collect sizes
																	// What to show would be part of camp setup
	Type 					string 			"camperType"			// Defines which class of camper this person is - Options defined in Camp
	TimeStamp 				string			"timeStamp"				// Date and time this camper was entered
	Comment 				string 			"comment"				// Comment field for each camper
	Sections 				[]bson.ObjectId "camperSections"		// IDs of sections this camper is assigned to
	WaitList 				[]bson.ObjectId "camperWaitListed"		// IDs of sections this camper is wait listed for
}

type Registration struct {
	RegId					bson.ObjectId 	`bson:"_id,omitempty"`	// ID of this registration
	RegCampId				bson.ObjectId 	"campId"				// ID of the camp this registration is connected to
	RegName					string 			"regName"				// Name for this registration
	RegPhone 				string 			"regPhone"				// Contact phone for this registration
	RegSite					string			"regSite"				// Home Site or other identifier for billing
	RegEmail				string 			"regEmail"				// Contact email for this registration
	RegTimeStamp			string 			"regTimeStamp"			// Date and time this entry was submitted
	RegComments				string 			"RegComments"			// Comment field		
	RegRequests				[]Request 		"regRequests"			// List of camp section ids and priority
	RegCampers				[]Camper 		"regCampers" 			// List of campers registered
}

// NewRegstration creates a new registration structure.
func NewRegistration() (Registration) {

	return Registration{bson.NewObjectId(), bson.NewObjectId(), "", "", "", "", "", "", nil, nil}
}

func NewCamper() (Camper) {
	return Camper{"", "", "", "", "", "", []bson.ObjectId{}, []bson.ObjectId{}}
}

func (r Registration) AddCamper(c Camper) Registration{
	r.RegCampers = append(r.RegCampers, c)
	return r 
}

func NewRequest(s bson.ObjectId, p int64) (Request) {
	// s = The section being request
	// p = The priority of this request
	return Request{s, p}
}

func ListRegistrations() ([]Registration, error) {
	// Setup Access to the database
	session, err := mgo.Dial(MONGODBSERVER)
	if err != nil {
		fmt.Println("Dialing the database presented the following error\n", err)
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)
	dbCollection := session.DB(DATABASE).C(COL_REG)

	defer session.Close()

	var theRegistrations []Registration
	// Find all the Camp Records
	err = dbCollection.Find(bson.M{}).All(&theRegistrations)

	return theRegistrations, err
}

func (r Registration)Save() (error) {
		// Setup Access to the database
	session, col, err := OpenRegistrationCollection()
	if err == nil {
		defer session.Close()
		err = col.Insert(r)
	}

	return err
}

func ListReservationsForCamp(campID bson.ObjectId) ([]Registration, error) {
	q := bson.M{"campId": campID}

		session, err := mgo.Dial(MONGODBSERVER)
	if err != nil {
		fmt.Println("Dialing the database presented the following error\n", err)
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)
	dbCollection := session.DB(DATABASE).C(COL_REG)

	defer session.Close()

	var theRegistrations []Registration

	err = dbCollection.Find(q).All(&theRegistrations)

	return theRegistrations, err
}

type ByTimeStamp []Registration

func (a ByTimeStamp) Len() int {
	return len(a)
}

func (a ByTimeStamp) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByTimeStamp) Less(i, j int) bool {
	return a[i].RegTimeStamp < a[i].RegTimeStamp
}

type ByTypeAndAge []Camper

func (a ByTypeAndAge) Len() int {
	return len(a)
}

func (a ByTypeAndAge) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByTypeAndAge) Less(i, j int) bool {
	if (a[i].Type == "Homeowner") && (a[j].Type == "Guest") {
		return true
	} else {
		if (a[i].Type == "Guest") && (a[j].Type == "Homeowner") {
			return false
		} else {
			ai, _ := strconv.ParseInt(a[i].Age, 10, 64)
			aj, _ := strconv.ParseInt(a[j].Age, 10, 64)
			return ai > aj
		}
	}
}
