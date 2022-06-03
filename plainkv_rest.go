package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"

	std "github.com/eaglebush/stdutil"
	"github.com/narsilworks/plainkv"
)

func PlainKVRequestHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			mime string
			err  error
			b    []byte
		)

		// connection string
		dsi := sb.Settings.GetDatabaseInfo("DEFAULT")

		// Request variables
		vars := std.GetRequestVars(r, *sb.Settings.JWTSecret)
		key := vars.Variables.Key
		cmd := vars.Variables.FirstCommand()

		if key == "" {
			w.WriteHeader(500)
			w.Write([]byte("Please specify key"))
			return
		}

		pkv := plainkv.NewPlainKV(dsi.ConnectionString)
		err = pkv.Open()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("PlainKV client: %s", err)))
			return
		}

		defer pkv.Close()

		if vars.IsGet() {

			if cmd == "" || cmd == "mime" {
				mime, err = pkv.GetMime(key)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte(fmt.Sprintf("PlainKV client: %s", err)))
					return
				}
				w.Header().Set("Content-Type", mime)
			}

			if cmd == "" {
				b, err = pkv.Get(key)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte(fmt.Sprintf("PlainKV client: %s", err)))
					return
				}
				w.Write(b)
			}

			if cmd == "list" {
				s, err := pkv.ListKeys(key)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte(fmt.Sprintf("PlainKV client: %s", err)))
					return
				}

				buf := &bytes.Buffer{}
				gob.NewEncoder(buf).Encode(s)
				bs := buf.Bytes()

				w.Header().Set("Content-Type", "application/json")
				w.Write(bs)
			}
		}

		if vars.IsPostOrPut() {
			if cmd == "" {
				err = pkv.Set(key, vars.Body)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte(fmt.Sprintf("PlainKV client: %s", err)))
					return
				}
			}

			mime = r.Header.Get("Content-Type")

			err = pkv.SetMime(key, mime)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(fmt.Sprintf("PlainKV client: %s", err)))
				return
			}
		}

		if vars.IsDelete() {
			err = pkv.Del(key)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(fmt.Sprintf("PlainKV client: %s", err)))
				return
			}
		}
	})
}
