package main

import (
	"sort"
	"time"
)

type positions struct {
	Car       int
	Lap       int
	Lastseen time.Time
}

type intsort struct {
	key   int64
	value int64
}
type sortbySecond []*intsort

func (d sortbySecond) Len() int {
	return len(d)
}
func (d sortbySecond) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
func (d sortbySecond) Less(i, j int) bool {
	return d[i].value < d[j].value
}

func sorty(cars []positions) []positions {
	// sort by lap then by last seen
	groups := groupy(cars)
	results := make([]positions, 0)
	// sort groups first!
	groupsort := make([]int, 0)
	for _, v := range groups {
		groupsort = append(groupsort, v[0].Lap)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(groupsort)))
	for _, group := range groupsort {
		lapsort := sortg(groups[group])
		for _, v := range lapsort {
			results = append(results, v)
		}
	}

	return results
}

func groupy(cars []positions) map[int][]positions {
	// group by lap
	groups := make(map[int][]positions)
	for _, v := range cars {
		// make group if we havent seen it before
		if _, ok := groups[v.Lap]; !ok {
			groups[v.Lap] = make([]positions, 0)
		}
		groups[v.Lap] = append(groups[v.Lap], v)
	}
	//return [positions{1,1,time.Now()}]
	return groups

}

func sortg(group []positions) []positions {
	// for each group sort by time
	sortbytime := make(sortbySecond, 0)
	results := make([]positions, 0)
	for k, v := range group {
		sortbytime = append(sortbytime, &intsort{int64(k), v.Lastseen.UnixNano()})
	}
	sort.Sort(sortbytime)
	for _, v := range sortbytime {
		results = append(results, group[v.key])
	}
	return results
}
