// +build !openbsd

package protector

// Source https://github.com/junegunn/fzf/blob/a1bcdc225e1c9b890214fcea3d19d85226fc552a/src/protector/protector.go

// Protect calls OS specific protections like pledge on OpenBSD
func Protect(doesntMatter string) {
	return
}
