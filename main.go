package main

import (
	"fmt"
	"math/rand"
	"time"
)

type board [][]rune

func Beat(b board, x, y int) board {
	tt := NewBoard(len(b))
	for i := 0; i < len(b); i++ {
		for j := 0; j < len(b[i]); j++ {
			tt[i][j] = b[i][j]
		}
	}
	tt[y][x] = 'Q'
	for i := 1; i < len(b); i++ {
		if y-i >= 0 && x-i > 0 {
			tt[y-i][x-i] = '.'
		}
		if y-i >= 0 && x+i < len(b) {
			tt[y-i][x+i] = '.'
		}
		if y+i < len(b) && x-i >= 0 {
			tt[y+i][x-i] = '.'
		}
		if y+i < len(b) && x+i < len(b) {
			tt[y+i][x+i] = '.'
		}
		if y-i >= 0 {
			tt[y-i][x] = '.'
		}
		if x-i >= 0 {
			tt[y][x-i] = '.'
		}
		if y+i < len(b) {
			tt[y+i][x] = '.'
		}
		if x+i < len(b) {
			tt[y][x+i] = '.'
		}
	}
	return tt
}

func placeQueen(b board, row int) (board, int) {

	if row >= len(b) {
		return b, row
	}
	// for i := 0; i < len(b); i++ {
	// 	for j := 0; j < len(b[i]); j++ {
	// 		fmt.Print(string(b[i][j]), " ")
	// 	}
	// 	fmt.Println()
	// }
	// fmt.Println()
	var bst board
	maxSc := -1
	for i := 0; i < len(b); i++ {
		if b[row][i] == 0 {
			tst, num := placeQueen(Beat(b, i, row), row+1)
			if num > maxSc {
				bst = tst
				maxSc = num
			}
		}
	}
	return bst, maxSc
}

func NewBoard(size int) board {
	b := make([][]rune, size)
	for i := 0; i < size; i++ {
		b[i] = make([]rune, size)
	}
	return board(b)
}

func Solve(size int) {
	b := NewBoard(size)
	res, val := placeQueen(b, 0)
	fmt.Println("Queens: ", val)
	for i := 0; i < len(res); i++ {
		for j := 0; j < len(res[i]); j++ {
			fmt.Print(string(res[i][j]), " ")
		}
		fmt.Println()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	Solve(12)
}
