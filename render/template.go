package render

import (
	"log"
	"text/template"

	"github.com/valyala/bytebufferpool"

	"github.com/spf13/afero"
)

type Template struct {
	fs    afero.Fs
	funcs map[string]interface{}
}

func (t *Template) funcInclude(templateName string, data interface{}) string {
	fp, err := afero.ReadFile(t.fs, templateName)
	if err != nil {
		log.Printf("failed to read template %v: %v", templateName, err)
		return ""
	}

	tmpl, err := template.New(templateName).Funcs(t.funcs).Parse(string(fp))
	if err != nil {
		log.Printf("failed to parse template %v: %v", templateName, err)
		return ""
	}

	out := bytebufferpool.Get()
	defer bytebufferpool.Put(out)

	err = tmpl.ExecuteTemplate(out, templateName, data)
	if err != nil {
		log.Printf("failed to execute template %v: %v", templateName, err)
		return ""
	}

	return out.String()
}

func NewTemplate(fs afero.Fs) *Template {
	out := new(Template)
	out.fs = fs
	out.funcs = map[string]interface{}{
		"include": out.funcInclude,
		"dict":    dict,
		"list":    list,
	}
	return out
}

func (t *Template) Render(name string, data interface{}) ([]byte, error) {
	fp, err := afero.ReadFile(t.fs, "base.html")
	if err != nil {
		return nil, err
	}

	baseTmpl, err := template.New("base.html").Funcs(t.funcs).Parse(string(fp))
	if err != nil {
		return nil, err
	}

	// --- NAME ---

	fp, err = afero.ReadFile(t.fs, name)
	if err != nil {
		return nil, err
	}

	tmpl, err := baseTmpl.New(name).Parse(string(fp))
	if err != nil {
		return nil, err
	}

	// --- EXECUTE ---

	out := bytebufferpool.Get()
	defer bytebufferpool.Put(out)

	err = tmpl.ExecuteTemplate(out, "base.html", data)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
