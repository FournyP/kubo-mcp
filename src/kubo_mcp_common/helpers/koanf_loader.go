package helpers

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const defaultConfigFile = "appsettings.json"

// LoadKoanfConfig reads an optional JSON file, then `SECTION__FIELD` env vars.
func LoadKoanfConfig(configFile string) (*koanf.Koanf, error) {
	k := koanf.New(".")

	if configFile == "" {
		configFile = defaultConfigFile
	}

	if _, err := os.Stat(configFile); err == nil {
		if err := k.Load(file.Provider(configFile), json.Parser()); err != nil {
			return nil, err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	envMap := mapEnvEntries(os.Environ())
	if len(envMap) == 0 {
		return k, nil
	}

	if err := k.Load(confmap.Provider(envMap, "."), nil); err != nil {
		return nil, err
	}

	return k, nil
}

func mapEnvEntries(entries []string) map[string]interface{} {
	envMap := map[string]interface{}{}
	for _, entry := range entries {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]
		if !strings.Contains(key, "__") {
			continue
		}

		segments := strings.Split(strings.ToLower(key), "__")
		updated, ok := insertEnvEntryValue(envMap, segments, value, false).(map[string]interface{})
		if ok {
			envMap = updated
		}
	}

	return envMap
}

func insertEnvEntryValue(current interface{}, segments []string, value string, previousNumeric bool) interface{} {
	if len(segments) == 0 {
		return current
	}

	segment := segments[0]
	segmentIsNumeric := isNumericSegment(segment)
	isLast := len(segments) == 1
	nextIsNumeric := len(segments) > 1 && isNumericSegment(segments[1])

	if segmentIsNumeric && previousNumeric {
		currentSlice, ok := current.([]interface{})
		if !ok {
			currentSlice = []interface{}{}
		}

		index, err := strconv.Atoi(segment)
		if err != nil || index < 0 {
			return currentSlice
		}

		for len(currentSlice) <= index {
			currentSlice = append(currentSlice, nil)
		}

		if isLast {
			currentSlice[index] = value
			return currentSlice
		}

		currentSlice[index] = insertEnvEntryValue(currentSlice[index], segments[1:], value, segmentIsNumeric)
		return currentSlice
	}

	currentMap, ok := current.(map[string]interface{})
	if !ok {
		currentMap = map[string]interface{}{}
	}

	if isLast {
		currentMap[segment] = value
		return currentMap
	}

	child, exists := currentMap[segment]
	if !exists || child == nil {
		if segmentIsNumeric && nextIsNumeric {
			child = []interface{}{}
		} else {
			child = map[string]interface{}{}
		}
	}

	currentMap[segment] = insertEnvEntryValue(child, segments[1:], value, segmentIsNumeric)
	return currentMap
}

func isNumericSegment(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil
}
