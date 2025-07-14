// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package main

import "net/http"

func main() {

}

func ValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API_Key") == "" {
			http.Error(w, "Missing API Key", http.StatusUnauthorized)
		}
		next.ServeHTTP(w, r)
	})
}
