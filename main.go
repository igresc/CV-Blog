package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

var path =: "/go/src/github.com/igresc/cv-blog/src"

var t *template.Template

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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := indexPage{
		Title: " Sergi Castro ",
	}

	if err := t.ExecuteTemplate(w, "home.html", p); err != nil {
		fmt.Println(err)
	}
}

func projectHandler(w http.ResponseWriter, r *http.Request) {
	resp, _ := http.Get("https://api.github.com/users/igresc/repos")

	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

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
	if err := t.ExecuteTemplate(w, "projects.html", p); err != nil {
		fmt.Println(err)
	}
}
func init() {
	t = template.New("")
	t.Funcs(template.FuncMap{"mod": func(i int) bool { return i%2 == 0 }})
	t = template.Must(t.ParseGlob("path/templates/*.html"))
}

func main() {
	cssHandler := http.FileServer(http.Dir(path+"/css/"))
	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/projects/", projectHandler)
	fmt.Println(http.ListenAndServe(":8000", nil))
}
