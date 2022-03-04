package cwlogger

import (
	aws "bott-the-pigeon/lib/aws/session"

	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/google/uuid"
)

// Parameters used to create a CloudWatch log.
type CloudWatch_Log struct {
	LogGroup string
	LogStream string
	Message string
}

// The CloudWatch client pointer is stored, and can be accessed later.
var cwsvc *cloudwatchlogs.CloudWatchLogs

// Returns the stored CloudWatch client or creates one if it doesn't exist,
// using the provided AWS session.
func getClient(awssess *session.Session) (*cloudwatchlogs.CloudWatchLogs) {
	if cwsvc != nil {
		return cwsvc
	} else {
		ssmsvc := cloudwatchlogs.New(awssess)
		return ssmsvc
	}
}

// Returns a UUID and creates a correlating log.
func Log(logGroup string, message string) (*string, error) {
	awssess, err := aws.GetSession()
	if err != nil {
		return nil, err
	}
	uuid := uuid.New().String()
	params := &CloudWatch_Log{
		LogGroup: logGroup,
		LogStream: uuid,
		Message: message,
	}
	err = createCWLog(awssess, params)
	if err != nil {
		return nil, err
	}
	return &uuid, nil
}

// Creates a CloudWatch log with the provided CloudWatch_Log parameters.
func createCWLog(awssess *session.Session, params *CloudWatch_Log) error {
	cwsvc := getClient(awssess)
	err := ensureLogGroupExists(cwsvc, params.LogGroup)
	if err != nil {
		return err
	}
	_, err = cwsvc.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{
		LogGroupName: &params.LogGroup,
		LogStreamName: &params.LogStream,
	})
	if err != nil {
		return fmt.Errorf("failed to create log stream; cannot create CloudWatch log: %v", err)
	}
	timestamp := time.Now().Unix() * 1000
	event := &cloudwatchlogs.InputLogEvent{
		Message: &params.Message,
		Timestamp: &timestamp,
	}
	events := make([]*cloudwatchlogs.InputLogEvent, 1)
	events[0] = event
	_, err = cwsvc.PutLogEvents(&cloudwatchlogs.PutLogEventsInput{
		LogEvents: events,
		LogGroupName: &params.LogGroup,
		LogStreamName: &params.LogStream,
	})
	if err != nil {
		return fmt.Errorf("failed to log to event stream: %v", err)
	}
	return nil
}

// Checks that the provided log group exists, and creates one if not.
func ensureLogGroupExists(cwsvc *cloudwatchlogs.CloudWatchLogs, name string) error {
	var limit int64 = 1
	lgs, err := cwsvc.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{
		Limit: &limit,
		LogGroupNamePrefix: &name,
	})
	if err != nil {
		return fmt.Errorf("failed to list existing log groups; cannot ensure log group exists: %v", err)
	}
	if len(lgs.LogGroups) < 1 {
		err := createLogGroup(cwsvc, name)
		if err != nil {
			return err
		}
	}
	return nil
}

// Creates a CloudWatch log group using the service and the log group name provided.
func createLogGroup(cwsvc *cloudwatchlogs.CloudWatchLogs, name string) error {
	_, err := cwsvc.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: &name,
	})
	if err != nil {
		return fmt.Errorf("failed to create log group: %v", err)
	}
	return nil
}
