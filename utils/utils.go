package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

/*
???    thanks to Mario Imperato for the snippets!!!package main

*/

var ConfigValueRegexp = regexp.MustCompile("(\\$\\{([a-zA-Z][a-zA-Z0-9_]*)\\})")

func FileSize(fn string) int64 {
	if fi, err := os.Stat(fn); err == nil {
		return fi.Size()

	} else if errors.Is(err, os.ErrNotExist) {
		return -1

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return -1
	}
}

func ResolveConfigValueToByteArray(v []byte) []byte {
	matches := ConfigValueRegexp.FindAllSubmatch(v, -1)

	sv := string(v)
	for _, m := range matches {
		env, ok := os.LookupEnv(string(m[2]))
		if ok {
			sv = strings.ReplaceAll(sv, string(m[1]), env)
		}
	}

	return []byte(sv)
}

func ResolveConfigValueToString(v string) string {

	matches := ConfigValueRegexp.FindAllSubmatch([]byte(v), -1)

	for _, m := range matches {
		env, ok := os.LookupEnv(string(m[2]))
		if ok {
			v = strings.ReplaceAll(v, string(m[1]), env)
		}
	}

	return v
}

func ReadFileAndResolveEnvVars(cfgFile string) ([]byte, error) {

	fsz := FileSize(cfgFile)
	if fsz < 0 {
		return nil, fmt.Errorf("error reading file %s", cfgFile)
	}

	file, err := os.Open(cfgFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sb bytes.Buffer
	sb.Grow(int(fsz))

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sb.WriteString(ResolveConfigValueToString(scanner.Text()))
		sb.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sb.Bytes(), nil
}

/*
func ReadConfig(fileConfigPathEnvVar string, defaultConfigContent []byte, resolveEnv bool) (string, []byte, error) {

	configPath := os.Getenv(fileConfigPathEnvVar)
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			cfgContent, rerr := ReadFileAndResolveEnvVars(configPath)
			if rerr != nil {
				return configPath, nil, err
			}
			return configPath, cfgContent, nil
		}
		return configPath, nil, fmt.Errorf("the %s env variable has been set but no file cannot be found at %s", fileConfigPathEnvVar, configPath)
	}

	if len(defaultConfigContent) > 0 {
		return "", ResolveConfigValueToByteArray(defaultConfigContent), nil
	}

	return "", nil, fmt.Errorf("the config path variable %s has not been set; please set", fileConfigPathEnvVar)
}
*/
