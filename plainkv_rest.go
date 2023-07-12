package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	std "github.com/eaglebush/stdutil"
	"github.com/narsilworks/plainkv"
)

func PlainKVRequestHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			mime string
			err  error
			b    []byte
			tlly int
		)

		// connection string
		dsi := sb.Settings.GetDatabaseInfo("DEFAULT")

		// Request variables
		vars, _ := std.GetRequestVars(r, *sb.Settings.JWTSecret, validate_token_times)
		key := vars.Variables.Key
		cmd := vars.Variables.FirstCommand()

		qs := vars.Variables.QueryString
		bucket, _ := qs.String("bucket")
		offset, _ := qs.Int("offset")

		writeError := func(err error) {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("PlainKV server: %s", err)))
		}

		if key == "" {
			writeError(fmt.Errorf("please specify key"))
			return
		}

		pkv := plainkv.NewPlainKV(dsi.ConnectionString, false)
		if bucket != "" {
			pkv.SetBucket(bucket)
		}

		if err = pkv.Open(); err != nil {
			writeError(err)
			return
		}
		defer pkv.Close()

		// Get method
		if vars.IsGet() {

			if cmd == "" {
				if b, err = pkv.Get(key); err != nil {
					writeError(err)
					return
				}
				w.Write(b)
			}

			if cmd == "" || cmd == "mime" {
				if mime, err = pkv.GetMime(key); err != nil {
					writeError(err)
					return
				}
				w.Header().Set("Content-Type", mime)

				return
			}

			if cmd == "list" {
				s, err := pkv.ListKeys(key)
				if err != nil {
					writeError(err)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				if len(s) == 0 {
					w.Write([]byte("[]"))
					return
				}

				bs, err := json.Marshal(s)
				if err != nil {
					writeError(err)
					return
				}

				w.Write(bs)
				return
			}

			if cmd == "tally" {
				if tlly, err = pkv.Tally(key, offset); err != nil {
					writeError(err)
					return
				}

				w.Write([]byte(strconv.Itoa(tlly)))
				return
			}
		}

		if vars.IsPostOrPut() {
			if vars.IsPost() {
				if cmd == "" {
					if err = pkv.Set(key, vars.Body); err != nil {
						writeError(err)
						return
					}
				}

				if cmd == "" || cmd == "mime" {
					mime = r.Header.Get("Content-Type")
					if mime != "" {
						if err = pkv.SetMime(key, mime); err != nil {
							writeError(err)
						}
					}

					return
				}
			}

			if vars.IsPut() {
				if cmd == "incr" {
					if tlly, err = pkv.Incr(key); err != nil {
						writeError(err)
						return
					}

					w.Write([]byte(strconv.Itoa(tlly)))
					return
				}

				if cmd == "decr" {
					if tlly, err = pkv.Decr(key); err != nil {
						writeError(err)
						return
					}

					w.Write([]byte(strconv.Itoa(tlly)))
					return
				}
			}
		}

		if vars.IsDelete() {
			if cmd == "" {
				if err = pkv.Del(key); err != nil {
					writeError(err)
				}
				return
			}

			if cmd == "tally" {
				if err = pkv.Reset(key); err != nil {
					writeError(err)
				}
				return
			}
		}
	})
}
