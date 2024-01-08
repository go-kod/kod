package rolling

import (
	"math"
	"sort"
	"sync"
)

// Count returns the number of elements in a window.
func Count(w Window) float64 {
	result := 0
	for _, bucket := range w {
		result += len(bucket)
	}
	return float64(result)
}

// Sum the values within the window.
func Sum(w Window) float64 {
	var result = 0.0
	for _, bucket := range w {
		for _, p := range bucket {
			result = result + p
		}
	}
	return result
}

// Avg the values within the window.
func Avg(w Window) float64 {
	var result = 0.0
	var count = 0.0
	for _, bucket := range w {
		for _, p := range bucket {
			result = result + p
			count = count + 1
		}
	}
	return result / count
}

// Min the values within the window.
func Min(w Window) float64 {
	var result = 0.0
	var started = true
	for _, bucket := range w {
		for _, p := range bucket {
			if started {
				result = p
				started = false
				continue
			}
			if p < result {
				result = p
			}
		}
	}
	return result
}

// Max the values within the window.
func Max(w Window) float64 {
	var result = 0.0
	var started = true
	for _, bucket := range w {
		for _, p := range bucket {
			if started {
				result = p
				started = false
				continue
			}
			if p > result {
				result = p
			}
		}
	}
	return result
}

// Percentile returns an aggregating function that computes the
// given percentile calculation for a window.
func Percentile(perc float64) func(w Window) float64 {
	var values []float64
	var lock = &sync.Mutex{}
	return func(w Window) float64 {
		lock.Lock()
		defer lock.Unlock()

		values = values[:0]
		for _, bucket := range w {
			values = append(values, bucket...)
		}
		if len(values) < 1 {
			return 0.0
		}
		sort.Float64s(values)
		var position = (float64(len(values))*(perc/100) + .5) - 1
		var k = int(math.Floor(position))
		var f = math.Mod(position, 1)
		if f == 0.0 {
			return values[k]
		}
		var plusOne = k + 1
		if plusOne > len(values)-1 {
			plusOne = k
		}
		return ((1 - f) * values[k]) + (f * values[plusOne])
	}
}

// FastPercentile implements the pSquare percentile estimation
// algorithm for calculating percentiles from streams of data
// using fixed memory allocations.
func FastPercentile(perc float64) func(w Window) float64 {
	perc = perc / 100.0
	return func(w Window) float64 {
		var initalObservations = make([]float64, 0, 5)
		var q [5]float64
		var n [5]int
		var nPrime [5]float64
		var dnPrime [5]float64
		var observations uint64
		for _, bucket := range w {
			for _, v := range bucket {

				observations = observations + 1
				// Record first five observations
				if observations < 6 {
					initalObservations = append(initalObservations, v)
					continue
				}
				// Before proceeding beyond the first five, process them.
				if observations == 6 {
					bubbleSort(initalObservations)
					for offset := range q {
						q[offset] = initalObservations[offset]
						n[offset] = offset
					}
					nPrime[0] = 0
					nPrime[1] = 2 * perc
					nPrime[2] = 4 * perc
					nPrime[3] = 2 + 2*perc
					nPrime[4] = 4
					dnPrime[0] = 0
					dnPrime[1] = perc / 2
					dnPrime[2] = perc
					dnPrime[3] = (1 + perc) / 2
					dnPrime[4] = 1
				}
				var k int // k is the target cell to increment
				switch {
				case v < q[0]:
					q[0] = v
					k = 0
				case q[0] <= v && v < q[1]:
					k = 0
				case q[1] <= v && v < q[2]:
					k = 1
				case q[2] <= v && v < q[3]:
					k = 2
				case q[3] <= v && v <= q[4]:
					k = 3
				case v > q[4]:
					q[4] = v
					k = 3
				}
				for x := k + 1; x < 5; x = x + 1 {
					n[x] = n[x] + 1
				}
				nPrime[0] = nPrime[0] + dnPrime[0]
				nPrime[1] = nPrime[1] + dnPrime[1]
				nPrime[2] = nPrime[2] + dnPrime[2]
				nPrime[3] = nPrime[3] + dnPrime[3]
				nPrime[4] = nPrime[4] + dnPrime[4]
				for x := 1; x < 4; x = x + 1 {
					var d = nPrime[x] - float64(n[x])
					if (d >= 1 && (n[x+1]-n[x]) > 1) ||
						(d <= -1 && (n[x-1]-n[x]) < -1) {
						var s = sign(d)
						var si = int(s)
						var nx = float64(n[x])
						var nxPlusOne = float64(n[x+1])
						var nxMinusOne = float64(n[x-1])
						var qx = q[x]
						var qxPlusOne = q[x+1]
						var qxMinusOne = q[x-1]
						var parab = q[x] + (s/(nxPlusOne-nxMinusOne))*((nx-nxMinusOne+s)*(qxPlusOne-qx)/(nxPlusOne-nx)+(nxPlusOne-nx-s)*(qx-qxMinusOne)/(nx-nxMinusOne))
						if qxMinusOne < parab && parab < qxPlusOne {
							q[x] = parab
						} else {
							q[x] = q[x] + s*((q[x+si]-q[x])/float64(n[x+si]-n[x]))
						}
						n[x] = n[x] + si
					}
				}

			}
		}

		if observations < 1 {
			return 0.0
		}
		// If we have less than five values then degenerate into a max function.
		// This is a reasonable value for data sets this small.
		if observations < 5 {
			bubbleSort(initalObservations)
			return initalObservations[len(initalObservations)-1]
		}
		return q[2]
	}
}

func sign(v float64) float64 {
	if v < 0 {
		return -1
	}
	return 1
}

// using bubblesort because we're only working with datasets of 5 or fewer
// elements.
func bubbleSort(s []float64) {
	for range s {
		for x := 0; x < len(s)-1; x = x + 1 {
			if s[x] > s[x+1] {
				s[x], s[x+1] = s[x+1], s[x]
			}
		}
	}
}
