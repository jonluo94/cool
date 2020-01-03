package password

import (
	"github.com/alexandrevicenzi/unchained"
	"github.com/jonluo94/cool/log"
)
var logger = log.GetLogger("password",log.ERROR)

const (
	Argon2Hasher       = "argon2"
	BCryptHasher       = "bcrypt"
	BCryptSHA256Hasher = "bcrypt_sha256"
	CryptHasher        = "crypt"
	MD5Hasher          = "md5"
	PBKDF2SHA1Hasher   = "pbkdf2_sha1"
	PBKDF2SHA256Hasher = "pbkdf2_sha256"
	SHA1Hasher         = "sha1"
	UnsaltedMD5Hasher  = "unsalted_md5"
	UnsaltedSHA1Hasher = "unsalted_sha1"

	DefaultHasher = PBKDF2SHA256Hasher
	DefaultSaltSize = 12
)

func Encode(password string,saltSize int,hasher string) string {
	hash, err := unchained.MakePassword(password, unchained.GetRandomString(saltSize), hasher)
	if err != nil {
		logger.Errorf("Error encoding password: %s\n", err)
	}
	return string(hash)
}

func Validate(password,cryto string) bool {
	valid, err := unchained.CheckPassword(password, cryto)
	if err != nil {
		logger.Errorf("Error decoding password: %s\n", err)
	}
	return valid
}