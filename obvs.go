package main

import (
	"fmt"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

var (
	region          = kingpin.Flag("region", "AWS region.").Default("us-east-1").OverrideDefaultFromEnvar("AWS_REGION").String()
	accessKeyId     = kingpin.Flag("access-key-id", "AWS access key ID.").Required().OverrideDefaultFromEnvar("AWS_ACCESS_KEY_ID").String()
	secretAccessKey = kingpin.Flag("secret-access-key", "AWS secret access key.").Required().OverrideDefaultFromEnvar("AWS_SECRET_ACCESS_KEY").String()

	groups = kingpin.Command("groups", "List log groups.")

	streams      = kingpin.Command("streams", "List log streams.")
	streamsGroup = streams.Arg("group", "Log group from which to list streams.").Required().String()

	events          = kingpin.Command("events", "List log events.").Default()
	eventsGroup     = events.Arg("group", "Log group from which to list events.").Required().String()
	eventsPattern   = events.Arg("pattern", "Filter events matching this pattern.").String()
	eventsStartTime = events.Flag("since", "Filter events since (gte) this time.").String()
	eventsEndTime   = events.Flag("until", "Filter events until (lte) this time.").String()
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("0.1")
	kingpin.CommandLine.Help = "CLI for AWS CloudWatch Logs."
	cmd := kingpin.Parse()

	creds := credentials.NewStaticCredentials(*accessKeyId, *secretAccessKey, "")
	config := session.New(&aws.Config{Region: region, Credentials: creds})
	client := cloudwatchlogs.New(config)

	switch cmd {
	case "groups":
		kingpin.FatalIfError(cmdGroups(client), "List log groups")
	case "streams":
		kingpin.FatalIfError(cmdStreams(client, *streamsGroup), "List log streams")
	case "events":
		kingpin.FatalIfError(cmdEvents(client, *eventsGroup, *eventsPattern, *eventsStartTime, *eventsEndTime), "List log events")
	}
}

func cmdGroups(client *cloudwatchlogs.CloudWatchLogs) error {
	req := cloudwatchlogs.DescribeLogGroupsInput{}

	handler := func(res *cloudwatchlogs.DescribeLogGroupsOutput, lastPage bool) bool {
		for _, group := range res.LogGroups {
			name := group.LogGroupName
			fmt.Printf("%s\n", *name)
		}

		return true // want more pages
	}

	err := client.DescribeLogGroupsPages(&req, handler)
	if err != nil {
		return err
	}

	return nil
}

func cmdStreams(client *cloudwatchlogs.CloudWatchLogs, group string) error {
	req := cloudwatchlogs.DescribeLogStreamsInput{LogGroupName: &group}

	handler := func(res *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) bool {
		for _, stream := range res.LogStreams {
			name := stream.LogStreamName
			fmt.Printf("%s\n", *name)
		}

		return true // want more pages
	}

	err := client.DescribeLogStreamsPages(&req, handler)
	if err != nil {
		return err
	}

	return nil
}

func parseTime(timeStr string) time.Time {
	loc, _ := time.LoadLocation("UTC")
	const timeFmt = "2006-01-02T15:04:05"

	t, _ := time.ParseInLocation(timeFmt, timeStr, loc)

	return t
}

func cmdEvents(client *cloudwatchlogs.CloudWatchLogs, group string, pattern string, startTime string, endTime string) error {
	interleaved := true
	req := cloudwatchlogs.FilterLogEventsInput{LogGroupName: &group, Interleaved: &interleaved}

	if pattern != "" {
		req.FilterPattern = &pattern
	}

	if startTime != "" {
		startTimeInt64 := parseTime(startTime).Unix()
		req.StartTime = &startTimeInt64
	}

	if endTime != "" {
		endTimeInt64 := parseTime(endTime).Unix()
		req.EndTime = &endTimeInt64
	}

	handler := func(res *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
		for _, event := range res.Events {
			message := event.Message
			streamName := event.LogStreamName
			eventId := event.EventId
			fmt.Printf("{%s/%s} %s\n", *streamName, *eventId, *message)
		}

		return true // want more pages
	}

	err := client.FilterLogEventsPages(&req, handler)
	if err != nil {
		return err
	}

	return nil
}

// API docs
// http://docs.aws.amazon.com/AmazonCloudWatch/latest/DeveloperGuide/FilterAndPatternSyntax.html
// http://docs.aws.amazon.com/AmazonCloudWatchLogs/latest/APIReference/API_FilterLogEvents.html
// http://docs.aws.amazon.com/sdk-for-go/api/service/cloudwatchlogs.html#type-FilterLogEventsInput
// http://golang.org/src/time/format.go
