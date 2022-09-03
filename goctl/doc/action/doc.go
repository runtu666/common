package action

import (
	"bytes"
	"fmt"
	"go/format"
	"strconv"
	"strings"
	"text/template"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/api/util"
)

func genDoc(api *spec.ApiSpec, dir, filename, mainTemplate, routeTemplate string) error {
	fp, _, err := util.MaybeCreateFile(dir, "", filename)
	if err != nil {
		return err
	}
	defer fp.Close()

	var builder strings.Builder
	err = parseMain(api, mainTemplate, &builder)
	if err != nil {
		return err
	}

	err = parseRoutes(api, routeTemplate, &builder)
	if err != nil {
		return err
	}

	_, err = fp.WriteString(strings.ReplaceAll(builder.String(), "&#34;", `"`))
	return err
}

func parseMain(api *spec.ApiSpec, tmp string, builder *strings.Builder) error {
	title := strings.Trim(api.Info.Properties["title"], `"`)
	author := strings.Trim(api.Info.Properties["author"], `"`)
	version := strings.Trim(api.Info.Properties["version"], `"`)
	service := api.Service.Name
	group := ""
	middleware := ""
	if len(api.Service.Groups) > 0 {
		group = api.Service.Groups[0].Annotation.Properties["group"]
		middleware = api.Service.Groups[0].Annotation.Properties["middleware"]
	}
	var buf bytes.Buffer
	t := template.Must(template.New("mainTemplate").Parse(tmp))
	err := t.Execute(&buf, map[string]string{
		"title":      title,
		"author":     author,
		"version":    version,
		"service":    service,
		"group":      group,
		"middleware": middleware,
	})
	if err != nil {
		return err
	}
	builder.Write(buf.Bytes())
	return nil
}

func parseRoutes(api *spec.ApiSpec, tmp string, builder *strings.Builder) error {
	types := parseTypes(api)

	for index, route := range api.Service.Routes() {
		routeDocs := strings.Trim(route.JoinedDoc(), `"`)
		if len(routeDocs) == 0 {
			routeDocs = "N/A"
		}
		var err error
		requestContent := ""
		responseContent := ""
		if route.RequestType != nil {
			requestContent, err = getTypes(route.RequestType.Name(), types)
		}
		if err != nil {
			return err
		}
		if route.ResponseType != nil {
			responseContent, err = getTypes(route.ResponseType.Name(), types)
		}
		if err != nil {
			return err
		}

		t := template.Must(template.New("routeTemplate").Parse(tmp))
		var buf bytes.Buffer
		err = t.Execute(&buf, map[string]string{
			"index":           strconv.Itoa(index + 1),
			"routeDocs":       routeDocs,
			"method":          strings.ToUpper(route.Method),
			"uri":             route.Path,
			"requestContent":  requestContent,
			"responseContent": responseContent,
		})
		if err != nil {
			return err
		}

		builder.Write(buf.Bytes())
	}
	return nil
}

func parseTypes(api *spec.ApiSpec) map[string]string {
	res := make(map[string]string, len(api.Types))
	for _, tp := range api.Types {
		s, err := parseType(tp)
		if err != nil {
			fmt.Println(err)
			continue
		}
		res[tp.Name()] = s
	}
	return res
}

func parseType(tp spec.Type) (string, error) {
	structType, ok := tp.(spec.DefineStruct)
	if !ok {
		return "", fmt.Errorf("unspport struct type: %s", tp.Name())
	}
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("type %s struct {\n", tp.Name()))
	for _, member := range structType.Members {
		if member.IsInline {
			buf.WriteString(fmt.Sprintf("%s\n", member.Type.Name()))
			continue
		}
		comment := member.GetComment()
		if len(comment) == 0 {
			buf.WriteString(fmt.Sprintf("%s %s %s\n", member.Name, member.Type.Name(), member.Tag))
		} else {
			comment := strings.TrimPrefix(comment, "//")
			comment = "//" + comment
			buf.WriteString(fmt.Sprintf("%s %s %s %s\n", member.Name, member.Type.Name(), member.Tag, comment))
		}
	}
	buf.WriteString("}")
	return buf.String(), nil
}

func getTypes(key string, types map[string]string) (string, error) {
	buf := bytes.NewBuffer(nil)
	seen := map[string]bool{}
	var help func(key string)
	help = func(key string) {
		if key == "" || seen[key] {
			return
		}
		tv, ok := types[key]
		if !ok || tv == "" {
			return
		}
		seen[key] = true
		buf.WriteString(tv)
		buf.WriteByte('\n')

		for _, v := range strings.Split(tv, "\n") {
			v = strings.TrimSpace(v)
			if strings.HasPrefix(v, "type") || strings.HasPrefix(v, "//") {
				continue
			}
			fields := strings.Fields(v)
			if len(fields) < 2 {
				continue
			}
			key := fields[1]
			key = strings.TrimPrefix(key, "[]") // []XXX
			key = strings.TrimPrefix(key, "*")  // *XXX, []*XXX
			key = strings.TrimPrefix(key, "map[string]")
			key = strings.TrimPrefix(key, "map[string]*")
			help(key)
		}
	}
	help(key)

	content := formatGoCode(buf.Bytes())
	content = strings.TrimSpace(content)
	return fmt.Sprintf("\n```golang\n%s\n```\n", content), nil
}

func formatGoCode(code []byte) string {
	if len(code) == 0 {
		return ""
	}
	const dummyCode = "package doc\n"
	code = append([]byte(dummyCode), code...)
	formated, err := format.Source(code)
	if err != nil {
		fmt.Println(err)
		return string(code)
	}
	res := string(formated)
	// delete dummy code
	res = res[len(dummyCode):]
	return strings.TrimSpace(res)
}
