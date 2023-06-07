package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RouteInfo struct {
	Package    *PackageInfo
	RefPackage map[string]*PackageInfo
	ExpPackage map[string]*PackageInfo
	Dir        string
	Group      string
	Methods    []*Method //有可能修改方法的信息，所以要用指针
}

type PackageInfo struct {
	Name  string
	Path  string
	Alias string
}

type Method struct {
	Name        string
	Recv        *Param
	Params      []*Param
	AliasName   string
	HttpMethods []string
	Results    []*Param
	Comments    []string
}

type Param struct {
	Name      string
	Selector  string
	ParamType string
	Pointer   bool //是否是指针
}

//读取源文件，解析代码结构
func ReadFile(base, fileName string) *RouteInfo {
	// Create the AST by parsing src.
	fset := token.NewFileSet() 
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	//  ast.Print(fset, f)
	var r RouteInfo
	r.Package = &PackageInfo{
		Name:  f.Name.Name,
		Path:  strings.TrimSuffix(base,"/"),
		Alias: f.Name.Name}
	r.RefPackage = make(map[string]*PackageInfo, 0)
	r.ExpPackage = make(map[string]*PackageInfo, 0)
	if f.Doc != nil {
		for _, doc := range f.Doc.List {
			test := strings.ToLower(doc.Text)
			if strings.HasPrefix(test, "//@methodgroup:") {
				r.Group = strings.TrimSpace(doc.Text[15:])
			}
		}
	}

	r.Methods = make([]*Method, 0)
	for _, n := range f.Decls {
		switch n.(type) {
		case *ast.GenDecl:
			x := n.(*ast.GenDecl)
			//处理文件的import信息
			if x.Tok == token.Lookup("import") {
				for _, v := range x.Specs {
					is := v.(*ast.ImportSpec)
					key := ""
					path := strings.Trim(is.Path.Value, "\"")
					if is.Name != nil {
						key = is.Name.Name
					} else {
						key = path[strings.LastIndex(path, "/")+1:]
					}
					r.RefPackage[key] = &PackageInfo{
						Name:  key,
						Path:  path,
						Alias: key,
					}
				}
			}

		case *ast.FuncDecl:
			x := n.(*ast.FuncDecl)
			//处理基于注释，没有注释的函数不用处理
			if x.Doc == nil {
				continue
			}
			method := Method{Name: x.Name.Name, AliasName: x.Name.Name}
			method.Params = make([]*Param, 0)
			method.Results = make([]*Param, 0)
			method.Comments = make([]string, 0)

			//逐行处理函数的注释，根据注释的前缀进行提取相应的信息,以后可以扩充更多的Tag
			for _, m := range x.Doc.List {
				test := strings.ToLower(m.Text)
				if strings.HasPrefix(test, "//@method:") {
					methods := strings.Replace(test, "//@method:", "", -1)
					method.HttpMethods = RemoveDuplicateElement(strings.Split(methods, ","))
				} else if strings.HasPrefix(test, "//@methodname:") {
					method.AliasName = strings.TrimSpace(m.Text[14:])
				} else {
					//这些是不需要处理的注释，放回函数的注释中
					method.Comments = append(method.Comments, m.Text)
				}
			}

			//只对标明了Http方法的函数进行处理
			if len(method.HttpMethods) > 0 {
				//如果函数有接收者
				if x.Recv != nil {
					p := Param{}
					reciver := x.Recv.List[0]
					p.Name = reciver.Names[0].Name
					switch reciver.Type.(type) {
					//引用类型，非基础类型的接收者
					case *ast.StarExpr:
						y := reciver.Type.(*ast.StarExpr).X.(*ast.SelectorExpr)
						p.ParamType = interface{}(y.Sel).(*ast.Ident).Name
						p.Selector = interface{}(y.X).(*ast.Ident).Name
						p.Pointer = true
					//非基础类型的接收者
					case *ast.SelectorExpr:
						y := interface{}(x.Recv.List[0].Type).(*ast.SelectorExpr)
						p.Selector = interface{}(y.X).(*ast.Ident).Name
						p.ParamType = interface{}(y.Sel).(*ast.Ident).Name
						p.Pointer = false
					//基础类型的接收者
					case *ast.Ident:
						p.ParamType = reciver.Type.(*ast.Ident).Name
						p.Pointer = false

					}
					method.Recv = &p
				}

				for _, p := range x.Type.Params.List {
					// fmt.Println(fmt.Sprintf("%T",p.Type),p)					
					switch p.Type.(type) {
					//基础类型	
					case *ast.Ident:
						for _,aParam := range p.Names {
							param := Param{Name: aParam.Name, Selector: "", ParamType: p.Type.(*ast.Ident).Name, Pointer: false}
							method.Params = append(method.Params, &param)
						}
					//非基础类型
					case *ast.SelectorExpr:
						rp := p.Type.(*ast.SelectorExpr)
						for _,aParam := range p.Names {
							param := Param{Name: aParam.Name, Selector: rp.X.(*ast.Ident).Name, ParamType: rp.Sel.Name, Pointer: false}
							method.Params = append(method.Params, &param)
						}
						//非基础类型，将定义此类型的包加入要引用的包列表中
						r.ExpPackage[rp.X.(*ast.Ident).Name] = &PackageInfo{
							Name:  rp.X.(*ast.Ident).Name,
							Path:  r.RefPackage[rp.X.(*ast.Ident).Name].Path,
							Alias: rp.X.(*ast.Ident).Name,
						}
					//指针类型
					case *ast.StarExpr:
						switch p.Type.(*ast.StarExpr).X.(type) {
						//基础类型的指针
						case *ast.Ident:
							for _,aParam := range p.Names {
								param := Param{Name: aParam.Name, Selector: "", ParamType: p.Type.(*ast.StarExpr).X.(*ast.Ident).Name, Pointer: true}							
								method.Params = append(method.Params, &param)
							}
						//非基础类型的指针
						case *ast.SelectorExpr:
							rp := p.Type.(*ast.StarExpr).X.(*ast.SelectorExpr)
							for _,aParam := range p.Names {
								param := Param{Name: aParam.Name}
								param.ParamType = interface{}(rp.Sel).(*ast.Ident).Name
								param.Selector = interface{}(rp.X).(*ast.Ident).Name
								param.Pointer = true	
								method.Params = append(method.Params, &param)
							}
							//非基础类型，将定义此类型的包加入要引用的包列表中
							r.ExpPackage[rp.X.(*ast.Ident).Name] = &PackageInfo{
								Name:  rp.X.(*ast.Ident).Name,
								Path:  r.RefPackage[rp.X.(*ast.Ident).Name].Path,
								Alias: rp.X.(*ast.Ident).Name,
							}
	
						}
					}

				}
				if x.Type.Results !=nil {
					for _, p := range x.Type.Results.List {						
						//返回值未必有名称
						//fmt.Println(fmt.Sprintf("%T",p.Type),p)					
						switch p.Type.(type) {
						//基础类型	
						case *ast.Ident:
							param := Param{ Selector: "", ParamType: p.Type.(*ast.Ident).Name, Pointer: false}
							if len(p.Names)>0 {
								param.Name = p.Names[0].Name
							}							
							method.Results = append(method.Results, &param)
						//非基础类型
						case *ast.SelectorExpr:
							rp := p.Type.(*ast.SelectorExpr)
							param := Param{Selector: rp.X.(*ast.Ident).Name, ParamType: rp.Sel.Name, Pointer: false}
							if len(p.Names)>0 {
								param.Name = p.Names[0].Name
							}
							method.Results = append(method.Results, &param)
							//非基础类型，将定义此类型的包加入要引用的包列表中
							// r.ExpPackage[rp.X.(*ast.Ident).Name] = &PackageInfo{
							// 	Name:  rp.X.(*ast.Ident).Name,
							// 	Path:  r.RefPackage[rp.X.(*ast.Ident).Name].Path,
							// 	Alias: rp.X.(*ast.Ident).Name,
							// }
						//指针类型
						case *ast.StarExpr:
							switch p.Type.(*ast.StarExpr).X.(type) {
							//基础类型的指针
							case *ast.Ident:								
								param := Param{Selector: "", ParamType: p.Type.(*ast.StarExpr).X.(*ast.Ident).Name, Pointer: true}
								if len(p.Names)>0 {
									param.Name = p.Names[0].Name
								}								
								method.Results = append(method.Results, &param)								
							//非基础类型的指针
							case *ast.SelectorExpr:
								rp := p.Type.(*ast.StarExpr).X.(*ast.SelectorExpr)
								param := Param{}
								if len(p.Names)>0 {
									param.Name = p.Names[0].Name
								}
								param.ParamType = interface{}(rp.Sel).(*ast.Ident).Name
								param.Selector = interface{}(rp.X).(*ast.Ident).Name
								param.Pointer = true								
								method.Results = append(method.Results, &param)
								//非基础类型，将定义此类型的包加入要引用的包列表中
								// r.ExpPackage[rp.X.(*ast.Ident).Name] = &PackageInfo{
								// 	Name:  rp.X.(*ast.Ident).Name,
								// 	Path:  r.RefPackage[rp.X.(*ast.Ident).Name].Path,
								// 	Alias: rp.X.(*ast.Ident).Name,
								// }
							}
						}

					}
				}
				r.Methods = append(r.Methods, &method)

			}

		}
	}
	//如果这个文件有解析出来的方法，则返回对应实体
	if len(r.Methods) > 0 {
		return &r
	} else {
		return nil
	}
}

//生成代码
func GenCode(r []*RouteInfo) {	
	//处理包的别名问题（相同包不同别名，相同别名不同包）
	ProcessPackageNameAlias(r)
	
	//生成Register文件
	GenRegister(r)	

	//生成Facade文件,目前没有处理参数为指针的情况，做为Facade函数，个人觉得参数不应该是指针
	GenFacade(r)

}

//处理包的别名
func ProcessPackageNameAlias(r []*RouteInfo) {
		
	for _, ri := range r {
		//梳理所有需要引用类型的包
		for key, pack := range ri.ExpPackage {
			bExists := false
			//先判断这个包是否已经在引用列表中了，应对不同别名对应相同包的场景
			for _, m := range allRefPackage {
				if pack.Path == m.Path {
					bExists = true
					break
				}
			}
			if bExists {	
				continue
			}
			old, ok := allRefPackage[key]
			if ok {
				//有相同别名的引用，看看引用的路径是否不同，不同就重新给个别名，相同略过
				if old.Path != pack.Path {
					//fmt.Println("Same key,", old.Path, ",", pack.Path)
					newKey := GetNewKey(key, allRefPackage)
					pack.Alias = newKey
					allRefPackage[newKey] = *pack
				}
			} else {
				allRefPackage[key] = *pack
			}
		}

		//梳理所有需要调用实际方法的包
		old, ok := allSrvPackage[ri.Package.Name]
		if ok {
			if old.Path != ri.Package.Path {
				newKey := GetNewKey(ri.Package.Name, allSrvPackage)
				ri.Package.Alias = newKey
				allSrvPackage[newKey] = *ri.Package
			}
		} else {
			allSrvPackage[ri.Package.Name] = *ri.Package
		}
	}

	//根据最终的引用列表，反向更新所有的包中发方法的参数类型别名，以及整个包的别名
	for _, ri := range r {
		for _, md := range ri.Methods {
			//暂时不处理有接收者的函数
			// if md.Recv != nil {
			// 	//md.Recv.Selector = GetKeyByContent(ri.ExpPackage[md.Recv.Selector].Path, allRefPackage)
			// }

			for _, pr := range md.Params {
				if pr.Selector != "" {
					//找到新的别名					
					pr.Selector = GetKeyByContent(ri.ExpPackage[pr.Selector].Path, allRefPackage)					
				}
			}
		}
		ri.Package.Alias = GetKeyByContent(ri.Package.Path, allSrvPackage)

	}
}

//生成Http方法的注册文件
func GenRegister(r []*RouteInfo) {
	var sb bytes.Buffer
	var f_router string
	
	if tplDir != "" {
		tr, _ := ioutil.ReadFile(tplDir + "/template_register.txt")
		template_register = string(tr)
	}

	f_router = strings.Replace(template_register, "{{Date}}", time.Now().Format("2006-01-02 15:04:05"), -1)
	f_router = strings.Replace(f_router, "{{ProjectRoot}}", strings.TrimSuffix(basePath,"/"), -1)
	f_router = strings.Replace(f_router, "{{RefPackage}}", sb.String(), -1)

	//开始生成路由表
	for _, ri := range r {
		if ri.Group != "" {
			sb.WriteString("r.Group(\"/" + ri.Group + "\",func(){\n")
		}

		for _, m := range ri.Methods {
			for _, v := range m.HttpMethods {
				for _, c := range m.Comments {
					sb.WriteString(c + "\n")
				}
				sb.WriteString(fmt.Sprintf("r.%s(\"/%s\",%s)\n", FormatMethodName(v), m.AliasName, ri.Package.Alias + m.Name + "Facade" + FormatMethodName(v) ))
			}
		}

		if ri.Group != "" {
			sb.WriteString("})\n")
		}
		sb.WriteString("\n")

	}

	f_router = strings.Replace(f_router, "{{CodeGenFunc}}", sb.String(), -1)

	str, err := format.Source([]byte(f_router))
	if err != nil {

		fmt.Println("Error:", err)
		fmt.Println(f_router)
	} else {
		WriteFile("register.go",str)	
	}
}

func GenFacade(r []*RouteInfo) {
	var sb bytes.Buffer
	var f_facade string

	if tplDir != "" {
	tf, _ := ioutil.ReadFile(tplDir + "/template_facade.txt")
	template_facade = string(tf)
	}

	//参数类型的类型定义包
	for k, v := range allRefPackage {
		//facade 应该只引用service下的包，但因为使用了Context参数，可能会把view包包含进来,其实这个包需要
		// if v.Path == basePath + "model/view" {
		// 	continue
		// }
		sb.WriteString(fmt.Sprintf("%s \"%s\"\n", k, v.Path))		
	}

	//需要调用函数的包
	for k, v := range allSrvPackage {
		sb.WriteString(fmt.Sprintf("%s \"%s\"\n", k, v.Path))
	}

	f_facade = strings.Replace(template_facade, "{{Date}}", time.Now().Format("2006-01-02 15:04:05"), -1)
	f_facade = strings.Replace(f_facade, "{{ProjectRoot}}", strings.TrimSuffix(basePath,"/"), -1)
	f_facade = strings.Replace(f_facade, "{{RefPackage}}", sb.String(), -1)

	sb.Reset()

	//生成函数体
	for _, ri := range r {
		for _, m := range ri.Methods {
			for _, c := range m.Comments {
				sb.WriteString(c + "\n")
			}

			for _,verb := range m.HttpMethods {
				sb.WriteString("func " + ri.Package.Alias + m.Name + "Facade"+ FormatMethodName(verb) + "(w http.ResponseWriter,r *http.Request) {\n")
				sb.WriteString(fmt.Sprintf(" _ticker_id,_ticker := StartTimer(\"%s\")\n",ri.Package.Alias + m.Name + "Facade"+ FormatMethodName(verb)))
				//根据函数的参数类型，从请求中获得对应的参数
				for _, v := range m.Params {
					//fmt.Println(m.Name, v)
					//内置的基础类型
					if v.Selector == "" {
						switch v.ParamType {
						case "bool":
							sb.WriteString(fmt.Sprintf("%s,err := %sBoolValue(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
							sb.WriteString("if err!= nil {\n w.WriteHeader(500)\n w.Write([]byte(err.Error()))\n return \n}\n")
						case "byte":
							sb.WriteString(fmt.Sprintf("%s,err := %sByteValue(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
							sb.WriteString("if err!= nil {\n w.WriteHeader(500)\n w.Write([]byte(err.Error()))\n return \n}\n")
						case "int8":
							sb.WriteString(fmt.Sprintf("%s,err := %sInt8Value(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
							sb.WriteString("if err!= nil {\n w.WriteHeader(500)\n w.Write([]byte(err.Error()))\n return \n}\n")
						case "int16":
							sb.WriteString(fmt.Sprintf("%s,err := %sInt16Value(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
							sb.WriteString("if err!= nil {\n w.WriteHeader(500)\n w.Write([]byte(err.Error()))\n return \n}\n")
						case "int32":
							sb.WriteString(fmt.Sprintf("%s,err := %sInt32Value(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
							sb.WriteString("if err!= nil {\n w.WriteHeader(500)\n w.Write([]byte(err.Error()))\n return \n}\n")
						case "int64":
							sb.WriteString(fmt.Sprintf("%s,err := %sInt64Value(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
							sb.WriteString("if err!= nil {\n w.WriteHeader(500)\n w.Write([]byte(err.Error()))\n return \n}\n")
						case "int":
							sb.WriteString(fmt.Sprintf("%s,err := %sIntValue(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
							sb.WriteString("if err!= nil {\n w.WriteHeader(500)\n w.Write([]byte(err.Error()))\n return \n}\n")
						case "unit":
							sb.WriteString(fmt.Sprintf("%s,err := %sUintValue(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
							sb.WriteString("if err!= nil {\n w.WriteHeader(500)\n w.Write([]byte(err.Error()))\n return \n}\n")
						case "uintptr":
							sb.WriteString(fmt.Sprintf("%s,err := %sUintptrValue(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
							sb.WriteString("if err!= nil {\n w.WriteHeader(500)\n w.Write([]byte(err.Error()))\n return \n}\n")
						case "float32":
							sb.WriteString(fmt.Sprintf("%s,err := %sFloat32Value(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
							sb.WriteString("if err!= nil {\n w.WriteHeader(500)\n w.Write([]byte(err.Error()))\n return \n}\n")
						case "float64":
							sb.WriteString(fmt.Sprintf("%s,err := %sFloat64Value(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
							sb.WriteString("if err!= nil {\n w.WriteHeader(500)\n w.Write([]byte(err.Error()))\n return \n}\n")
						case "complex64":
							sb.WriteString(fmt.Sprintf("%s,err := %sComplex64Value(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
							sb.WriteString("if err!= nil {\n w.WriteHeader(500)\n w.Write([]byte(err.Error()))\n return \n}\n")
						case "complex128":
							sb.WriteString(fmt.Sprintf("%s,err := %sComplex128Value(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
							sb.WriteString("if err!= nil {\n w.WriteHeader(500)\n w.Write([]byte(err.Error()))\n return \n}\n")
						case "string":
							sb.WriteString(fmt.Sprintf("%s := %sStringValue(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
						}
					} else if v.Selector == "time" {
						if v.Pointer {
							sb.WriteString(fmt.Sprintf("var %s *time.Time\n", v.Name))
							sb.WriteString(fmt.Sprintf("%s,err = &%sTimeValue(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
							sb.WriteString("if err!= nil {\n w.WriteHeader(500)\n w.Write([]byte(err.Error()))\n return \n}\n")
						} else {
							sb.WriteString(fmt.Sprintf("var %s time.Time\n", v.Name))
							sb.WriteString(fmt.Sprintf("%s,err = %sTimeValue(r,\"%s\")\n", v.Name, FormatMethodName(verb), v.Name))
							sb.WriteString("if err!= nil {\n w.WriteHeader(500)\n w.Write([]byte(err.Error()))\n return \n}\n")
						}						
						
					} else if v.Selector != "" && v.ParamType == "Context" {
						if v.Pointer {
							sb.WriteString(fmt.Sprintf("var %s *%s.Context\n %s = CreateContext(w,r)\n", v.Name,v.Selector,v.Name))
						} else {
							sb.WriteString(fmt.Sprintf("var %s %s.Context\n %s = *CreateContext(w,r)\n", v.Name,v.Selector,v.Name))
						}						
					} else {						
						sb.WriteString(fmt.Sprintf("tmp_%s := %s.%s{}\n", v.Name, v.Selector, v.ParamType))
						sb.WriteString(fmt.Sprintf(" if %sObjectValue(r,\"%s\",&tmp_%s) != nil {\n w.WriteHeader(500)\n w.Write([]byte(%sObjectValue(r,\"%s\",&tmp_%s).Error()))\n return\n }\n", FormatMethodName(verb), v.Name,v.Name, FormatMethodName(verb), v.Name,v.Name))						
						if v.Pointer {
							sb.WriteString(fmt.Sprintf(" %s := &tmp_%s\n", v.Name,v.Name))
						} else {
							sb.WriteString(fmt.Sprintf(" %s := tmp_%s\n", v.Name,v.Name))
						}						
					}
				}
				if m.Recv != nil {
					//暂时不支持接收器参数				
				} else {
					if len(m.Results) == 1 {
						sb.WriteString(" ret := " + ri.Package.Alias + "." + m.Name + "(")
					} else {
						sb.WriteString( ri.Package.Alias + "." + m.Name + "(")
					}
				}

				params := ""
				for _, v := range m.Params {
					params = params + v.Name + ","
				}
				if params != "" {
					params = strings.TrimSuffix(params, ",")
				}
				sb.WriteString(params + ")\n") 
				
				if len(m.Results) == 1{
					sb.WriteString(" RenderResult(w,ret)\n")
				} 

				sb.WriteString(" EndTimer(_ticker_id,_ticker)\n}\n\n")
			}

		}
	}

	f_facade = strings.Replace(f_facade, "{{CodeGenFunc}}", sb.String(), -1)

	str, err := format.Source([]byte(f_facade))
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(f_facade)
	} else {
		WriteFile("facade.go",str)
	}			
}
func WriteFile(fileName string,content []byte){
	f := filepath.Join(outDir,fileName)	
	ioutil.WriteFile(f, content, 0666)		
	fmt.Println("生成文件",f,",共",strings.Count(string(content),"\n"),"行")	
}

//将Http方法的首字符大写
func FormatMethodName(n string) string {
	switch strings.ToLower(n) {
	case "get":
		return "Get"
	case "post":
		return "Post"
	default:
		return "Get"
	}
}

func RemoveDuplicateElement(arr []string) (newArr []string) {
    newArr = make([]string, 0)
    for i := 0; i < len(arr); i++ {
        repeat := false
        for j := i + 1; j < len(arr); j++ {
            if FormatMethodName(arr[i]) == FormatMethodName(arr[j]) {
                repeat = true
                break
            }
        }
        if !repeat {
            newArr = append(newArr, FormatMethodName(arr[i]))
        }
    }
    return
}

func GetNewKey(base string, m map[string]PackageInfo) string {
	i := 2
	for {
		ret := fmt.Sprintf("%s%d", base, i)
		_, ok := m[ret]
		if !ok {
			return ret
		} else {
			i++
		}
	}
}

func GetKeyByContent(content string, m map[string]PackageInfo) string {
	for k, v := range m {
		if v.Path == content {
			return k
		}
	}
	return ""
}

var basePath, baseDir, tplDir,outDir string
var template_facade = `
/*
此文件为自动生成，生成时间:{{Date}}
*/
package router

import (
	"net/http"	
	{{RefPackage}}
)

{{CodeGenFunc}}

`
var template_register = `
/*
此文件为自动生成，生成时间:{{Date}}
*/
package router

import (
	"{{ProjectRoot}}/lib/router"
	{{RefPackage}}
)

func Register() (*router.Router,error) {
	r := router.New()
	r.BasePath("/api")

	{{CodeGenFunc}}

	return r,nil
}
`
//类型引用需要的包
var	allRefPackage map[string]PackageInfo
//调用函数需要的包
var	allSrvPackage map[string]PackageInfo

func init() {
	const (
		defaultBasePath = ""
		defaultBaseDir  = "."
		defaultOutDir = "router"
		defaultTmpDir = ""
	)
	flag.StringVar(&basePath, "ProjectRoot", defaultBasePath, "当前路径在项目路径或go path中的位置,或go mod的项目名")
	flag.StringVar(&basePath, "r", defaultBasePath, "当前路径在项目路径或go path中的位置,或go mod的项目名")
	flag.StringVar(&baseDir, "Dir", defaultBaseDir, "要扫描的路径")
	flag.StringVar(&baseDir, "d", defaultBaseDir, "要扫描的路径")
	flag.StringVar(&tplDir, "TemplateDir", defaultTmpDir, "文件模版路径")
	flag.StringVar(&tplDir, "t", defaultTmpDir, "文件模版路径")
	flag.StringVar(&outDir, "OutDir", defaultOutDir, "文件输出路径")
	flag.StringVar(&outDir, "o", defaultOutDir, " 文件输出路径")
	//类型引用需要的包
	allRefPackage = make(map[string]PackageInfo, 0)
	//调用函数需要的包
	allSrvPackage = make(map[string]PackageInfo, 0)
}

func main() {

	flag.Parse()
	
	list := make([]*RouteInfo, 0)
	tplDir = strings.TrimSuffix(tplDir, "/")


	basePath = strings.TrimSuffix(basePath,"/") +  "/"
	ListDir(baseDir)
	Sep := string(os.PathSeparator)	
	baseDir := strings.TrimSuffix(baseDir, Sep) + Sep
	bp := filepath.Dir(baseDir)
	
	for _, v := range files {
		path := filepath.Dir(v)
		path = strings.TrimPrefix(path, bp)
		path = strings.TrimSuffix(strings.TrimPrefix(path, Sep),Sep) + Sep
		path = strings.Replace(path, Sep, "/", -1)

		x := ReadFile(strings.TrimSuffix(basePath+path, Sep), v)
		if x != nil {
			list = append(list, x)
			
		}
	}

	
	// fmt.Println(ToJson(list))
	GenCode(list)

	fmt.Println("Done")

}

var files []string


func ListDir(dirPath string) error {
	dirPath = strings.TrimSuffix(dirPath, string(os.PathSeparator))
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	pathSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() {
			ListDir(dirPath + pathSep + fi.Name())
		} else {
			if filepath.Ext(fi.Name()) == ".go" {
				files = append(files, dirPath+pathSep+fi.Name())
			}
		}
	}
	return nil
}

//序列化
func ToJson(obj interface{}) string {
	json, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("序列化失败" + err.Error())
		return ""
	}
	return string(json)
}
