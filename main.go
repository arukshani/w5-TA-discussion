package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type GradeInfo struct {
	Name       string
	Grade      string
	UploadTime string
}

var students []GradeInfo

func handleIndex(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, students); err != nil {
			http.Error(w, "Template Error", http.StatusInternalServerError)
		}
	}
}

func handleUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now()
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Template Error", http.StatusInternalServerError)
			return
		}
		studentGrade := GradeInfo{Name: r.FormValue("student_name"), Grade: r.FormValue("student_grade"),
			UploadTime: currentTime.Format("2006.01.02 15:04:05")}
		students = append(students, studentGrade)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func findStudent(name string) (GradeInfo, bool) {
	for _, gradeinfo := range students {
		if gradeinfo.Name == name {
			return gradeinfo, true // Return the person and true if found
		}
	}
	return GradeInfo{}, false // Return an empty Person and false if not found
}

func handleStudentInfo(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Invalid Request!", http.StatusBadRequest)
			return
		}

		foundStudent, found := findStudent(id)

		if found {
			if err := tmpl.Execute(w, foundStudent); err != nil {
				http.Error(w, "Template error", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Student record not found", http.StatusNotFound)
			return
		}
	}
}

func main() {

	tmplIndex := template.Must(template.ParseFiles("index.html"))
	tmplInfo := template.Must(template.ParseFiles("gradeinfo.html"))

	http.HandleFunc("/", handleIndex(tmplIndex))
	http.HandleFunc("/upload", handleUpload())
	http.HandleFunc("/student", handleStudentInfo(tmplInfo))

	addr := ":8080"
	log.Printf("Server started on http://localhost:%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
