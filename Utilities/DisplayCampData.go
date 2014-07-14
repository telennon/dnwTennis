package main
/* ************************************************************************
*	Author		TELennon
*	Created		May 2014
*
*	Copyright 2014 - Tom Lennon.  All rights reserved.
*	Use of this source code is governed by a MIT-style
*	license that can be found in the LICENSE.md file.
*	
*	DisplayCampData.go
*		SetupRegistrationArchive will install all the registration and camper
*		records from previous years camps. 2012 - 2014 
*
*		The new format consists of two collections Camps and Registrations
*		The old archive data comes in two forms - One is an old spreadsheet
*		that has been edited to separate registrations from campers. The other
*		is a CSV export from the existing SQL database
*		
************************************************************************* */
import (
	"fmt"
	"sort"

	"code.google.com/p/gofpdf"
	"github.com/telennon/dnwTennis/dnwcamp"
)

var CurrentHeaderText string

func main() {
	// Basic Algorithm
	// Initialize PDF
	pdf := gofpdf.New("P", "in", "Letter", "")
	pdf.SetHeaderFunc(func() {
		pdf.SetFont("Arial", "B", 24)
		pdf.SetFillColor(204, 204, 204)
		//pdf.CellFormat(0, 1, allCamps[j].Title, "", 2, "LM", true, 0, "")
		pdf.CellFormat(0, 1, CurrentHeaderText, "", 2, "LM", true, 0, "")
	})
	pdf.SetFooterFunc(func() {
		pdf.SetY(-.5)
		pdf.SetFont("Arial", "", 8)
		pdf.CellFormat(0, .5, fmt.Sprintf("Page %d of {nb}", pdf.PageNo()), "", 0, "C", false, 0, "")
	})
	pdf.SetFontLocation("/Users/tom/gopkg/src/code.google.com/p/gofpdf/font")
	pdf.SetMargins(.25, .5, .25)
	pdf.SetAutoPageBreak(true, .75)
	pdf.AliasNbPages("")
	// Generate a title
	//pdf.SetFont("Arial", "B", 18)
	//pdf.AddPage()
	//pdf.SetMargins(.25, .5, .25)
	//pdf.SetAutoPageBreak(true, 0)
	//pdf.SetXY(.75,.75)

	// Fetch all the camps
	allCamps, err := dnwcamp.ListCamps()
	if err != nil {
		fmt.Println("Error getting camps")
	}
	// For each camp
	for j := range allCamps {															// Camps
		//		fetch all reservations
		allCampReservations, err := dnwcamp.ListReservationsForCamp(allCamps[j].ID)
		sort.Sort(dnwcamp.ByTimeStamp(allCampReservations))

		if err != nil {
			fmt.Println("Error getting reservations")
		}
		CurrentHeaderText = allCamps[j].Title

		//		generate a PDF report
		for i := range allCamps[j].Sections {											// Sections
			// Create camp section header
			//	loop through reservations looking for campers who are in this section
			pdf.AddPage()
			pdf.CellFormat(0, .2, "", "", 2, "", false, 0, "")
			pdf.SetFont("Arial", "B", 14)
			pdf.SetFillColor(3, 172, 229)

			pdf.CellFormat(0, .25, allCamps[j].Sections[i].Name, "1", 2, "LM", true, 0, "")
			for y := range allCampReservations {															// Registrations
				printRegistrationHeader := true
				sort.Sort(dnwcamp.ByTypeAndAge(allCampReservations[y].RegCampers))
				for x := range allCampReservations[y].RegCampers {										// Campers
					for z := range allCampReservations[y].RegCampers[x].Sections {						// Camper Sections
						if allCamps[j].Sections[i].ID == allCampReservations[y].RegCampers[x].Sections[z] {

							// LEFT OFF - Need to write the routine below but realizing that the check done above
							//  Will need to be done agian in the routine so probably just want to pass in
							//	general stuff and have the routine figure out who gets printed
							if printRegistrationHeader {	
								if pdf.GetY() > 9 {
									pdf.AddPage()
								}
								pdf.SetFont("Arial", "", 12)
								pdf.SetFillColor(197, 221, 229)
								pdf.CellFormat(2.25, .25, allCampReservations[y].RegName + " (" + allCampReservations[y].RegSite + ")", "1", 0, "LM", true, 0, "")
								pdf.CellFormat(2.75, .25, "Email: " + allCampReservations[y].RegEmail, "1", 0, "LM", true, 0, "")
								pdf.CellFormat(3, .25, "Phone: " + allCampReservations[y].RegPhone, "1", 1, "LM", true, 0, "")
								printRegistrationHeader = false
							}

							pdf.SetFont("Arial", "", 12)
							pdf.CellFormat(.12, .25, "", "1", 0, "LM", false, 0, "")
							pdf.CellFormat(2.13, .25, allCampReservations[y].RegCampers[x].Name, "1", 0, "LM", false, 0, "")
							pdf.CellFormat(.75, .25, allCampReservations[y].RegCampers[x].Age, "1", 0, "LM", false, 0, "")
							pdf.CellFormat(2, .25, allCampReservations[y].RegCampers[x].ShirtSize, "1", 0, "LM", false, 0, "")
							pdf.CellFormat(3, .25, allCampReservations[y].RegCampers[x].Type, "1", 1, "LM", false, 0, "")
						}
					}
				}
			}
		}
	}
	err = pdf.OutputFileAndClose("./campData.pdf")
}

//	pdf := gofpdf.New("P", "in", "Letter", "")
//	pdf.SetFontLocation("/Users/tom/gopkg/src/code.google.com/p/gofpdf/font")
//	pdf.SetFont("Arial", "B", 18)
//	pdf.AddPage()
//	pdf.SetMargins(1, 2, 1)
//	pdf.SetAutoPageBreak(false, 0)
///	pdf.SetXY(.75,.75)
//	pdf.CellFormat(0, .5, "DNW Tennis Camp Summary", "", 2, "TC", false, 0, "")
//	pdf.SetFont("Arial", "", 12)
//	pdf.CellFormat(0, .5, "Reservation 1", "", 2, "TL", false, 0, "")
//
//	fn := "/Users/tom/gocode/src/github.com/telennon/dnwTennis/testcamps/testfile.pdf"
//
//	err := pdf.OutputFileAndClose(fn)
//	if err != nil {
//		fmt.Println(err)
//	}
//}