package filter

import (
	"testing"
)

var caseStr = "abcde"

func TestRSHash(t *testing.T) {
	println(rsHash([]byte(caseStr)))
}

func TestJSHash(t *testing.T) {
	println(jsHash([]byte(caseStr)))
}

func TestELFHash(t *testing.T) {
	println(elfHash([]byte(caseStr)))
}

func TestBKDRHash(t *testing.T) {
	println(bkdRHash([]byte(caseStr)))
}

func TestAPHash(t *testing.T) {
	println(apHash([]byte(caseStr)))
}

func TestDJBHash(t *testing.T) {
	println(djbHash([]byte(caseStr)))
}

func TestSDBMHash(t *testing.T) {
	println(sdbMHash([]byte(caseStr)))
}

func TestPJWHash(t *testing.T) {
	println(pjwHash([]byte(caseStr)))
}