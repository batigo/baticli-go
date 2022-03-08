package baticli

import (
	"crypto/md5"
	"fmt"

	"github.com/hashicorp/go-uuid"
)

func Genmsgid() string {
	s, _ := uuid.GenerateUUID()
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))[:12]
}
