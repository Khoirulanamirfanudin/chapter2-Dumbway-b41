package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var Data = map[string]interface{}{
	"Title": "Personal Web",
}

type Card struct {
	Title     string
	Post_date string
	Author    string
	Content   string
}

var Cards = []Card{
	{
		Title:     "Terserah Kamu mau isi apa",
		Post_date: "21 October 2022 22:20 WIB",
		Author:    "Khoirul Anam Irfanudin",
		Content:   "Isi apapun yang kamu suka",
	},
}

func main() {
	route := mux.NewRouter()

	// route path folder untuk public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	//routing
	route.HandleFunc("/hello", helloWorld).Methods("GET")
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/addProject", addProject).Methods("GET")
	route.HandleFunc("/card-detail/{index}", cardDetail).Methods("GET")
	route.HandleFunc("/form-card", formAddCard).Methods("GET")
	route.HandleFunc("/add-card", addCard).Methods("POST")
	route.HandleFunc("/delete/{index}", delete).Methods("GET")

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

	respData := map[string]interface{}{
		"Cards": Cards,
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

	var CardDetail = Card{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range Cards {
		if index == i {
			CardDetail = Card{
				Title:     data.Title,
				Content:   data.Content,
				Post_date: data.Post_date,
				Author:    data.Author,
			}
		}
	}

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

	fmt.Println("Title : " + r.PostForm.Get("inputTitle")) // value berdasarkan dari tag input name
	fmt.Println("Content : " + r.PostForm.Get("inputContent"))

	var title = r.PostForm.Get("inputTitle")
	var content = r.PostForm.Get("inputContent")

	var newCard = Card{
		Title:     title,
		Content:   content,
		Author:    "Khoirul Anam Irfanudin",
		Post_date: time.Now().String(),
	}

	Cards = append(Cards, newCard)
	// fmt.Println(Cards)

	http.Redirect(w, r, "/addProject", http.StatusMovedPermanently)
}

func delete(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	fmt.Println(index)

	Cards = append(Cards[:index], Cards[index+1:]...)
	fmt.Println(Cards)

	http.Redirect(w, r, "/addProject", http.StatusFound)
}
