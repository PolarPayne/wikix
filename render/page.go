package render

import (
	"regexp"
	"text/template"

	"github.com/valyala/bytebufferpool"

	"github.com/spf13/afero"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

var (
	regexLink  = regexp.MustCompile(`\[\[ *([^\]]+?) *\| *([^\]#]+?)(#[^\]#]+?)? *\]\]`)
	regexLink2 = regexp.MustCompile(`\[\[ *([^\]#\|]+?)(#[^\]#]+?)? *\]\]`)
)

func wikilinksToMarkdown(b []byte) []byte {
	out := regexLink.ReplaceAll(b, []byte(`[$1](/p/$2$3)`))
	return regexLink2.ReplaceAll(out, []byte(`[$1](/p/$1$2)`))
}

type Page struct {
	fs    afero.Fs
	funcs map[string]interface{}
}

func (p *Page) funcMacro(macro string, data interface{}) (string, error) {
	fp, err := afero.ReadFile(p.fs, macro+".md")
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(macro).Funcs(p.funcs).Parse(string(fp))
	if err != nil {
		return "", err
	}

	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	err = tmpl.ExecuteTemplate(b, macro, data)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func NewPage(fs afero.Fs) *Page {
	out := new(Page)
	out.fs = fs
	out.funcs = map[string]interface{}{
		"include": out.funcMacro,
		"dict":    dict,
		"list":    list,
	}

	return out
}

func (p *Page) Render(data []byte, unsafe bool) (string, error) {
	tmpl, err := template.New("main").Funcs(p.funcs).Parse(string(data))
	if err != nil {
		return "", err
	}

	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	err = tmpl.ExecuteTemplate(b, "main", nil)
	if err != nil {
		return "", err
	}

	markdownParser := parser.NewWithExtensions(parser.AutoHeadingIDs | parser.BackslashLineBreak | parser.Footnotes | parser.SuperSubscript)

	md := wikilinksToMarkdown(b.Bytes())
	html := markdown.ToHTML(md, markdownParser, nil)

	if !unsafe {
		return string(bluemonday.UGCPolicy().SanitizeBytes(html)), nil
	}
	return string(html), nil
}
