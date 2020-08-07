package comic

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	// config data file name
	defaultConfigFileName = "config"
	// config data file path
	defaultConfigFilePath = "."
	// placeholder command name used for single command applications
	defaultSingleCommandAppName = "main"
	// separator of nested keys in env vars
	defaultEnvVarNestedKeySeparator = "_"
	// separator of nested keys in Viper
	viperNestedKeySeparator = "."
	// pattern of the path of keys used to set required config variables in config data file
	commandKeyPattern = "required.%s."
	// separator of nested command names
	commandNameSeparator = " "
)

var c *Comic

func init() {
	c = New()
}

// Comic contains all relevant info of a Comic instance
type Comic struct {
	Options
	vip comicViper
}

// Options contains all configurable options of Comic
type Options struct {
	ConfigFileName, ConfigFilePath, SingleCommandAppName, EnvVarNestedKeySeparator string
}

// New creates a new instance of Comic with it's own instance of Viper and default options
func New() *Comic {
	return NewWithOptions(defaultOptions())
}

// NewWithOptions creates a new instance of Comic with it's own instance of Viper and supplied options
func NewWithOptions(opts Options) *Comic {
	defOpts := defaultOptions()

	if opts.ConfigFileName == "" {
		opts.ConfigFileName = defOpts.ConfigFileName
	}

	if opts.ConfigFilePath == "" {
		opts.ConfigFilePath = defOpts.ConfigFilePath
	}

	if opts.SingleCommandAppName == "" {
		opts.SingleCommandAppName = defOpts.SingleCommandAppName
	}

	if opts.EnvVarNestedKeySeparator == "" {
		opts.EnvVarNestedKeySeparator = defOpts.EnvVarNestedKeySeparator
	}

	return &Comic{
		Options: opts,
		vip:     viper.New(),
	}
}

// defaultOptions returns the default options of Comic
func defaultOptions() Options {
	return Options{
		ConfigFileName:           defaultConfigFileName,
		ConfigFilePath:           defaultConfigFilePath,
		SingleCommandAppName:     defaultSingleCommandAppName,
		EnvVarNestedKeySeparator: defaultEnvVarNestedKeySeparator,
	}
}

// FromCommandPath returns the full command name of a command path i.e.
// including all parent commands separated by spaces
// excluding the binary name
// e.g. ./binary command sub-command => command sub-command
func FromCommandPath(commandPath string) string {
	return strings.Join(strings.Split(commandPath, commandNameSeparator)[1:], commandNameSeparator)
}

// Viper returns the Viper instance in use by Comic
func Viper() *viper.Viper { return c.Viper() }
func (c *Comic) Viper() *viper.Viper {
	return c.vip.(*viper.Viper)
}

// MustLoad:
// - verifies the required config variables of the application
// - loads all config variables into the passed struct
// a panic is thrown in case of a failure
//
// note: cfg *must* be a pointer
func MustLoad(cfg interface{}) { c.MustLoad(cfg) }
func (c *Comic) MustLoad(cfg interface{}) {
	c.MustLoadForCommand(cfg, c.SingleCommandAppName)
}

// MustLoadForCommand:
// - verifies the required config variables of the passed command
// - loads all config variables into the passed struct
// a panic is thrown in case of a failure
//
// note: cfg *must* be a pointer
func MustLoadForCommand(cfg interface{}, commandName string) { c.MustLoadForCommand(cfg, commandName) }
func (c *Comic) MustLoadForCommand(cfg interface{}, commandName string) {
	if err := c.LoadForCommand(cfg, commandName); err != nil {
		panic(err)
	}
}

// Load:
// - verifies the required config variables of the application
// - loads all config variables into the passed struct
// an error is returned in case of a failure
//
// note: cfg *must* be a pointer
func Load(cfg interface{}) error { return c.Load(cfg) }
func (c *Comic) Load(cfg interface{}) error {
	return c.LoadForCommand(cfg, c.SingleCommandAppName)
}

// LoadForCommand:
// - verifies the required config variables of the passed command
// - loads all config variables into the passed struct
// an error is returned in case of a failure
//
// note: cfg *must* be a pointer
func LoadForCommand(cfg interface{}, commandName string) error {
	return c.LoadForCommand(cfg, commandName)
}
func (c *Comic) LoadForCommand(cfg interface{}, commandName string) error {
	if commandName == "" {
		return errors.New("command name empty")
	}

	c.vip.SetConfigName(c.ConfigFileName)
	c.vip.AddConfigPath(c.ConfigFilePath)

	c.vip.AutomaticEnv()
	c.vip.SetEnvKeyReplacer(strings.NewReplacer(viperNestedKeySeparator, c.EnvVarNestedKeySeparator))

	if err := c.vip.ReadInConfig(); err != nil {
		return fmt.Errorf("config not loaded: %s", err)
	}

	if err := c.checkRequiredVars(commandName); err != nil {
		return fmt.Errorf("required config for command '%s' missing: %s", commandName, err)
	}

	if err := c.vip.Unmarshal(cfg); err != nil {
		return fmt.Errorf("config not parsed: %s", err)
	}

	return nil
}

// checkRequiredVars verifies that all required config variables are present (i.e. have values)
// for the passed command name
func (c *Comic) checkRequiredVars(commandName string) error {
	for _, varName := range c.getRequiredVarNames(commandName) {
		if !c.vip.IsSet(varName) {
			return fmt.Errorf("config not present: %s", varName)
		}
	}

	return nil
}

// getRequiredVarNames returns the key names of all the required config variables of the passed command name
func (c *Comic) getRequiredVarNames(commandName string) (requiredVarNames []string) {
	for _, key := range c.vip.AllKeys() {
		if requiredKeyName, ok := getRequiredKeyName(key, commandName); ok {
			requiredVarNames = append(requiredVarNames, requiredKeyName)
		}
	}

	return
}

// getRequiredKeyName checks if the passed key is a requirement key of the passed command name
// i.e. a key of a required config variable of the command
// if so, it returns the requirement key name and true, otherwise, it returns an empty string and false
func getRequiredKeyName(key, commandName string) (requiredKeyName string, ok bool) {
	if !strings.HasPrefix(key, fmt.Sprintf(commandKeyPattern, commandName)) {
		return
	}

	if requiredKeyName = strings.TrimPrefix(key, fmt.Sprintf(commandKeyPattern, commandName)); requiredKeyName == "" {
		return
	}

	ok = true

	return
}
