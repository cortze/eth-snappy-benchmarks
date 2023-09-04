package main

import (
    "os"
    "context"

	"github.com/cortze/eth-snappy-benchmarks/cmd"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
    snappyBenchmark := &cli.App{
        Name:   "snappy-benchmark",
        Usage:  "executes the Snappy compression benchmarks with the blocks provided in a given directory",
        Commands: []*cli.Command{
            cmd.RunCmd,
            cmd.FetchBlocksCmd,
        },
    }
    err := snappyBenchmark.RunContext(context.Background(), os.Args);
    if err != nil {
        log.Errorf("error: %v\n", err)
        os.Exit(1)
    }
    os.Exit(0)
}
