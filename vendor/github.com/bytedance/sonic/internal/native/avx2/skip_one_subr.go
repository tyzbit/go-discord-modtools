// +build !noasm !appengine
// Code generated by asm2asm, DO NOT EDIT.

package avx2

import (
	`github.com/bytedance/sonic/loader`
)

const (
    _entry__skip_one = 432
)

const (
    _stack__skip_one = 120
)

const (
    _size__skip_one = 10292
)

var (
    _pcsp__skip_one = [][2]uint32{
        {1, 0},
        {4, 8},
        {6, 16},
        {8, 24},
        {10, 32},
        {12, 40},
        {13, 48},
        {9740, 120},
        {9744, 48},
        {9745, 40},
        {9747, 32},
        {9749, 24},
        {9751, 16},
        {9753, 8},
        {9757, 0},
        {10292, 120},
    }
)

var _cfunc_skip_one = []loader.CFunc{
    {"_skip_one_entry", 0,  _entry__skip_one, 0, nil},
    {"_skip_one", _entry__skip_one, _size__skip_one, _stack__skip_one, _pcsp__skip_one},
}