package main

import (
	"os"
	"context"
	"time"
	"fmt"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	// "github.com/prometheus/common/model"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "prom",
	Short: "prom",
	Long:  `Run Prometheus Queries`,
}


func run(query string, host string)  {
	// 1 create prometheus api
	ctx := context.Background()
	raw_client, err := api.NewClient(api.Config{
		Address: host,
	})

	if err != nil {
		panic(err)
	}

	client := v1.NewAPI(raw_client)

	result, err := client.Query(ctx, query, time.Now())

	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

func configPFlag(cmd *cobra.Command, name string, def string, description string) {
	cmd.PersistentFlags().String(name, def, description)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}

func init() {
	queryCmd := &cobra.Command{
		Use: "query",
		Short: "Query prom",
		Run: func(cmd *cobra.Command, args []string) {
			query := viper.GetString("query")
			host := viper.GetString("host")
			run(query, host)
		},
	}

	rootCmd.AddCommand(queryCmd)
}

func main() {
	viper.SetEnvPrefix("prom")
	viper.AutomaticEnv()

	configPFlag(rootCmd, "query", "prometheus_build_info{}", "the query")
	configPFlag(rootCmd, "host",  "http://localhost:9090", "Prom cluster")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
