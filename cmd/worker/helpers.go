package worker

import (
	"bytes"
	"encoding/json"
	"errors"

	"example.com/taskservice/internal/errs"
)

func parseRawConfig[T IntervalConfig | OddEvenConfig | YearlyDateConfig](raw json.RawMessage) (T, error) {
	var cfg T
	decoder := json.NewDecoder(bytes.NewBuffer(raw))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&cfg); err != nil {
		return cfg, errors.Join(errs.ErrInvalidInput, err)
	}
	if decoder.More() {
		return cfg, errors.New("unexpected trailing data")
	}
	return cfg, nil
}
