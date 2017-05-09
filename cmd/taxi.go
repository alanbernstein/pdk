package cmd

import (
	"io"
	"log"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pilosa/pdk/usecase/taxi"
	"github.com/spf13/cobra"
	"github.com/uber/jaeger-client-go/config"
)

var TaxiMain *taxi.Main

func configureTracer() {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	tracer, _, _ := cfg.New("pdk-taxi")
	opentracing.SetGlobalTracer(tracer)
}

func NewTaxiCommand(stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	configureTracer()
	TaxiMain = taxi.NewMain()
	taxiCommand := &cobra.Command{
		Use:   "taxi",
		Short: "taxi - import taxi data to pilosa",
		Long:  `TODO`,
		RunE: func(cmd *cobra.Command, args []string) error {
			start := time.Now()
			err := TaxiMain.Run()
			if err != nil {
				return err
			}
			log.Println("Done: ", time.Since(start))
			select {}
		},
	}
	flags := taxiCommand.Flags()
	flags.IntVarP(&TaxiMain.Concurrency, "concurrency", "c", 1, "Number of goroutines fetching and parsing")
	flags.IntVarP(&TaxiMain.BufferSize, "buffer-size", "b", 10000000, "Size of buffer for importers - heavily affects memory usage")
	flags.StringVarP(&TaxiMain.PilosaHost, "pilosa", "p", "localhost:10101", "Pilosa host")
	flags.StringVarP(&TaxiMain.Index, "index", "i", "taxi", "Pilosa db to write to")
	flags.StringVarP(&TaxiMain.URLFile, "url-file", "f", "usecase/taxi/urls-short.txt", "File to get raw data urls from. Urls may be http or local files.")

	return taxiCommand
}

func init() {
	subcommandFns["taxi"] = NewTaxiCommand
}
