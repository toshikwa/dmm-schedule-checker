package dmm

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	teacherTableName  = os.Getenv("TEACHER_TABLE_NAME")
	scheduleTableName = os.Getenv("SCHEDULE_TABLE_NAME")
)

type Teacher struct {
	Id string `json:"id" dynamodbav:"id"`
}

type DynamodbPutItemApi interface {
	PutItem(
		ctx context.Context,
		params *dynamodb.PutItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.PutItemOutput, error)
}

func AddTeacher(
	ctx context.Context,
	api DynamodbPutItemApi,
	teacherId string,
) error {
	item := map[string]ddbTypes.AttributeValue{
		"id": &ddbTypes.AttributeValueMemberS{
			Value: teacherId,
		},
	}
	_, err := api.PutItem(ctx, &dynamodb.PutItemInput{Item: item, TableName: aws.String(teacherTableName)})
	if err != nil {
		return fmt.Errorf("failed to put item: %s", err)
	}
	return nil
}

type DynamodbDeleteItemApi interface {
	DeleteItem(
		ctx context.Context,
		params *dynamodb.DeleteItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.DeleteItemOutput, error)
}

func DeleteTeacher(
	ctx context.Context,
	api DynamodbDeleteItemApi,
	teacherId string,
) error {
	key := map[string]ddbTypes.AttributeValue{
		"id": &types.AttributeValueMemberS{
			Value: teacherId,
		},
	}
	_, err := api.DeleteItem(ctx, &dynamodb.DeleteItemInput{Key: key, TableName: aws.String(teacherTableName)})
	if err != nil {
		return fmt.Errorf("failed to put item: %s", err)
	}
	return nil
}

type DynamodbScanApi interface {
	Scan(
		ctx context.Context,
		params *dynamodb.ScanInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.ScanOutput, error)
}

func ListTeachers(
	ctx context.Context,
	api DynamodbScanApi,
) ([]Teacher, error) {
	res, err := api.Scan(ctx, &dynamodb.ScanInput{TableName: aws.String(teacherTableName)})
	teachers := []Teacher{}
	var t Teacher
	for _, item := range res.Items {
		_ = attributevalue.UnmarshalMap(item, &t)
		teachers = append(teachers, t)
	}
	return teachers, err
}

type DynamodbBatchWriteItemApi interface {
	BatchWriteItem(
		ctx context.Context,
		params *dynamodb.BatchWriteItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.BatchWriteItemOutput, error)
}

func AddNewSlots(
	ctx context.Context,
	api DynamodbBatchWriteItemApi,
	teacherId string,
	slots []Slot,
) error {
	// write 25 or less items at once
	maxItems := 2
	var j int
	for i := 0; i < len(slots); i += maxItems {
		j = i + maxItems
		if j > len(slots) {
			j = len(slots)
		}
		reqs := []ddbTypes.WriteRequest{}
		for _, s := range slots[i:j] {
			item, err := attributevalue.MarshalMap(s)
			if err != nil {
				return fmt.Errorf("failed to marshal slot: %s", err)
			}
			reqs = append(reqs, ddbTypes.WriteRequest{PutRequest: &ddbTypes.PutRequest{Item: item}})
		}
		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]ddbTypes.WriteRequest{
				scheduleTableName: reqs,
			},
		}
		_, err := api.BatchWriteItem(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to batch write: %v", err)
		}
	}
	return nil
}

type DynamodbBatchGetItemApi interface {
	BatchGetItem(
		ctx context.Context,
		params *dynamodb.BatchGetItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.BatchGetItemOutput, error)
}

func ListExistingSlots(
	ctx context.Context,
	api DynamodbBatchGetItemApi,
	teacherId string,
	slots []Slot,
) ([]Slot, error) {
	keys := []map[string]ddbTypes.AttributeValue{}
	for _, s := range slots {
		item, err := attributevalue.MarshalMap(s)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal slot: %s", err)
		}
		keys = append(keys, item)
	}
	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]ddbTypes.KeysAndAttributes{
			scheduleTableName: {Keys: keys},
		},
	}
	res, err := api.BatchGetItem(ctx, input)
	exists := []Slot{}
	var s Slot
	for _, item := range res.Responses[scheduleTableName] {
		_ = attributevalue.UnmarshalMap(item, &s)
		exists = append(exists, s)
	}
	return exists, err
}
