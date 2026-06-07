package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go-printos-backend-quickstart/api"

	mongo "github.azc.ext.hp.com/3DSoftware/go-gravity-mongo-driver/v6"
	featureflagsdk "github.azc.ext.hp.com/3DSoftware/gravity-feature-flag/v2"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// Command-line action
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start the go-gravity-quick server",
	Long:  `starts a http server and serves the api`,
	Run: func(cmd *cobra.Command, args []string) {
		initServer()
		server()
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func initServer() {
	log.SetPrefix("go-quickstart: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func server() {

	viper.SetDefault("LAUNCH_DARKLY_OFFLINE", false)

	errs := startServer(ServerConfig{
		Version:         viper.GetString("QST_SERVICE_VERSION"),
		LogFilePath:     viper.GetString("QST_LOG_FILE_PATH"),
		Host:            viper.GetString("QST_HOST"),
		ExecuteAsDryRun: viper.GetBool("EXECUTE_AS_DRY_RUN"),
		HTTPPort:        viper.GetString("QST_HTTP_PORT"),
		FeatureFlagsConfig: featureflagsdk.Config{
			SDKkey:          viper.GetString("LAUNCHDARKLY_SDK_KEY"),
			Offline:         viper.GetBool("LAUNCHDARKLY_OFFLINE"),
			FileDatasources: viper.GetString("LAUNCHDARKLY_LOADDATA_FROM_FILES"),
		},
		//REMOVE IF NOT USING MONGO
		MongoConfig: mongo.Config{
			EnableDNSSeedlist:    viper.GetBool("MONGODB_ENABLE_DNS_SEEDLIST"),
			Username:             viper.GetString("MONGODB_USERNAME"),
			Password:             viper.GetString("MONGODB_PASSWORD"),
			Host:                 viper.GetString("MONGODB_HOST"),
			Parameter:            viper.GetString("MONGODB_PARAMETER"),
			Name:                 viper.GetString("MONGODB_NAME"),
			AuthMechanism:        mongo.AuthMechanism(viper.GetString("MONGODB_AUTH_MECHANISM")),
			Socks5Proxy:          viper.GetString("MONGODB_SOCKS5_PROXY"),
			AssumedRoleARN:       viper.GetString("AWS_ROLE_ARN"),
			SessionName:          viper.GetString("AWS_SESSION_NAME"),
			Region:               viper.GetString("AWS_REGION"),
			WebIdentityTokenFile: viper.GetString("AWS_WEB_IDENTITY_TOKEN_FILE"),
		},
	})

	for {
		<-errs
		os.Exit(-1)
	}
}

// ServerConfig contains the configuration for the HTTP server.
type ServerConfig struct {
	Version            string
	LogFilePath        string
	Host               string
	HTTPPort           string
	ExecuteAsDryRun    bool
	FeatureFlagsConfig featureflagsdk.Config
	MongoConfig        mongo.Config
}

func addressFromHostAndPort(host string, port string) string {
	return fmt.Sprint(host, ":", port)
}

func logServerRunning(address string) {
	log.Println("service running on " + address)
}

func startServer(c ServerConfig) chan error {
	bootstrap := api.NewBootstrap(api.BootstrapConfig{
		ServiceVersion: c.Version,
		LogWritter: io.MultiWriter(os.Stdout, &lumberjack.Logger{
			Filename:   c.LogFilePath,
			MaxSize:    50,
			MaxAge:     1,
			MaxBackups: 0,
			Compress:   false,
		}),
		ExecuteAsDryRun: c.ExecuteAsDryRun,
		MongoConfig:     c.MongoConfig,
		FeatureFlagsConfig: featureflagsdk.Config{
			SDKkey:          c.FeatureFlagsConfig.SDKkey,
			Offline:         c.FeatureFlagsConfig.Offline,
			FileDatasources: c.FeatureFlagsConfig.FileDatasources,
		},
	})

	mux := api.NewHandler(bootstrap)

	errs := make(chan error)

	go func() {
		address := addressFromHostAndPort(c.Host, c.HTTPPort)
		logServerRunning(address)
		if err := http.ListenAndServe(address, mux); err != nil {
			bootstrap.Logger.Error("EXITING! ", err.Error())
			errs <- err
		}
	}()

	return errs
}
