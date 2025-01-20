package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/maintc/rustmaps-cli/pkg/rustmaps"
	"github.com/maintc/rustmaps-cli/pkg/types"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	generator *rustmaps.Generator
	logger    *zap.Logger
	logLevel  string
)

func GetGenerator() *rustmaps.Generator {
	return generator
}

// initLogger initializes the Zap logger with proper configuration
func initLogger() error {
	// Create a zapcore.EncoderConfig with your preferred encoding configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create an encoder (used by both cores)
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Open the log file for writing
	logFilePath := generator.GetLogPath()
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file %q: %v", logFilePath, err)
	}
	fileWriteSyncer := zapcore.AddSync(logFile)

	// Create the stdout write syncer
	stdoutWriteSyncer := zapcore.Lock(os.Stdout)

	// Parse the log level for stdout
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return fmt.Errorf("invalid log level %q: %v", logLevel, err)
	}

	// Create two zapcore.Core:
	// 1. For stdout, using the specified log level
	stdoutCore := zapcore.NewCore(encoder, stdoutWriteSyncer, zap.NewAtomicLevelAt(level))

	// 2. For file, always logging at debug level (or all levels)
	fileCore := zapcore.NewCore(encoder, fileWriteSyncer, zap.DebugLevel)

	// Combine the two cores using zapcore.NewTee
	core := zapcore.NewTee(stdoutCore, fileCore)

	// Build the logger from the combined core
	logger = zap.New(core, zap.AddCaller())

	// Ensure we sync the logger on exit
	defer logger.Sync()

	return nil
}

func loadFromParams(csv, seed string, size int, savedConfig string, staging, force, random bool) {
	if csv != "" {
		if err := generator.LoadCSV(logger, csv); err != nil {
			fmt.Println("Error validating map file, check logs for more info")
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		if random {
			seed = generator.GetRandomSeed()
		}
		m := types.NewMap(seed, size, savedConfig, staging)
		generator.AddMap(m)
	}

	if err := generator.Import(logger, force); err != nil {
		fmt.Println("Failed to import file, check logs for more info")
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "rustmaps",
	Short: "RustMaps CLI",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		generator, err = rustmaps.NewGenerator(nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing generator: %v\n", err)
			os.Exit(1)
		}

		if err := generator.InitDirs(); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing directories: %v\n", err)
			os.Exit(1)
		}

		if err := generator.LoadConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		if err := initLogger(); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		fmt.Println()
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  Resource\tPath") // Indented header
		fmt.Fprintln(w, "  --------\t----")
		fmt.Fprintf(w, "  Downloads directory\t%s\n", generator.GetDownloadsDir())
		fmt.Fprintf(w, "  Imports directory\t%s\n", generator.GetImportDir())
		fmt.Fprintf(w, "  Config file\t%s\n", generator.GetConfigPath())
		fmt.Fprintf(w, "  Log file\t%s\n", generator.GetLogPath())
		w.Flush()
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "fatal",
		"Log level (debug, info, warn, error, dpanic, panic, fatal)")
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(generateCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
