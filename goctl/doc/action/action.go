package action

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
)

const (
	defaultTemplate = `# {{ .title }}
author: {{ .author }}

version: {{ .version }}

## {{ .service }}

group: {{ .group }}

middleware: {{ .middleware }}
`
	defaultRoutesTemplate = `
### {{ .index }}. {{ .routeDocs }}

##### {{ .method }} {{ .uri }}

Request:
{{ .requestContent }}

Response:
{{ .responseContent }}
`
)

// DocAction generate Markdown doc file
func DocAction(c *cli.Context) error {
	api, defaultOutputFile, err := getApi(c)
	if err != nil {
		return err
	}

	outputFile, err := getOutputFile(c, defaultOutputFile)
	if err != nil {
		return nil
	}
	outputDir := filepath.Dir(outputFile)
	outputFile = outputFile[len(outputDir):]

	mainTemplate, routesTemplate, err := getTemplates(c)
	if err != nil {
		return err
	}

	return genDoc(api, outputDir, outputFile, mainTemplate, routesTemplate)
}

func getApi(c *cli.Context) (*spec.ApiSpec, string, error) {
	apiFile := c.String("api")
	if len(apiFile) == 0 {
		api, err := getApiFromPlugin()
		if err != nil {
			return nil, "", fmt.Errorf("%s, please check the -api option", err.Error())
		}
		return api, "api.md", nil
	}
	apiFile, err := filepath.Abs(apiFile)
	if err != nil {
		return nil, "", err
	}
	api, err := parser.Parse(apiFile)
	if err != nil {
		return nil, "", fmt.Errorf("parse file: %s, err: %w", apiFile, err)
	}
	api.Service = api.Service.JoinPrefix()
	defaultOutputFile := apiFile[len(filepath.Dir(apiFile)):]
	defaultOutputFile = strings.Replace(defaultOutputFile, ".api", ".md", 1)
	return api, defaultOutputFile, nil
}

func getApiFromPlugin() (*spec.ApiSpec, error) {
	p, err := plugin.NewPlugin()
	if err != nil {
		return nil, err
	}
	if p.Api == nil {
		return nil, errors.New("no api")
	}
	return p.Api, nil
}

func getOutputFile(c *cli.Context, defaultOutName string) (string, error) {
	outputFile := c.String("o")
	if len(outputFile) == 0 {
		var err error
		outputDir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		outputFile = filepath.Join(outputDir, defaultOutName)
	}
	outputFile, err := filepath.Abs(outputFile)
	if err != nil {
		return "", err
	}

	return outputFile, nil
}

func getTemplates(c *cli.Context) (string, string, error) {
	mainTemplate, err := getFileContentOrDefault(c.String("mainTemplate"), defaultTemplate)
	if err != nil {
		return "", "", err
	}

	routesTemplate, err := getFileContentOrDefault(c.String("routesTemplate"), defaultRoutesTemplate)
	if err != nil {
		return "", "", err
	}

	return mainTemplate, routesTemplate, nil
}

func getFileContentOrDefault(path, defaultContent string) (string, error) {
	if path != "" {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}

		return string(data), nil
	}

	return defaultContent, nil
}
