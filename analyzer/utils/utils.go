// Copyrigth (c) 2024 Erik Kassubek
//
// File: utils.go
// Brief: Utility function to check if an slice contains a value
//
// Author: Erik Kassubek
// Created: 2024-04-06
//
// License: BSD-3-Clause

package utils

import "strings"

/*
* Check if a slice ContainsString an element
* Args:
*   s: slice to check
*   e: element to check
* Returns:
*   bool: true is e in s, false otherwise
 */
func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsInt(slice []int, elem int) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}

/*
 * Split the string into two parts at the last occurrence of the separator
 * Args:
 *   str (string): string to split
 *   sep (string): separator to split at
 * Returns:
 *   []string: If sep in string: list with two elements split at the sep,
 *     if not then list containing str
 */
func SplitAtLast(str string, sep string) []string {
	if sep == "" {
		return []string{str}
	}

	i := strings.LastIndex(str, sep)
	if i == -1 {
		return []string{str}
	}
	return []string{str[:i], str[i+1:]}
}
