package common

import (
	"github.com/HaoyuHu/gosimhash"
	"github.com/HaoyuHu/gosimhash/utils"
)

type HashType uint8

const (
	SipHash     = iota
	JenkinsHash
)

var simhasher *gosimhash.Simhasher

func InitializeSimhasher(hashType HashType, dict string, hmm string, userDict string, idf string, stopWords string) {
	var inner utils.Hasher
	switch hashType {
	case SipHash:
		inner = utils.NewSipHasher([]byte(gosimhash.DEFAULT_HASH_KEY))
		break
	case JenkinsHash:
		inner = utils.NewJenkinsHasher()
		break
	default:
		inner = utils.NewJenkinsHasher()
	}
	simhasher = gosimhash.NewSimhasher(inner, dict, hmm, userDict, idf, stopWords)
}

func MakeSimhash(doc *string, topN int) uint64 {
	return simhasher.MakeSimhash(doc, topN)
}

func Diff(simhash uint64, another uint64, limit int) bool {
	distance := gosimhash.CalculateDistanceBySimhash(simhash, another)
	// different documents when distance > limit
	return distance > limit
}

func Distance(simhash uint64, another uint64) int {
	return gosimhash.CalculateDistanceBySimhash(simhash, another)
}

func Free() {
	simhasher.Free()
}
