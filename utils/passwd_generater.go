package utils

import (
	"math/rand"
	"time"
)

// GeneratePassword 密码生成器
func GeneratePassword(length int, lower, upper, number, symbols bool) string {
	m := make(map[int]string)

	if !lower && !upper && !number && !symbols {
		lower = true
		upper = true
		number = true
	}

	k := 0
	if lower {
		m[k] = "abcdefghijklmnopqrstuvwxyz"
		k += 1
	}

	if upper {
		m[k] = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		k += 1
	}

	if number {
		m[k] = "0123456789"
		k += 1
	}

	if symbols {
		m[k] = "!@#$%^&*()_+{}|:<>?~"
		k += 1
	}

	randArray := generateRandomArray(length, k)
	ret := make([]byte, 0)

	for i, alternative := range m {
		for _ = range randArray[i] {
			ret = append(ret, byte(alternative[rand.Intn(len(alternative))]))
		}
	}

	return shuffleString(string(ret))
}

func shuffleString(s string) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	runes := []rune(s)
	r.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	return string(runes)
}

func generateRandomArray(n, k int) []int {
	ret := make([]int, k)
	for i := range ret {
		ret[i] = 1
	}

	remaining := n - k
	rand.New(rand.NewSource(time.Now().UnixNano()))

	for remaining > 0 {
		index := rand.Intn(k)
		ret[index]++
		remaining--
	}

	return ret
}
