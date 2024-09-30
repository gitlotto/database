package database

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const compositeRecordsTableName = "gitlotto.compositeRecords"

type CompositeRecord struct {
	PartitionKey string `dynamodbav:"partition_key"`
	SortKey      string `dynamodbav:"sort_key"`
}

func (record CompositeRecord) Id() RecordIdentifier {
	return RecordIdentifier{
		TableName: compositeRecordsTableName,
		PrimaryKey: PrimaryKey{
			PartitionKey: DynamodbKey{
				Name:  "partition_key",
				Value: record.PartitionKey,
			},
			SortKey: &DynamodbKey{
				Name:  "sort_key",
				Value: record.SortKey,
			},
		},
	}
}

func Test_CompositeRecord_should_be_stored_in_correct_form(t *testing.T) {
	var err error

	partitionKey := uuid.New().String()
	sortKey := uuid.New().String()

	record := CompositeRecord{
		PartitionKey: partitionKey,
		SortKey:      sortKey,
	}

	err = Persist(dynamodbClient, record)
	assert.NoError(t, err)

	getEventOutput, err := dynamodbClient.GetItem(
		&dynamodb.GetItemInput{
			TableName: aws.String(compositeRecordsTableName),
			Key: map[string]*dynamodb.AttributeValue{
				"partition_key": {
					S: aws.String(partitionKey),
				},
				"sort_key": {
					S: aws.String(sortKey),
				},
			},
		},
	)

	assert.NoError(t, err)
	actualItems := getEventOutput.Item

	expectedItems := map[string]*dynamodb.AttributeValue{
		"partition_key": {
			S: aws.String(partitionKey),
		},
		"sort_key": {
			S: aws.String(sortKey),
		},
	}

	assert.Equal(t, expectedItems, actualItems)

	expectedRecord := record

	actualRecord := CompositeRecord{}
	err = dynamodbattribute.UnmarshalMap(getEventOutput.Item, &actualRecord)

	assert.NoError(t, err)
	assert.Equal(t, expectedRecord, actualRecord)

}

func Test_CompositeRecord_should_be_fetch_in_correct_form(t *testing.T) {
	var err error

	partitionKey := uuid.New().String()
	sortKey := uuid.New().String()

	record := CompositeRecord{
		PartitionKey: partitionKey,
		SortKey:      sortKey,
	}

	expectedItems := map[string]*dynamodb.AttributeValue{
		"partition_key": {
			S: aws.String(partitionKey),
		},
		"sort_key": {
			S: aws.String(sortKey),
		},
	}

	actualItems, err := dynamodbattribute.MarshalMap(record)
	assert.NoError(t, err)
	assert.Equal(t, expectedItems, actualItems)

	_, err = dynamodbClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(compositeRecordsTableName),
		Item:      actualItems,
	})
	assert.NoError(t, err)

	actualRecord, err := Fetch(dynamodbClient, record)
	assert.NotNil(t, actualRecord)
	assert.NoError(t, err)

	expectedRecord := record
	assert.Equal(t, expectedRecord, *actualRecord)

}

func Test_Fetch_should_return_nil_if_the_composite_record_does_not_exist(t *testing.T) {
	var err error

	partitionKey := uuid.New().String()
	sortKey := uuid.New().String()

	record := CompositeRecord{
		PartitionKey: partitionKey,
		SortKey:      sortKey,
	}

	actualRecord, err := Fetch(dynamodbClient, record)
	assert.NoError(t, err)
	assert.Nil(t, actualRecord)

}
