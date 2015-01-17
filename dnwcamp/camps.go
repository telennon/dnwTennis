package dnwcamp
// Copyright 2014 - Tom Lennon.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//
// This package manages the concept of a camp.
//
// Original created April 2014


import (
	"errors"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Section struct {
	ID 					bson.ObjectId 	`bson:"Id"`			// ID of this section
	Name				string 			"Name"				// Friendly name for this section
	Start 				string 			"Start"				// DateTime camp section starts
	End 				string			"End"				// DateTime camp section ends
	CostDifferential	int64 			"CostDifferential"	// Discount or premium paid for this section
	// ToDo -- Add a section registration limit
}

type Camp struct {
	ID 					bson.ObjectId	`bson:"_id,omitempty"`
	Cid 				int64			"Cid"				// Defines a logical Camp
	Title 				string			"Title"				// Generic name for camp
	Active				bool 			"Active"			// Flips to no when camp is past signup dates
	Cost				int64 			"Cost"				// Base cost of camp - Applies to all sections
	RegStart			string			"RegStart"			// DateTime signup can begin
	RegEnd				string 			"RegEnd"			// DateTime signup can end
	RefundDeadline 		string			"RefundDeadline"	// Last date to cancel and get a refund
	BillingDate			string			"BillingDate"		// Date on which charge will be billed -- 1st time use in 2015
	CampOver			string			"CampOver"			// Date after which the site should no longer allow registrations and flip to a see you next year page
	CamperTypes			[]string 		"CamperTypes"		// Classes of campers - For DNW - Homeowner / Guest
	Sections			[]Section 		"Sections"			// Array of sections for this camp
}
// NewCamp returns a new camp structure.
func NewCamp() (*Camp) {
	return &Camp{bson.NewObjectId(), "", false, 0.00, "", "", "", "", "", nil, nil}
}

// AddSection adds a new camp section to this camp
func (s *Camp) AddSection() error {

	s.Sections = append(s.Sections, *NewCampSection())
	return nil
}
func NewCampSection() (*Section) {
	return &Section{bson.NewObjectId(), "", "", "", 0}
}

// Save - Saves changes made to a camp
func (s *Camp) Save() error {
	sess, col, err := OpenCampCollection()
	defer sess.Close()
	if err == nil {
		err = col.Insert(s)
	}
	return err
}

// ListCurrentCamp - returns only the current camp
func ListCurrentCamp() (*Camp, error) {
	// List out all the camps
	allCamps, err := ListCamps()
	if err != nil {
		return NewCamp(), err
	}
	
	currentCamp := -1
	for i := range allCamps {
		if allCamps[i].Active {
			if currentCamp < 0 {
				currentCamp = i
			} else {
				return NewCamp(), errors.New("More than one active camp!")
			}
		} 
	}
	if currentCamp < 0 {
		return NewCamp(), errors.New("No current camp found")
	}
	return &allCamps[currentCamp], err
}

// ListCamps - returns a slice containing all camp documents
func ListCamps() ([]Camp, error) {
	// Setup Access to the database
	sess, col, err := OpenCampCollection()
	defer sess.Close()

	var theCamps []Camp
	// Find all the Camp Records
	err = col.Find(nil).Sort("RegStart").All(&theCamps)

	return theCamps, err
}

// CreateCampIndex - Is a setup function called to create the index set required on the camps collection
func CreateCampIndex() error {

	index := mgo.Index{
		Key:		[]string{"Title"},
		Unique:		false,
		DropDups:	false,
	}

	index2 := mgo.Index{
		Key:		[]string{"RegStart"},
		Unique:		false,
		DropDups:	false,
	}
	sess, col, err := OpenCampCollection()
	defer sess.Close()
	if err == nil {
		err := col.EnsureIndex(index)
		if err == nil {
			err = col.EnsureIndex(index2)
		}
	}

	return err
}

func GetActiveCamp(cid string) (*Camp, error) {
	// Setup Access to the database
	sess, col, err := OpenCampCollection()
	defer sess.Close()

	var activeCamp Camp
	// Find all the Camp Records
	err = col.Find(bson.M{"Cid": cid, "Active": "true"}).One(&activeCamp)

	return activeCamp, err
}
