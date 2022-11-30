package risingwave

import (
	"github.com/bamboovir/postgres/lib/risingwave"
	"github.com/spf13/cobra"
)

type RootArgs struct {
	Verbose     bool
	ConnStr     string
	ForceFlush  bool
	Random      bool
	QueryFactor float64
	InsertNum   int
}

func NewRootCMD() *cobra.Command {
	args := &RootArgs{}

	cmd := &cobra.Command{
		Use:   "risingwave-benchmark",
		Short: "risingwave-benchmark",
		Long:  "risingwave-benchmark",
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			benchmark, err := risingwave.New(args.ConnStr)
			benchmark.WithVerbose(args.Verbose)
			benchmark.WithForceFlush(args.ForceFlush)
			benchmark.WithInsertNum(args.InsertNum)
			benchmark.WithQueryFactor(args.QueryFactor)
			benchmark.WithRandom(args.Random)
			if err != nil {
				return err
			}
			err = benchmark.Benchmark()
			return err
		},
	}

	cmd.Flags().StringVar(&args.ConnStr, "conn-str", "", "connection string")
	cmd.MarkFlagRequired("conn-str")

	cmd.Flags().BoolVar(&args.Verbose, "verbose", false, "set verbose output")
	cmd.Flags().BoolVar(&args.Random, "random", false, "insert random data")
	cmd.Flags().BoolVar(&args.ForceFlush, "force-flush", false, "false flush")
	cmd.Flags().IntVar(&args.InsertNum, "insert-num", 10, "insert rows number")
	cmd.Flags().Float64Var(&args.QueryFactor, "query-factor", 1.0, "query factor")

	return cmd
}
