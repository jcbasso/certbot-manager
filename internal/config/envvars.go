package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"reflect"
	"strings"
)

// bindEnvsRecursive DFS traverses a struct, building Viper keys and ENV var names, then binds them.
// viperKeyPrefix: e.g., "globals", "certificates.0" (accumulated dot-separated path for Viper)
// envVarPrefix: e.g., "CERTBOT_MANAGER_GLOBALS_", "CERTBOT_MANAGER_CERTIFICATES_0_" (accumulated underscore-separated for ENV)
// val: The reflect.Value of the struct/field to inspect.
func bindEnvsRecursive(viperKeyPrefix string, val reflect.Value, v *viper.Viper) {
	envVarPrefix := fmt.Sprintf("%s_%s", v.GetEnvPrefix(), strings.ToUpper(viperKeyPrefix))
	bindEnvsRecursive2(viperKeyPrefix, envVarPrefix, val, v)
}

// bindEnvsRecursive2 is a helper function that does the actual work of binding ENV vars to Viper keys.
func bindEnvsRecursive2(viperKeyPrefix string, envVarPrefix string, val reflect.Value, v *viper.Viper) {
	// Dereference pointer if val is a pointer to a struct
	if val.Kind() == reflect.Ptr {
		if val.IsNil() { // Important: If pointer is nil, can't proceed
			return
		}
		val = val.Elem()
	}

	// Ensure we are working with a struct
	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !field.IsExported() {
			continue // Skip unexported fields
		}

		// --- Determine Field Name for Viper Key and ENV Var Segment ---
		tag := field.Tag.Get("mapstructure")
		tagParts := strings.Split(tag, ",")
		mapKeyName := tagParts[0] // Name from mapstructure tag (e.g., "webroot_path")
		if mapKeyName == "" {
			// If no mapstructure tag, default to lowercase field name for Viper key
			mapKeyName = strings.ToLower(field.Name)
		}
		if mapKeyName == "-" { // Skip fields explicitly ignored by mapstructure
			continue
		}

		// Determine if this field is squashed
		isSquashed := false
		for _, part := range tagParts {
			if part == "squash" {
				isSquashed = true
				break
			}
		}

		// --- Construct Paths and Recurse or Bind ---
		var nextViperKeyPath string
		var nextEnvVarName string

		if isSquashed {
			// For squashed structs, the fields of the embedded struct are treated as if
			// they are part of the parent. So, we don't add the squashed field's name
			// to the path/prefix, but recurse on its value.
			if fieldVal.Kind() == reflect.Struct || (fieldVal.Kind() == reflect.Ptr && !fieldVal.IsNil() && fieldVal.Elem().Kind() == reflect.Struct) {
				bindEnvsRecursive2(viperKeyPrefix, envVarPrefix, fieldVal, v)
			}
		} else {
			// For regular (non-squashed) fields:
			// Viper key path segment is the mapKeyName
			if viperKeyPrefix == "" { // Top-level field of the initial struct
				nextViperKeyPath = mapKeyName
			} else {
				nextViperKeyPath = viperKeyPrefix + "." + mapKeyName
			}

			// Environment variable segment is the UPPPERCASE mapKeyName
			// (or uppercased field name if no mapstructure tag, though tag is preferred)
			envSegment := strings.ToUpper(mapKeyName)

			if envVarPrefix == "" { // Top-level field, use only APP_PREFIX + SEGMENT
				// This case is for fields directly in the `Config` struct if they were bound directly
				// For our setup, we usually start recursion with a base prefix.
				// The global `v.GetEnvPrefix()` should be the app's main prefix.
				nextEnvVarName = v.GetEnvPrefix() + "_" + envSegment // Assumes EnvPrefix is set
			} else {
				nextEnvVarName = envVarPrefix + "_" + envSegment
			}

			// Recurse for nested structs (that are not squashed)
			if fieldVal.Kind() == reflect.Struct || (fieldVal.Kind() == reflect.Ptr && !fieldVal.IsNil() && fieldVal.Elem().Kind() == reflect.Struct) {
				bindEnvsRecursive2(nextViperKeyPath, nextEnvVarName, fieldVal, v)
			} else if fieldVal.Kind() == reflect.Slice && fieldVal.Type().Elem().Kind() == reflect.Struct {
				// --- Advanced: Handle Slices of Structs ---
				// This would bind CERTBOT_MANAGER_CERTIFICATES_0_CMD, CERTBOT_MANAGER_CERTIFICATES_1_CMD etc.
				// We need to iterate a few potential indices, as Viper doesn't auto-expand array env vars.
				// This is a common limitation. For simplicity, we might bind a few (e.g., 0-9)
				// or users would typically override the whole slice via a JSON/YAML string in one ENV var.
				// For now, let's just log a placeholder or skip.
				// log.Printf("Info: Binding for slice '%s' elements. Max 10 elements supported via individual ENV vars.", nextViperKeyPath)
				// for j := 0; j < 10; j++ { // Arbitrary limit for example
				// 	elemViperPath := fmt.Sprintf("%s.%d", nextViperKeyPath, j)
				// 	elemEnvPrefix := fmt.Sprintf("%s_%d", nextEnvVarName, j)
				// 	// Need a dummy struct of the slice element type to recurse
				// 	elemType := fieldVal.Type().Elem()
				// 	if elemType.Kind() == reflect.Ptr { elemType = elemType.Elem() } // Handle ptr to struct
				//  if elemType.Kind() == reflect.Struct {
				//      dummyElem := reflect.New(elemType).Elem() // Create new instance of element type
				//      bindEnvsRecursive(elemViperPath, elemEnvPrefix, dummyElem)
				//  }
				// }
				log.Printf("Skipping automatic ENV binding for slice elements in '%s'. Configure via TOML or a single ENV var with JSON/YAML string.", nextViperKeyPath)

			} else {
				// Bind ENV for simple (non-struct, non-slice) fields
				_ = v.BindEnv(nextViperKeyPath, nextEnvVarName)
			}
		}
	}
}
