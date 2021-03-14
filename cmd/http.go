package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
		// we ignore zap logger sync error
		//nolint: errcheck
		defer logger.Sync()
		logger = logger.With(zap.String("service", "sensoplas-api"))
		r := mux.NewRouter()
		r.Handle("/uvi-prediction", uvindex.MakeHTTPHandler(&uvindex.RNGService{}, logger))

		// This is very much not needed because of endless. Will catch signals either way.
		errs := make(chan error, 2)
		go func() {
			logger.Info("server starting", zap.String("transport", "http"))
			errs <- endless.ListenAndServe("localhost:"+port, r)
		}()
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT)
			errs <- fmt.Errorf("%s", <-c)
		}()

		logger.Info("terminated", zap.Error(<-errs))

		return nil
	},
}

func init() {
	httpServerCmd.Flags().StringVarP(&port, "port", "p", "8080", "specifies the port for the http server to run on, defaults to 8080")
	rootCmd.AddCommand(httpServerCmd)
}
