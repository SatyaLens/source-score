//go:build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/onsi/ginkgo/v2"
	_ "github.com/onsi/gomega"
	// _ "github.com/maxbrunsfeld/counterfeiter/v6"
)