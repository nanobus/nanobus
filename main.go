package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const header = `/*
Copyright 2022 The NanoBus Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/`

func main() {
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			fmt.Println(path, info.Size())

			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if bytes.HasPrefix(data, []byte(header)) {
				data = bytes.TrimLeft(data[len(header):], " \n\t")
				err = os.WriteFile(path, data, 0644)
				if err != nil {
					return err
				}
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
