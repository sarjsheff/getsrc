package main

import (
	"log"
	"net/http"
	"text/template"

	"sheff.online/getsrc/internal/getsrc"
)

func main() {
	config, err := getsrc.NewConfig("./getsrc.yaml")
	if err != nil {
		log.Fatal(err)
		return
	}

	for k, v := range *config.Repos {
		getsrc.RegDumbHTTPRepo(k, v.Path)
	}

	http.Handle("/css/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpls, err := template.ParseFiles("./tmpl/list.go.html", "./tmpl/icons.go.html", "./tmpl/common.go.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		err = tmpls.Execute(w, getsrc.NewHttpObject(config.Repos, nil))
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
