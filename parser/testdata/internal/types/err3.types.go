// Code generated by goctl-types plugin. DO NOT EDIT.
package types

import (
	"time"
)

var (
	_ = time.Now()
)

type Error3 struct {
	PointerType *EmbededError3 `json:"pointerType"`
}

type EmbededError3 struct {
	Message string `json:"message"`
}
