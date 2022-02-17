// Package edkey is a copy of https://github.com/mikesmitty/edkey, which creates ED25519 keys in the OpenSSH format.
//
// Instead of importing it as a dependency, we have copied it here. This is because this is critical security functionality, and
// we cannot rely on a single private person to
// - not remove the code at some random point in time
// - not inject the code with security compromising code. We won't be affected by such an eventuality before we bump the version
// of the dependency, but it's not unlikely that a future developer will bump a version without reading the contents of the code
// they are bumping to.
//
// There also shouldn't be a need to bump this code at all, as it's working.
//
// Next question, is it secure? Yes. The only thing this package does is wrap a ed25519 key in the OpenSSH format. Details:
// - The code uses the official ed25519 package.
// - Line 56, "pk1.Pub = pubKey", shows that the code is just wrapping a pubKey into a "pk1" struct that contains various
// metadata. edkey_test.go -> TestEdkeyContainsOriginalKey verifies this wrapping.
//
// We should replace this whole package when Golang adds proper support for OpenSSH format. See:
// https://github.com/golang/go/issues/37132
package edkey
