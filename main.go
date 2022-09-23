package main

import (
	"fmt"
	"net/http"

	"github.com/NarsilWorks-Inc/servicebase"
	std "github.com/eaglebush/stdutil"
)

var (
	sb *servicebase.ServiceBase
	use_valid_token,
	validate_token_times,
	validate_token_app bool
)

const APPLICATION_ID string = `PlainKV`

func main() {

	var (
		ok bool
	)

	sb, ok = servicebase.CreateService(
		std.NameValue[string]{Name: "id", Value: APPLICATION_ID},
		std.NameValue[string]{Name: "name", Value: "PlainKV Service"},
		std.NameValue[string]{Name: "version", Value: "1.0"},
	)
	if !ok {
		for _, m := range sb.GetMessages() {
			fmt.Println(m)
		}
	}

	if tmp := sb.Settings.Flag("debug").Bool(); tmp != nil {
		sb.ProductionMode = *tmp
	}

	if tmp := sb.Settings.Flag("use_valid_token").Bool(); tmp != nil {
		use_valid_token = *tmp
	}

	if tmp := sb.Settings.Flag("validate_token_times").Bool(); tmp != nil {
		validate_token_times = *tmp
	}

	if tmp := sb.Settings.Flag("validate_token_app").Bool(); tmp != nil {
		validate_token_app = *tmp
	}

	sb.SetMime(std.NameValue[string]{Name: ".txt", Value: "text/plain"})
	sb.Router.PathPrefix("/api/").Handler(PlainKVRequestHandler())

	sb.Router.Use(CORSMiddleware)
	if use_valid_token {
		sb.Router.Use(JWTMiddleware)
	} else {
		sb.DelayedLog(std.MsgWarn, "Token validation is disabled. Certain problems maybe encountered if this service depends on token payload data.")
	}

	sb.Serve()
}

// CORSMiddleware handles CORS request
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sb.PreflightRespond(&w, *r)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// JWTMiddleware handles the JWT token validity
func JWTMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		if sb.Settings.JWTSecret == nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`API Secret not set`))
			return
		}

		ji, err := std.ValidateJWT(r, *sb.Settings.JWTSecret, validate_token_times)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(fmt.Sprintf(`%s`, err)))
			return
		}

		if !ji.ValidAuthToken {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`Invalid access token`))
			return
		}

		if validate_token_app {
			if ji.TokenApplicationID != APPLICATION_ID {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`Invalid application from token`))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
