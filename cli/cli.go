package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/samuelstevens/gocaption/util"
)

const (
	writeHelp     = "Writes any captions found in original .html documents."
	silentHelp    = "Doesn't report any captions to stdout"
	thresholdHelp = "Specifies a minimum confidence threshold."
	configHelp    = "Specify a config file for API keys."
	cacheHelp     = "Specify a json file to cache captions"
	fileTypesHelp = "Specify a comma-separated list of file types to label"
	apiKeyHelp    = "Specify an API key for MS Azure"
	endpointHelp  = "Specfiy an endpoint for MS Azure"
	loudHelp      = "Writes to stdout when getting a new description"

	writeDefault     = false
	silentDefault    = false
	thresholdDefault = 0.7
	configDefault    = "~/.labelrc.json"
	cacheDefault     = "~/.label_captions.json"
	fileTypesDefault = ""
	apiKeyDefault    = ""
	endpointDefault  = ""
	loudDefault      = false
)

type Options struct {
	Write      bool
	Silent     bool
	Files      []string
	ConfigFile string
	CacheFile  string
	Endpoint   string
	APIKey     string
	Threshold  float64
	Loud       bool
}

type ConfigFile struct {
	Endpoint  string  `json:"endpoint"`
	APIKey    string  `json:"key"`
	Threshold float64 `json:"threshold"`
}

func shorthandHelp(help string) string {
	return help + " (shorthand)"
}

func parseConfig(fileName string) *ConfigFile {
	config := ConfigFile{}

	contents, err := ioutil.ReadFile(fileName)

	if err != nil {
		return &config
	}

	err = json.Unmarshal(contents, &config)

	if err != nil {
		log.Fatalf("Cannot parse config file: %s", err.Error())
	}

	return &config
}

func argsToFiles(args []string, validFileTypes *util.StringSet) []string {
	filepaths := []string{}

	for _, path := range args {
		info, err := os.Stat(path)

		if err != nil {
			fmt.Printf("Not parsing %s; %s.\n", path, err.Error())
			continue
		}

		if info.IsDir() {
			err := filepath.Walk(path, func(nestedPath string, nestedInfo os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if nestedInfo.IsDir() {
					// @TODO: prevent git from being captioned
					return nil
				}

				if validFileTypes.Empty() {
					filepaths = append(filepaths, nestedPath)
				} else {

					ext := filepath.Ext(nestedPath)

					if len(ext) != 0 {
						ext = ext[1:]
					}

					if validFileTypes.Contains(ext) {
						filepaths = append(filepaths, nestedPath)
					}

				}

				return nil
			})

			if err != nil {
				fmt.Printf("Not parsing %s; %s.\n", path, err.Error())
			}
		} else {
			filepaths = append(filepaths, path)
		}
	}

	return filepaths
}

func betterConfigString(configValue string, flagValue string) string {
	if flagValue != "" {
		return flagValue
	}

	return configValue
}

func betterConfigFloat(configValue float64, flagValue float64, defaultValue float64) float64 {
	if flagValue != defaultValue {
		return flagValue
	}

	return configValue
}

func Cli() *Options {
	opts := Options{}

	flag.BoolVar(&opts.Write, "write", writeDefault, writeHelp)
	flag.BoolVar(&opts.Write, "w", writeDefault, shorthandHelp(writeHelp))

	flag.BoolVar(&opts.Silent, "silent", silentDefault, silentHelp)
	flag.BoolVar(&opts.Silent, "s", silentDefault, shorthandHelp(silentHelp))
	flag.BoolVar(&opts.Silent, "quiet", silentDefault, silentHelp)
	flag.BoolVar(&opts.Silent, "q", silentDefault, shorthandHelp(silentHelp))

	flag.BoolVar(&opts.Loud, "loud", loudDefault, loudHelp)
	flag.BoolVar(&opts.Loud, "l", loudDefault, shorthandHelp(loudHelp))

	flag.Float64Var(&opts.Threshold, "threshold", thresholdDefault, thresholdHelp)
	flag.Float64Var(&opts.Threshold, "t", thresholdDefault, shorthandHelp(thresholdHelp))

	flag.StringVar(&opts.ConfigFile, "config", configDefault, configHelp)
	flag.StringVar(&opts.ConfigFile, "c", configDefault, shorthandHelp(configHelp))
	opts.ConfigFile = util.ExpandUserDirectory(opts.ConfigFile)

	flag.StringVar(&opts.APIKey, "key", apiKeyDefault, apiKeyHelp)
	flag.StringVar(&opts.APIKey, "k", apiKeyDefault, shorthandHelp(apiKeyHelp))

	flag.StringVar(&opts.Endpoint, "endpoint", endpointDefault, endpointHelp)
	flag.StringVar(&opts.Endpoint, "e", endpointDefault, shorthandHelp(endpointHelp))

	flag.StringVar(&opts.CacheFile, "cache", cacheDefault, cacheHelp)
	opts.CacheFile = util.ExpandUserDirectory(opts.CacheFile)

	var fileTypesFlag string

	flag.StringVar(&fileTypesFlag, "filetypes", fileTypesDefault, fileTypesHelp)
	flag.StringVar(&fileTypesFlag, "f", fileTypesDefault, shorthandHelp(fileTypesHelp))

	flag.Parse()

	fileTypes := util.NewStringSet(strings.Split(fileTypesFlag, ","))

	fileTypes.Remove("") // in case the flag was empty

	args := flag.Args()

	opts.Files = argsToFiles(args, fileTypes)
	config := parseConfig(opts.ConfigFile)

	opts.APIKey = betterConfigString(config.APIKey, opts.APIKey)
	opts.Endpoint = betterConfigString(config.Endpoint, opts.Endpoint)
	opts.Threshold = betterConfigFloat(config.Threshold, opts.Threshold, thresholdDefault)

	return &opts
}
