package utils

import (
	"crypto/sha256"
	"fmt"

	"github.com/prakashjegan/review-processor/app/review-services/utils/uid64"
)

func GetUID() uint64 {
	uid, _ := uid64.New()
	return uint64(uid)

}

func GenerateChecksum(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}
