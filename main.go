package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var path string = "./"

var t *template.Template

// Repo struct for Github repositories
type Repo struct {
	Name        string `json:"name"`
	URL         string `json:"html_url"`
	Fork        bool   `json:"fork"`
	Description string `json:"description"`
}

type indexPage struct {
	Title string
}

type projectPage struct {
	Title string
	Repos []Repo
}

type cvPage struct {
	Title string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := indexPage{
		Title: " Sergi Castro ",
	}

	if err := t.ExecuteTemplate(w, "home.gohtml", p); err != nil {
		fmt.Println(err)
	}
}

func projectHandler(w http.ResponseWriter, r *http.Request) {
	var bytes []byte
	// resp, _ := http.Get("https://api.github.com/users/igresc/repos")

	// bytes, _ := ioutil.ReadAll(resp.Body)
	// resp.Body.Close()

	var rList []Repo
	json.Unmarshal(bytes, &rList)
	for i, repo := range rList {
		if repo.Fork {
			copy(rList[i:], rList[i+1:])
			rList[len(rList)-1] = Repo{} // or the zero value of T
			rList = rList[:len(rList)-1]
		}
	}
	p := projectPage{
		Title: " Sergi Castro Projects ",
		Repos: rList,
	}
	if err := t.ExecuteTemplate(w, "projects.gohtml", p); err != nil {
		fmt.Println(err)
	}
}

func cvHandler(w http.ResponseWriter, r *http.Request) {
	p := cvPage{
		Title: " Sergi Castro CV ",
	}

	if err := t.ExecuteTemplate(w, "cv.gohtml", p); err != nil {
		log.Fatal(err)
	}
}

func init() {
	t = template.New("")
	t.Funcs(template.FuncMap{"mod": func(i int) bool { return i%2 == 0 }})
	t = template.Must(t.ParseGlob(path + "/templates/*.gohtml"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>ERROR 404. Page not found!</h1>")
}

func main() {

	r := mux.NewRouter()

	// imgHandler := http.FileServer(http.Dir(path + "/img/"))
	// r.Handle("/img/", http.StripPrefix("/img/", imgHandler))
	// cssHandler := http.FileServer(http.Dir(path + "/css/"))
	// r.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	r.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("."+"/img/"))))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("."+"/css/"))))

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/projects", projectHandler)
	r.HandleFunc("/cv", cvHandler)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	fmt.Println(http.ListenAndServe(":3000", r))
}
