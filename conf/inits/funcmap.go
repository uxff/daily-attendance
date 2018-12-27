package inits

import (
	"fmt"
	"html/template"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/beego/i18n"
	"github.com/ikeikeikeike/gopkg/convert"
	"github.com/mattn/go-runewidth"

	"github.com/astaxie/beego"
)

func init() {

	beego.AddFuncMap("i18n", i18n.Tr)

	beego.AddFuncMap("i18nja", func(format string, args ...interface{}) string {
		return i18n.Tr("ja-JP", format, args...)
	})

	beego.AddFuncMap("datenow", func(format string) string {
		return time.Now().Add(time.Duration(9) * time.Hour).Format(format)
	})

	beego.AddFuncMap("dateformatJst", func(in time.Time) string {
		in = in.Add(time.Duration(9) * time.Hour)
		return in.Format("2006/01/02 15:04")
	})

	beego.AddFuncMap("timefmt", func(in time.Time) string {
		return in.Format("2006-01-02 15:04:05")
	})

	beego.AddFuncMap("timefmtm", func(in time.Time) string {
		return in.Format("2006-01-02 15:04")
	})

	beego.AddFuncMap("qescape", func(in string) string {
		return url.QueryEscape(in)
	})

	beego.AddFuncMap("nl2br", func(in string) string {
		return strings.Replace(in, "\n", "<br>", -1)
	})

	beego.AddFuncMap("tostr", func(in interface{}) string {
		return convert.ToStr(reflect.ValueOf(in).Interface())
	})

	beego.AddFuncMap("first", func(in interface{}) interface{} {
		return reflect.ValueOf(in).Index(0).Interface()
	})

	beego.AddFuncMap("last", func(in interface{}) interface{} {
		s := reflect.ValueOf(in)
		return s.Index(s.Len() - 1).Interface()
	})

	beego.AddFuncMap("truncate", func(in string, length int) string {
		return runewidth.Truncate(in, length, "...")
	})

	beego.AddFuncMap("noname", func(in string) string {
		if in == "" {
			return "(未入力)"
		}
		return in
	})

	beego.AddFuncMap("cleanurl", func(in string) string {
		return strings.Trim(strings.Trim(in, " "), "　")
	})

	beego.AddFuncMap("append", func(data map[interface{}]interface{}, key string, value interface{}) template.JS {
		if _, ok := data[key].([]interface{}); !ok {
			data[key] = []interface{}{value}
		} else {
			data[key] = append(data[key].([]interface{}), value)
		}
		return template.JS("")
	})

	beego.AddFuncMap("appendmap", func(data map[interface{}]interface{}, key string, name string, value interface{}) template.JS {
		v := map[string]interface{}{name: value}

		if _, ok := data[key].([]interface{}); !ok {
			data[key] = []interface{}{v}
		} else {
			data[key] = append(data[key].([]interface{}), v)
		}
		return template.JS("")
	})

	beego.AddFuncMap("mod", func(big int, lit int, zeroEcho, nonZeroEcho string) string {
		if lit > 0 && big%lit == 0 {
			return zeroEcho
		}
		return nonZeroEcho
	})

	beego.AddFuncMap("tri", func(b bool, trueEcho, falseEcho string) string {
		if b {
			return trueEcho
		}
		return falseEcho
	})

	// needle=23 see=12,45,23,67 return trueEcho
	beego.AddFuncMap("inarr", func(needle interface{}, sea string, trueEcho, falseEcho string) string {
		needleStr := fmt.Sprintf("%v", needle)
		seas := strings.Split(sea, ",")
		for _, s := range seas {
			if s == needleStr {
				return trueEcho
			}
		}
		return falseEcho
	})

	// genlist 2 3 4 return [3, 7]; genlist 4 5 6 return [5,11,17,23]
	beego.AddFuncMap("genlist", func(num, start, eachStep int) []int {
		dss := make([]int, num)
		for i := 0; i < num; i++ {
			dss[i] = i*eachStep + start
		}
		return dss
	})

	beego.AddFuncMap("plus", func(num1, num2 int) string {
		return fmt.Sprintf("%d", num1+num2)
	})

}
