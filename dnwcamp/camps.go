package dnwcamp
// Copyright 2014 - Tom Lennon.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//
// This package manages the concept of a camp.
//
// Original created April 2014


import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Camp struct {
	ID 					bson.ObjectId	`bson:"_id,omitempty"`
	Title 				string			"Title"				// Generic name for camp
	Active				bool 			"Active"			// True: Camp will show up in lists; False: Camp does not show up 
	Cost				int64 			"Cost"				// Base cost of camp - Applies to all sections
	RegStart			string			"RegStart"			// DateTime signup begins
	RegEnd				string 			"RegEnd"			// DateTime signup ends
	RefundDeadline 		string			"RefundDeadline"	// Last date to cancel and get a refund
	CamperTypes			[]string 		"CamperTypes"		// Classes of campers - For DNW - Homeowner / Guest
	Sections			[]Section 		"Sections"			// Array of sections for this camp
}

type Section struct {
	ID 					bson.ObjectId 	`bson:"Id"`			// ID of this section
	Name				string 			"Name"				// Friendly name for this section
	Start 				string 			"Start"				// DateTime camp section starts
	End 				string			"End"				// DateTime camp section ends
	CostDifferential	int64 			"CostDifferential"	// Discount or premium paid for this section
	// ToDo -- Add a section registration limit
}

// NewCamp returns a new camp structure.
func NewCamp() (*Camp) {
	return &Camp{bson.NewObjectId(), "", false, 0.00, "", "", "", nil, nil}
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
