package registry

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/hxia043/nest/pkg/util/crypto"
)

const argsLimitCount int = 1

func readCertFrom(caCertPath string) string {
	plaintext, err := os.ReadFile(caCertPath)
	if err != nil {
		fmt.Fprint(os.Stderr, "Warning: "+err.Error()+"\n")
		return ""
	}

	ciphertext := make([]byte, len(plaintext))
	crypto.EncryptAES(ciphertext, plaintext, crypto.Key, crypto.Iv)

	return base64.StdEncoding.EncodeToString(ciphertext)
}

func getRegistryName(args []string) (string, error) {
	if len(args) != argsLimitCount {
		return "", fmt.Errorf("incorrect registry number, the required number of registry is %d", argsLimitCount)
	}

	return args[0], nil
}
