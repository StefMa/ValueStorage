package apicommon

import (
	"errors"
	"net/http"
)

const authorizationToken = "Token ABC123"

func CheckToken(r *http.Request) error {
	authorization := r.Header.Get("Authorization")
	if authorization != authorizationToken {
		return errors.New("Authorization failed")
	}
	return nil
}
