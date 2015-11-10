package lines

import (
	"bufio"
	"io"
	"strings"
	"text/template"

	"github.com/gogap/spirit"
)

const (
	linesTranslatorOutURN = "urn:spirit:translator:out:lines"
)

var _ spirit.OutputTranslator = new(LinesOutputTranslator)

type TemplateDelims struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}

type LinesOutputTranslatorConfig struct {
	Template string         `json:"template"`
	Delims   TemplateDelims `json:"delims"`
}

type LinesOutputTranslator struct {
	tmpl *template.Template
	conf LinesOutputTranslatorConfig
}

func init() {
	spirit.RegisterOutputTranslator(linesTranslatorOutURN, NewLinesOutputTranslator)
}

func NewLinesOutputTranslator(options spirit.Options) (translator spirit.OutputTranslator, err error) {
	conf := LinesOutputTranslatorConfig{}

	if err = options.ToObject(&conf); err != nil {
		return
	}

	if conf.Template == "" {
		conf.Template = "{{getJSON .delivery.Payload.GetData}}\n"
	}

	var tmpl *template.Template
	tmpl = template.New(linesTranslatorOutURN).Option("missingkey=error")

	conf.Delims.Left = strings.TrimSpace(conf.Delims.Left)
	conf.Delims.Right = strings.TrimSpace(conf.Delims.Right)

	if conf.Delims.Left != "" && conf.Delims.Right != "" {
		tmpl = tmpl.Delims(conf.Delims.Left, conf.Delims.Right)
	}

	tmpl = tmpl.Funcs(funcMap)

	if tmpl, err = tmpl.Parse(conf.Template); err != nil {
		return
	}

	translator = &LinesOutputTranslator{
		conf: conf,
		tmpl: tmpl,
	}
	return
}

func (p *LinesOutputTranslator) Out(w io.WriteCloser, delivery spirit.Delivery) (err error) {
	newWriter := bufio.NewWriter(w)

	if err = p.tmpl.Execute(newWriter, map[string]interface{}{"delivery": delivery}); err != nil {
		return
	}

	err = newWriter.Flush()
	return
}