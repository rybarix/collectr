// Copyright 2024 Sandro Rybárik

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 		http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collectr

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type HttpHandler struct {
	Logger    *slog.Logger
	Collector *Collectr
}

func (c *HttpHandler) PostAppend(w http.ResponseWriter, r *http.Request) {
	// Declare a new Person struct.
	var payload map[string]any

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c.Logger.Error(err.Error())
		return
	}

	err = c.Collector.Append(payload)

	if err != nil {
		// FIXME: Append emits more error types...
		c.Logger.Error(err.Error())
		http.Error(w, "Unable to process request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
