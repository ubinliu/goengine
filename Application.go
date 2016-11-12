package engine

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

type Hooks map[string]RequestHandler

type Application struct{
	name string
	listen string
	assets map[string]string
	router *Router
	request *Request
	response *Response
	ctx *Context
	template *Template
	hooks Hooks
}

func NewApplication(name string, listen string) (*Application){
	app := &Application{}
	app.name = name
	app.listen = listen
	app.router = &Router{}
	app.template = &Template{}
	app.assets = make(map[string]string, 0)
	app.hooks = make(map[string]RequestHandler, 0)
	return app
}

func (app *Application) Run(){
	err := http.ListenAndServe(app.listen, app)
	if err != nil {
		fmt.Println(err)
	}
}

func (app *Application) Registe4xxHandler(handler RequestHandler){
	app.router.Registe4xxHandler(handler)
}

func (app *Application) Registe5xxHandler(handler RequestHandler){
	app.router.Registe5xxHandler(handler)
}

func (app *Application) InitRouterHandlers(routerHandlers RouterHandlerMap){
	app.router.InitRouterHandlers(routerHandlers)
}

func (app *Application) RegisteTemplateFuncMap(funcMap template.FuncMap){
	app.template.RegisteFuncMap(funcMap)
}

func (app *Application) ParseAllTemplates(tplDir string, suffix string){
	app.template.ParseAllTemplates(tplDir, suffix)
}

func (app *Application) RegisteAssetFiles(prefix string, assetsDir string){
	app.assets[prefix] = assetsDir
} 

func (app *Application) RegisteBeforeRequestHooks(handler RequestHandler){
	app.hooks["before"] = handler
}

func (app *Application) RegisteAfterRequestHooks(handler RequestHandler){
	app.hooks["after"] = handler
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request){
	//for assets route
	currentPath := r.URL.Path
	if len(app.assets) > 0 {
		for prefix, assetsDir := range(app.assets){
			if ok, _ := regexp.MatchString(prefix, currentPath); ok {
		        http.StripPrefix(prefix, http.FileServer(http.Dir(assetsDir))).ServeHTTP(w, r)
				return
		    }
		}
	}
	
	app.request = NewRequest(r)
	app.response = NewResponse(w)
	app.response.template = app.template
	app.ctx = NewContext(app.request, app.response)
	
	beforeHook, ok := app.hooks["before"]
	if ok {
		beforeHook(app.ctx)
	}
	
	app.router.Route(app.ctx)

	afterHook, ok := app.hooks["after"]
	if ok {
		afterHook(app.ctx)
	}
}