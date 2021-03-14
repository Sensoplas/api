package cmd

import (
	"github.com/Sensoplas/api/internal/uvindex"
	"github.com/fvbock/endless"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var port string

var httpServerCmd = &cobra.Command{
	Use: "http",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := zap.NewProductionConfig()
		switch loggingLevel {
		case "info":
			conf.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		case "debug":
			conf.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		case "warn":
			conf.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
		}
		logger, _ := conf.Build()
		defer func() {
			err := logger.Sync()
			if err != nil {
				panic(err)
			}
		}()
		logger = logger.With(zap.String("service", "sensoplas-api"))
		r := mux.NewRouter()
		r.Handle("/uvi-prediction", uvindex.MakeHTTPHandler(&uvindex.RNGService{}, logger))
		return endless.ListenAndServe("localhost:"+port, r)
	},
}

func init() {
	httpServerCmd.Flags().StringVarP(&port, "port", "p", "8080", "specifies the port for the http server to run on, defaults to 8080")
	rootCmd.AddCommand(httpServerCmd)
}
