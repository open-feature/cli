// Package manifestutils contains useful functions for loading the flag manifest.
package manifestutils

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/open-feature/cli/internal/filesystem"
	"github.com/open-feature/cli/internal/flagkeys"
	"github.com/open-feature/cli/internal/generate/types"
	flagmanifest "github.com/open-feature/cli/schema/v0"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// LoadData loads the data from the flag manifest.
func LoadData(manifestPath string, supportedFlagTypes map[types.FlagType]bool) (*types.BaseTmplData, error) {
	fs := filesystem.FileSystem()
	data, err := afero.ReadFile(fs, manifestPath)
	if err != nil {
		return nil, fmt.Errorf("error reading contents from file %q", manifestPath)
	}
	unfilteredData, err := unmarshalFlagManifest(data)
	if err != nil {
		return nil, err
	}

	filteredData := filterUnsupportedFlags(unfilteredData, supportedFlagTypes)

	return filteredData, nil
}

func filterUnsupportedFlags(unfilteredData *types.BaseTmplData, supportedFlagTypes map[types.FlagType]bool) *types.BaseTmplData {
	filteredData := &types.BaseTmplData{
		OutputPath: unfilteredData.OutputPath,
	}
	for _, flagData := range unfilteredData.Flags {
		if supportedFlagTypes[flagData.Type] {
			filteredData.Flags = append(filteredData.Flags, flagData)
		}
	}
	return filteredData
}

var stringToFlagType = map[string]types.FlagType{
	"string":  types.StringType,
	"boolean": types.BoolType,
	"float":   types.FloatType,
	"integer": types.IntType,
	"object":  types.ObjectType,
}

func getDefaultValue(defaultValue interface{}, flagType types.FlagType) string {
	switch flagType {
	case types.BoolType:
		return strconv.FormatBool(defaultValue.(bool))
	case types.IntType:
		//the conversion to float64 instead of integer typically occurs
		//due to how JSON is parsed in Go. In Go's encoding/json package,
		//all JSON numbers are unmarshaled into float64 by default when decoding into an interface{}.
		fallthrough
	case types.FloatType:
		return strconv.FormatFloat(defaultValue.(float64), 'g', -1, 64)
	case types.StringType:
		return defaultValue.(string)
	default:
		return ""
	}
}

func unmarshalFlagManifest(data []byte) (*types.BaseTmplData, error) {
	dynamic := make(map[string]interface{})
	err := json.Unmarshal(data, &dynamic)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	sch, err := jsonschema.CompileString(flagmanifest.SchemaPath, flagmanifest.Schema)
	if err != nil {
		return nil, fmt.Errorf("error compiling JSON schema: %v", err)
	}
	if err = sch.Validate(dynamic); err != nil {
		return nil, fmt.Errorf("error validating JSON schema: %v", err)
	}
	// All casts can be done directly since the JSON is already validated by the schema.
	iFlags := dynamic["flags"]
	flags := iFlags.(map[string]interface{})
	btData := types.BaseTmplData{
		OutputPath: viper.GetString(flagkeys.OutputPath),
	}
	for flagKey, iFlagData := range flags {
		flagData := iFlagData.(map[string]interface{})
		flagTypeString := flagData["flagType"].(string)
		flagType := stringToFlagType[flagTypeString]
		docs := flagData["description"].(string)
		defaultValue := getDefaultValue(flagData["defaultValue"], flagType)
		btData.Flags = append(btData.Flags, &types.FlagTmplData{
			Name:         flagKey,
			Type:         flagType,
			DefaultValue: defaultValue,
			Docs:         docs,
		})
	}
	// Ensure consistency of order of flag generation.
	sort.Slice(btData.Flags, func(i, j int) bool {
		return btData.Flags[i].Name < btData.Flags[j].Name
	})
	return &btData, nil
}
