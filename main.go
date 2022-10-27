package main

import (
	"context"
	"day-10/connection"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var Data = map[string]interface{}{
	"Title": "Personal Web",
}

type Card struct {
	Id         int
	Name       string
	Start_date string
	End_date   string
	// Technologies string
	Author      string
	Description string
}

var Cards = []Card{
	// {
	// 	Title:     "Terserah Kamu mau isi apa",
	// 	Post_date: "21 October 2022 22:20 WIB",
	// 	Author:    "Khoirul Anam Irfanudin",
	// 	Content:   "Isi apapun yang kamu suka",
	// },
}

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	// route path folder untuk public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	//routing
	route.HandleFunc("/hello", helloWorld).Methods("GET")
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/addProject", addProject).Methods("GET")
	route.HandleFunc("/card-detail/{id}", cardDetail).Methods("GET")
	route.HandleFunc("/form-card", formAddCard).Methods("GET")
	route.HandleFunc("/add-card", addCard).Methods("POST")
	route.HandleFunc("/edit/{id}", edit).Methods("GET")
	route.HandleFunc("/delete/{id}", delete).Methods("GET")

	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", route)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text-html;charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Println(Cards)

	var tmpl, err = template.ParseFiles("views/addProject.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	rows, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description FROM tb_projects")

	var result []Card

	for rows.Next() {
		var each = Card{}

		var err = rows.Scan(&each.Id, &each.Name, &each.Start_date, &each.End_date, &each.Description)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Author = "Khoirul Anam Irfanudin"
		// each.Format_date = each.Post_date.Format("3 Maret 2008")

		result = append(result, each)
	}

	fmt.Println(result)

	respData := map[string]interface{}{
		"Cards": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func cardDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/cardDetail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var CardDetail = Card{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description FROM tb_projects WHERE id=$1", id).Scan(
		&CardDetail.Id, &CardDetail.Name, &CardDetail.Start_date, &CardDetail.End_date, &CardDetail.Description,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// for i, data := range Cards {
	// 	if index == i {
	// 		CardDetail = Card{
	// 			Name:        data.Name,
	// 			Description: data.Description,
	// 			Start_date:  data.Start_date,
	// 			End_date:    data.End_date,
	// 			Author:      data.Author,
	// 		}
	// 	}
	// }

	data := map[string]interface{}{
		"Card": CardDetail,
	}

	fmt.Println(data)

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func formAddCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/add-card.html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func addCard(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Name : " + r.PostForm.Get("inputName")) // value berdasarkan dari tag input name
	fmt.Println("Description : " + r.PostForm.Get("inputDescriptiom"))
	fmt.Println("Start_date : " + r.PostForm.Get("inputStartDate"))
	fmt.Println("End_date : " + r.PostForm.Get("inputEndDate"))
	// fmt.Println("Technologies : " + r.PostForm.Get("inputTechnologies"))

	var name = r.PostForm.Get("inputName")
	var description = r.PostForm.Get("inputDescription")
	var start_date = r.PostForm.Get("inputStartDate")
	var end_date = r.PostForm.Get("inputEndDate")
	// var technologies = r.PostForm.Get("inputTechnologies")

	// var newCard = Card{
	// 	Name:        name,
	// 	Description: description,
	// 	Start_date:  start_date,
	// 	End_date:    end_date,
	// 	Author:      "Khoirul Anam Irfanudin",
	// }

	// Cards = append(Cards, newCard)
	// // fmt.Println(Cards)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name, start_date, end_date, description) VALUES ($1, $2, $3, $4) ", name, start_date, end_date, description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/addProject", http.StatusMovedPermanently)
}

func edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	http.Redirect(w, r, "/form-card", http.StatusFound)

}

func delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)

	// Cards = append(Cards[:index], Cards[index+1:]...)
	// fmt.Println(Cards)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	http.Redirect(w, r, "/addProject", http.StatusFound)
}
