package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"time"
	"time"
	"errors"
	"log"
)

type Database struct {
	session *mgo.Session
}

type DataSess struct {
	session  *mgo.Session
	database *mgo.Database
}

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Email    string
	Username string
	Password string
	Roles    []string
}

// DATABASE struct functions

func (db *Database) Connect() {
	session, err := mgo.DialWithTimeout(config.DatabaseURL, 10)

	if err != nil {
		panic("Failed to connect to the database at URL: " + config.DatabaseURL)
	}

	db.session = session

	db.session.SetSocketTimeout(1 * time.Hour)
}

func (db *Database) EnsureConnected() {
	defer func() {
		if r := recover(); r != nil {
			db.Connect()
		}
	}()

	//Ping panics if session is closed. (see mgo.Session.Panic())
	db.session.Ping()
}

func (db *Database) GetSession() *DataSess {
	db.EnsureConnected()

	sess := DataSess{}
	sess.session = db.session.Copy()
	sess.database = sess.session.DB(config.DatabaseName)

	return &sess;
}

func (db *Database) CheckDefaultUser() {
	ds := db.GetSession()
	defer ds.session.Close()

	defaultUser := User{
		Email:    "admin@example.com",
		Username: "admin",
		Password: "",
		Roles:    []string{"admin"},
	}

	if val, err := ds.UserExists(defaultUser.Email, defaultUser.Username); err != nil {
		panic(err)
	} else if val {
		log.Println("Default user exists.. nothing to do")
		return
	} else if ds.AddUser(&defaultUser) != nil {
		panic(err)
	} else {
		log.Println("No default user found. Added user")
	}
}

func (db *Database) Close() {
	db.session.Close()
}

// DATASESS struct functions

func (db *DataSess) AddUser(user *User) error {
	var err error

	defer func() {
		if r := recover(); r != nil {
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			// return the modified err and rep
		}
	}()

	c := db.database.C("users")

	err = c.Insert(user)

	return err
}

func (db *DataSess) RemoveUser(email string) error {
	var err error = nil

	defer func() {
		if r := recover(); r != nil {
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			// return the modified err and rep
		}
	}()

	c := db.database.C("users")

	err = c.Remove(bson.M{"email": email})

	return err
}

func (db *DataSess) UpdateUser(user *User) error {
	var err error = nil

	defer func() {
		if r := recover(); r != nil {
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			// return the modified err and rep
		}
	}()

	c := db.database.C("users")

	err = c.UpdateId(user.ID, user)

	return err
}

func (db *DataSess) GetUser(email string, username string) (User, error) {
	var ret User
	var err error = nil

	defer func() {
		if r := recover(); r != nil {
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			// return the modified err and rep
		}
	}()

	c := db.database.C("users")

	result := User{}
	err = c.Find(bson.M{"email": email}).One(&result)

	if err != nil {
		err = c.Find(bson.M{"username": username}).One(&result)
		if (err != nil) {
			err = nil
		}
	}

	if result.Email != "" {
		ret = result;
	}

	return ret, err
}

func (db *DataSess) UserExists(email string, username string) (bool, error) {
	var ret bool = false
	var err error = nil

	defer func() {
		if r := recover(); r != nil {
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			// return the modified err and rep
		}
	}()

	c := db.database.C("users")

	result := User{}
	err = c.Find(bson.M{"email": email}).One(&result)

	if (err != nil) {
		err = nil
	}

	if result.Email == "" {
		err = c.Find(bson.M{"username": username}).One(&result)
	}

	if (err != nil) {
		err = nil
	}

	if result.Email != "" {
		ret = true
	}

	return ret, err
}

