package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"

	"github.com/polarpayne/wikix/render"
	"github.com/polarpayne/wikix/types"

	"github.com/husobee/vestigo"
)

func getNormalizedName(r *http.Request) string {
	return vestigo.Param(r, "name")
	// return strings.ReplaceAll(page, "_", " ")
}

type server struct {
	*render.Template
	types.StorePage
	types.StoreStatic
}

func (s *server) Run(addr string) {
	r := vestigo.NewRouter()

	r.Get("/", s.handleIndex(*varRootPage))
	r.Get("/p", s.handlePageAll())
	r.Get("/p/:name", s.handlePage())
	r.Get("/t", s.handleTagAll())
	r.Get("/t/:name", s.handleTag())
	r.Get("/s", s.handleStaticAll())
	r.Get("/s/:name", s.handleStatic())

	log.Print("starting server on ", addr)
	log.Fatal(http.ListenAndServe(addr, middlewareLogging(middlewareTrailingSlash(r))))
}

// errPage renders the error.html page to `w`.
// `err` must be one of `string`, `[]byte`, `error`, or `fmt.Stringer`, panics otherwise.
// If `err` is a string, it may be followed by `args` that will then be given to `fmt.Sprintf` to format the string.
func (s *server) errPage(w http.ResponseWriter, code int, err interface{}, args ...interface{}) {
	var errString string

	switch e := err.(type) {
	case string:
		errString = fmt.Sprintf(e, args...)
	case []byte:
		errString = string(e)
	case error:
		errString = e.Error()
	case fmt.Stringer:
		errString = e.String()
	default:
		panic("err is of invalid type, it must be one of `string`, `[]byte`, `error`, `fmt.Stringer`")
	}

	out, er := s.Template.Render("error.html", types.Error{Error: errString, Code: code})
	if er != nil {
		w.WriteHeader(500)
		io.WriteString(w, er.Error())
		return
	}

	w.WriteHeader(code)
	w.Write(out)
}

func (s *server) handleIndex(rootPage string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request to get index page: %v", rootPage)

		if !s.PageExists(rootPage) {
			s.errPage(w, 404, "page '%v' does not exist", rootPage)
			return
		}

		p, err := s.Page(rootPage)
		if err != nil {
			s.errPage(w, 500, err)
			return
		}

		out, err := s.Render("page.html", p)
		if err != nil {
			s.errPage(w, 500, err)
			return
		}

		w.Write(out)
	}
}

func (s *server) handlePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := getNormalizedName(r)

		if !s.PageExists(name) {
			s.errPage(w, 404, "page '%v' does not exist", name)
			return
		}

		p, err := s.Page(name)
		if err != nil {
			s.errPage(w, 500, err)
			return
		}

		out, err := s.Render("page.html", p)
		if err != nil {
			s.errPage(w, 500, err)
			return
		}

		w.Write(out)
	}
}

func (s *server) handlePageAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages, err := s.PageAll()
		if err != nil {
			s.errPage(w, 500, err)
		}

		out, err := s.Render("page-all.html", pages)
		if err != nil {
			s.errPage(w, 500, err)
			return
		}

		w.Write(out)
	}
}

func (s *server) handleTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := getNormalizedName(r)

		tag, err := s.Tag(name)
		if err != nil {
			s.errPage(w, 500, err)
			return
		}

		out, err := s.Render("tag.html", tag)
		if err != nil {
			s.errPage(w, 500, err)
			return
		}

		w.Write(out)
	}
}

func (s *server) handleTagAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags, err := s.TagAll()
		if err != nil {
			s.errPage(w, 500, err)
			return
		}

		out, err := s.Render("tag-all.html", tags)
		if err != nil {
			s.errPage(w, 500, err)
			return
		}

		w.Write(out)
	}
}

func (s *server) handleStatic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := getNormalizedName(r)

		if !s.StaticExists(name) {
			s.errPage(w, 404, "static file '%v' does not exist", name)
			return
		}

		mediatype, params, err := mime.ParseMediaType(name)
		if err == nil {
			w.Header().Set("Content-Type", mime.FormatMediaType(mediatype, params))
		} else {
			log.Printf("failed to parse mimetype for %v: %v", name, err)
		}

		f, err := s.Static(name)
		if err != nil {
			s.errPage(w, 500, err)
			return
		}

		io.Copy(w, f)
		f.Close()
	}
}

func (s *server) handleStaticAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statics, err := s.StaticAll()
		if err != nil {
			s.errPage(w, 500, err)
			return
		}

		out, err := s.Render("static-all.html", statics)
		if err != nil {
			s.errPage(w, 500, err)
			return
		}

		w.Write(out)
	}
}
