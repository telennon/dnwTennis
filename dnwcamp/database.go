package dnwcamp
// Copyright 2014 - Tom Lennon.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//
// This package manages database connections for DNW Camps
//
// Original created May 2014

import (
	"fmt"
	"labix.org/v2/mgo"
	"os"
)

const (
	// MONGODBSERVER = "mongodb://dnwUser:dnw0001@ds051459.mongolab.com:51459/campmaster" // Database URL
	DATABASE = "campmaster"					// Database name containing the app registration collection
	COL_CAMPS = "Camps"			// Name of camp configuration collection
	COL_REG = "Registrations"	// All camp registrations
)

	


func BuildDBSession() (*mgo.Session, error) {
	mongoURI := os.Getenv("MONGOURI")
	if mongoURI == "" {
		fmt.Println("Database connection string was not found in the environment - Quitting")
		os.Exit(12)
	}

	session, err := mgo.Dial(mongoURI)
	if err != nil {
		fmt.Println("Dialing the database presented the following error..", err)
		os.Exit(10)
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
		// TODO Add some code here to check that the collection exists prior to dropping
		err = col.DropCollection()
	}
	sess.Close()
	return err
}