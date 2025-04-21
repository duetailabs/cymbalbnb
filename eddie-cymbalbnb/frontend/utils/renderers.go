// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"strconv"
	"strings"
)

func RenderMoney(amount float32) string {
	return "$" + strconv.FormatFloat(float64(amount), 'f', 2, 32)
}

func RangeMax[T comparable](slice []T, maxLength int) []T {
	max := len(slice)
	if max > maxLength {
		max = maxLength
	}
	return slice[:max]
}

func FirstTwoWords(s string) string {
	words := strings.Split(s, " ")
	if len(words) < 3 {
		return s
	}
	return strings.Join(words[:2], " ") + "..."
}
