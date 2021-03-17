// Package rotater embeds the lambda function code zip file
package rotater

// We import embed like this because it is required by the embed
// functionality in go 1.16
// nolint: golint
import _ "embed"

// LambdaFunctionZip contains the embedded data for the
// rotater lambda
// nolint: gochecknoglobals
//go:embed lambda_function.zip
var LambdaFunctionZip []byte
