package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
)

var (
	rootFlagSet = flag.NewFlagSet("sqs-pub", flag.ExitOnError)
	cfg         = SQSMessageReplayConfig{}
	replayer    *SQSMessageReplayer
)

func init() {
	rootFlagSet.StringVar(&cfg.region, "region", "eu-west-1", "AWS region")
	rootFlagSet.StringVar(&cfg.from, "from", "queue-name-source", "sqs queue from where messages will be sourced from")
	rootFlagSet.StringVar(&cfg.to, "to", "queue-name-destination", "sqs queue where messages will be pushed to")
	rootFlagSet.StringVar(&cfg.filters, "filters", "10104211111292", "comma separted text that can be used a message body filter")
	rootFlagSet.BoolVar(&cfg.deleteFromSource, "delete", true, "delete messages from source after successfuly pushed to destination queue")
	rootFlagSet.BoolVar(&cfg.dryrun, "dryrun", false, "a flag to run the replay in dry run mode.")

	replayer = NewSQSMessageReplayer(&cfg)
}

func main() {
	root := &ffcli.Command{
		Name:       "replay",
		ShortUsage: "sqs_pub [-from queue1 - to queue2 -filter text1,text2,...]",
		ShortHelp:  "Source message from the given queue then push to the destination queue",
		FlagSet:    rootFlagSet,
		Exec:       replayer.replay,
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

}
