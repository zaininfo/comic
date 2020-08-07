package comic

import (
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	expected := &Comic{
		Options: Options{
			ConfigFileName:           "config",
			ConfigFilePath:           ".",
			SingleCommandAppName:     "main",
			EnvVarNestedKeySeparator: "_",
		},
		vip: viper.New(),
	}

	assert.Equal(t, expected, New())
}

func TestNewWithOptions(t *testing.T) {
	cases := []struct {
		opts     Options
		expected *Comic
	}{
		{
			opts: Options{},
			expected: &Comic{
				Options: Options{
					ConfigFileName:           "config",
					ConfigFilePath:           ".",
					SingleCommandAppName:     "main",
					EnvVarNestedKeySeparator: "_",
				},
				vip: viper.New(),
			},
		},
		{
			opts: Options{
				ConfigFileName: "configuration",
			},
			expected: &Comic{
				Options: Options{
					ConfigFileName:           "configuration",
					ConfigFilePath:           ".",
					SingleCommandAppName:     "main",
					EnvVarNestedKeySeparator: "_",
				},
				vip: viper.New(),
			},
		},
		{
			opts: Options{
				ConfigFilePath: "..",
			},
			expected: &Comic{
				Options: Options{
					ConfigFileName:           "config",
					ConfigFilePath:           "..",
					SingleCommandAppName:     "main",
					EnvVarNestedKeySeparator: "_",
				},
				vip: viper.New(),
			},
		},
		{
			opts: Options{
				SingleCommandAppName: "app",
			},
			expected: &Comic{
				Options: Options{
					ConfigFileName:           "config",
					ConfigFilePath:           ".",
					SingleCommandAppName:     "app",
					EnvVarNestedKeySeparator: "_",
				},
				vip: viper.New(),
			},
		},
		{
			opts: Options{
				EnvVarNestedKeySeparator: "::",
			},
			expected: &Comic{
				Options: Options{
					ConfigFileName:           "config",
					ConfigFilePath:           ".",
					SingleCommandAppName:     "main",
					EnvVarNestedKeySeparator: "::",
				},
				vip: viper.New(),
			},
		},
		{
			opts: Options{
				ConfigFileName: "configuration",
				ConfigFilePath: "..",
			},
			expected: &Comic{
				Options: Options{
					ConfigFileName:           "configuration",
					ConfigFilePath:           "..",
					SingleCommandAppName:     "main",
					EnvVarNestedKeySeparator: "_",
				},
				vip: viper.New(),
			},
		},
		{
			opts: Options{
				ConfigFileName:           "configuration",
				ConfigFilePath:           "..",
				SingleCommandAppName:     "app",
				EnvVarNestedKeySeparator: "::",
			},
			expected: &Comic{
				Options: Options{
					ConfigFileName:           "configuration",
					ConfigFilePath:           "..",
					SingleCommandAppName:     "app",
					EnvVarNestedKeySeparator: "::",
				},
				vip: viper.New(),
			},
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, NewWithOptions(c.opts))
	}
}

func TestDefaultOptions(t *testing.T) {
	expected := Options{
		ConfigFileName:           "config",
		ConfigFilePath:           ".",
		SingleCommandAppName:     "main",
		EnvVarNestedKeySeparator: "_",
	}

	assert.Equal(t, expected, defaultOptions())
}

func TestFromCommandPath(t *testing.T) {
	cases := []struct {
		commandPath, expected string
	}{
		{
			commandPath: "",
			expected:    "",
		},
		{
			commandPath: "./binary",
			expected:    "",
		},
		{
			commandPath: "./binary ",
			expected:    "",
		},
		{
			commandPath: "./binary command",
			expected:    "command",
		},
		{
			commandPath: "./binary command sub-command",
			expected:    "command sub-command",
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, FromCommandPath(c.commandPath))
	}
}

func TestViper(t *testing.T) {
	assert.Equal(t, viper.New(), Viper())
}

func TestComic_Viper(t *testing.T) {
	c := &Comic{
		vip: viper.New(),
	}

	assert.Equal(t, viper.New(), c.Viper())
}

type sampleConfig struct {
	name   string
	server struct {
		port int
	}
}

func loadTestCases() []struct {
	comic          *Comic
	cfg            *sampleConfig
	cmd            string
	expectedOutput *sampleConfig
	expectedError  error
} {
	return []struct {
		comic          *Comic
		cfg            *sampleConfig
		cmd            string
		expectedOutput *sampleConfig
		expectedError  error
	}{
		{
			comic:          &Comic{},
			cfg:            &sampleConfig{},
			cmd:            "",
			expectedOutput: &sampleConfig{},
			expectedError:  errors.New("command name empty"),
		},
		{
			comic: &Comic{
				vip: &mockViper{},
			},
			cfg:            &sampleConfig{},
			cmd:            "run",
			expectedOutput: &sampleConfig{},
			expectedError:  nil,
		},
		{
			comic: &Comic{
				vip: &mockViper{
					cfg: &sampleConfig{
						name: "app",
						server: struct{ port int }{
							port: 123,
						},
					},
					keys: map[string]bool{
						"name":        false,
						"server.port": false,
					},
				},
			},
			cfg: &sampleConfig{},
			cmd: "run",
			expectedOutput: &sampleConfig{
				name: "app",
				server: struct{ port int }{
					port: 123,
				},
			},
			expectedError: nil,
		},
		{
			comic: &Comic{
				vip: &mockViper{
					cfg: &sampleConfig{
						name: "app",
						server: struct{ port int }{
							port: 123,
						},
					},
					keys: map[string]bool{
						"name":                     false,
						"server.port":              false,
						"required.run.name":        false,
						"required.run.server.port": false,
					},
				},
			},
			cfg:            &sampleConfig{},
			cmd:            "run",
			expectedOutput: &sampleConfig{},
			expectedError:  errors.New("required config for command 'run' missing: config not present: name"),
		},
		{
			comic: &Comic{
				vip: &mockViper{
					cfg: &sampleConfig{
						name: "app",
						server: struct{ port int }{
							port: 123,
						},
					},
					keys: map[string]bool{
						"name":                     false,
						"server.port":              false,
						"required.run.name":        false,
						"required.run.server.port": false,
					},
				},
			},
			cfg: &sampleConfig{},
			cmd: "schedule",
			expectedOutput: &sampleConfig{
				name: "app",
				server: struct{ port int }{
					port: 123,
				},
			},
			expectedError: nil,
		},
		{
			comic: &Comic{
				vip: &mockViper{
					cfg: &sampleConfig{
						name: "app",
						server: struct{ port int }{
							port: 123,
						},
					},
					keys: map[string]bool{
						"name":                     true,
						"server.port":              true,
						"required.run.name":        false,
						"required.run.server.port": false,
					},
				},
			},
			cfg: &sampleConfig{},
			cmd: "run",
			expectedOutput: &sampleConfig{
				name: "app",
				server: struct{ port int }{
					port: 123,
				},
			},
			expectedError: nil,
		},
	}
}

func runAndRecover(f func()) (panic interface{}) {
	defer func() {
		panic = recover()
	}()

	f()

	return
}

func TestMustLoad(t *testing.T) {
	for _, tc := range loadTestCases() {
		tc.comic.SingleCommandAppName = tc.cmd
		c = tc.comic

		err := runAndRecover(func() {
			MustLoad(tc.cfg)
		})

		assert.Equal(t, tc.expectedOutput, tc.cfg)
		assert.Equal(t, tc.expectedError, err)
	}
}

func TestComic_MustLoad(t *testing.T) {
	for _, c := range loadTestCases() {
		c.comic.SingleCommandAppName = c.cmd

		err := runAndRecover(func() {
			c.comic.MustLoad(c.cfg)
		})

		assert.Equal(t, c.expectedOutput, c.cfg)
		assert.Equal(t, c.expectedError, err)
	}
}

func TestMustLoadForCommand(t *testing.T) {
	for _, tc := range loadTestCases() {
		c = tc.comic

		err := runAndRecover(func() {
			MustLoadForCommand(tc.cfg, tc.cmd)
		})

		assert.Equal(t, tc.expectedOutput, tc.cfg)
		assert.Equal(t, tc.expectedError, err)
	}
}

func TestComic_MustLoadForCommand(t *testing.T) {
	for _, c := range loadTestCases() {
		err := runAndRecover(func() {
			c.comic.MustLoadForCommand(c.cfg, c.cmd)
		})

		assert.Equal(t, c.expectedOutput, c.cfg)
		assert.Equal(t, c.expectedError, err)
	}
}

func TestLoad(t *testing.T) {
	for _, tc := range loadTestCases() {
		tc.comic.SingleCommandAppName = tc.cmd
		c = tc.comic

		err := Load(tc.cfg)

		assert.Equal(t, tc.expectedOutput, tc.cfg)
		assert.Equal(t, tc.expectedError, err)
	}
}

func TestComic_Load(t *testing.T) {
	for _, c := range loadTestCases() {
		c.comic.SingleCommandAppName = c.cmd

		err := c.comic.Load(c.cfg)

		assert.Equal(t, c.expectedOutput, c.cfg)
		assert.Equal(t, c.expectedError, err)
	}
}

func TestLoadForCommand(t *testing.T) {
	for _, tc := range loadTestCases() {
		c = tc.comic

		err := LoadForCommand(tc.cfg, tc.cmd)

		assert.Equal(t, tc.expectedOutput, tc.cfg)
		assert.Equal(t, tc.expectedError, err)
	}
}

func TestComic_LoadForCommand(t *testing.T) {
	for _, c := range loadTestCases() {
		err := c.comic.LoadForCommand(c.cfg, c.cmd)

		assert.Equal(t, c.expectedOutput, c.cfg)
		assert.Equal(t, c.expectedError, err)
	}
}

func TestComic_checkRequiredVars(t *testing.T) {
	cases := []struct {
		comic         *Comic
		commandName   string
		expectedError error
	}{
		{
			comic: &Comic{
				vip: &mockViper{},
			},
			commandName:   "",
			expectedError: nil,
		},
		{
			comic: &Comic{
				vip: &mockViper{},
			},
			commandName:   "run",
			expectedError: nil,
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":        false,
						"server.port": false,
					},
				},
			},
			commandName:   "run",
			expectedError: nil,
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":              false,
						"server.port":       false,
						"required.run.name": false,
					},
				},
			},
			commandName:   "run",
			expectedError: errors.New("config not present: name"),
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":              true,
						"server.port":       false,
						"required.run.name": false,
					},
				},
			},
			commandName:   "run",
			expectedError: nil,
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":                     false,
						"server.port":              false,
						"required.run.name":        false,
						"required.run.server.port": false,
					},
				},
			},
			commandName:   "run",
			expectedError: errors.New("config not present: name"),
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":                     true,
						"server.port":              false,
						"required.run.name":        false,
						"required.run.server.port": false,
					},
				},
			},
			commandName:   "run",
			expectedError: errors.New("config not present: server.port"),
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":                     false,
						"server.port":              true,
						"required.run.name":        false,
						"required.run.server.port": false,
					},
				},
			},
			commandName:   "run",
			expectedError: errors.New("config not present: name"),
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":                     true,
						"server.port":              true,
						"required.run.name":        false,
						"required.run.server.port": false,
					},
				},
			},
			commandName:   "run",
			expectedError: nil,
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":                         false,
						"server.port":                  false,
						"required.run job.name":        false,
						"required.run job.server.port": false,
					},
				},
			},
			commandName:   "run job",
			expectedError: errors.New("config not present: name"),
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":                         true,
						"server.port":                  false,
						"required.run job.name":        false,
						"required.run job.server.port": false,
					},
				},
			},
			commandName:   "run job",
			expectedError: errors.New("config not present: server.port"),
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":                         false,
						"server.port":                  true,
						"required.run job.name":        false,
						"required.run job.server.port": false,
					},
				},
			},
			commandName:   "run job",
			expectedError: errors.New("config not present: name"),
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":                         true,
						"server.port":                  true,
						"required.run job.name":        false,
						"required.run job.server.port": false,
					},
				},
			},
			commandName:   "run job",
			expectedError: nil,
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"required.run job.name":        false,
						"required.run job.server.port": false,
					},
				},
			},
			commandName:   "run job",
			expectedError: errors.New("config not present: name"),
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":                         true,
						"required.run job.name":        false,
						"required.run job.server.port": false,
					},
				},
			},
			commandName:   "run job",
			expectedError: errors.New("config not present: server.port"),
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"server.port":                  true,
						"required.run job.name":        false,
						"required.run job.server.port": false,
					},
				},
			},
			commandName:   "run job",
			expectedError: errors.New("config not present: name"),
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.expectedError, c.comic.checkRequiredVars(c.commandName))
	}
}

func TestComic_getRequiredVarNames(t *testing.T) {
	cases := []struct {
		comic       *Comic
		commandName string
		expected    []string
	}{
		{
			comic: &Comic{
				vip: &mockViper{},
			},
			commandName: "",
			expected:    nil,
		},
		{
			comic: &Comic{
				vip: &mockViper{},
			},
			commandName: "run",
			expected:    nil,
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":        false,
						"server.port": false,
					},
				},
			},
			commandName: "run",
			expected:    nil,
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":              false,
						"server.port":       false,
						"required.run.name": false,
					},
				},
			},
			commandName: "run",
			expected:    []string{"name"},
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":                     false,
						"server.port":              false,
						"required.run.name":        false,
						"required.run.server.port": false,
					},
				},
			},
			commandName: "run",
			expected:    []string{"name", "server.port"},
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"name":                         false,
						"server.port":                  false,
						"required.run job.name":        false,
						"required.run job.server.port": false,
					},
				},
			},
			commandName: "run job",
			expected:    []string{"name", "server.port"},
		},
		{
			comic: &Comic{
				vip: &mockViper{
					keys: map[string]bool{
						"required.run job.name":        false,
						"required.run job.server.port": false,
					},
				},
			},
			commandName: "run job",
			expected:    []string{"name", "server.port"},
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.comic.getRequiredVarNames(c.commandName))
	}
}

func TestGetRequiredKeyName(t *testing.T) {
	cases := []struct {
		key, commandName, expectedOutput string
		expectedStatus                   bool
	}{
		{
			key:            "required",
			commandName:    "run",
			expectedOutput: "",
			expectedStatus: false,
		},
		{
			key:            "required.",
			commandName:    "run",
			expectedOutput: "",
			expectedStatus: false,
		},
		{
			key:            "required.run",
			commandName:    "run",
			expectedOutput: "",
			expectedStatus: false,
		},
		{
			key:            "required.run.",
			commandName:    "run",
			expectedOutput: "",
			expectedStatus: false,
		},
		{
			key:            "required.run.port",
			commandName:    "",
			expectedOutput: "",
			expectedStatus: false,
		},
		{
			key:            "required.run.port",
			commandName:    "run",
			expectedOutput: "port",
			expectedStatus: true,
		},
		{
			key:            "required.run job.port",
			commandName:    "run job",
			expectedOutput: "port",
			expectedStatus: true,
		},
		{
			key:            "required.run job.port_number",
			commandName:    "run job",
			expectedOutput: "port_number",
			expectedStatus: true,
		},
	}

	for _, c := range cases {
		requiredKeyName, ok := getRequiredKeyName(c.key, c.commandName)

		assert.Equal(t, c.expectedOutput, requiredKeyName)
		assert.Equal(t, c.expectedStatus, ok)
	}
}
