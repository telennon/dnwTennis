package dnwcamp
// Copyright 2014 - Tom Lennon.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//
// This package manages the concept of a camp.
// I may decide to merge camp and sections packages as they are linked, a section
// must be connected to a camp to make sense and a camp with no sections is nonsensicle as well
//
// Original created April 2014


import (
	//"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)



type Camp struct {
	ID 					bson.ObjectId	`bson:"_id,omitempty"`
	Title 				string			"Title"				// Generic name for camp
	Active				bool 			"Active"			// Flips to no when camp is past signup dates
	Cost				int64 			"Cost"				// Base cost of camp - Applies to all sections
	RegStart			string			"RegStart"			// DateTime signup can begin
	RegEnd				string 			"RegEnd"			// DateTime signup can end
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

// NewCamp creates a new camp structure.
func NewCamp() (*Camp) {
	return &Camp{bson.NewObjectId(), "", false, 0.00, "", "", "", nil, nil}
}

func (s *Camp) AddSection() error {

	s.Sections = append(s.Sections, *NewCampSection())
	return nil
}

func (s *Camp) Save() error {
	sess, col, err := OpenCampCollection()
	defer sess.Close()
	if err == nil {
		err = col.Insert(s)
	}
	return err
}



func NewCampSection() (*Section) {
	return &Section{bson.NewObjectId(), "", "", "", 0}
}

func ListCamps() ([]Camp, error) {
	// Setup Access to the database
	sess, col, err := OpenCampCollection()
	defer sess.Close()

	var theCamps []Camp
	// Find all the Camp Records
	err = col.Find(nil).Sort("RegStart").All(&theCamps)

	return theCamps, err
}

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
