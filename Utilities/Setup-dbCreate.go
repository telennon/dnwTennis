package main
/* ************************************************************************
*	Author		TELennon
*	Created		Jan 2014
*
*	Copyright 2014 - Tom Lennon.  All rights reserved.
*	Use of this source code is governed by a MIT-style
*	license that can be found in the LICENSE.md file.

*	Setup-dbCreate.go
*		Setup-dbCreate will convert past tennis databases to the new mongoDB form.
* 		Past Data exists for the years 2012 - 2014 and all will be converted into a
*		single database
*
*		The new format consists of two collections camps and registrations
*		Camp documents are new and will be created from data structures.
*		Registration and camper data will be read from CSV files stored in the _dataImport
*		directory.
*
************************************************************************* */

import (
	"fmt"
	"os"

	"github.com/telennon/dnwTennis/dnwcamp"
)

func main() {
	
	// Delete the entire Database
	err := dnwcamp.DeleteDNWTennisDB()
	if err != nil {
		fmt.Println("Error occured trying to delete the DNWTennis Database\n", err)
		os.Exit(1)
	}

	// Create the camps archive
	err = createCampsArchive()
	if err != nil {
		fmt.Println("Error occured trying to create the Camps Collection\n", err, "\nStopping the application")
		os.Exit(2)
	}

	// Create the registrations archive
	err = createRegistrationsArchive()
	if err != nil {
		fmt.Println("Error occured trying to create registration archive\n", err, "\nStopping the application")
	}
}