// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cli // import "onebill.app/internal/cli"

import (
	"errors"
	"fmt"

	"onebill.app/internal/model"
	"onebill.app/internal/storage"
	"onebill.app/internal/validator"
)

func resetPassword(store *storage.Storage) {
	username, password := askCredentials()
	user, err := store.UserByUsername(username)
	if err != nil {
		printErrorAndExit(err)
	}

	if user == nil {
		printErrorAndExit(errors.New("user not found"))
	}

	userModificationRequest := &model.UserModificationRequest{
		Password: &password,
	}
	if validationErr := validator.ValidateUserModification(store, user.ID, userModificationRequest); validationErr != nil {
		printErrorAndExit(validationErr.Error())
	}

	user.Password = password
	if err := store.UpdateUser(user); err != nil {
		printErrorAndExit(err)
	}

	fmt.Println("Password changed!")
}
