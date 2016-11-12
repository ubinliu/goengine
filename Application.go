package goengine

import (
	"fmt"
	"strings"
	"html/template"
	"net"
	"net/http"
	"regexp"
	"bytes"
    "encoding/base64"
)

type Hooks map[string]RequestHandler

type FuncMap template.FuncMap

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
	AuthName string
	AuthPwd string
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
	
	ln, err := net.Listen("tcp4", app.listen)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(http.Serve(ln, app))
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

func (app *Application) RegisteTemplateFuncMap(funcMap FuncMap){
	app.template.RegisteFuncMap(template.FuncMap(funcMap))
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

func (app *Application) UseBasicAuth(username string, password string){
	app.AuthName = username
	app.AuthPwd = password
}

func (app *Application)basicAuth() bool{
    basicAuthPrefix := "Basic "
	
    // 获取 request header
    auth := app.ctx.Request.GetHeader("Authorization")
    // 如果是 http basic auth
    if strings.HasPrefix(auth, basicAuthPrefix) {
        // 解码认证信息
        payload, err := base64.StdEncoding.DecodeString(
            auth[len(basicAuthPrefix):],
        )
        if err == nil {
            pair := bytes.SplitN(payload, []byte(":"), 2)
            if len(pair) == 2 && bytes.Equal(pair[0], []byte(app.AuthName)) &&
                bytes.Equal(pair[1], []byte(app.AuthPwd)) {
                	return true
            }
        }
    }

    // 认证失败，提示 401 Unauthorized
    // Restricted 可以改成其他的值
    app.ctx.Response.SetHeader("WWW-Authenticate", `Basic realm="Restricted"`)
    app.ctx.Response.SetHttpStatus(401)
    app.ctx.Response.RenderText("401 Authorization Required")
    return false
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
	
	if app.AuthName != "" {
		if ! app.basicAuth() {
			fmt.Println("unauthed")
			return
		}
	}
	
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