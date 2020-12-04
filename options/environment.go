package options

import (
	"fmt"
	"strings"
)

type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentStaging     Environment = "staging"
	EnvironmentProduction  Environment = "production"
)

func (e Environment) IsDevelopment() bool {
	return e == EnvironmentDevelopment
}

func (e Environment) IsStaging() bool {
	return e == EnvironmentStaging
}

func (e Environment) IsProduction() bool {
	return e == EnvironmentProduction
}

func ParseEnvironment(str string) Environment {
	switch strings.ToLower(str) {
	case string(EnvironmentDevelopment), "dev":
		return EnvironmentDevelopment

	case string(EnvironmentStaging):
		return EnvironmentStaging

	case string(EnvironmentProduction), "prod":
		return EnvironmentProduction
	}

	panic(fmt.Sprintf("invalid environment '%s'", str))
}
