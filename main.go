package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Course struct {
	Courseid    string  `json:"courseid"`
	Coursename  string  `json:"coursename"`
	Courseprice int     `json:"courseprice"`
	Author      *Author `json:"author"`
}

type Author struct {
	Fullname string `json:"fullname"`
	Website  string `json:"website"`
}

var courses []Course //db

func (c *Course) IsEmpty() bool {
	if c.Courseid == "" {
		return true
	}
	return false
}

func servehome(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Welcome to the api server"))

}

func getallcourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllCourses")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getcoursebyid(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getCourseByID")

	w.Header().Set("Content-Type", "application/json")
	Parms := mux.Vars(r)
	for _, course := range courses {
		if course.Courseid == Parms["id"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	json.NewEncoder(w).Encode("Course not found")

}

func createonecourse(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint Hit: createOneCourse")
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send a request body")
		return
	}
	//check title is dupilcate
   

	//what if the body is - {}
	
	var course Course
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	course.Courseid = strconv.Itoa(randGen.Intn(100))

	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
	return

}

func updatecourse(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint Hit: updateCourse")
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send a request body")
		return
	}

	Parms := mux.Vars(r)
	var course Course
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	for i, item := range courses {
		if item.Courseid == Parms["id"] {
			courses[i] = course
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	json.NewEncoder(w).Encode("Course not found")
}

func deletecourse(w http.ResponseWriter, r *http.Request) {
	
	fmt.Println("Endpoint Hit: deleteCourse")
	w.Header().Set("Content-Type", "application/json")

	Parms := mux.Vars(r)
	for i, item := range courses {
		if item.Courseid == Parms["id"] {
			courses = append(courses[:i], courses[i+1:]...)
			json.NewEncoder(w).Encode("Course deleted")
			return
		}
	}
	json.NewEncoder(w).Encode("Course not found")
}
func main() {
	
	fmt.Println("Starting the application...")
	courses = append(courses, Course{Courseid: "1", Coursename: "Python", Courseprice: 1000, Author: &Author{Fullname: "John Doe", Website: "www.johndoe.com"}})
	courses = append(courses, Course{Courseid: "2", Coursename: "Java", Courseprice: 2000, Author: &Author{Fullname: "Jane Doe", Website: "www.janedoe.com"}})
	
	r := mux.NewRouter()
	r.HandleFunc("/", servehome).Methods("GET")
	r.HandleFunc("/courses", getallcourses).Methods("GET")
	r.HandleFunc("/courses/{id}", getcoursebyid).Methods("GET")
	r.HandleFunc("/courses", createonecourse).Methods("POST")
	r.HandleFunc("/courses/{id}", updatecourse).Methods("PUT")
	r.HandleFunc("/courses/{id}", deletecourse).Methods("DELETE")


	
	log.Fatal(http.ListenAndServe(":8080", r))
}
