package rolling

import (
	"fmt"
	"testing"
)

// https://gist.github.com/cevaris/bc331cbe970b03816c6b
var epsilon = 0.00000001

func floatEquals(a float64, b float64) bool {
	return (a-b) < epsilon && (b-a) < epsilon
}

var largeEpsilon = 0.001

func floatMostlyEquals(a float64, b float64) bool {
	return (a-b) < largeEpsilon && (b-a) < largeEpsilon

}

func TestCount(t *testing.T) {
	var numberOfPoints = 100
	var w = NewWindow(numberOfPoints)
	var p = NewPointPolicy(w)
	for x := 1; x <= numberOfPoints; x = x + 1 {
		p.Append(float64(x))
	}
	var result = p.Reduce(Count)

	var expected = 100.0
	if !floatEquals(result, expected) {
		t.Fatalf("count calculated incorrectly: %f versus %f", expected, result)
	}
}

func TestCountPreallocatedWindow(t *testing.T) {
	var numberOfPoints = 100
	var w = NewPreallocatedWindow(numberOfPoints, 100)
	var p = NewPointPolicy(w)
	for x := 1; x <= numberOfPoints; x = x + 1 {
		p.Append(float64(x))
	}
	var result = p.Reduce(Count)

	var expected = 100.0
	if !floatEquals(result, expected) {
		t.Fatalf("count with prealloc window calculated incorrectly: %f versus %f", expected, result)
	}
}

func TestSum(t *testing.T) {
	var numberOfPoints = 100
	var w = NewWindow(numberOfPoints)
	var p = NewPointPolicy(w)
	for x := 1; x <= numberOfPoints; x = x + 1 {
		p.Append(float64(x))
	}
	var result = p.Reduce(Sum)

	var expected = 5050.0
	if !floatEquals(result, expected) {
		t.Fatalf("avg calculated incorrectly: %f versus %f", expected, result)
	}
}

func TestAvg(t *testing.T) {
	var numberOfPoints = 100
	var w = NewWindow(numberOfPoints)
	var p = NewPointPolicy(w)
	for x := 1; x <= numberOfPoints; x = x + 1 {
		p.Append(float64(x))
	}
	var result = p.Reduce(Avg)

	var expected = 50.5
	if !floatEquals(result, expected) {
		t.Fatalf("avg calculated incorrectly: %f versus %f", expected, result)
	}
}

func TestMax(t *testing.T) {
	var numberOfPoints = 100
	var w = NewWindow(numberOfPoints)
	var p = NewPointPolicy(w)
	for x := 1; x <= numberOfPoints; x = x + 1 {
		p.Append(100.0 - float64(x))
	}
	var result = p.Reduce(Max)

	var expected = 99.0
	if !floatEquals(result, expected) {
		t.Fatalf("max calculated incorrectly: %f versus %f", expected, result)
	}
}

func TestMin(t *testing.T) {
	var numberOfPoints = 100
	var w = NewWindow(numberOfPoints)
	var p = NewPointPolicy(w)
	for x := 1; x <= numberOfPoints; x = x + 1 {
		p.Append(float64(x))
	}
	var result = p.Reduce(Min)

	var expected = 1.0
	if !floatEquals(result, expected) {
		t.Fatalf("Min calculated incorrectly: %f versus %f", expected, result)
	}
}

func TestPercentileAggregateInterpolateWhenEmpty(t *testing.T) {
	var numberOfPoints = 0
	var w = NewWindow(numberOfPoints)
	var p = NewPointPolicy(w)
	var perc = 99.9
	var a = Percentile(perc)
	var result = p.Reduce(a)
	if !floatEquals(result, 0) {
		t.Fatalf("percentile should be zero but got %f", result)
	}
}

func TestPercentileAggregateInterpolateWhenInsufficientData(t *testing.T) {
	var numberOfPoints = 100
	var w = NewWindow(numberOfPoints)
	var p = NewPointPolicy(w)
	for x := 1; x <= numberOfPoints; x = x + 1 {
		p.Append(float64(x))
	}
	var perc = 99.9
	var a = Percentile(perc)
	var result = p.Reduce(a)

	// When there are insufficient values to satisfy the precision then the
	// percentile algorithm degenerates to a max function. In this case, we need
	// 1000 values in order to select a 99.9 but only have 100. 100 is also the
	// maximum value and will be selected as k and k+1 in the linear
	// interpolation.
	var expected = 100.0
	if !floatEquals(result, expected) {
		t.Fatalf("%f percentile calculated incorrectly: %f versus %f", perc, expected, result)
	}
}

func TestPercentileAggregateInterpolateWhenSufficientData(t *testing.T) {
	var numberOfPoints = 1000
	var w = NewWindow(numberOfPoints)
	var p = NewPointPolicy(w)
	for x := 1; x <= numberOfPoints; x = x + 1 {
		p.Append(float64(x))
	}
	var perc = 99.9
	var a = Percentile(perc)
	var result = p.Reduce(a)
	var expected = 999.5
	if !floatEquals(result, expected) {
		t.Fatalf("%f percentile calculated incorrectly: %f versus %f", perc, expected, result)
	}
}

func TestFastPercentileAggregateInterpolateWhenEmpty(t *testing.T) {
	var numberOfPoints = 0
	var w = NewWindow(numberOfPoints)
	var p = NewPointPolicy(w)
	var perc = 99.9
	var a = FastPercentile(perc)
	var result = p.Reduce(a)
	if !floatEquals(result, 0) {
		t.Fatalf("fast percentile should be zero but got %f", result)
	}
}

func TestFastPercentileAggregateInterpolateWhenSufficientData(t *testing.T) {
	// Using a larger dataset so that the algorithm can converge on the
	// correct value. Smaller datasets where the value might be interpolated
	// linearly in the typical percentile calculation results in larger error
	// in the result. This is acceptable so long as the estimated value approaches
	// the correct value as more data are given.
	var numberOfPoints = 10000
	var w = NewWindow(numberOfPoints)
	var p = NewPointPolicy(w)
	for x := 1; x <= numberOfPoints; x = x + 1 {
		p.Append(float64(x))
	}
	var perc = 99.9
	var a = FastPercentile(perc)
	var result = p.Reduce(a)
	var expected = 9990.0
	if !floatEquals(result, expected) {
		t.Fatalf("%f percentile calculated incorrectly: %f versus %f", perc, expected, result)
	}
}

func TestFastPercentileAggregateUsingPSquaredDataSet(t *testing.T) {
	var numberOfPoints = 20
	var w = NewWindow(numberOfPoints)
	var p = NewPointPolicy(w)
	p.Append(0.02)
	p.Append(0.15)
	p.Append(0.74)
	p.Append(0.83)
	p.Append(3.39)
	p.Append(22.37)
	p.Append(10.15)
	p.Append(15.43)
	p.Append(38.62)
	p.Append(15.92)
	p.Append(34.60)
	p.Append(10.28)
	p.Append(1.47)
	p.Append(0.40)
	p.Append(0.05)
	p.Append(11.39)
	p.Append(0.27)
	p.Append(0.42)
	p.Append(0.09)
	p.Append(11.37)
	var perc = 50.0
	var a = FastPercentile(perc)
	var result = p.Reduce(a)
	var expected = 4.44
	if !floatMostlyEquals(result, expected) {
		t.Fatalf("%f percentile calculated incorrectly: %f versus %f", perc, expected, result)
	}
}

var aggregateResult float64

type policy interface {
	Append(float64)
	Reduce(func(Window) float64) float64
}
type aggregateBench struct {
	inserts       int
	policy        policy
	aggregate     func(Window) float64
	aggregateName string
}

func BenchmarkAggregates(b *testing.B) {
	var baseCases = []*aggregateBench{
		{aggregate: Sum, aggregateName: "sum"},
		{aggregate: Min, aggregateName: "min"},
		{aggregate: Max, aggregateName: "max"},
		{aggregate: Avg, aggregateName: "avg"},
		{aggregate: Count, aggregateName: "count"},
		{aggregate: Percentile(50.0), aggregateName: "p50"},
		{aggregate: Percentile(99.9), aggregateName: "p99.9"},
		{aggregate: FastPercentile(50.0), aggregateName: "fp50"},
		{aggregate: FastPercentile(99.9), aggregateName: "fp99.9"},
	}
	var insertions = []int{1, 1000, 10000, 100000}
	var benchCases = make([]*aggregateBench, 0, len(baseCases)*len(insertions))
	for _, baseCase := range baseCases {
		for _, inserts := range insertions {
			var w = NewWindow(inserts)
			var p = NewPointPolicy(w)
			for x := 1; x <= inserts; x = x + 1 {
				p.Append(float64(x))
			}
			benchCases = append(benchCases, &aggregateBench{
				inserts:       inserts,
				aggregate:     baseCase.aggregate,
				aggregateName: baseCase.aggregateName,
				policy:        p,
			})
		}
	}

	for _, benchCase := range benchCases {
		b.Run(fmt.Sprintf("Aggregate:%s-DataPoints:%d", benchCase.aggregateName, benchCase.inserts), func(bt *testing.B) {
			var result float64
			bt.ResetTimer()
			for n := 0; n < bt.N; n = n + 1 {
				result = benchCase.policy.Reduce(benchCase.aggregate)
			}
			aggregateResult = result
		})
	}
}
