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

/*
* Check if a slice Contains an element
* Args:
*   s: slice to check
*   e: element to check
 */
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
