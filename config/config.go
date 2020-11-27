package config

import "time"

type Config interface {
	Get(path ...string) Value
	Del(path ...string) error
	Set(val interface{}, path ...string) error
}

type Value interface {
	Bool(def bool) bool
	Int(def int) int
	String(def string) string
	Float64(def float64) float64
	Duration(def time.Duration) time.Duration
	Strings(def []string) []string
	StringMap(def map[string]string) map[string]string
	Scan(val interface{}) error
}
