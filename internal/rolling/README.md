<a id="markdown-rolling" name="rolling"></a>
# rolling
[![GoDoc](https://godoc.org/github.com/asecurityteam/rolling?status.svg)](https://godoc.org/github.com/asecurityteam/rolling)
[![Build Status](https://travis-ci.com/asecurityteam/rolling.png?branch=master)](https://travis-ci.com/asecurityteam/rolling)
[![codecov.io](https://codecov.io/github/asecurityteam/rolling/coverage.svg?branch=master)](https://codecov.io/github/asecurityteam/rolling?branch=master)

**A rolling/sliding window implementation for Google-golang**

<!-- TOC -->

- [rolling](#rolling)
    - [Usage](#usage)
        - [Point Window](#point-window)
        - [Time Window](#time-window)
    - [Aggregating Windows](#aggregating-windows)
            - [Custom Aggregations](#custom-aggregations)
    - [Contributors](#contributors)
    - [License](#license)

<!-- /TOC -->

<a id="markdown-usage" name="usage"></a>
## Usage

<a id="markdown-point-window" name="point-window"></a>
### Point Window

```golang
var p = rolling.NewPointPolicy(rolling.NewWindow(5))

for x := 0; x < 5; x = x + 1 {
  p.Append(x)
}
p.Reduce(func(w Window) float64 {
  fmt.Println(w) // [ [0] [1] [2] [3] [4] ]
  return 0
})
w.Append(5)
p.Reduce(func(w Window) float64 {
  fmt.Println(w) // [ [5] [1] [2] [3] [4] ]
  return 0
})
w.Append(6)
p.Reduce(func(w Window) float64 {
  fmt.Println(w) // [ [5] [6] [2] [3] [4] ]
  return 0
})
```

The above creates a window that always contains 5 data points and then fills
it with the values 0 - 4. When the next value is appended it will overwrite
the first value. The window continuously overwrites the oldest value with the
latest to preserve the specified value count. This type of window is useful
for collecting data that have a known interval on which they are capture or
for tracking data where time is not a factor.

<a id="markdown-time-window" name="time-window"></a>
### Time Window

```golang
var p = rolling.NewTimeWindow(rolling.NewWindow(3000), time.Millisecond)
var start = time.Now()
for range time.Tick(time.Millisecond) {
  if time.Since(start) > 3*time.Second {
    break
  }
  p.Append(1)
}
```

The above creates a time window that contains 3,000 buckets where each bucket
contains, at most, 1ms of recorded data. The subsequent loop populates each
bucket with exactly one measure (the value 1) and stops when the window is full.
As time progresses, the oldest values will be removed such that if the above
code performed a `time.Sleep(3*time.Second)` then the window would be empty
again.

The choice of bucket size depends on the frequency with which data are expected
to be recorded. On each increment of time equal to the given duration the window
will expire one bucket and purge the collected values. The smaller the bucket
duration then the less data are lost when a bucket expires.

This type of bucket is most useful for collecting real-time values such as
request rates, error rates, and latencies of operations.

<a id="markdown-aggregating-windows" name="aggregating-windows"></a>
## Aggregating Windows

Each window exposes a `Reduce(func(w Window) float64) float64` method that can
be used to aggregate the data stored within. The method takes in a function
that can compute the contents of the `Window` into a single value. For
convenience, this package provides some common reductions:

```golang
fmt.Println(p.Reduce(rolling.Count))
fmt.Println(p.Reduce(rolling.Avg))
fmt.Println(p.Reduce(rolling.Min))
fmt.Println(p.Reduce(rolling.Max))
fmt.Println(p.Reduce(rolling.Sum))
fmt.Println(p.Reduce(rolling.Percentile(99.9)))
fmt.Println(p.Reduce(rolling.FastPercentile(99.9)))
```

The `Count`, `Avg`, `Min`, `Max`, and `Sum` each perform their expected
computation. The `Percentile` aggregator first takes the target percentile and
returns an aggregating function that works identically to the `Sum`, et al.

For cases of very large datasets, the `FastPercentile` can be used as a
replacement for the standard percentile calculation. This alternative version
uses the p-squared algorithm for estimating the percentile by processing
only one value at a time, in any order. The results are quite accurate but can
vary from the *actual* percentile by a small amount. It's a tradeoff of accuracy
for speed when calculating percentiles from large data sets. For more on the
p-squared algorithm see: <http://www.cs.wustl.edu/~jain/papers/ftp/psqr.pdf>.

<a id="markdown-custom-aggregations" name="custom-aggregations"></a>
#### Custom Aggregations

Any function that matches the form of `func(rolling.Window)float64` may be given
to the `Reduce` method of any window policy. The `Window` type is a named
version of `[][]float64`. Calling `len(window)` will return the number of
buckets. Each bucket is, itself, a slice of floats where `len(bucket)` is the
number of values measured within that bucket. Most aggregate will take the form
of:

```golang
func MyAggregate(w rolling.Window) float64 {
  for _, bucket := range w {
    for _, value := range bucket {
      // aggregate something
    }
  }
}
```

<a id="markdown-contributors" name="contributors"></a>
## Contributors

Pull requests, issues and comments welcome. For pull requests:

*   Add tests for new features and bug fixes
*   Follow the existing style
*   Separate unrelated changes into multiple pull requests

See the existing issues for things to start contributing.

For bigger changes, make sure you start a discussion first by creating
an issue and explaining the intended change.

Atlassian requires contributors to sign a Contributor License Agreement,
known as a CLA. This serves as a record stating that the contributor is
entitled to contribute the code/documentation/translation to the project
and is willing to have it used in distributions and derivative works
(or is willing to transfer ownership).

Prior to accepting your contributions we ask that you please follow the appropriate
link below to digitally sign the CLA. The Corporate CLA is for those who are
contributing as a member of an organization and the individual CLA is for
those contributing as an individual.

*   [CLA for corporate contributors](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=e1c17c66-ca4d-4aab-a953-2c231af4a20b)
*   [CLA for individuals](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=3f94fbdc-2fbe-46ac-b14c-5d152700ae5d)

<a id="markdown-license" name="license"></a>
## License

Copyright (c) 2017 Atlassian and others.
Apache 2.0 licensed, see [LICENSE.txt](LICENSE.txt) file.
