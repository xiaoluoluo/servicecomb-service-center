/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package validator

import (
	"errors"
	"reflect"

	"github.com/apache/servicecomb-service-center/pkg/validate"
)

var createAccountValidator = &validate.Validator{}
var updateAccountValidator = &validate.Validator{}
var createRoleValidator = &validate.Validator{}

var changePWDValidator = &validate.Validator{}
var accountLoginValidator = &validate.Validator{}

func init() {
	createAccountValidator.AddRule("Name", &validate.Rule{Max: 64, Regexp: nameRegex})
	createAccountValidator.AddRule("Roles", &validate.Rule{Min: 1, Max: 5, Regexp: nameRegex})
	createAccountValidator.AddRule("Password", &validate.Rule{Regexp: &validate.PasswordChecker{}})
	createAccountValidator.AddRule("Status", &validate.Rule{Regexp: accountStatusRegex})

	updateAccountValidator.AddRule("Roles", createAccountValidator.GetRule("Roles"))
	updateAccountValidator.AddRule("Status", createAccountValidator.GetRule("Status"))

	createRoleValidator.AddRule("Name", &validate.Rule{Max: 64, Regexp: nameRegex})

	changePWDValidator.AddRule("Password", &validate.Rule{Regexp: &validate.PasswordChecker{}})
	changePWDValidator.AddRule("Name", &validate.Rule{Regexp: nameRegex})

	accountLoginValidator.AddRule("TokenExpirationTime", &validate.Rule{Regexp: &validate.TokenExpirationTimeChecker{}})
}

func baseCheck(v interface{}) error {
	if v == nil {
		return errors.New("data is nil")
	}
	sv := reflect.ValueOf(v)
	if sv.Kind() == reflect.Ptr && sv.IsNil() {
		return errors.New("pointer is nil")
	}
	return nil
}
