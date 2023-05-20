package getsrc

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/dustin/go-humanize"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	htmlgold "github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/mermaid"
)

type HttpObject struct {
	Repo    *Repo
	Repos   *map[string]ConfigRepo
	Start   time.Time
	Config  *Config
	SubPath string
	Path    string
}

func NewHttpObject(Repos *map[string]ConfigRepo, Repo *Repo, config *Config, subpath string, rawpath string) *HttpObject {
	ret := &HttpObject{
		Start:   time.Now(),
		Repos:   Repos,
		Repo:    Repo,
		Config:  config,
		SubPath: subpath,
		Path:    rawpath,
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

// Конвертация markdown
func (o *HttpObject) ToHtml(content string) string {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Strikethrough,
			extension.Linkify,
			extension.TaskList,
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			extension.Typographer,
			extension.CJK,
			meta.Meta,
			&mermaid.Extender{},
			// &pikchr.Extender{},
			highlighting.NewHighlighting(
				highlighting.WithStyle("manni"),
				highlighting.WithFormatOptions(
					html.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			htmlgold.WithUnsafe(),
		),
	)
	context := parser.NewContext()
	var buf bytes.Buffer
	if err := markdown.Convert([]byte(content), &buf, parser.WithContext(context)); err != nil {
		return ""
	}

	return buf.String()
}

type HTTP struct {
	singleTmpl *template.Template
	listTmpl   *template.Template
	config     *Config
}

func NewHTTP(config *Config) (*HTTP, error) {
	ret := &HTTP{config: config}

	tmpls, err := template.ParseFiles("./tmpl/single.go.html", "./tmpl/icons.go.html", "./tmpl/common.go.html", "./tmpl/gen.go.html")
	if err != nil {
		return nil, err
	}
	ret.singleTmpl = tmpls

	tmpls, err = template.ParseFiles("./tmpl/list.go.html", "./tmpl/icons.go.html", "./tmpl/common.go.html", "./tmpl/gen.go.html")
	if err != nil {
		return nil, err
	}
	ret.listTmpl = tmpls

	for k, v := range *config.Repos {
		ret.RegDumbHTTPRepo(k, v.Path, config)
	}

	http.Handle("/css/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/", ret.ListHandler)

	return ret, nil
}

func (h *HTTP) ListHandler(w http.ResponseWriter, r *http.Request) {
	err := h.listTmpl.Execute(w, &HttpObject{
		Start:   time.Now(),
		Repos:   h.config.Repos,
		Repo:    nil,
		Config:  h.config,
		SubPath: "",
		Path:    r.URL.Path,
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
}

func (h *HTTP) RegDumbHTTPRepo(name string, repopath string, config *Config) {
	repo, err := NewRepo(name, repopath)
	if err != nil {
		log.Println(err)
		return
	}

	contextUrl := "/git/" + name + "/"

	http.HandleFunc(contextUrl, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Header.Get("User-Agent"), "git") {
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

				}
			}
		} else {
			if p, ok := strings.CutPrefix(r.URL.Path, contextUrl); ok {
				err = h.singleTmpl.Execute(w, &HttpObject{
					Start:   time.Now(),
					Repos:   nil,
					Repo:    repo,
					Config:  config,
					SubPath: p,
					Path:    r.URL.Path,
				})
				if err != nil {
					log.Println(err)
					w.WriteHeader(500)
					return
				}
			}
		}
	})

}
