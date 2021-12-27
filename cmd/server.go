/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nicolastakashi/cole/internal/cole"
	"github.com/nicolastakashi/cole/internal/command"
	"github.com/nicolastakashi/cole/internal/k8sclient"
	"github.com/nicolastakashi/cole/internal/logging"
	"github.com/nicolastakashi/cole/internal/loghandler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if err := logging.Configure(scmd.LogLevel); err != nil {
			os.Exit(1)
		}

		logrus.Info("Welcome to cole...")

		ctx, cancel := context.WithCancel(context.Background())
		wg, ctx := errgroup.WithContext(ctx)

		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)

		srv := createHttpServer(serverPort)

		err := prometheus.DefaultRegisterer.Register(version.NewCollector("cole"))

		if err != nil {
			logrus.Errorf("error to register version collector %v", err)
		}

		logrus.Info("listen on " + serverPort)

		wg.Go(func() error {
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				logrus.Errorf("http server error %v", err)
				return err
			}
			return nil
		})

		wg.Go(func() error {
			client, err := k8sclient.New(ctx, scmd.KubeConfig)

			if err != nil {
				return err
			}

			lastSinceTime := time.Now()
			cole := cole.Cole{
				Ctx:           ctx,
				Scmd:          *scmd,
				Client:        client,
				LogHandler:    loghandler.New(*scmd),
				LastSinceTime: &lastSinceTime,
				Timer:         time.NewTimer(1 * time.Millisecond),
			}

			if err := cole.Start(); err != nil {
				return err
			}
			return nil
		})

		select {
		case <-term:
			logrus.Info("received SIGTERM, exiting gracefully...")
		case <-ctx.Done():
		}

		if err := srv.Shutdown(ctx); err != nil {
			logrus.Errorf("server shutdown error %v", err)
		}

		cancel()

		if err := wg.Wait(); err != nil {
			logrus.Errorf("unhandled error received. Exiting... %v", err)
			os.Exit(1)
		}

		os.Exit(0)

	},
}

var serverPort = ""
var scmd = &command.Server{}

func init() {
	serverCmd.Flags().StringVar(&serverPort, "http.port", ":9754", "listem port for http endpoints")
	serverCmd.Flags().StringVar(&scmd.LogLevel, "log.level", logrus.InfoLevel.String(), "listem port for http endpoints")
	serverCmd.Flags().StringVar(&scmd.KubeConfig, "kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	serverCmd.Flags().StringVar(&scmd.Namespace, "namespace", "default", "namespace where Grafana is running")
	serverCmd.Flags().BoolVar(&scmd.IncludeUname, "metrics.includeUname", false, "Include user name to metrics (disabled by default)")
	serverCmd.Flags().StringVar(&scmd.LabelSelector, "grafana.podLabelselector", "", "Grafana pod label selector")
	serverCmd.Flags().StringVar(&scmd.Container, "grafana.containerName", "grafana", "Grafana container name (default grafana)")
	serverCmd.Flags().StringVar(&scmd.LogFormat, "grafana.log.format", "", "Grafana pod log format")
	rootCmd.AddCommand(serverCmd)
}

func createHttpServer(port string) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/-/health", func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/json")

		resp := map[string]string{
			"message": "Healthy",
		}

		jsonResp, err := json.Marshal(resp)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}

		rw.Write(jsonResp)
	})
	mux.HandleFunc("/-/ready", func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/json")

		resp := map[string]string{
			"message": "Ready",
		}

		jsonResp, err := json.Marshal(resp)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}

		rw.Write(jsonResp)
	})

	srv := &http.Server{
		Addr:     port,
		Handler:  mux,
		ErrorLog: &log.Logger{},
	}
	return srv
}
