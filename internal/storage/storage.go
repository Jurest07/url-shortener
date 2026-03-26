package storage

import "errors"

var (
	ErrURLNotFound 		= errors.New("url not found")
	ErrURLExsits   		= errors.New("url exists")
	ErrAliasNotFound 	= errors.New("alias not fould")
	ErrAliasExists 		= errors.New("alias exists")
)
