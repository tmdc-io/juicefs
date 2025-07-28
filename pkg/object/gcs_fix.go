//go:build !nos3
// +build !nos3

/*
 * JuiceFS, Copyright 2018 Juicedata, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package object

import (
	"reflect"
	"unsafe"
)

// Hack to fix GCS interoperability issue
// The AWS SDK includes accept-encoding in signature calculation, but GCS expects it to be ignored
type Rule interface {
	IsValid(value string) bool
}
type Rules []Rule

//go:linkname __ignoredHeaders github.com/aws/aws-sdk-go-v2/aws/signer/internal/v4.IgnoredHeaders
var __ignoredHeaders unsafe.Pointer

// Avoids "go.info.github.com/aws/aws-sdk-go-v2/aws/signer/internal/v4.IgnoredHeaders:
// relocation target go.info.github.com/xxx/xxx/xxx.Rules not defined"
// refer https://github.com/pkujhd/goloader/blob/09f36c84ac85502eb5df4670f1aa7472934ba03a/iface.1.10.go#L31-L36
var ignoredHeaders = (*Rules)(unsafe.Pointer(&__ignoredHeaders))

func init() {
	// Add accept-encoding to ignored headers to fix GCS interoperability
	reflect.ValueOf((*ignoredHeaders)[0]).FieldByName("Rule").Elem().SetMapIndex(
		reflect.ValueOf("Accept-Encoding"), reflect.ValueOf(struct{}{}))
} 