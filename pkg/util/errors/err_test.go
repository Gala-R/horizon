// Copyright © 2023 Horizoncd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errSome = errors.New("hello world")

const ErrCodeDuplicate = ErrorCode("Duplicate")

func TestError_WrapError(t *testing.T) {
	err := E("A", errSome)
	err = E("B", err)
	assert.Equal(t, "B - A - hello world", err.Error())
	assert.Equal(t, string(ErrCodeInternalError), Code(err))
}

func TestError_TranslateError(t *testing.T) {
	err := E("A", ErrCodeInternalError, "fatal error")
	err = E("B", 409, ErrCodeDuplicate, err, "duplicate")
	assert.Equal(t, "B - A - <409 Duplicate - duplicate> <500 InternalError - fatal error> ", err.Error())
	assert.Equal(t, string(ErrCodeDuplicate), Code(err))
	assert.Equal(t, "duplicate", Message(err))
}

func TestError_Status(t *testing.T) {
	assert.Equal(t, StatusInternalError, Status(errSome))
	err := E("A")
	assert.Equal(t, StatusInternalError, Status(err))

	err = E("B", 409, err)
	assert.Equal(t, 409, Status(err))

	err = E("C", 405, err)
	assert.Equal(t, 405, Status(err))
}

func TestError_ChangeCode(t *testing.T) {
	err := E("A", ErrCodeInternalError, errSome)
	err = E("B", 409, ErrCodeDuplicate, err)
	assert.Equal(t, "B - A - <409 Duplicate - > <500 InternalError - > hello world", err.Error())
	assert.Equal(t, string(ErrCodeDuplicate), Code(err))
	assert.Equal(t, "hello world", Message(err))
}

func TestError_Message(t *testing.T) {
	err := E("A", ErrorCode("BadRequest"))
	assert.Equal(t, "An internal error has occurred. Please contact technical support.", Message(err))

	err = E("B")
	assert.Equal(t, "An internal error has occurred. Please contact technical support.", Message(err))

	err = E("B", "hello world")
	assert.Equal(t, "hello world", Message(err))
}
