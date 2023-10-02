package main

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"

	util "book-store-be/utils"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Log           LogOptions `yaml:"log" mapstructure:"log" json:"log"`
	Database      DbOptions  `yaml:"database" mapstructure:"database" json:"database"`
	Observability ObsOptions `yaml:"observability" mapstructure:"observability" json:"observability"`
}

type DbOptions struct {
	ConnectionStringPostgres DbPostgresOptions `yaml:"connectionStringPostgres" mapstructure:"connectionStringPostgres" json:"connectionStringPostgres"`
}

type DbPostgresOptions struct {
	Host     string `yaml:"host" mapstructure:"host" json:"host"`
	Port     int    `yaml:"port" mapstructure:"port" json:"port"`
	User     string `yaml:"user" mapstructure:"user" json:"user"`
	Password string `yaml:"password" mapstructure:"password" json:"password"`
	DbName   string `yaml:"dbname" mapstructure:"dbname" json:"dbname"`
	SslMode  string `yaml:"sslMode" mapstructure:"sslMode" json:"sslMode"`
}

type ObsOptions struct {
	ServiceName string `yaml:"serviceName" mapstructure:"serviceName" json:"serviceName"`
	Endpoint    string `yaml:"endpoint" mapstructure:"endpoint" json:"endpoint"`
	Enable      bool   `yaml:"enable" mapstructure:"enable" json:"enable"`
}
type LogOptions struct {
	Level      int  `yaml:"level" mapstructure:"level" json:"level"`
	EnableJSON bool `yaml:"enableJson" mapstructure:"enableJson" json:"enableJson"`
}

// Default Config il file is not found
var DefaultConfig = Config{
	Log: LogOptions{
		Level:      -1,
		EnableJSON: false,
	},
	Database: DbOptions{
		ConnectionStringPostgres: DbPostgresOptions{},
	},
	Observability: ObsOptions{
		ServiceName: "",
		Endpoint:    "",
		Enable:      false,
	},
}

// Default config file.
//
//go:embed config.yaml
var projectConfigFile []byte

const ConfigFileEnvVar = "BOOK_STORE_BE_FILE_PATH"
const ConfigurationName = "BOOK_STORE_BE"

/*
???    thanks to Mario Imperato for the snippets!!!
*/
func ReadConfig() (*Config, error) {

	configPath := os.Getenv(ConfigFileEnvVar)
	var cfgContent []byte
	var err error
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			log.Info().Str("cfg-file-name", configPath).Msg("reading config")
			cfgContent, err = util.ReadFileAndResolveEnvVars(configPath)
			log.Info().Msg("++++CFG:" + string(cfgContent))
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("the %s env variable has been set but no file cannot be found at %s", ConfigFileEnvVar, configPath)
		}
	} else {
		log.Warn().Msgf("The config path variable %s has not been set. Reverting to bundled configuration", ConfigFileEnvVar)
		cfgContent = util.ResolveConfigValueToByteArray(projectConfigFile)
		// return nil, fmt.Errorf("the config path variable %s has not been set; please set", ConfigFileEnvVar)
	}

	appCfg := DefaultConfig
	err = yaml.Unmarshal(cfgContent, &appCfg)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if !appCfg.Log.EnableJSON {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.SetGlobalLevel(zerolog.Level(appCfg.Log.Level))

	return &appCfg, nil
}
