package main

import (
	"flag"
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name  string
	Phone string
}

func main() {
	server := flag.String("serverdb", "localhost", "address mongoDB")
	flag.Parse()
	fmt.Println("server: ", *server)

	session, err := mgo.Dial(*server)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")
	err = c.Insert(&Person{"Robot1", "+00 09 8765 4321"},
		&Person{"Robot2", "+12 34 5678 9000"})
	if err != nil {
		log.Fatal(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Robot1"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
}
