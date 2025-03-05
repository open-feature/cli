package generators

import (
	"strconv"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
)

func defaultFuncs() template.FuncMap {
	// Update the contributing doc when adding a new function
	return template.FuncMap{
		// TODO rename to ToPascal
		"ToCamel": strcase.ToCamel,
		// TODO rename to ToCamel
		"ToLowerCamel": strcase.ToLowerCamel,
		"ToKebab": strcase.ToKebab,
		"ToScreamingKebab": strcase.ToScreamingKebab,
		"ToSnake": strcase.ToSnake,
		"ToScreamingSnake": strcase.ToScreamingSnake,
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
		"Title": cases.Title,
		"Quote": strconv.Quote,
		"QuoteString": func (input any) any {
			if str, ok := input.(string); ok {
				return strconv.Quote(str)
			}
			return input
		},
	}
}

func init() {
	// results in "Api" using ToCamel("API")
	// results in "api" using ToLowerCamel("API")
	strcase.ConfigureAcronym("API", "api")
}