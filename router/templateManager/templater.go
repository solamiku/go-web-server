package templater

import (
	"html/template"
	"regexp"
	"sync"
)

const (
	TEMPLATE_NIL_123 = "TEMPLATE_NIL_123"
)

type veiwTemplate struct {
	files     []string
	pTemplate *template.Template
}

var templates map[string]*veiwTemplate
var lock *sync.RWMutex
var tagRegxep *regexp.Regexp
var debugFlag bool = true

func init() {
	templates = make(map[string]*veiwTemplate)
	lock = new(sync.RWMutex)
	nilT, _ := template.New(TEMPLATE_NIL_123).Parse("template nil")
	templates[TEMPLATE_NIL_123] = &veiwTemplate{
		files:     []string{},
		pTemplate: nilT,
	}
	tagRegxep, _ = regexp.Compile(`\w*\.html`)
}

func ParseTemplate(name string, coms []string) error {
	files := append([]string{name}, coms...)
	pTemplate, err := template.ParseFiles(files...)
	if err != nil {
		return err
	}
	lock.Lock()
	defer lock.Unlock()
	templates[tagRegxep.FindString(name)] = &veiwTemplate{
		files:     files,
		pTemplate: pTemplate,
	}
	return nil
}

func GetTemplate(sName string) *template.Template {
	lock.RLock()
	defer lock.RUnlock()
	if pvt, ok := templates[sName]; ok {
		if debugFlag && len(pvt.files) > 0 {
			if nt, err := template.ParseFiles(pvt.files...); err != nil {
				pvt.pTemplate, _ = template.New(sName).Parse(err.Error())
			} else {
				pvt.pTemplate = nt
			}
		}
		return templates[sName].pTemplate
	}
	return templates[TEMPLATE_NIL_123].pTemplate
}

// set debug flag.
// if set false, you need to reload the template content manual  in times of need.
// arguments - bool
func Debug(f bool) {
	debugFlag = f
}
