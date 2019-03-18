package handler

import (
	"context"
	"net/http"

	"azure.com/ecovo/trip-service/cmd/middleware/auth"
)

// Auth validates a request's authorization header using the given validator
// to ensure that the user is authorized to access an endpoint and extracts the
// authenticated user's information.
//
// The authenticated user's information placed in the request's context and can
// be accessed by using the auth.FromContext utility function.
func Auth(validators []auth.Validator, next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		header := r.Header.Get("Authorization")

		var err error
		var userInfo *auth.UserInfo
		for _, v := range validators {
			userInfo, err = v.Validate(header)
			if err == nil {
				break
			}
		}

		if err != nil {
			return err
		}

		ctx := context.WithValue(r.Context(), auth.UserInfoContextKey, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))

		return nil
	}
}
