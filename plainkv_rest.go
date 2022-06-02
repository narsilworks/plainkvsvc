package main

import (
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

		res := std.ResultAny{
			Result: std.InitResult(),
		}

		// connection string
		dsi := sb.Settings.GetDatabaseInfo("DEFAULT")

		// Request variables
		vars := std.GetRequestVars(r, *sb.Settings.JWTSecret)
		key := vars.Variables.Key
		cmd := vars.Variables.FirstCommand()

		if key == "" {
			res.AddError("Please specify key")
			sb.Respond(&res, &w, *r)
			return
		}

		pkv := plainkv.NewPlainKV(dsi.ConnectionString)
		err = pkv.Open()
		if err != nil {
			res.AddErrorf("PlainKV client: %s", err)
			sb.Respond(&res, &w, *r)
			return
		}

		defer pkv.Close()

		if vars.IsGet() {

			if cmd == "" {
				b, err = pkv.Get(key)
				if err != nil {
					res.AddErrorf("PlainKV client: %s", err)
					sb.Respond(&res, &w, *r)
					return
				}
				res.Data = b
			}

			if cmd == "" || cmd == "mime" {
				mime, err = pkv.GetMime(key)
				if err != nil {
					res.AddErrorf("PlainKV client: %s", err)
					sb.Respond(&res, &w, *r)
					return
				}
			}

			if cmd == "list" {
				s, err := pkv.ListKeys(key)
				if err != nil {
					res.AddErrorf("PlainKV client: %s", err)
					sb.Respond(&res, &w, *r)
					return
				}
				res.Data = s
			}

			if mime == "application/json" {
				res.Return(std.OK)
				sb.Respond(&res, &w, *r)
			} else {
				sb.RespondBytes(b, ".txt", &w, *r)
			}
		}

		if vars.IsPostOrPut() {
			if cmd == "" {
				err = pkv.Set(key, vars.Body)
				if err != nil {
					res.AddErrorf("PlainKV client: %s", err)
					sb.Respond(&res, &w, *r)
					return
				}
			}

			if cmd == "" {
				mime = r.Header.Get("Content-Type")
			} else {
				mime = string(vars.Body)
			}

			err = pkv.SetMime(key, mime)
			if err != nil {
				res.AddErrorf("PlainKV client: %s", err)
				sb.Respond(&res, &w, *r)
				return
			}

			res.Return(std.OK)
			sb.Respond(&res, &w, *r)
		}

		if vars.IsDelete() {

			err = pkv.Del(key)
			if err != nil {
				res.AddErrorf("PlainKV client: %s", err)
				sb.Respond(&res, &w, *r)
				return
			}

			res.Return(std.OK)
			sb.Respond(&res, &w, *r)
		}
	})
}
