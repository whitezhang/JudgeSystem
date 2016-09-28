package encrypt

import (
	"crypto/md5"
	"encoding/hex"
)

func DoEncryption(pwd string) (encryptedPwd string) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(pwd))
	cipherStr := md5Ctx.Sum(nil)
	encryptedPwd = hex.EncodeToString(cipherStr)
	return
}
