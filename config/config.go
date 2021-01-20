package config

import "time"

type Config interface {
	Get(path ...string) Value
	Del(path ...string) error
	Set(val interface{}, path ...string) error
}

type Value interface {
	Error() error
	Raw() string
	Scan(v interface{}) error
	String() (string, bool)
	StringOr(def string) string
	Int32() (int32, bool)
	Int32Or(def int32) int32
	Int64() (int64, bool)
	Int64Or(def int64) int64
	Bool() (bool, bool)
	BoolOr(def bool) bool
	Float64() (float64, bool)
	Float64Or(def float64) float64
	Duration() (time.Duration, bool)
	DurationOr(def time.Duration) time.Duration
}
