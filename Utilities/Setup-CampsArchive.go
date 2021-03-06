package main
/* ************************************************************************
*	Author		TELennon
*	Created		Jan 2014
*
*	Copyright 2014 - Tom Lennon.  All rights reserved.
*	Use of this source code is governed by a MIT-style
*	license that can be found in the LICENSE.md file.

*	convertdb.go
*		convertdb will convert past tennis databases to the new mongoDB form.
* 		Past Data exists for the years 2012 - 2014 and all will be converted into a
*		single database
*
*		The new format consists of two collections campConfig and campRegistrations
*		CampConfig is new and those entries need to be created manually
*		The other data will be read from CSV files. Data conversion to attach to 
*		the concept of camps and sections will be done and then the records written
*		to the new Mongodb database.
*
************************************************************************* */
import (
	"fmt"
	"labix.org/v2/mgo/bson"

	"github.com/telennon/dnwTennis/dnwcamp"
)


func createCampsArchive() (error) {
	// Past Camp Configuration Data
	//		This data is created manually for all past camps for which there is data.
	// Todo:
	//		Add section ids to each camp section that will need to be assigned programatically
	//		For symmetry consider carrying the campId in the sections - Probably not
	pastcamps := []dnwcamp.Camp {
		{ID: bson.NewObjectId(),
		Title: "2012 DNW Tennis Camp", 
		Active: false,
		Cost: 145,
		RegStart: "2012/03/15 00:00:00",
		RegEnd: "2012/04/28 23:59:59",
		RefundDeadline: "2012/06/15 00:00:00",
		CamperTypes: []string { "Homeowner", "Guest"},
		Sections: []dnwcamp.Section {
			{
			Name: "Camp 1",
			Start: "2012/07/16 00:00:00",
			End: "2012/07/20 00:00:00",
			CostDifferential: 0},
			{
			Name: "Camp 2",
			Start: "2012/08/01 00:00:00",
			End: "2012/08/05 00:00:00",
			CostDifferential: 0},
			{
			Name: "Camp 3",
			Start: "2012/08/15 00:00:00",
			End: "2014/08/19 00:00:00",
			CostDifferential: 0},
		},
		},
		{ID: bson.NewObjectId(),
		Title: "2013 DNW Tennis Camp", 
		Active: false,
		Cost: 145,
		RegStart: "2013/03/15 00:00:00",
		RegEnd: "2013/04/28 23:59:59",
		RefundDeadline: "2013/06/15 00:00:00",
		CamperTypes: []string { "Homeowner", "Guest"},
		Sections: []dnwcamp.Section{
			{
			Name: "Camp 1",
			Start: "2013/07/29 00:00:00",
			End: "2014/08/02 00:00:00",
			CostDifferential: 0},
			{
			Name: "Camp 2",
			Start: "2014/08/05 00:00:00",
			End: "2014/08/09 00:00:00",
			CostDifferential: 0},
			{
			Name: "Camp 3",
			Start: "2014/08/12 00:00:00",
			End: "2014/08/16 00:00:00",
			CostDifferential: 0},
		},
		},
		{ID: bson.NewObjectId(),
		Title: "2014 DNW Tennis Camp", 
		Active: true,
		Cost: 150,
		RegStart: "2014/03/15 00:00:00",
		RegEnd: "2014/04/28 23:59:59",
		RefundDeadline: "2014/06/15 00:00:00",
		CamperTypes: []string { "Homeowner", "Guest"},
		Sections: []dnwcamp.Section{
			{
			Name: "Camp 1",
			Start: "2014/07/28 00:00:00",
			End: "2014/08/01 00:00:00",
			CostDifferential: 0},
			{
			Name: "Camp 2",
			Start: "2014/08/04 00:00:00",
			End: "2014/08/08 00:00:00",
			CostDifferential: 0},
			{
			Name: "Camp 3",
			Start: "2014/08/11 00:00:00",
			End: "2014/08/15 00:00:00",
			CostDifferential: 0},
		},
		},
	}

	// Create camp documents
	for i := 0; i < len(pastcamps); i++ {
		aCamp := dnwcamp.NewCamp()
		aCamp.Title = pastcamps[i].Title
		aCamp.Active = pastcamps[i].Active
		aCamp.Cost = pastcamps[i].Cost
		aCamp.RegStart = pastcamps[i].RegStart
		aCamp.RegEnd = pastcamps[i].RegEnd
		aCamp.RefundDeadline = pastcamps[i].RefundDeadline
		aCamp.CamperTypes = pastcamps[i].CamperTypes
		for y := 0; y < len(pastcamps[i].Sections); y++ {
			aCamp.AddSection()
			aCamp.Sections[y].Name = pastcamps[i].Sections[y].Name
			aCamp.Sections[y].Start = pastcamps[i].Sections[y].Start
			aCamp.Sections[y].End = pastcamps[i].Sections[y].End
			aCamp.Sections[y].CostDifferential = pastcamps[i].Sections[y].CostDifferential
		}
		err := aCamp.Save()
		if err != nil {
			fmt.Println("Error - Creating camp ", pastcamps[i].Title, "\n", err)
			return err
		}
	}

	err := dnwcamp.CreateCampIndex()
	if err != nil {
		fmt.Println("Error - Creating camp indexes\n", err)
	}

	return err
}