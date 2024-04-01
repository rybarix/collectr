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

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/rybarix/email-collectr/internal/collectr"
)

func setupLogger(pth string) (*slog.Logger, *os.File) {
	f, err := os.OpenFile(pth, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{AddSource: true}))

	return logger, f
}

func main() {
	// Parse port number, logging file
	portNumber := flag.Int("port", 8000, "server port number")
	logFilePath := flag.String("logfile", "json.log", "path to app log file")
	flag.Parse()

	if !(*portNumber >= 0 && *portNumber <= 0xFFFF) {
		panic(errors.New("invalid port number, use range 0-65535"))
	}

	// Setup logger
	logger, logf := setupLogger(*logFilePath)
	defer logf.Close()

	coll, err := collectr.New()
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
	defer coll.Close()

	state := collectr.HttpHandler{
		Logger:    logger,
		Collector: coll,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /collect", state.PostAppend)

	logger.Info("starting server")
	// Start the server
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *portNumber), mux); err != nil {
		logger.Error(err.Error())
		panic(err)
	}
}
