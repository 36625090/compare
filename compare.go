/*
 * Copyright 2022 The Go Authors<36625090@qq.com>. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 */

package compare

import (
	"fmt"
	"reflect"
)

type Comparable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

type CmpMode int

const (
	CmpModeLess         CmpMode = -1
	CmpModeLessEqual    CmpMode = -10
	CmpModeEqual        CmpMode = 0
	CmpModeGreater      CmpMode = 1
	CmpModeGreaterEqual CmpMode = 10
)

type Serializer interface {
	HashCode() int
}

func Compare[T any](i, j T, mode CmpMode) bool {
	iValue := reflect.ValueOf(i)
	jValue := reflect.ValueOf(j)

	switch iValue.Kind() {
	case reflect.Pointer:
		if iValue.Elem().Kind() == reflect.Struct {
			switch ip := iValue.Interface().(type) {
			case Serializer:
				switch jp := jValue.Interface().(type) {
				case Serializer:
					return Cmp(ip.HashCode(), jp.HashCode(), mode)
				}
			}
			switch ip := iValue.Interface().(type) {
			case fmt.Stringer:
				switch jp := jValue.Interface().(type) {
				case fmt.Stringer:
					return Cmp(ip.String(), jp.String(), mode)
				}
			}
		}
		return Compare(iValue.Elem().Interface(), jValue.Elem().Interface(), mode)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Cmp(iValue.Int(), jValue.Int(), mode)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Cmp(iValue.Uint(), jValue.Uint(), mode)

	case reflect.Float32, reflect.Float64:
		return Cmp(iValue.Float(), jValue.Float(), mode)

	case reflect.Complex128, reflect.Complex64:
		return Complex128Cmp(iValue.Complex(), jValue.Complex(), mode)

	case reflect.String:
		return Cmp(iValue.String(), jValue.String(), mode)

	case reflect.Bool:
		if mode == CmpModeEqual {
			return iValue.Bool() == jValue.Bool()
		}

	default:
		return false
	}

	return false
}

func Cmp[T Comparable](i, j T, mode CmpMode) bool {
	switch mode {
	case CmpModeLess:
		return i < j
	case CmpModeLessEqual:
		return i <= j
	case CmpModeEqual:
		return i == j
	case CmpModeGreater:
		return i > j
	case CmpModeGreaterEqual:
		return i >= j
	}
	return false
}

func Complex64Cmp(i, j complex64, mode CmpMode) bool {
	switch mode {
	case CmpModeLess:
		return real(i) < real(j) || (real(i) == real(j) && imag(i) < imag(j))
	case CmpModeEqual:
		return real(i) == real(j) && imag(i) == imag(j)
	case CmpModeGreater:
		return real(i) > real(j) || (real(i) == real(j) && imag(i) > imag(j))
	}
	return false
}

func Complex128Cmp(i, j complex128, mode CmpMode) bool {
	switch mode {
	case CmpModeLess:
		return real(i) < real(j) || (real(i) == real(j) && imag(i) < imag(j))
	case CmpModeLessEqual:
		return real(i) <= real(j) || (real(i) == real(j) && imag(i) <= imag(j))
	case CmpModeEqual:
		return real(i) == real(j) && imag(i) == imag(j)
	case CmpModeGreater:
		return real(i) > real(j) || (real(i) == real(j) && imag(i) > imag(j))
	case CmpModeGreaterEqual:
		return real(i) >= real(j) || (real(i) == real(j) && imag(i) >= imag(j))
	}
	return false
}
