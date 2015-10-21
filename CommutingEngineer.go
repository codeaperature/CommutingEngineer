/*

https://www.codeeval.com/open_challenges/90/

CHALLENGE DESCRIPTION:

Commuters in the bay area who commute to and from South Bay spend on
average 2-3 hours of valuable time getting to and from work every day.
That's why startups like Mashery, Flurry, New Relic and Glassdoor
have called the San Francisco Peninsula their home.

Today, we're visiting some of area's fastest growing startups and
would like to find the shortest possible distance to visit each
company once starting from the CodeEval offices at 1355 Market
Street.

Solving the following challenge means finding the shortest
possible route which visits each coordinate once starting
from the point 1. You may read more about the Travelling
Salesman problem

On the map you can see the best route for 6 coordinates.
But we've added 4 more for the offices of Mashery, Flurry,
New Relic and Glassdoor. So you need to find the best
route for 10 coordinates.


input sample:

1 | CodeEval 1355 Market St, SF (37.7768016, -122.4169151)
2 | Yelp 706 Mission St, SF (37.7860105, -122.4025377)
3 | Square 110 5th St, SF (37.7821494, -122.4058960)
4 | Airbnb 99 Rhode Island St, SF (37.7689269, -122.4029053)
5 | Dropbox 185 Berry St, SF (37.7768800, -122.3911496)
6 | Zynga 699 8th St, SF (37.7706628, -122.4040139)

7 | Mashery 717 Market St, SF (37.7870361, -122.4039444)
8 | Flurry 3060 3rd St, SF (37.7507903, -122.3877184)
9 | New Relic 188 Spear St, SF (37.7914417, -122.3927229)
10 | Glassdoor 1 Harbor Drive, Sausalito (37.8672841, -122.5010216)


output for 1st 6
1
3
2
5
6
4



*/

package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

/*
var data =

`1 | CodeEval 1355 Market St, SF (37.7768016, -122.4169151)`
2 | Yelp 706 Mission St, SF (37.7860105, -122.4025377)
3 | Square 110 5th St, SF (37.7821494, -122.4058960)
4 | Airbnb 99 Rhode Island St, SF (37.7689269, -122.4029053)
5 | Dropbox 185 Berry St, SF (37.7768800, -122.3911496)
6 | Zynga 699 8th St, SF (37.7706628, -122.4040139)
`
*/

type LocType struct {
	id int
	i  int
	y  int
	x  int
	d  []uint64
}

type SolutionType struct {
	totalCost uint64
	cnum      []int
}

//var solution []SolutionType
var companies int
var locs []LocType

func main() {

	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		return
	}

	lines := strings.Split(string(b), "\n")

	companies = len(lines)

	//	fmt.Println("lines = ", lines)

	if len(lines[companies-1]) == 0 {
		lines = lines[0 : companies-1]
		companies--
	}
	//	fmt.Println("lines = ", lines)

	locs = make([]LocType, companies, companies)
	//	solution := make([]SolutionType, companies, companies)
	//	fmt.Println("n = ", companies)

	for i, locstr := range lines {
		//		fmt.Println(locstr)
		junk1 := strings.Split(locstr, " | ")
		//	fmt.Printf("1 - >%s<\n", junk1[0])
		junk2 := strings.Split(junk1[1], "(37.")
		//	fmt.Printf("2 - >%s<\n", junk2[1])
		junk3 := strings.Split(junk2[1], ", -122.")
		junk4 := strings.Split(junk3[1], ")")
		//	fmt.Printf("3 - >%s<, >%s<\n", junk3[0], junk4[0])
		locs[i].id, err = strconv.Atoi(junk1[0])
		locs[i].i = i
		locs[i].y, err = strconv.Atoi(junk3[0])
		locs[i].x, err = strconv.Atoi(junk4[0])
		if err != nil {
			fmt.Println("error = ", err)
			return
		}

		locs[i].d = make([]uint64, companies, companies)
		//		fmt.Println("locs: ", locs[i])
	}

	// let's go with some approximations -- the eng is in SF bay area where
	// 1 deg lat/y  => ~69.17 miles
	// 1 deg long/x => ~54.73 miles
	// everything is in bounds: -122 >= x > -123 &&  37 <= y < 38
	city := make([]int, companies, companies)
	for i := 0; i < companies; i++ {
		city[i] = i
		for j := 0; j < i; j++ {
			if i == j {
				locs[i].d[j] = 0
			} else {
				var uy, ux uint64
				if dx := locs[i].x - locs[j].x; dx < 0 {
					ux = uint64(-dx)
				} else {
					ux = uint64(dx)
				}
				ux *= 5473
				if dy := locs[i].y - locs[j].y; dy < 0 {
					uy = uint64(-dy)
				} else {
					uy = uint64(dy)
				}
				uy *= 6917
				locs[i].d[j] = uint64(math.Sqrt(float64(ux*ux + uy*uy)))
				//				fmt.Println("i =", i, "j =", j, " => ", locs[i].x, "-", locs[j].x, "dx =", ux, " => ", locs[i].y, "-", locs[j].y, "dy =", uy, "l=", locs[i].d[j])
				locs[j].d[i] = locs[i].d[j]
			}

			//			fmt.Println(locs[i])

		}
	}
	/*
		for i := 0; i < companies; i++ {
			fmt.Println(locs[i])
		}

		for i := 0; i < companies; i++ {
			fmt.Printf("%2d: ", i)
			for j := 0; j < companies; j++ {
				fmt.Printf("%10d ", locs[i].d[j])

			}
			fmt.Println()
		}
	*/
	city = append(city[:0], city[1:]...)
	//	fmt.Println(city)
	fsol := getMinDistRt(0, city)
	//	fmt.Println(fsol)
	//	fmt.Println(locs)

	for i := companies - 1; i >= 0; i-- {
		fmt.Println(locs[fsol.cnum[i]].id)
	}

}

func getMinDistRt(city int, remaining []int) SolutionType {
	// at end?
	var sol SolutionType
	locLen := len(remaining)
	//sol.cnum = make([]int, companies, companies)
	//	fmt.Println("remaining [", locLen, "]: ", remaining)
	if locLen == 0 {
		sol.cnum = make([]int, companies, companies)
		sol.cnum[locLen] = city
		sol.totalCost = 0

		//		fmt.Println("isol: ", sol)
		return sol
	}

	// loop over remaining cities
	minCostTotal := uint64(0xffffffffffffffff)
	var temSol SolutionType
	var temCost uint64
	for i, dest := range remaining {
		remainingLoc := append([]int{}, remaining...)
		remainingLoc = append(remainingLoc[:i], remainingLoc[i+1:]...)
		//	fmt.Println("rml: ", remainingLoc)

		temSol = getMinDistRt(dest, remainingLoc)
		remainingLoc = nil
		//		fmt.Println("sol: ", sol)
		temCost = temSol.totalCost + locs[city].d[dest]

		//		fmt.Println("city: ", city, "tsol: ", temSol)
		if minCostTotal > temCost {
			minCostTotal = temCost
			sol = temSol
		}

	}
	sol.totalCost = minCostTotal
	sol.cnum[locLen] = city

	//	fmt.Println("city: ", city, "xsol: ", sol)

	return sol
}
