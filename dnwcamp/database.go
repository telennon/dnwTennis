package dnwcamp
// Copyright 2014 - Tom Lennon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This package manages database connections for DNW Camps
//
// Original created May 2014

import (
	"fmt"
	"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"
)

const (
	MONGODBSERVER = "localhost"			// DataBase URL
	DATABASE = "CampMaster"					// Database name containing the app registration collection
	COL_CAMPS = "Camps"			// Name of camp configuration collection
	COL_REG = "Registrations"	// All camp registrations
)

	


func BuildDBSession() (*mgo.Session, error) {
	
	session, err := mgo.Dial(MONGODBSERVER)
	if err != nil {
		fmt.Println("Dialing the database presented the following error\n", err)
		return session, err
	}

	session.SetMode(mgo.Monotonic, true)
	
	return session, nil
}

func ConnectCollection(s *mgo.Session, db string, collection string) (*mgo.Collection) {
		c := s.DB(db).C(collection)
		return c
}

func CloseConnection(s *mgo.Session) {
	defer s.Close()
}

func OpenCampCollection() (*mgo.Session, *mgo.Collection, error) {
		sess, err := BuildDBSession()
		if err != nil {
			return nil, nil, err
		}

		col := ConnectCollection(sess, DATABASE, COL_CAMPS)

		// defer sess.Close()

		return sess, col, err
}

func DeleteDNWTennisDB() error {
	// TODO - Check for database before dropping
	sess, err := BuildDBSession()
	if err == nil {
		err = sess.DB(DATABASE).DropDatabase()
	}
	sess.Close()
	return err
}

func OpenRegistrationCollection() (*mgo.Session, *mgo.Collection, error) {
		sess, err := BuildDBSession()
		if err != nil {
			return nil, nil, err
		}

		col := ConnectCollection(sess, DATABASE, COL_REG)

		// defer sess.Close()

		return sess, col, err
}

func DeleteRegistrationCollection() error {
	sess, col, err := OpenRegistrationCollection()
	if err == nil {
		err = col.DropCollection()
	}
	sess.Close()
	return err
}

	// Setup the database connection
	//sess, err := dnwcamp.BuildDBSession()
	//if err != nil {
	//	return 											// No database then just give up
	//}
	//col := dnwcamp.ConnectCollection(sess, dnwcamp.DATABASE, dnwcamp.COL_CAMPS)
	//defer sess.Close()

	//err = sess.DB(dnwcamp.DATABASE).DropDatabase()
	//if err != nil {
	//	fmt.Println("Failed to delete the database with error ..", err)
	//}