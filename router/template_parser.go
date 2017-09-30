/*
	view template parser
*/
package router

import (
	"html/template"
	"webserver/config"

	"github.com/cihub/seelog"
)

var tNilTemplate *template.Template

func InitTemplate() error {
	t, err := template.New("nil").Parse("template parse err")
	if err != nil {
		return err
	}
	tNilTemplate = t
	return nil
}

//template parse for router
func templateParse(sTemplateName ...string) *template.Template {
	var tTemplate *template.Template
	var err error
	if len(sTemplateName) == 0 {
		tTemplate, err = template.ParseFiles(config.G.Template.Viewroot + config.G.Template.Default)
	} else {
		//add deafult header_start and header_end for any template
		files := []string{}
		for _, f := range sTemplateName {
			files = append(files, config.G.Template.Viewroot+f)
		}
		for _, c := range config.G.Template.Components {
			if len(c) > 0 {
				files = append(files, config.G.Template.Viewroot+c)
			}
		}
		tTemplate, err = template.ParseFiles(files...)
	}
	if err != nil {
		seelog.Errorf("parse template:%v err:%v", sTemplateName, err)
		return tNilTemplate
	}
	return tTemplate
}
