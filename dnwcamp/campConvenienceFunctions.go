package dnwcamp
// Copyright 2014 - Tom Lennon.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.
//
// This package manages the concept of a camp.
// I may decide to merge camp and sections packages as they are linked, a section
// must be connected to a camp to make sense and a camp with no sections is nonsensicle as well
//
// Original created April 2014


import (
	"fmt"
	//"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"
)

func TestLinkUp() {
	fmt.Println("Checking some connections")
	fmt.Println(MONGODBSERVER)
}
