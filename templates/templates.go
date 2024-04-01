package templates

import (
	"errors"
	"html/template"
	"log"
)

// type TemplateData map[string]interface{}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid template dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("template dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

var Tmpl = template.New("master").Funcs(template.FuncMap{"dict": dict})

func init() {
	_, err := Tmpl.ParseGlob("templates/**/*.html")
	if err != nil {
		log.Fatalln("Error loading templates:", err)
	}

	// Print the templates that were loaded
	log.Println("Templates loaded:")
	for _, t := range Tmpl.Templates() {
		log.Println(" -", t.Name())
	}
}
