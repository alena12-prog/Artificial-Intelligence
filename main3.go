package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var (
	correct [3][3]int = [3][3]int{[3]int{1, 2, 3}, [3]int{4, 5, 6}, [3]int{7, 8, 0}}
	step    int

	succcesfull int

	shuffleCircles = 2
	shuffleUsed    [][3][3]int
	SolveArr       []int
	timeArr        []time.Duration
)

type field struct {
	field [3][3]int
	H     func([3][3]int) float64
}

func (f *field) H1(fld [3][3]int) float64 {
	dst := 0.0
	for i := 1; i < 9; i++ {
		x, y := (*f).Locate(i, fld)
		x1, y1 := (*f).Locate(i, correct)
		if x != x1 || y != y1 {
			dst++
		}
	}
	return dst
}

func (f *field) SetH(s string) {
	if s == "H1" {
		(*f).H = (*f).H1
	} else if s == "H2" {
		(*f).H = (*f).H2
	}
}

func (f *field) H2(fld [3][3]int) float64 {
	dst := 0.0
	for i := 1; i < 9; i++ {
		x, y := (*f).Locate(i, fld)
		x1, y1 := (*f).Locate(i, correct)
		dst += math.Sqrt(float64((x1-x)*(x1-x) + (y1-y)*(y1-y)))
	}
	return dst
}

func (f *field) Create() {
	(*f).field = correct
}

func (f *field) HILL(fld [3][3]int, stopSearchChance float64) [3][3]int {
	start := fld
	for {
		if start == correct {
			return start
		}
		neigh := (*f).RecieveNeighbours(start)
		var nextNode [3][3]int
		nextEval := math.Inf(1)
		step++
		for _, v := range neigh {
			if (*f).H(v) < nextEval {
				nextEval = (*f).H(v)
				nextNode = v
				stp := rand.Float64()
				if stp < stopSearchChance {
					break
				}
			}
		}
		if (*f).H(start) <= (*f).H(nextNode) {
			break
		}
		start = nextNode
	}
	return [3][3]int{[3]int{-1, -1, -1}, [3]int{-1, -1, -1}, [3]int{-1, -1, -1}}
}

func DoHILL(s string, tries int, stopChance float64) {
	SolveArr = make([]int, 0)
	timeArr = make([]time.Duration, 0)
	fmt.Println()
	for i := 0; i < 20; i++ {
		fmt.Println("\n\nЕКСПЕРИМЕНТ №", i)
		f := Start(s)
		fmt.Println("\tПочатковий стан:")
		(*f).Print()
		t := time.Now()
		step = 0
		for i := 0; i < tries; i++ {
			res := (*f).HILL((*f).field, stopChance)
			if res == correct {
				letter := "ів"
				st := step % 10
				if st == 1 {
					letter = ""
				} else if st <= 3 {
					letter = "и"
				}
				tt := time.Since(t)
				fmt.Printf("\n\tСпроба %v\nРішення знайдено за %v крок%s; час пошуку: %v \n", i+1, step, letter, tt)
				SolveArr = append(SolveArr, step)
				timeArr = append(timeArr, tt)
				break
			} else {

				fmt.Printf("\n\tСпроба %v\nРішення не знайдено", i+1)

			}
		}
	}
	fmt.Println("\nМетод HILL:")
	fmt.Printf("\tВирішено для %v з %v\n", len(SolveArr), 20)
	if len(SolveArr) > 0 {
		median := 0
		var tMed time.Duration
		for i := 0; i < len(SolveArr); i++ {
			median += SolveArr[i]
			tMed += timeArr[i]
		}
		median /= len(SolveArr)
		tMed /= time.Duration(len(SolveArr))

		fmt.Printf("\tСередній час при вирішенні: %v, Середня кількість кроків: %v", tMed, median)
	}
}

func (f *field) Shuffle() {
	rep := true
	(*f).Create()
	for rep {
		for i := 0; i < shuffleCircles; i++ {
			x, y := 0, 0
			for {
				x, y = rand.Int()%2, rand.Int()%2
				if (*f).field[y][x] == 0 || (*f).field[y][x+1] == 0 || (*f).field[y+1][x] == 0 || (*f).field[y+1][x+1] == 0 {
					break
				}
			}

			tmp := (*f).field[y][x]
			(*f).field[y][x] = (*f).field[y][x+1]
			(*f).field[y][x+1] = (*f).field[y+1][x+1]
			(*f).field[y+1][x+1] = (*f).field[y+1][x]
			(*f).field[y+1][x] = tmp
		}
		rep = false
		for i := 0; i < len(shuffleUsed); i++ {
			if (*f).field == shuffleUsed[i] {
				rep = true
			}
		}
	}
	shuffleUsed = append(shuffleUsed, (*f).field)
}

func (f *field) RecieveNeighbours(fld [3][3]int) [][3][3]int {
	nilY, nilX := (*f).Locate(0, fld)
	res := make([][3][3]int, 0)
	cField := fld
	if nilX > 0 {
		fld[nilY][nilX], fld[nilY][nilX-1] = fld[nilY][nilX-1], fld[nilY][nilX]
		res = append(res, fld)
	}
	fld = cField
	if nilY > 0 {
		fld[nilY][nilX], fld[nilY-1][nilX] = fld[nilY-1][nilX], fld[nilY][nilX]
		res = append(res, fld)
	}
	fld = cField
	if nilX < 2 {
		fld[nilY][nilX], fld[nilY][nilX+1] = fld[nilY][nilX+1], fld[nilY][nilX]
		res = append(res, fld)
	}
	fld = cField
	if nilY < 2 {
		fld[nilY][nilX], fld[nilY+1][nilX] = fld[nilY+1][nilX], fld[nilY][nilX]
		res = append(res, fld)
	}
	return res
}

func (f *field) Locate(num int, fld [3][3]int) (int, int) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if fld[i][j] == num {
				return i, j
			}
		}
	}
	fmt.Println("nil not found:", fld)
	return -1, -1
}

func Start(s string) *field {
	f := field{}
	f.Create()
	f.Shuffle()
	f.SetH(s)
	return &f
}

func (f field) Print() {
	fmt.Println()
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if f.field[i][j] != 0 {
				fmt.Print("\t", f.field[i][j], "  ")
			} else {
				fmt.Print("\t_  ")
			}
		}
		fmt.Println()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	shuffleCircles = 2

	DoHILL("H1", 3, 0.5)

}
