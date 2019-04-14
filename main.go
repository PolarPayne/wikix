package main

import (
	"flag"
	"fmt"

	"github.com/polarpayne/wikix/page"
	"github.com/polarpayne/wikix/render"
	"github.com/polarpayne/wikix/static"

	"github.com/spf13/afero"
)

var (
	varUnsafe       = flag.Bool("unsafe", false, "allow arbitary html in markdown")
	varPathTemplate = flag.String("template", "template", "path to template directory")
	varPathStatic   = flag.String("static", "static", "path to static directory")
	varPathPage     = flag.String("page", "page", "path to page directory")
	varRootPage     = flag.String("root-page", "Welcome", "page that is used for the root path (homepage)")
	varHost         = flag.String("host", "localhost", "host that the application is started on")
	varPort         = flag.Int("port", 8080, "port that the application is started on")
)

func main() {
	flag.Parse()

	af := func(p string) afero.Fs {
		return afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), p))
	}

	s := server{
		render.NewTemplate(af(*varPathTemplate)),
		page.NewFS(af(*varPathPage), *varUnsafe),
		static.NewFS(af(*varPathStatic)),
	}

	addr := fmt.Sprintf("%v:%v", *varHost, *varPort)
	s.Run(addr)
}
