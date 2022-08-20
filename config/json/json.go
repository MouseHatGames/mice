package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/MouseHatGames/mice/config"
	"github.com/MouseHatGames/mice/options"
)

func Config(filePath string) options.Option {
	return func(o *options.Options) {
		o.Config = &jsonConfig{
			filePath: filePath,
		}
	}
}

type jsonConfig struct {
	filePath string
	data     map[string]interface{}
}

func (c *jsonConfig) load() error {
	f, err := os.ReadFile(c.filePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	if err := json.Unmarshal(f, &c.data); err != nil {
		return fmt.Errorf("parse json: %w", err)
	}

	return nil
}

func (c *jsonConfig) save() error {
	b, err := json.Marshal(c.data)
	if err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	if err := os.WriteFile(c.filePath, b, 0); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (c *jsonConfig) loadIfNotLoaded() error {
	if c.data == nil {
		if err := c.load(); err != nil {
			return fmt.Errorf("load config: %w", err)
		}
	}

	return nil
}

func (c *jsonConfig) Get(path ...string) config.Value {
	if err := c.loadIfNotLoaded(); err != nil {
		return &jsonValue{err: err}
	}

	var obj interface{} = c.data

	for _, p := range path {
		if m, ok := obj.(map[string]interface{}); ok {
			obj = m[p]
		} else {
			return cannotIndexValue
		}
	}

	return &jsonValue{val: obj}
}

func (c *jsonConfig) Delete(path ...string) error {
	if err := c.loadIfNotLoaded(); err != nil {
		return err
	}

	var obj interface{} = c.data

	for i, p := range path {
		if m, ok := obj.(map[string]interface{}); ok {
			if i == len(path)-1 {
				delete(m, p)
			} else {
				obj = m[p]
			}
		} else {
			return ErrCannotIndexValue
		}
	}

	return c.save()
}

func (c *jsonConfig) Set(val interface{}, path ...string) error {
	if err := c.loadIfNotLoaded(); err != nil {
		return err
	}

	var obj interface{} = c.data

	for i, p := range path {
		if m, ok := obj.(map[string]interface{}); ok {
			if i == len(path)-1 {
				m[p] = val
			} else {
				obj = m[p]
			}
		} else {
			return ErrCannotIndexValue
		}
	}

	return c.save()
}

// ErrCannotIndexValue is returned when you try to access a property of a non-object value
var ErrCannotIndexValue = errors.New("tried to index non-indexable value")
var cannotIndexValue = &jsonValue{err: ErrCannotIndexValue}

type jsonValue struct {
	err error
	val interface{}
}

func (v *jsonValue) Error() error {
	return v.err
}

func (v *jsonValue) Raw() string {
	b, _ := json.Marshal(v.val)
	return string(b)
}

func (v *jsonValue) Scan(out interface{}) error {
	b, _ := json.Marshal(v.val)
	return json.Unmarshal(b, out)
}

func (v *jsonValue) String() (string, bool) {
	o, ok := v.val.(string)
	return o, ok
}

func (v *jsonValue) StringOr(def string) string {
	if val, ok := v.String(); ok {
		return val
	}
	return def
}

func (v *jsonValue) Int32() (int32, bool) {
	if o, ok := v.val.(float64); ok {
		return int32(o), true
	}
	return 0, false
}

func (v *jsonValue) Int32Or(def int32) int32 {
	if val, ok := v.Int32(); ok {
		return val
	}
	return def
}

func (v *jsonValue) Int64() (int64, bool) {
	if o, ok := v.val.(float64); ok {
		return int64(o), true
	}
	return 0, false
}

func (v *jsonValue) Int64Or(def int64) int64 {
	if val, ok := v.Int64(); ok {
		return val
	}
	return def
}

func (v *jsonValue) Bool() (bool, bool) {
	o, ok := v.val.(bool)
	return o, ok
}

func (v *jsonValue) BoolOr(def bool) bool {
	if val, ok := v.Bool(); ok {
		return val
	}
	return def
}

func (v *jsonValue) Float64() (float64, bool) {
	o, ok := v.val.(float64)
	return o, ok
}

func (v *jsonValue) Float64Or(def float64) float64 {
	if val, ok := v.Float64(); ok {
		return val
	}
	return def
}

func (v *jsonValue) Duration() (time.Duration, bool) {
	if val, ok := v.val.(string); ok {
		dur, err := time.ParseDuration(val)
		if err != nil {
			return 0, false
		}
		return dur, true
	}
	return 0, false
}

func (v *jsonValue) DurationOr(def time.Duration) time.Duration {
	if val, ok := v.Duration(); ok {
		return val
	}
	return def
}
