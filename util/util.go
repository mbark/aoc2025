package util

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func ReadInput(in, splitBy string) []string {
	trimmed := strings.Trim(in, "\n")
	return strings.Split(trimmed, splitBy)
}

func Str2IntSlice(in []string) []int {
	var list []int
	for _, s := range in {
		list = append(list, Str2Int(s))
	}

	return list
}

func NumberList(in string, separator string) []int {
	var list []int
	for _, s := range strings.Split(in, separator) {
		if s == "" {
			continue
		}
		list = append(list, Str2Int(s))
	}

	return list
}

func Str2Int(in string) int {
	i, _ := strconv.Atoi(in)
	return i
}

func NewBoolMatrix(width, height int) map[int]map[int]bool {
	m := make(map[int]map[int]bool, width)
	for i := 0; i < width; i++ {
		m[i] = make(map[int]bool, height)
		for j := 0; j < height; j++ {
			m[i][j] = false
		}
	}

	return m
}
func CopyList[T any](l []T) []T {
	cp := make([]T, len(l))
	copy(cp, l)
	return cp
}

func CopyMap[K comparable, V any](m map[K]V) map[K]V {
	cp := make(map[K]V, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

func WithTime() func() {
	now := time.Now()

	return func() { fmt.Printf("time taken: %v\n", time.Now().Sub(now)) }
}

func WithProfiling() func() {
	f, err := os.Create("profile.out")
	if err != nil {
		panic(err)
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		panic(err)
	}
	return pprof.StopCPUProfile
}

func ParseInt[T int64 | int32 | int](s string) T {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("unable to convert string %s to int", s))
	}
	return T(i)
}

func Btoi[T int64 | int32 | int](s string) T {
	i, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		panic(fmt.Sprintf("unable to convert string %s to binary", s))
	}
	return T(i)
}

func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}
