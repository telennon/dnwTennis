package main
/* ************************************************************************
*	Author		TELennon
*	Created		Jan 2014
*
*	Copyright 2014 - Tom Lennon.  All rights reserved.
*	Use of this source code is governed by a MIT-style
*	license that can be found in the LICENSE.md file.
*
*	SetupRegistrationArchive.go
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
	"strconv"
	"strings"
	"errors"
	"time"
	//"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"encoding/csv"
	//"io"
	"os"

	"github.com/telennon/dnwTennis/dnwcamp"
)

const (
	//MongoSrv = "localhost"			// DataBase URL
	//DB = "CampMaster"					// Database name containing the app registration collection
	//COL_CONF = "Camps"			// Name of camp configuration collection
	//COL_REG = "Registrations"	// All camp registrations
	// Indexes into camper record
	C_ID = 0
	C_LOTTERYNUM = 1
	C_NAME = 2
	C_AGE = 3
	C_TSHIRT = 4
	C_HOMEOWNER = 5
	C_CAMPS = 6
	C_WAITLISTS = 7
	C_PHPTIMESTAMP = 8
	C_ORGCOMMENTS = 9
	C_YEAR = 10
	I_RECORDTYPE = 0
	S1_YEAR = 1
	S1_SECTION = 2
	R_NAME = 1
)

type campFiles struct {
	style int64
	year string
	payorcombinedfile	string
	payorFields int
	camperfile	string
	camperFields int
}

func main() {

	camperDataArchive := []campFiles {
		{
			style: 1,
			year: "2012",
			payorcombinedfile: "/Users/tom/gocode/src/github.com/telennon/dnwtennis/_dataImport/2012TennisCamp1.csv",
			payorFields: 8,
			camperfile: "",
			camperFields: 0},
		{
			style: 2,
			year: "2013",
			payorcombinedfile: "/Users/tom/gocode/src/github.com/telennon/dnwtennis/_dataImport/payors2013.csv",
			payorFields: 10,
			camperfile: "/Users/tom/gocode/src/github.com/telennon/dnwtennis/_dataImport/campers2013.csv",
			camperFields: 11},
		{
			style: 2,
			year: "2014",
			payorcombinedfile: "/Users/tom/gocode/src/github.com/telennon/dnwtennis/_dataImport/payors2014.csv",
			payorFields: 10,
			camperfile: "/Users/tom/gocode/src/github.com/telennon/dnwtennis/_dataImport/campers2014.csv",
			camperFields: 11},		
	}

	// Drop any existing collection
	err := dnwcamp.DeleteRegistrationCollection
	if err != nil {
		fmt.Println("Attempt to delete the collection", dnwcamp.COL_REG, "returned the following error ..", err)
		// TODO - Add some code to check for the existance of the collection before attempting to drop it
	}


	for i := 0; i < len(camperDataArchive); i++ {
		switch {
		case camperDataArchive[i].style == 1:
			err := importStyle1Archive(camperDataArchive[i])
			if err != nil {
				fmt.Println("Error importing archive from", camperDataArchive[i].year)
			}
		case camperDataArchive[i].style == 2:
			err := importStyle2Archive(camperDataArchive[i])
			if err != nil {
				fmt.Println("Error importing archive from", camperDataArchive[i].year)
			}
		case camperDataArchive[i].style > 2:
			fmt.Println("Whoa - Some unexpected import style showed up")
		}
	}
}

func importStyle2Archive(campArchive campFiles) error {
	fmt.Println("Importing a style 2 archive")

	// Read in the registration records
	rawRegistrations, err := readCSV(campArchive.payorcombinedfile, campArchive.payorFields)
	if err != nil {
		fmt.Println("Error reading registration file ", campArchive.payorcombinedfile)
		fmt.Println("With Error:", err)
		return err
	}

	// Read in the camper records
	rawCampers, err := readCSV(campArchive.camperfile, campArchive.camperFields)
	if err != nil {
		fmt.Println("Error reading camper file ", campArchive.camperfile)
		fmt.Println("With Error:", err)
		return err
	}	

	// Get the current camp to add registrations to
	currentCamp, err := findJustTheCamp(campArchive.year)
	if err != nil {
		fmt.Println("Could not find a camp with title containing year: ", campArchive.year)
		return err
	}

	for i := 1; i < len(rawRegistrations); i++ {
		// Create a registration record
		currentRegistration, oldId, err := regFromRaw(rawRegistrations[i], currentCamp)
		if err != nil {
			fmt.Println("Error processing this registration")
			fmt.Println(rawRegistrations[i])
			return err
		}
		for y := 1; y < len(rawCampers); y++ {
			// If the record belongs
			rawCamperCampId, err := strconv.ParseInt(rawCampers[y][1], 10, 64)
			if err != nil {
				fmt.Println("Encountered a bad ID value for a camper")
				return err
			}
			if rawCamperCampId == oldId {
				// create a camper record
				currentCamper := camperFromRaw(rawCampers[y], currentRegistration, currentCamp)
				// append to current registration
				currentRegistration.RegCampers = append(currentRegistration.RegCampers, currentCamper)

			}
		}
		// save the registration record
		err = currentRegistration.Save()
		if err != nil {
			fmt.Println("Error saving registration record")
			fmt.Printf("%+v\n", currentRegistration)
			return err
		}
	}
	return nil
}

func regFromRaw(r []string, c dnwcamp.Camp) (dnwcamp.Registration, int64, error) {

	aNewReg := dnwcamp.NewRegistration()
	
	// Copy across the simple strings
	aNewReg.RegName = r[1]
	aNewReg.RegPhone = r[2]
	aNewReg.RegSite = r[3]
	aNewReg.RegEmail = r[4]
	
	// Insert reference to camp record
	aNewReg.RegCampId = c.ID

	// Convert PHP timestamp to string date / time
	iTime, err := strconv.ParseInt(r[8], 10, 64)
	if err != nil {
		fmt.Println("Attempting to convert ", r[8], "and had error", err)
	}
	t := time.Unix(iTime, 0)
	aNewReg.RegTimeStamp = t.Format(time.RFC3339)
	
	// Convert Camp requests to something for filling in later
	iFirst, err := strconv.ParseInt(r[5], 10, 64)
	if err != nil {
		fmt.Println("Error converting first Choice..", err)
		return aNewReg, 0, err
	}
	if ((iFirst >= 1) && (iFirst <= 3))  {
		// 		Create a request structure
		// 		Fill with Priority 1 and ID of camp requested
		// 		Append to request array in registration
		aNewRequest := dnwcamp.NewRequest(c.Sections[iFirst-1].ID, 1)
		aNewReg.RegRequests = append(aNewReg.RegRequests, aNewRequest)

	} else {
		fmt.Println("First Choice was not in range - Expected a number between 1 and 3 but got ", r[6])
	}
	
	// repeat all the stuff done for first choice but using priority 2
	// Second choices are not always filled in in the data so need to check for nil then force a 4 = No Choice
	if r[6] == "" {
		r[6] = "4"
	}

	iSecond, err := strconv.ParseInt(r[6], 10, 64)
	if err != nil {
		fmt.Println("Error converting second Choice..", err)
		return aNewReg, 0, err
	}
	if ((iSecond >= 1) && (iSecond <= 3))  {
		// 		Create a request structure
		// 		Fill with Priority 2 and ID of camp requested
		// 		Append to request array in registration
		aNewRequest := dnwcamp.NewRequest(c.Sections[iSecond-1].ID, 2)
		aNewReg.RegRequests = append(aNewReg.RegRequests, aNewRequest)
	} else {
		if iSecond != 4 {
			fmt.Println("second Choice was not in range - Expected a number between 1 and 3 but got ", r[6])
		}
	}

	// Collect the old ID for the return value
	iOldId, err := strconv.ParseInt(r[0], 10, 64)
	if err != nil {
		fmt.Println("Error converting old id to Int - Value is .. ", r[0])
		iOldId = 0
	}

	return aNewReg, iOldId, nil
}

func camperFromRaw(record []string, reg dnwcamp.Registration, theCamp dnwcamp.Camp) dnwcamp.Camper {
	// Create a new camper structure
	aCamper := dnwcamp.NewCamper()

	// Simple Data moves for name, age, T-Shirt
	aCamper.Name = record[C_NAME]
	aCamper.Age = record[C_AGE]
	aCamper.ShirtSize = record[C_TSHIRT]
	
	// Convert Homeowner flag to index into camper type
	if record[C_HOMEOWNER] == "1" {
		aCamper.Type = "Homeowner"
	} else {
		if record[C_HOMEOWNER] == "0" {
			aCamper.Type = "Guest"
		} else {
			aCamper.Type = "Oops"
		}
	}

	// Convert Camp value to an array of references to camp sections

	iCampValue, err := strconv.ParseInt(record[C_CAMPS], 10, 64)
	if err != nil {
		fmt.Println("Error converting Camp Value to integer .. ", record[C_CAMPS])
		fmt.Println(err)
		return aCamper
	}
	aCamper.Sections = convertSections(iCampValue, theCamp)
	// Camp value was defined as follows. 1 = Camp 1, 10 = Camp 2; 100 = Camp 3
	// 	You could add them together to be in multiple camps so 11 would mean camp 1 and 2

	// Convert waitlist value to an array of references to camp sections

	iWaitListValue, err := strconv.ParseInt(record[C_WAITLISTS], 10, 64)
	if err != nil {
		fmt.Println("Error converting Camp Value to integer .. ", record[C_WAITLISTS])
		fmt.Println(err)
		return aCamper
	}
	aCamper.WaitList = convertSections(iWaitListValue, theCamp)

	// Convert time stamp to a huma readable date/time
	iTime, err := strconv.ParseInt(record[C_PHPTIMESTAMP], 10, 64)
	if err != nil {
		fmt.Println("Attempting to convert camper timestamp ", record[C_PHPTIMESTAMP], "and had error", err)
	}
	t := time.Unix(iTime, 0)
	aCamper.TimeStamp = t.Format(time.RFC3339)

	// Org comments will be all empty but we will create anyway and stick a copy
	//		of the old record into the comments
	aCamper.Comment = strings.Join(record, ",")

	// Year is discarded because we are linked to the camp through the registration

	return aCamper
}

func convertSections(i int64, c dnwcamp.Camp) ([]bson.ObjectId) {
	sectionlist := []bson.ObjectId{}
	idx := 0
	for i > 0 {
		if ((i % 2) != 0) {
			sectionlist = append(sectionlist, c.Sections[idx].ID)
		}
		idx++
		i = (i / 10)
	}
	return sectionlist
}


func importStyle1Archive(campArchive campFiles) error {
	rawData, err := readCSV(campArchive.payorcombinedfile, campArchive.payorFields)
	if err != nil {
		fmt.Println("Error trying to read data with error")
		fmt.Println(err)
		return err
	}

	iYearRecords := 0
	iBlankRecords := 0
	iRegistrationRecords := 0
	iCamperRecords := 0
	alreadyBuildingAReg := false

	var currentCampId bson.ObjectId
	var currentSectionId bson.ObjectId
	//var currentRegistrationId bson.ObjectId
	var currentRegistration dnwcamp.Registration

	for i := 0; i < len(rawData); i++ {
		switch {
		case rawData[i][I_RECORDTYPE] == "Y":
			// The year record encodes the camp year and the camp section all 
			//	entries that follow are related to. When a year record is seen 
			//  currentCampId and currentSectionId variables are updated
			iYearRecords++
			currentCampId, currentSectionId, err = findTheCamp(rawData[i][S1_YEAR], rawData[i][S1_SECTION])
			if err != nil {
				fmt.Println("Searching for a unique camp and section using:")
				fmt.Println("Camp keyword", rawData[i][S1_YEAR])
				fmt.Println("Section Keyword", rawData[i][S1_SECTION])
				fmt.Println("With the error: ", err)
				return errors.New("Not able to process the entire Style 1 Data file")
			}
		case rawData[i][I_RECORDTYPE] == "B":
			// B records are simply Blank - They may or may not have data in them 
			// but are to be ignored
			iBlankRecords++
		case rawData[i][I_RECORDTYPE] == "R":
			// R Records are registration records - When one is seen it creates a new 
			// registration and updates the currentRegistration variable
		//reg := registration.NewRegistration()
		//err := buildRegistration(*reg, rawData[i])
			if alreadyBuildingAReg {
				// Save the registration we have been building 
				// Set alreadyBui.. to false
				err := currentRegistration.Save()
				if err != nil {
					fmt.Printf("Error saving this registration\n%+v\nWith the following err\n%v", currentRegistration, err)
					return err
				}
				alreadyBuildingAReg = false
			}

			currentRegistration = dnwcamp.NewRegistration()
			currentRegistration.RegName = rawData[i][R_NAME]
			currentRegistration.RegCampId = currentCampId
			currentRegistration.RegRequests = append(currentRegistration.RegRequests, dnwcamp.NewRequest(currentSectionId, 1))
			// For the old spreadsheet style type 1 conversion the rest of the information in a registration
			//		record must be blank because that data was not tracked in the old spreadsheet

			iRegistrationRecords++
			alreadyBuildingAReg = true

		case rawData[i][I_RECORDTYPE] == "C":
			// C records are camper records - When one is seen the Camper is
			// added to the currentRegistration
			currentCamper := dnwcamp.NewCamper()
			// Camper records look like the following
			// RecordType=C,Empty,Name,Age,HomeownerFlag,GuestFLag,ShirtSize,Comment
			currentCamper.Name = rawData[i][2]
			currentCamper.Age = rawData[i][3]
			if rawData[i][4] == "1" {
				currentCamper.Type = "Homeowner"
			} else {
				if rawData[i][5] == "1" {
					currentCamper.Type = "Guest"
				} else {
					currentCamper.Type = "Unknown"
				}
			}
			// Camper TimeStamp was not maintained in the old data - Nil here
			currentCamper.ShirtSize = rawData[i][6]
			currentCamper.Sections = append(currentCamper.Sections, currentSectionId)

			currentRegistration.RegCampers = append(currentRegistration.RegCampers, currentCamper)
			iCamperRecords++
		}
	}

	// When we exit the loop we might have one left over unsaved registration that needs to be save.
	if alreadyBuildingAReg {
		// Save the registration we have been building 
		// Set alreadyBui.. to false
		err := currentRegistration.Save()
		if err != nil {
			fmt.Printf("Error saving this registration\n%+v\nWith the following err\n%v", currentRegistration, err)
			return err
		}
		alreadyBuildingAReg = false
	}

	fmt.Println("Processed ", iYearRecords, "Year Records")
	fmt.Println("Processed ", iBlankRecords, "Blank Records")
	fmt.Println("Processed ", iRegistrationRecords, "Registration Records")
	fmt.Println("Processed ", iCamperRecords, "Camper Records")
	return nil
}

func findJustTheCamp(yr string) (dnwcamp.Camp, error) {
	allCamps, _ := dnwcamp.ListCamps()

	for i := 0; i < len(allCamps); i++ {
		if strings.Contains(allCamps[i].Title, yr) {
			// This is just a conversion routine and will return the first camp found
			return allCamps[i], nil
		}
	}
	return allCamps[0], errors.New("No match found")
}

func findTheCamp(yr, sec string) (bson.ObjectId, bson.ObjectId, error) {
	// findTheCamp is a search routine that scans the list of camps looking for a
	// camp title that contains the string supplied in yr and then a section that 
	// contains the string in sec
	// Matching more than one camp or section will return an error
	//	Not matching any will also be an error
	var campMatches []bson.ObjectId
	var sectionMatches []bson.ObjectId

	allCamps, _ := dnwcamp.ListCamps()
	//fmt.Println(allCamps)

	for i := 0; i < len(allCamps); i++ {
		if strings.Contains(allCamps[i].Title, yr) {
			campMatches = append(campMatches, allCamps[i].ID)
			for y := 0; y < len(allCamps[i].Sections); y++ {
				if strings.Contains(allCamps[i].Sections[y].Name, sec) {
					sectionMatches = append(sectionMatches, allCamps[i].Sections[y].ID)
				}
			} 
		}
	}
	// Outside the evaluation loop we do some error checking on what we found
	// camp or section length = 0 -- Error, no camp or section found
	// camp or section length > 1 -- Error, no unique camp or section was found
	if (len(campMatches) + len(sectionMatches)) == 0 {
		return bson.NewObjectId(), bson.NewObjectId(), errors.New("No matching camp or section could be found")
	}
	if (len(campMatches) + len(sectionMatches)) > 2 {
		return bson.NewObjectId(), bson.NewObjectId(), errors.New("No unique camp or section could be found")
	}
	return campMatches[0], sectionMatches[0], nil
}

//func buildRegistration(r *registration, d []string) (error) {
//
//}

func readCSV(fName string, fCount int) ([][]string, error) {
	cFile, err := os.Open(fName)
	if err != nil {
		fmt.Println("Error opening file ", fName, " with error ", err)
		return nil, err
	}
	defer cFile.Close()

	reader := csv.NewReader(cFile)
	reader.FieldsPerRecord = fCount
	v, err := reader.ReadAll()
	return v, err
}
