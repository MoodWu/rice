package router

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"strings"
	"regexp"
	"rice/lib/avl"
	"sort"
)

type (
	//MiddlewareFunc 中间件处理函数
	MiddlewareFunc func(http.ResponseWriter, *http.Request, ServeHTTPFunc)
	//ServeHTTPFunc http处理函数
	ServeHTTPFunc func(http.ResponseWriter, *http.Request)
	//Router 实例
	Router struct {
		basePath          string
		globalMiddlewares []MiddlewareFunc
		groupStack        []middlewareItem
		routeItems        map[string]routeItem
		dynRouteItems	  *avl.Node
		//static resource
		staticEnable  bool
		staticPath    string
		staticHandler http.Handler
	}
	middlewareItem struct {
		prefix      string
		middlewares []MiddlewareFunc
	}
	routeItem struct {
		methods  []string
		callback ServeHTTPFunc
	}
	routeItemEx struct{
		routeItem
		rawUri string //包含参数定义的URI
		segment int // URI中/ 的个数
		pattern string //正则
		prefix string //固定的URI前缀，可以为空
		suffix string //固定的URI后缀，可以为空
	}

	routeNode struct {
		avl.Node
	}

)

//New 创建新路由对象
func New() *Router {
	//创建路由对象,给路由表默认100长度
	return &Router{routeItems: make(map[string]routeItem, 100),dynRouteItems : nil,basePath: ""}
}

//ServeStatic 静态资源
func (r *Router) ServeStatic(staticPath, physicsPath string) {
	r.staticEnable = true
	r.staticPath = staticPath
	r.staticHandler = http.StripPrefix(staticPath, http.FileServer(http.Dir(physicsPath)))
}

//BasePath 设置view路由的基础路径[不影响静态资源]
func (r *Router) BasePath(basePath string) {
	r.basePath = basePath
}

//Use 添加全局中间件
func (r *Router) Use(middlewares ...MiddlewareFunc) error {
	if len(middlewares) == 0 {
		return errors.New("middlewares length can't not be zero")
	}
	for _, middleware := range middlewares {
		if middleware == nil {
			return errors.New("middleware can't not be nil")
		}
	}
	r.globalMiddlewares = append(r.globalMiddlewares, middlewares...)
	return nil
}

//Group 分组路由
func (r *Router) Group(prefix string, callback func(), middlewares ...MiddlewareFunc) error {
	if len(prefix) == 0 {
		return errors.New("prefix can't not be nil or empty")
	}
	for _, middleware := range middlewares {
		if middleware == nil {
			return errors.New("middleware can't not be nil")
		}
	}
	r.groupStack = append(r.groupStack, middlewareItem{prefix, middlewares}) //push item
	callback()
	r.groupStack = r.groupStack[:len(r.groupStack)-1] //pop item
	return nil
}

//Get 添加GET请求路由
func (r *Router) Get(uri string, handle ServeHTTPFunc, middlewares ...MiddlewareFunc) error {
	return r.AddRoute(uri, []string{http.MethodGet}, handle, middlewares...)
}

//Post 添加POST请求路由
func (r *Router) Post(uri string, handle ServeHTTPFunc, middlewares ...MiddlewareFunc) error {
	return r.AddRoute(uri, []string{http.MethodPost}, handle, middlewares...)
}

//Any 添加GET,POST请求路由
func (r *Router) Any(uri string, handle ServeHTTPFunc, middlewares ...MiddlewareFunc) error {
	return r.AddRoute(uri, []string{http.MethodGet, http.MethodPost}, handle, middlewares...)
}

//AddRoute 添加路由
func (r *Router) AddRoute(uri string, methods []string, handle ServeHTTPFunc, middlewares ...MiddlewareFunc) error {
	if len(uri) == 0 {
		return errors.New("prefix can't not be nil or empty")
	}
	if len(methods) == 0 {
		return errors.New("methods can't not be nil or empty")
	}
	if len(methods) > 2 {
		return errors.New("method support GET or POST")
	}
	for _, method := range methods {
		if method != http.MethodGet && method != http.MethodPost {
			return errors.New("method support GET or POST")
		}
	}
	if handle == nil {
		return errors.New("handle can't not be nil")
	}
	for _, middleware := range middlewares {
		if middleware == nil {
			return errors.New("middleware can't not be nil")
		}
	}
	for i := len(r.groupStack) - 1; i >= 0; i-- {
		groupItem := r.groupStack[i]
		uri = groupItem.prefix + uri
		middlewares = append(groupItem.middlewares, middlewares...)
	}
	if r.basePath != "" {
		uri = r.basePath + uri
	}

	middlewares = append(r.globalMiddlewares, middlewares...)
	callback := r.pipeline(handle, middlewares...)

	if strings.Contains(uri,":") {		
		prefix,suffix,pattern := processUri(uri)
		routeList  := make([]routeItemEx,0)
		node := routeItemEx{routeItem:routeItem{methods, callback},rawUri:uri,segment:strings.Count(uri,"/"),prefix:prefix,suffix:suffix,pattern:pattern }
		if r.dynRouteItems == nil {
			routeList = append(routeList,node)
			x := avl.InitTree(node.prefix,routeList,nil)			
			r.dynRouteItems = x			
		} else {			
			n := r.dynRouteItems.Search(node.prefix)
			if n == nil {
				routeList = append(routeList,node)	
				r.dynRouteItems = r.dynRouteItems.Put(node.prefix,routeList,nil)
			} else {
				//判断是否有想同Segment和suffix的路由项存在，有就报错
				for _,m := range n.Data.([]routeItemEx) {
					if m.segment == node.segment && m.suffix == node.suffix {
						panic("路由重复了："+ m.rawUri + "," + node.rawUri)
					}
				}
				n.Data = append(n.Data.([]routeItemEx),node)
				//对有想同前缀的路由项目进行倒排，确保最长路径 和 最长后缀的排在前面
				sort.Slice(n.Data,func(i,j int) bool {
					list := n.Data.([]routeItemEx)
					if list[i].segment == list[j].segment {
						if  strings.Count(list[i].suffix,"/") == strings.Count(list[j].suffix,"/") {
							return list[i].suffix > list[j].suffix
						} else {
							return strings.Count(list[i].suffix,"/") > strings.Count(list[j].suffix,"/")
						}
					} else {
						return list[i].segment > list[j].segment
					}
				})
			}			
		}
	} else {
		r.routeItems[uri] = routeItem{methods, callback}
	}
	return nil
}

func processUri(uri string )(prefix,suffix,pattern string) {		
	suffix = ""
	pattern = ""
	iPos := strings.Index(uri,":")
	prefix = uri[:iPos]
	c := strings.Count(uri,":")
	uri = uri[iPos+1:]
	
	iPos = strings.LastIndex(uri,":")
	if iPos >0 {
		uri = uri[iPos +1:]
	}
	
	iPos = strings.Index(uri,"/")
	if iPos > 0 {
		suffix = uri[iPos+1:]
	}
		
	if suffix == ""{
		pattern = fmt.Sprintf("%s%s{%d}[^/]+$", prefix , "([^/]+/)",c-1)			
	} else {
		pattern = fmt.Sprintf("%s%s{%d}%s$", prefix , "([^/]+/)",c,suffix)	
	}
	return 
}


//ServeHTTP http请求处理函数 实现http.Handler接口
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	uri := html.EscapeString(req.URL.Path)

	if r.staticEnable && strings.HasPrefix(uri, r.staticPath) {
		r.staticHandler.ServeHTTP(w, req)
		return
	}

	routeItem, exist := r.routeItems[uri]
	if exist {
		processRouteItem(routeItem,uri,method,w,req)
		return
	} else {
		// fmt.Println("Dyn Router,uri:",uri)
		// r.dynRouteItems.Print()
		//检查动态路由
		//从最长路劲开始匹配
		
		path := uri
		for {	
			// fmt.Println("Search Path:",path)
			path = strings.TrimSuffix(path,"/")
			path = strings.TrimSuffix(path,GetSuffix(path,"/"))
			//已经匹配完了，还没有找到合适的
			if path == "/" {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "uri %s not found", uri)
				return
			}			 

			dynItem := r.dynRouteItems.Search(path)
			if dynItem != nil {
				//遍历所有路由项，找到参数想同
				routeList := dynItem.Data.([]routeItemEx)
				for _,item := range routeList{
					if strings.Count(uri,"/") == item.segment && strings.HasSuffix(uri,item.suffix) {
						processDynRouteItem(item,uri,method,w,req)
						return
					}
				}
			}
		}
	}
}

func processRouteItem(item routeItem,uri,method string,w http.ResponseWriter, req *http.Request) {
	supportMethod := false
	for _, m := range item.methods {
		if m == method {
			supportMethod = true
			break
		}
	}
	if !supportMethod {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "uri %s with method %s not support", uri, method)
		return
	}
	item.callback(w, req)
}

func processDynRouteItem(item routeItemEx,uri,method string,w http.ResponseWriter, req *http.Request) {
	supportMethod := false
	for _, m := range item.methods {
		if m == method {
			supportMethod = true
			break
		}
	}
	if !supportMethod {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "uri %s with method %s not support", uri, method)
		return
	}
	//处理动态参数，并根据请求方法压入栈中
	rawList := strings.Split(item.rawUri,"/")[strings.Count(item.prefix,"/"):]
	detailList := strings.Split(uri,"/")[strings.Count(item.prefix,"/"):]
	for k,m  := range rawList{
		if strings.HasPrefix(m,":") {
			req.ParseForm()
			switch strings.ToUpper(method) {
			case "GET":				
				req.Form.Set(strings.TrimPrefix(m,":"),detailList[k])
			case "POST":
				req.PostForm.Set(strings.TrimPrefix(m,":"),detailList[k])
			}			
		}
	}
	fmt.Println("Call")
	item.callback(w, req)
}

func (r *Router) pipeline(handle ServeHTTPFunc, middlewares ...MiddlewareFunc) ServeHTTPFunc {
	if len(middlewares) <= 0 {
		return handle
	}
	middleware := middlewares[len(middlewares)-1]
	callback := func(w http.ResponseWriter, r *http.Request) {
		middleware(w, r, handle)
	}
	callback = r.pipeline(callback, middlewares[:len(middlewares)-1]...)

	return callback
}

func (node *routeNode)Search(uri string) *routeNode {	
	if node == nil {
		return node
	}

	// fmt.Println("Search：",uri,",Pattern:",node.Data.(routeItemEx).pattern)
	match,_ := regexp.MatchString(node.Data.(routeItemEx).pattern, uri)
	//如果符合正则，并且参数段数相同则为，并且后缀想同
	if match && strings.Count(node.Data.(routeItemEx).rawUri,"/") == strings.Count(uri,"/") && (node.Data.(routeItemEx).suffix == "" || GetSuffix(uri,"/") == node.Data.(routeItemEx).suffix) {		
		//当前节点匹配，为了实现最长的前缀匹配，要看下个节点是否匹配，直到下个节点不匹配了，才返回
		fmt.Println("Pending,",node)
		if node.LNode == nil || node.LNode.Search(uri) == nil {
			fmt.Println("Found:",node.Data.(routeItemEx).rawUri)
			return node
		} else {
			n := &routeNode{*node.LNode}
			return n.Search(uri)			
		}
	} else if uri < node.Data.(routeItemEx).prefix {
		if node.LNode != nil {
			n := &routeNode{*node.LNode}
			return n.Search(uri)			
		} 
	} else {
		if node.RNode != nil {
		n := &routeNode{*node.RNode}		
		return  n.Search(uri)			
		} 
	}

	return nil
}

func GetSuffix(uri,sep string) string {
	iPos := strings.LastIndex(uri,sep)
	
	if iPos >=0 {
		// fmt.Println("GetSuffix:",uri[iPos+1:])
		return uri[iPos+1:]
	} else {
		// fmt.Println("GetSuffix:",uri)
		return uri
	}
}