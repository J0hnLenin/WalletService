package pgstorage

import "fmt"

func tableWithBucket(bucket bucketNum) string {
	return fmt.Sprintf("%s%d.%s", bucketPrefix, bucket, tableName)
}