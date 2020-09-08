package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/robfig/cron.v2"
)

var path string = "./"
var rList []Repo
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
	Path  string
}

type projectPage struct {
	Title string
	Repos []Repo
	Path  string
}

type cvPage struct {
	Title string
	Path  string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := indexPage{
		Title: "Sergi Castro",
		Path:  r.URL.Path,
	}

	if err := t.ExecuteTemplate(w, "home.gohtml", p); err != nil {
		fmt.Println(err)
	}
}

func requestGithubProjects() {
	fmt.Println("Cron executed")
	resp, _ := http.Get("https://api.github.com/users/igresc/repos")
	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	json.Unmarshal(bytes, &rList)
	for i, repo := range rList {
		if repo.Fork {
			copy(rList[i:], rList[i+1:])
			rList[len(rList)-1] = Repo{} // or the zero value of T
			rList = rList[:len(rList)-1]
		}
	}
}

func projectHandler(w http.ResponseWriter, r *http.Request) {
	// var bytes []byte
	// resp, _ := http.Get("https://api.github.com/users/igresc/repos")
	// bytes, _ := ioutil.ReadAll(resp.Body)
	// resp.Body.Close()

	// var rList []Repo
	// json.Unmarshal(bytes, &rList)
	// for i, repo := range rList {
	// 	if repo.Fork {
	// 		copy(rList[i:], rList[i+1:])
	// 		rList[len(rList)-1] = Repo{} // or the zero value of T
	// 		rList = rList[:len(rList)-1]
	// 	}
	// }
	p := projectPage{
		Title: "Sergi Castro Projects",
		Repos: rList,
		Path:  r.URL.Path,
	}
	if err := t.ExecuteTemplate(w, "projects.gohtml", p); err != nil {
		fmt.Println(err)
	}
}

func cvHandler(w http.ResponseWriter, r *http.Request) {
	p := cvPage{
		Title: "Sergi Castro CV",
		Path:  r.URL.Path,
	}

	if err := t.ExecuteTemplate(w, "cv.gohtml", p); err != nil {
		log.Fatal(err)
	}
}

func init() {
	requestGithubProjects()

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

	c := cron.New()
	c.AddFunc("@daily", requestGithubProjects)
	c.Start()

	r := mux.NewRouter()

	r.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("."+"/img/"))))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("."+"/css/"))))

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/projects", projectHandler)
	r.HandleFunc("/cv", cvHandler)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	fmt.Println(http.ListenAndServe(":80", r))
}
