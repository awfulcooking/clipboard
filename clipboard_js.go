// Copyright 2021 @awfulcooking. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build js

package clipboard

import (
	"errors"
	"syscall/js"
)

func readAll() (string, error) {
	v, err := awaitPromise(_clipboard.Call("readText"))
	return v.String(), err
}

func writeAll(text string) error {
	_, err := awaitPromise(_clipboard.Call("writeText", text))
	return err
}

var _clipboard js.Value

func init() {
	_clipboard = js.Global().Get("navigator").Get("clipboard")
}

func awaitPromise(p js.Value) (js.Value, error) {
	resolve := make(chan js.Value, 1)
	reject  := make(chan js.Value, 1)

	p.Call("then", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		resolve <- args[0]
		return nil
	})).Call("catch", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		reject <- args[0]
		return nil
	}))

	select {
	case result := <-resolve:
		return result, nil
	case rejection := <-reject:
		return js.Null(), errors.New(rejection.String())
	}
}

