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
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
)

// Structure of json configuration.
type CollectrConf struct {
	Fields map[string]string
	File   string
}

type Collectr struct {
	Writer io.Writer
	f      *os.File
	Conf   CollectrConf
	mu     sync.Mutex // file access mutex
}

// Loads json config of collectr.
func (c *Collectr) LoadConfJson(r io.Reader) error {
	by, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	var conf CollectrConf
	err = json.Unmarshal(by, &conf)
	if err != nil {
		return err
	}
	c.Conf = conf
	return nil
}

// Handles validation rules "required|nonempty|text|number|regexp:".
func (c *Collectr) validate(cmap map[string]any) error {
	// cmap needs to have these types to be valid
	for field, validationRuleStr := range c.Conf.Fields {
		value := cmap[field]

		valRules := strings.Split(validationRuleStr, "|")
		// Iterates over validation rules
		for _, rule := range valRules {
			if rule == "required" {
				_, ok := cmap[field]
				if !ok {
					return errors.New(field + " is required")
				}
			} else if rule == "text" {
				_, ok := value.(string)
				if !ok {
					return errors.New("invalid type of " + field + ", should have text type")
				}
			} else if rule == "nonempty" {
				s, ok := value.(string)
				if ok && len(s) == 0 {
					return errors.New(field + " must be nonempty string")
				}
			} else if rule == "number" {
				_, okInt := value.(int)
				_, okFloat := value.(int)

				if !(okInt || okFloat) {
					return errors.New("invalid type of " + field + ", should have number type")
				}
			} else if strings.HasPrefix(rule, "regexp:") {
				if len(rule) <= len("regexp:") {
					return errors.New("invalid regexp definition, use regexp:regexp_here")
				}

				regexpStr := rule[7:]
				rx, err := regexp.Compile(regexpStr)

				if err != nil {
					return err
				}

				v, okStr := value.(string)
				if !okStr {
					return errors.New("regexp can be applied only on strings")
				}

				ok := rx.Match([]byte(string(v)))
				if !ok {
					return errors.New("invalid type of " + field + ", should have number type")
				}
			}
		}
	}

	return nil
}

// Inits new collectr instance with file backend.
// Use this for production environment.
func New() (*Collectr, error) {
	coll := &Collectr{}

	by, err := os.ReadFile("collectr.json")
	if err != nil {
		return nil, err
	}

	var conf CollectrConf
	err = json.Unmarshal(by, &conf)
	if err != nil {
		return nil, err
	}
	coll.Conf = conf

	f, err := os.OpenFile(coll.Conf.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	coll.Writer = f
	coll.f = f

	return coll, nil
}

// Use to close file when initialized using New() function.
func (c *Collectr) Close() {
	c.f.Close()
}

// Append data into file specified in collectr.json as "file" field.
// Handles slight validation before appending.
// Aims to be simple "write to the disk"data from API call.
func (c *Collectr) Append(cmap map[string]any) error {
	// check if cmap is valid
	if err := c.validate(cmap); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Marshal map to JSON string
	jb, err := json.Marshal(cmap)

	if err != nil {
		return err
	}

	// Log in JSONL format
	if _, err := io.WriteString(c.Writer, string(jb)+"\n"); err != nil {
		return err
	}

	return nil
}
