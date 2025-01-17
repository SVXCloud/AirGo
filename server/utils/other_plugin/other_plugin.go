package other_plugin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
)

// 约定前端需要匹配的类型,暂时约定4个类型，date->日期，string->字符串，num->数字，bool->布尔
var typeComparisonTable = map[string]string{
	"time.Time":  "date",
	"*time.Time": "date",

	"string":    "string",
	"uuid.UUID": "string",

	"int32":   "num",
	"int64":   "num",
	"int":     "num",
	"float64": "num",
	"float32": "num",

	"bool": "bool",
}

// 对长度不足n的数字后面补0
func Sup(i, n int64) string {
	m := fmt.Sprintf("%d", i)
	for int64(len(m)) < n {
		m = fmt.Sprintf("%s0", m)
	}
	return m
}

// struct转map
func StructToMap(data interface{}) map[string]interface{} {

	m := make(map[string]interface{})

	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem() //需要用struct的Type，不能用指针的Type
	}
	if v.Kind() != reflect.Struct {
		return m
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {

		//field.Name,            //变量名称
		//field.Offset,          //相对于结构体首地址的内存偏移量，string类型会占据16个字节
		//field.Anonymous,       //是否为匿名成员
		//field.Type,            //数据类型，reflect.Type类型
		//field.IsExported(),    //包外是否可见（即是否以大写字母开头）
		//field.Tag.Get("json")) //获取成员变量后面``里面定义的tag

		name := t.Field(i).Name
		tag := t.Field(i).Tag.Get("json")
		if tag == "-" || name == "-" {
			continue
		}
		if tag != "" {
			index := strings.Index(tag, ",")
			if index == -1 {
				name = tag
			} else {
				name = tag[:index]
			}
		}
		m[name] = v.Field(i).Interface()
	}
	return m
}

// m3=获取结构体字段数组，如  [name:aaa,age:18]
// m1=获取结构体字段英文名和中文名map，如  name:名字
// m2=获取结构体字段类型map，如  name:string
func GetStructFieldMap(data interface{}) ([]string, map[string]interface{}, map[string]interface{}) {

	m1 := make([]string, 0)
	m2 := make(map[string]interface{})
	m3 := make(map[string]interface{})
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem() //需要用struct的Type，不能用指针的Type
	}
	if v.Kind() != reflect.Struct {
		return m1, m2, m3
	}
	t := v.Type() //Value转Type
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct {
			//嵌套结构体，递归
			mm1, mm2, mm3 := GetStructFieldMap(v.Field(i).Interface())
			for _, v := range mm1 {
				m1 = append(m1, v)
			}
			for k, v := range mm2 {
				if _, ok := m2[k]; !ok {
					m2[k] = v
				}
			}
			for k, v := range mm3 {
				if _, ok := m3[k]; !ok {
					m3[k] = v
				}
			}

			continue
		}
		//json
		jsonName := t.Field(i).Tag.Get("json")
		if jsonName == "-" {
			continue
		}
		index := strings.Index(jsonName, ",")
		if index != -1 {
			jsonName = jsonName[:index]
		}
		//fmt.Println("jsonName:", jsonName)
		//gorm
		comment := t.Field(i).Tag.Get("gorm")
		if comment == "-" {
			continue
		}
		index = strings.Index(comment, "comment:")
		if index == -1 {
			continue
		}
		comment = comment[index+8:]
		index = strings.Index(comment, ";")
		if index != -1 {
			comment = comment[:index]
		}
		//fmt.Println("comment:", comment)
		m2[jsonName] = comment
		if mv, ok := typeComparisonTable[v.Field(i).Type().String()]; ok {
			m3[jsonName] = mv
		} else {
			m3[jsonName] = typeComparisonTable["string"]
		}
		m1 = append(m1, jsonName)

	}
	//fmt.Println(m)
	return m1, m2, m3

}

// golang struct 动态创建
func RegisterType(elem ...interface{}) map[string]reflect.Type {
	var typeRegistry = make(map[string]reflect.Type)
	for i := 0; i < len(elem); i++ {
		t := reflect.TypeOf(elem[i])
		typeRegistry[t.Name()] = t
	}
	return typeRegistry
}
func NewStruct(name string, typeRegistry map[string]reflect.Type) (interface{}, bool) {
	elem, ok := typeRegistry[name]
	fmt.Println("elem", elem)
	if !ok {
		return nil, false
	}
	return reflect.New(elem).Elem().Interface(), true

}

// gin.Context中获取user id
func GetUserIDFromGinContext(ctx *gin.Context) (int64, bool) {
	userID, ok := ctx.Get("uID")
	return userID.(int64), ok
}

// gin.Context中获取user id
func GetUserNameFromGinContext(ctx *gin.Context) (string, bool) {
	userName, ok := ctx.Get("uName")
	return userName.(string), ok
}
