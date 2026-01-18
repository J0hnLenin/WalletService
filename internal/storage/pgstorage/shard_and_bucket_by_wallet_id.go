package pgstorage

import (
	"hash/fnv"

	"github.com/google/uuid"
)

func (pg *PGStorage) shardAndBucketByWalletID(id uuid.UUID) (int, bucketNum) {
	h := fnv.New32a()
	h.Write(id[:]) 
	hash := h.Sum32()

	shardIndex := (hash) % uint32(len(pg.shards))

	bucketPerShard := uint32(pg.numberOfBuckets) / uint32(len(pg.shards))
	bucketIndex := bucketNum(shardIndex * bucketPerShard + (hash % bucketPerShard))

	return int(shardIndex), bucketIndex
}