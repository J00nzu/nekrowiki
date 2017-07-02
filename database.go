package main

import (
	"fmt"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"time"
)

//var session mgo.Session

type Person struct {
	Name  string
	Phone string
}

var (
	Session  mgo.Session
	Database mgo.Database
)

func Connect() {
	session, err := mgo.Dial("server1.example.com,server2.example.com")

	Session = *session

	if err != nil {
		panic(err)
	}
	defer Session.Close()

	c := Session.DB("test").C("people")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
}

func Close() {

}
