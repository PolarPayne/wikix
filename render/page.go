package render

import (
	"html/template"
	"regexp"

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

func Page(data []byte, unsafe bool) template.HTML {
	markdownParser := parser.NewWithExtensions(parser.AutoHeadingIDs | parser.BackslashLineBreak | parser.Footnotes | parser.SuperSubscript)

	md := wikilinksToMarkdown(data)
	html := markdown.ToHTML(md, markdownParser, nil)

	if !unsafe {
		return template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))
	}
	return template.HTML(html)
}
