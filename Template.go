package engine

import (
	"html/template"
	"path/filepath"
	"os"
	"strings"
	"fmt"
)

type Template struct{
	AllTemplates *template.Template
	FuncMap template.FuncMap
}

func (t *Template) RegisteFuncMap(funcMap template.FuncMap){
	t.FuncMap = funcMap
}

func (t *Template)ParseAllTemplates(viewPath string, suffix string){
	t.AllTemplates = template.New("AllTemplates")

	if len(t.FuncMap) != 0 {
		t.AllTemplates.Funcs(t.FuncMap)
	}

	suffix = strings.ToUpper(suffix)
	pathLen := len(viewPath)
	err := filepath.Walk(viewPath, func(path string, f os.FileInfo, err error) error {
        if (f == nil) {
			return err
		}
        if f.IsDir() {
			return nil
		}
		
		if strings.HasSuffix(strings.ToUpper(f.Name()), suffix) {
			name := path[pathLen:len(path)]
			
			t.AllTemplates, err = t.AllTemplates.ParseFiles(path)
			if err != nil {
				fmt.Println("template parse files failed", err.Error())
				return err
			}else{
				fmt.Println("template parse file", name, path)
			}
			
		}
        return nil
    })
	
    if err != nil {
		fmt.Printf("ParseAllTemplates returned %v\n", err)
		os.Exit(1)
    }
} 
