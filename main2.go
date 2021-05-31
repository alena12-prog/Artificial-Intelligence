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

func DoLDFS(s string, maxDepth int) {
	SolveArr = make([]int, 0)
	timeArr = make([]time.Duration, 0)
	fmt.Println("Максимальна глибина пошуку:", maxDepth)
	fmt.Println()
	for i := 0; i < 20; i++ {
		fmt.Println("\nЕКСПЕРИМЕНТ №", i+1)
		f := Start(s)
		fmt.Println("\tПочатковий стан:")
		(*f).Print()
		t := time.Now()
		step = 0
		res := (*f).DLS((*f).field, correct, maxDepth)
		if res == correct {
			letter := "ів"
			st := step % 10
			if st <= 1 {
				letter = ""
			} else if st <= 3 {
				letter = "и"
			}
			tt := time.Since(t)
			fmt.Printf("\n\tРішення знайдено за %v крок%s; час пошуку: %v \n", step, letter, tt)
			SolveArr = append(SolveArr, step)
			timeArr = append(timeArr, tt)
		} else {
			fmt.Println("\n\tРішення не знайдено")
		}
	}
	fmt.Println("\nМетод LDFS:")
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

func (f *field) DLS(root, goal [3][3]int, depth int) [3][3]int {
	step++
	if step > 10000000000 {
		return [3][3]int{[3]int{-1, -1, -1}, [3]int{-1, -1, -1}, [3]int{-1, -1, -1}}
	}
	//
	if depth == 0 && root == correct {
		return root
	} else if depth > 0 {
		neighbours := (*f).RecieveNeighbours(root)
		for _, n := range neighbours {
			res := (*f).DLS(n, goal, depth-1)
			if res == correct {
				return res
			}
		}
	}
	return [3][3]int{[3]int{-1, -1, -1}, [3]int{-1, -1, -1}, [3]int{-1, -1, -1}}
}

func (f *field) AStar(fld [3][3]int, maxSteps int) [3][3]int {
	openSet := make(map[[3][3]int][3][3]int)
	openSet[fld] = [3][3]int{[3]int{-1, -1, -1}, [3]int{-1, -1, -1}, [3]int{-1, -1, -1}}
	gScore := make(map[[3][3]int]float64)
	gScore[fld] = 0
	fScore := make(map[[3][3]int]float64)
	fScore[fld] = (*f).H(fld)
	for len(openSet) > 0 {
		step++
		if step > maxSteps {
			fmt.Println("Перевищено ліміт кроків; відміна пошуку")
			return [3][3]int{[3]int{-1, -1, -1}, [3]int{-1, -1, -1}, [3]int{-1, -1, -1}}
		}
		// if step%1000 == 0 {
		// 	fmt.Print(".")
		// }
		var curr [3][3]int
		mind := math.MaxFloat64
		for n, v := range fScore {
			if _, ok := openSet[n]; ok {
				if v < mind {
					curr = n
					mind = v
				}
			}
		}
		//fmt.Println(curr)
		if curr == correct {
			return curr
		}
		delete(openSet, curr)
		//fmt.Println(openSet)
		//fmt.Println("+a")
		neigh := (*f).RecieveNeighbours(curr)
		for _, v := range neigh {
			//fmt.Println(v)
			tentScore := gScore[curr] + 1
			if _, ok := gScore[v]; !ok {
				if _, ok := openSet[v]; !ok {
					if _, ok := fScore[v]; !ok {
						openSet[v] = curr
					}
				}
				gScore[v] = tentScore
				fScore[v] = gScore[v] + (*f).H(v)
				//fmt.Println("+c")
			} else if tentScore < gScore[v] {
				if _, ok := openSet[v]; !ok {
					if _, ok := fScore[v]; !ok {
						openSet[v] = curr
					}
				}
				gScore[v] = tentScore
				fScore[v] = gScore[v] + (*f).H(v)
				//fmt.Println("+d")

			}
		}
	}
	return [3][3]int{[3]int{-1, -1, -1}, [3]int{-1, -1, -1}, [3]int{-1, -1, -1}}
}

func DoAStar(s string, maxSteps int) {
	SolveArr = make([]int, 0)
	timeArr = make([]time.Duration, 0)
	fmt.Println()
	for i := 0; i < 20; i++ {
		fmt.Println("\nЕКСПЕРИМЕНТ №", i+1)
		f := Start(s)
		fmt.Println("\tПочатковий стан:")
		(*f).Print()
		t := time.Now()
		step = 0
		res := (*f).AStar((*f).field, maxSteps)
		if res == correct {
			letter := "ів"
			st := step % 10
			if st <= 1 {
				letter = ""
			} else if st <= 3 {
				letter = "и"
			}
			tt := time.Since(t)
			fmt.Printf("\n\tРішення знайдено за %v крок%s; час пошуку: %v \n", step, letter, tt)
			SolveArr = append(SolveArr, step)
			timeArr = append(timeArr, tt)
		}
	}
	fmt.Println("\nМетод A*:")
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

	//DoLDFS("H1", 20)

	DoAStar("H1", 25000)

}
