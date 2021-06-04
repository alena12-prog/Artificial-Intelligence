package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

type node struct {
	name string
	to   map[*node]int
}

func Create(name string) *node {
	return &node{name, map[*node]int{}}
}

func (n *node) Connect(n1 *node, distance int) *node {
	(*n).to[n1] = distance
	(*n1).to[n] = distance
	return n1
}

func (n *node) ToChildren(prevStr string, name string, dist, depth int) (int, string) {
	if depth > 0 {
		prevStr += " -> " + (*n).name
		if (*n).name == name {
			return dist, prevStr
		}
		mind := math.MaxInt32
		ps := ""
		for nm, v := range (*n).to {
			r, s := (*nm).ToChildren(prevStr, name, dist+v, depth-1)
			//fmt.Println(strings.Repeat("\t", 15-depth), (*n).name, v, (*nm).name)
			if r < mind {
				mind = r
				ps = s
			}
		}
		return mind, ps
	} else {
		return math.MaxInt32, ""
	}
}

func (n *node) To(name string, depth int) {
	r, s := (*n).ToChildren("", name, 0, depth)
	fmt.Println(s)
	fmt.Println("Відстань:", r)
}

func (n *node) Find(name string) *node {
	exc := []*node{}
	exc = append(exc, n)
	return (*n).FindInChildren(name, &exc)
}

func (n *node) FindInChildren(name string, exceptionArr *[]*node) *node {
	for nn := range (*n).to {
		if (*nn).name == name {
			return nn
		}
		b := false
		for _, v := range *exceptionArr {
			if nn == v {
				b = true
			}
		}
		if !b {
			(*exceptionArr) = append((*exceptionArr), nn)
			res := (*nn).FindInChildren(name, exceptionArr)
			if res != nil {
				return res
			}
		}
	}
	return nil
}

func FromTxt() *node {
	str, err := ioutil.ReadFile("data.txt")
	if err != nil {
		panic(err)
	}
	dt := strings.Split(string(str), "\n")
	var start *node
	if len(dt) > 0 {
		dt[0] = strings.Trim(dt[0], "\r")
		dtt := strings.Split(dt[0], " ")
		start = Create(dtt[0])
		l, _ := strconv.Atoi(dtt[1])
		start.Connect(Create(dtt[2]), l)
	}
	for _, v := range dt {
		v = strings.Trim(v, "\r")
		dtt := strings.Split(v, " ")
		c1 := start.Find(dtt[0])
		if c1 == nil {
			c1 = Create(dtt[0])
		}
		c2 := start.Find(dtt[2])
		if c2 == nil {
			c2 = Create(dtt[2])
		}
		l, _ := strconv.Atoi(dtt[1])
		c1.Connect(c2, l)
		fmt.Println((*c1).name, l, (*c2).name)
	}
	return start
}

func BuildMap() *node {
	return FromTxt()
}

func main() {
	start := BuildMap()
	fmt.Println(start.Find("Луганcьк\r"))
	start.Find("Луганcьк").To("Львів", 15)
	s := ""
	fmt.Scanln(&s)
}
