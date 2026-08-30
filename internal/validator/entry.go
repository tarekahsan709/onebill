// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package validator // import "onebill.app/internal/validator"

import (
	"errors"
	"fmt"

	"onebill.app/internal/model"
)

// ValidateEntriesStatusUpdateRequest validates a status update for a list of entries.
func ValidateEntriesStatusUpdateRequest(request *model.EntriesStatusUpdateRequest) error {
	if len(request.EntryIDs) == 0 {
		return errors.New(`the list of entries cannot be empty`)
	}

	return ValidateEntryStatus(request.Status)
}

// ValidateEntryStatus makes sure the entry status is valid.
func ValidateEntryStatus(status string) error {
	switch status {
	case model.EntryStatusRead, model.EntryStatusUnread:
		return nil
	}

	return fmt.Errorf(`invalid entry status, valid status values are: %q and %q`, model.EntryStatusRead, model.EntryStatusUnread)
}

// ValidateEntryOrder makes sure the sorting order is valid.
func ValidateEntryOrder(order string) error {
	switch order {
	case "id", "status", "changed_at", "published_at", "created_at", "title", "author":
		return nil
	}

	return errors.New(`invalid entry order, valid order values are: "id", "status", "changed_at", "published_at", "created_at", "title", "author"`)
}
