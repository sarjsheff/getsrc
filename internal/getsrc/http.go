package getsrc

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/go-git/go-git/v5/plumbing"
)

type HttpObject struct {
	Repo   *Repo
	Repos  *map[string]ConfigRepo
	Start  time.Time
	Config *Config
}

func NewHttpObject(Repos *map[string]ConfigRepo, Repo *Repo, config *Config) *HttpObject {
	ret := &HttpObject{
		Start:  time.Now(),
		Repos:  Repos,
		Repo:   Repo,
		Config: config,
	}
	return ret
}

func (o *HttpObject) MD5(text string) string {
	hash := md5.Sum([]byte(text))
	return fmt.Sprintf("%x", hash)
}

func (o *HttpObject) Gravatar(email string) string {
	hash := md5.Sum([]byte(email))

	return fmt.Sprintf("https://www.gravatar.com/avatar/%x", hash)
}

func (o *HttpObject) Now() time.Time {
	return time.Now()
}

func (o *HttpObject) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (o *HttpObject) SinceHuman(t time.Time) string {
	return humanize.Time(t)
}

func (o *HttpObject) ExecTime() time.Duration {
	return time.Since(o.Start)
}

func RegDumbHTTPRepo(name string, repopath string, config *Config) {
	repo, err := NewRepo(name, repopath)
	if err != nil {
		log.Println(err)
		return
	}

	contextUrl := "/git/" + name + "/"

	http.HandleFunc(contextUrl, func(w http.ResponseWriter, r *http.Request) {
		if p, ok := strings.CutPrefix(r.URL.Path, contextUrl); ok {
			if p == "info/refs" {
				it, err := repo.Repo.References()
				if err == nil {
					ref, err := it.Next()
					for err == nil {
						if ref.Type() == 1 {
							fmt.Fprintf(w, "%s\t%s\n", ref.Hash().String(), ref.Name())
						}
						ref, err = it.Next()
					}
				}
			} else if p == "HEAD" {
				h, err := repo.Repo.Head()
				if err == nil {
					fmt.Fprintf(w, "ref: %s", h.Name())
				}
			} else if strings.HasPrefix(p, "objects") {
				if cfg, err := repo.Repo.Config(); err == nil {
					prefix := ".git"
					if cfg.Core.IsBare {
						prefix = ""
					}
					if ss, ok := strings.CutPrefix(p, "objects/"); ok {
						if repo.Repo.Storer.HasEncodedObject(plumbing.NewHash(strings.ReplaceAll(ss, "/", ""))) == nil {
							bt, err := os.ReadFile(path.Join(repopath, prefix, p))
							if err == nil {
								w.Write(bt)
							} else {
								log.Println(err)
							}
						}

					}

				}

			} else {
				tmpls, err := template.ParseFiles("./tmpl/single.go.html", "./tmpl/icons.go.html", "./tmpl/common.go.html", "./tmpl/gen.go.html")
				if err != nil {
					log.Println(err)
					w.WriteHeader(500)
					return
				}

				err = tmpls.Execute(w, NewHttpObject(nil, repo, config))
				if err != nil {
					log.Println(err)
					w.WriteHeader(500)
					return
				}
			}
		}
	})

}
