// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package validator // import "onebill.app/internal/validator"

import (
	"testing"

	"onebill.app/internal/model"
)

func TestValidateRange(t *testing.T) {
	if err := ValidateRange(-1, 0); err == nil {
		t.Error(`An invalid offset should generate a error`)
	}

	if err := ValidateRange(0, -1); err == nil {
		t.Error(`An invalid limit should generate a error`)
	}

	if err := ValidateRange(0, model.MaxEntryLimit+1); err == nil {
		t.Error(`A limit above MaxEntryLimit should generate an error`)
	}

	if err := ValidateRange(42, 42); err != nil {
		t.Error(`A valid offset and limit should not generate any error`)
	}

	if err := ValidateRange(0, model.MaxEntryLimit); err != nil {
		t.Error(`A limit equal to MaxEntryLimit should not generate an error`)
	}
}

func TestValidateDirection(t *testing.T) {
	for _, status := range []string{"asc", "desc"} {
		if err := ValidateDirection(status); err != nil {
			t.Error(`A valid direction should not generate any error`)
		}
	}

	if err := ValidateDirection("invalid"); err == nil {
		t.Error(`An invalid direction should generate a error`)
	}
}
