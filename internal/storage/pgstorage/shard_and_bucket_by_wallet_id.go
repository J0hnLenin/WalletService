package pgstorage

import (
	"hash/fnv"

	"github.com/google/uuid"
)

func (pg *PGStorage) shardAndBucketByWalletID(id uuid.UUID) (int, bucketNum) {
	h := fnv.New32a()
	h.Write(id[:]) 
	hash := h.Sum32()
	bucketIndex := bucketNum(hash % uint32(pg.numberOfBuckets))
	shardIndex := int(bucketIndex) % len(pg.shards)

	return shardIndex, bucketIndex
}