package dmm

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
	item := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{
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
	key := map[string]types.AttributeValue{
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

type DynamodbQueryApi interface {
	Query(
		ctx context.Context,
		params *dynamodb.QueryInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.QueryOutput, error)
}

func ListSlots(
	ctx context.Context,
	api DynamodbQueryApi,
	teacherId string,
) ([]Slot, error) {
	keyExpr := expression.Key("teacherId").Equal(expression.Value(teacherId))
	expr, err := expression.NewBuilder().WithKeyCondition(keyExpr).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression for query: %v", err)
	}
	res, err := api.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(scheduleTableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query slots: %v", err)
	}
	slots := []Slot{}
	err = attributevalue.UnmarshalListOfMaps(res.Items, &slots)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal slots: %v", err)
	}
	return slots, err
}

type DynamodbBatchWriteItemApi interface {
	BatchWriteItem(
		ctx context.Context,
		params *dynamodb.BatchWriteItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.BatchWriteItemOutput, error)
}

func AddSlots(
	ctx context.Context,
	api DynamodbBatchWriteItemApi,
	teacherId string,
	slots []Slot,
) error {
	// write 25 or less items at once
	maxItems := 25
	var j int
	for i := 0; i < len(slots); i += maxItems {
		j = i + maxItems
		if j > len(slots) {
			j = len(slots)
		}
		reqs := []types.WriteRequest{}
		for _, s := range slots[i:j] {
			item, err := attributevalue.MarshalMap(
				SlotWithTTL{
					TeacherId: s.TeacherId,
					DateTime:  s.DateTime,
					Ttl:       time.Now().Add(7 * 24 * time.Hour).Unix(),
				})
			if err != nil {
				return fmt.Errorf("failed to marshal slot: %s", err)
			}
			reqs = append(reqs, types.WriteRequest{PutRequest: &types.PutRequest{Item: item}})
		}
		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
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

func DeleteSlots(
	ctx context.Context,
	api DynamodbBatchWriteItemApi,
	teacherId string,
	slots []Slot,
) error {
	// write 25 or less items at once
	maxItems := 25
	var j int
	for i := 0; i < len(slots); i += maxItems {
		j = i + maxItems
		if j > len(slots) {
			j = len(slots)
		}
		reqs := []types.WriteRequest{}
		for _, s := range slots[i:j] {
			key, err := attributevalue.MarshalMap(s)
			if err != nil {
				return fmt.Errorf("failed to marshal slot: %s", err)
			}
			reqs = append(reqs, types.WriteRequest{DeleteRequest: &types.DeleteRequest{Key: key}})
		}
		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
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
