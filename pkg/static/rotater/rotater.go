// Package rotater embeds the lambda function code zip file
package rotater

import _ "embed"

// LambdaFunctionZip contains the embedded data for the
// rotater lambda
// nolint: gochecknoglobals
//go:embed lambda_function.zip
var LambdaFunctionZip []byte
