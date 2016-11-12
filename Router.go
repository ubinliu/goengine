package goengine

import (
	"fmt"
	"regexp"
)

type RequestHandler func(ctx *Context)

type RouterHandlerMap map[string]RequestHandler
type RouterRegMap map[string]*regexp.Regexp

type Router struct{
	RouterHandlers RouterHandlerMap
	RouterRegs RouterRegMap
	Handle4xx RequestHandler
	Handle5xx RequestHandler
}

func (router *Router) PreCompilePathReg(routerHandlers RouterHandlerMap){
	if len(router.RouterRegs) == 0 {
		router.RouterRegs = make(RouterRegMap, 0)
	}

	var err error
	for pathExpr, _ := range(routerHandlers){
		(router.RouterRegs)[pathExpr], err = regexp.Compile(pathExpr)
		if err != nil {
			panic(err)
		}
	}
} 

func (router *Router) InitRouterHandlers(routerHandlers RouterHandlerMap){
	if len(routerHandlers) == 0 {
		return
	}
	
	router.RouterHandlers = routerHandlers
	router.PreCompilePathReg(routerHandlers)
}

func (router *Router) Registe4xxHandler(handler RequestHandler){
	router.Handle4xx = handler
}

func (router *Router) Registe5xxHandler(handler RequestHandler){
	router.Handle5xx = handler
}

func (router *Router) Route(ctx *Context){
	path := ctx.Request.r.URL.Path
	//handle panic error
	defer func(){
        if err := recover(); err != nil{
            if router.Handle5xx != nil {
				router.Handle5xx(ctx)
			}else{
				ctx.Response.RenderText("<h1>500 server internal error</h1>")
				fmt.Println(err)
			}
        }
    }()
	//complete match first
	for p, handler := range(router.RouterHandlers) {
		if p == path {
			handler(ctx)
			return
		}
	}	
	//regexp complete match second
	for pathExpr, regexp := range(router.RouterRegs) {
		if pathExpr == "/" {
			continue
		}
		if regexp.MatchString(path) {
			router.RouterHandlers[pathExpr](ctx)
			return
		}
	}	
	//not found
	if router.Handle4xx != nil {
		router.Handle4xx(ctx)
		return
	}
	ctx.Response.RenderText("<h1>404 not found</h1>")
}

