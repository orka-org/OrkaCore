package tokens

import (
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
)

func TokenFromHeader(header string) string {
	splits := strings.Split(header, " ")
	if len(splits) != 2 {
		return ""
	}
	return splits[1]
}

func ValidateAuthorizationHeader(header string) (string, error) {
	if header == "" {
		return "", errors.Unauthorized("Authorization Header Missing", "The Authorization header is missing")
	}

	if !strings.HasPrefix(header, "Bearer ") {
		errMsg := fmt.Sprintf("The Authorization header is invalid, expected Bearer <token>, got %s", header)
		return "", errors.Unauthorized("Invalid Authorization Header", errMsg)
	}
	token := TokenFromHeader(header)
	if token == "" {
		return "", errors.Unauthorized("Invalid Authorization Header", "The Authorization header is invalid, got empty token")
	}
	return token, nil
}
