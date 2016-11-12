package engine

type Context struct{
	Request *Request
	Response *Response
	Template *Template
	dict map[string]string
}

func NewContext(request *Request, response *Response) (*Context){
	context := &Context{}
	context.dict = make(map[string]string, 10)
	context.Request = request
	context.Response = response 
	return context
}

func (ctx *Context)Set(key string, value string){
	ctx.dict[key] = value
}

func (ctx *Context)Get(key string) string{	
	value, ok := ctx.dict[key]
	if ok {
		return value
	}
	return ""
}

