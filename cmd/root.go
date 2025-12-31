package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	baseDir string
	verbose bool
	dryRun  bool
)

// version is set at build time via ldflags
var version = "dev"

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "composectl",
	Short: "Manage Docker Compose stacks",
	Long: `Composectl is a CLI tool for managing multiple Docker Compose stacks 
with centralized configuration and automated execution.`,
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		// Skip validation for commands that don't need it (like help, version, completion)
		if cmd.Name() == "help" || cmd.Name() == "completion" || cmd.Name() == "__complete" {
			return nil
		}

		// Validate base directory exists
		if _, err := os.Stat(GetBaseDir()); os.IsNotExist(err) {
			return fmt.Errorf("base directory does not exist: %s", GetBaseDir())
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $COMPOSECTL_DIR/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&baseDir, "base-dir", "",
		"base directory for stacks (default: $COMPOSECTL_DIR or /volmain/.@docker_compose)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "print commands without executing (implies --verbose)")

	// Bind flags to viper
	//nolint:errcheck // Flag binding errors are non-critical
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	//nolint:errcheck // Flag binding errors are non-critical
	viper.BindPFlag("base_dir", rootCmd.PersistentFlags().Lookup("base-dir"))
	//nolint:errcheck // Flag binding errors are non-critical
	viper.BindPFlag("dry_run", rootCmd.PersistentFlags().Lookup("dry-run"))
}

func initConfig() {
	// Get base directory
	if baseDir == "" {
		baseDir = os.Getenv("DOCKER_LOADER_DIR")
		if baseDir == "" {
			baseDir = "/volmain/.@docker_compose"
		}
	}
	viper.Set("base_dir", baseDir)

	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in base directory
		viper.AddConfigPath(baseDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}

// GetBaseDir returns the configured base directory
func GetBaseDir() string {
	return viper.GetString("base_dir")
}

// IsVerbose returns whether verbose mode is enabled
// Dry-run mode automatically enables verbose
func IsVerbose() bool {
	return viper.GetBool("verbose") || viper.GetBool("dry_run")
}

// IsDryRun returns whether dry-run mode is enabled
func IsDryRun() bool {
	return viper.GetBool("dry_run")
}

// GetLogFile returns the log file path
func GetLogFile() string {
	return filepath.Join(GetBaseDir(), "docker-loader.log")
}
