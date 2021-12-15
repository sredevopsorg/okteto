// Copyright 2021 The Okteto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"context"
	"os"
	"strconv"

	"github.com/okteto/okteto/pkg/log"
	"github.com/okteto/okteto/pkg/types"
)

// HasAccessToNamespace checks if the user has access to a namespace/preview
func HasAccessToNamespace(ctx context.Context, namespace string, oktetoClient types.OktetoInterface) (bool, error) {

	nList, err := oktetoClient.Namespaces().List(ctx)
	if err != nil {
		return false, err
	}

	for i := range nList {
		if nList[i].ID == namespace {
			return true, nil
		}
	}

	previewList, err := oktetoClient.Previews().List(ctx)
	if err != nil {
		return false, err
	}

	for i := range previewList {
		if previewList[i].ID == namespace {
			return true, nil
		}
	}

	return false, nil
}

func LoadBoolean(k string) bool {
	v := os.Getenv(k)
	if v == "" {
		v = "false"
	}

	h, err := strconv.ParseBool(v)
	if err != nil {
		log.Yellow("'%s' is not a valid value for environment variable %s", v, k)
	}

	return h
}
