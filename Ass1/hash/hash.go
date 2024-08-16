package hash
import(
	"crypto/md5"
	"encoding/hex"
)

// Helper function to hash a password using MD5
func HashMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}