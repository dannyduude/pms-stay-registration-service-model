package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dimiro1/banner"
	"github.com/mattn/go-colorable"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, cfgDir string

const defaultCExt = "json"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "go-gravity-quick",
	Short: "go-gravity-quick template",
	Long: `go-gravity-quick template service.
	                 Used as base for developing new go webservices within HP3D Software organization.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		showBanner()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config_appname.yaml)")
	RootCmd.PersistentFlags().StringVar(&cfgDir, "configDir", "", "config directory. Loads all valid config files in this directory (default is none")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// first of all load env variables
	viper.AutomaticEnv()

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}

		return
	}

	if cfgDir != "" {
		fmt.Println("Using config dir:", cfgDir)
		if err := loadConfigDir(cfgDir); err != nil {
			fmt.Println("error reading config dir:", err)
		}

		return
	}

	envConfigDir := os.Getenv("CONFIG_DIR")
	if envConfigDir != "" {
		fmt.Println("Using config dir:", envConfigDir)
		if err := loadConfigDir(envConfigDir); err != nil {
			fmt.Println("error reading config dir:", err)
		}

		return
	}

	// Search config in home directory with name ".config.yml".
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	viper.SetConfigName(".env")
	viper.SetConfigName(".config")
	viper.SetConfigName("config.yml")

	viper.SetConfigType("yml")
	viper.AddConfigPath(home)
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../../config")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func loadConfigDir(configDir string) error {
	files, err := os.ReadDir(configDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			fmt.Println("reading config file:", file.Name())
			err := loadConfigFile(configDir, file.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func loadConfigFile(cPath, cFileName string) error {
	cExt := filepath.Ext(cFileName)
	if len(cExt) > 2 {
		viper.SetConfigType(cExt[1:])
	}

	if cExt == "" {
		fmt.Printf("using default configuration extension: %s", defaultCExt)
		viper.SetConfigType(defaultCExt)
	}

	viper.AddConfigPath(cPath)
	viper.SetConfigName(cFileName)
	return viper.MergeInConfig()
}

func showBanner() {

	templ := `{{ .Title "Go-Gravity-Template" "" 4 }}
   {{ .AnsiColor.BrightCyan }}The title for the application{{ .AnsiColor.Default }}
   GO_Version: {{ .GoVersion }}
   GO_OS: {{ .GOOS }}
   GO_ARCH: {{ .GOARCH }}
   NumCPU: {{ .NumCPU }}
   GO_PATH: {{ .GOPATH }}
   GO_ROOT: {{ .GOROOT }}
   Compiler: {{ .Compiler }}
   ENV: {{ .Env "GOPATH" }}
   HTTP_PROXY: {{ .Env "http_proxy" }}
   HTTPS_PROXY: {{ .Env "https_proxy" }}   
   Now: {{ .Now "Monday, 2 Jan 2006" }}

   {{ .AnsiColor.BrightGreen }}Version: ` + Version + `{{ .AnsiColor.Default }}
   
   `
	banner.InitString(colorable.NewColorableStdout(), true, true, templ)
}
