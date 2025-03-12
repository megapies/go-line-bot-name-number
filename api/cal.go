package handler

import "fmt"

func getAlphabetNumber(a rune) (int, error) {
	// Define Thai character to number mapping
	thaiNumMap := map[rune]int{
		'ก': 1,
		'ด': 1,
		'ท': 1,
		'ถ': 1,
		'ภ': 1,
		'ฤ': 1,
		'ฦ': 1,
		'่': 1,
		'ุ': 1,
		'า': 1,
		'ำ': 1,
		'ข': 2,
		'ช': 2,
		'ง': 2,
		'บ': 2,
		'ป': 2,
		'้': 2,
		'เ': 2,
		'แ': 2,
		'ู': 2,
		'ฆ': 3,
		'ต': 3,
		'ฑ': 3,
		'ฒ': 3,
		'๋': 3,
		'ค': 4,
		'ธ': 4,
		'ญ': 4,
		'ร': 4,
		'ษ': 4,
		'ะ': 4,
		'ิ': 4,
		'โ': 4,
		'ั': 4,
		'ฉ': 5,
		'ณ': 5,
		'ฌ': 5,
		'น': 5,
		'ม': 5,
		'ห': 5,
		'ฎ': 5,
		'ฬ': 5,
		'ฮ': 5,
		'ึ': 5,
		'จ': 6,
		'ล': 6,
		'ว': 6,
		'อ': 6,
		'ใ': 6,
		'ซ': 7,
		'ศ': 7,
		'ส': 7,
		'๊': 7,
		'ี': 7,
		'ื': 7,
		'ผ': 8,
		'ฝ': 8,
		'พ': 8,
		'ฟ': 8,
		'ย': 8,
		'็': 8,
		'ฏ': 9,
		'ฐ': 9,
		'ไ': 9,
		'์': 9,
	}

	number, ok := thaiNumMap[a]
	if !ok {
		return 0, fmt.Errorf("invalid character: %c", a)
	}

	return number, nil
}

func CalNameNumber(name string) (int, string, error) {
	// Define Thai character to number mapping

	return 0, "", nil
}
