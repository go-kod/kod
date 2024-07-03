package ratelimit

import (
	"context"
	"errors"
	"math"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/samber/lo"
	"github.com/shirou/gopsutil/v3/process"

	"github.com/go-kod/kod/internal/rolling"
)

// ErrLimitExceed is returned when the rate limiter is
// triggered and the request is rejected due to limit exceeded.
var ErrLimitExceed = errors.New("rate limit exceeded")

var (
	gCPU  int64
	decay = 0.95
)

type (
	cpuGetter func() int64
)

func init() {
	go cpuproc()
}

// cpu = cpuᵗ⁻¹ * decay + cpuᵗ * (1 - decay)
func cpuproc() {
	ticker := time.NewTicker(time.Millisecond * 500) // same to cpu sample rate
	defer func() {
		ticker.Stop()
		if err := recover(); err != nil {
			go cpuproc()
		}
	}()

	p := lo.Must(process.NewProcess(int32(os.Getpid())))
	// EMA algorithm: https://blog.csdn.net/m0_38106113/article/details/81542863
	for range ticker.C {
		usage, _ := p.Percent(0)
		prevCPU := atomic.LoadInt64(&gCPU)
		curCPU := int64(float64(prevCPU)*decay + float64(usage)*(1.0-decay))
		atomic.StoreInt64(&gCPU, curCPU)
	}
}

// stat contains the metrics snapshot of bbr.
type stat struct {
	CPU         int64
	InFlight    int64
	MaxInFlight int64
	MinRt       int64
	MaxPass     int64
}

// counterCache is used to cache maxPASS and minRt result.
// Value of current bucket is not counted in real time.
// Cache time is equal to a bucket duration.
type counterCache struct {
	val  int64
	time time.Time
}

// ratelimitOptions of bbr limiter.
type ratelimitOptions struct {
	// WindowSize defines time duration per window
	Window time.Duration
	// BucketNum defines bucket number for each window
	Bucket int
	// CPUThreshold
	CPUThreshold int64
	// CPUQuota
	CPUQuota float64
}

// Ratelimit implements bbr-like limiter.
// It is inspired by sentinel.
// https://github.com/alibaba/Sentinel/wiki/%E7%B3%BB%E7%BB%9F%E8%87%AA%E9%80%82%E5%BA%94%E9%99%90%E6%B5%81
type Ratelimit struct {
	cpu             cpuGetter
	passStat        *rolling.TimePolicy
	rtStat          *rolling.TimePolicy
	inFlight        int64
	bucketPerSecond int64
	bucketDuration  time.Duration

	// prevDropTime defines previous start drop since initTime
	prevDropTime atomic.Value
	maxPASSCache atomic.Value
	minRtCache   atomic.Value

	opts ratelimitOptions
}

// NewLimiter returns a bbr limiter
func NewLimiter(_ context.Context, _ string) *Ratelimit {
	opt := ratelimitOptions{
		Window:       time.Second * 10,
		Bucket:       100,
		CPUThreshold: 80,
	}

	bucketDuration := opt.Window / time.Duration(opt.Bucket)
	passStat := rolling.NewTimePolicy(rolling.NewWindow(opt.Bucket), opt.Window)
	rtStat := rolling.NewTimePolicy(rolling.NewWindow(opt.Bucket), opt.Window)

	limiter := &Ratelimit{
		opts:            opt,
		passStat:        passStat,
		rtStat:          rtStat,
		bucketDuration:  bucketDuration,
		bucketPerSecond: int64(time.Second / bucketDuration),
		cpu:             func() int64 { return atomic.LoadInt64(&gCPU) },
	}

	if opt.CPUQuota != 0 {
		// if cpuQuota is set, use new cpuGetter,Calculate the real CPU value based on the number of CPUs and Quota.
		limiter.cpu = func() int64 {
			return int64(float64(atomic.LoadInt64(&gCPU)) * float64(runtime.NumCPU()) / opt.CPUQuota)
		}
	}

	return limiter
}

func (l *Ratelimit) maxPASS() int64 {
	passCache := l.maxPASSCache.Load()
	if passCache != nil {
		ps := passCache.(*counterCache)
		if l.timespan(ps.time) < 1 {
			return ps.val
		}
	}
	rawMaxPass := int64(l.passStat.Reduce(func(iterator rolling.Window) float64 {
		result := 1.0
		for i := 1; i <= len(iterator) && i < l.opts.Bucket; i++ {
			bucket := iterator[i-1]
			count := 0.0
			for _, p := range bucket {
				count += p
			}
			result = math.Max(result, count)
		}
		return result
	}))
	l.maxPASSCache.Store(&counterCache{
		val:  rawMaxPass,
		time: time.Now(),
	})
	return rawMaxPass
}

// timespan returns the passed bucket count
// since lastTime, if it is one bucket duration earlier than
// the last recorded time, it will return the BucketNum.
func (l *Ratelimit) timespan(lastTime time.Time) int {
	v := int(time.Since(lastTime) / l.bucketDuration)
	if v > -1 {
		return v
	}
	return l.opts.Bucket
}

func (l *Ratelimit) minRT() int64 {
	rtCache := l.minRtCache.Load()
	if rtCache != nil {
		rc := rtCache.(*counterCache)
		if l.timespan(rc.time) < 1 {
			return rc.val
		}
	}
	rawMinRT := int64(math.Ceil(l.rtStat.Reduce(func(iterator rolling.Window) float64 {
		result := math.MaxFloat64
		for i := 1; i <= len(iterator) && i < l.opts.Bucket; i++ {
			bucket := iterator[i-1]
			if len(bucket) == 0 {
				continue
			}
			total := 0.0
			for _, p := range bucket {
				total += p
			}
			avg := total / float64(len(bucket))
			result = math.Min(result, avg)
		}
		return result
	})))
	if rawMinRT <= 0 {
		rawMinRT = 1
	}
	l.minRtCache.Store(&counterCache{
		val:  rawMinRT,
		time: time.Now(),
	})
	return rawMinRT
}

func (l *Ratelimit) maxInFlight() int64 {
	return int64(math.Floor(float64(l.maxPASS()*l.minRT()*l.bucketPerSecond)/1000.0) + 0.5)
}

func (l *Ratelimit) shouldDrop() bool {
	now := time.Duration(time.Now().UnixNano())
	if l.cpu() < l.opts.CPUThreshold {
		// current cpu payload below the threshold
		prevDropTime, _ := l.prevDropTime.Load().(time.Duration)
		if prevDropTime == 0 {
			// haven't start drop,
			// accept current request
			return false
		}
		if time.Duration(now-prevDropTime) <= time.Second {
			// just start drop one second ago,
			// check current inflight count
			inFlight := atomic.LoadInt64(&l.inFlight)
			return inFlight > 1 && inFlight > l.maxInFlight()
		}
		l.prevDropTime.Store(time.Duration(0))
		return false
	}
	// current cpu payload exceeds the threshold
	inFlight := atomic.LoadInt64(&l.inFlight)
	drop := inFlight > 1 && inFlight > l.maxInFlight()
	if drop {
		prevDrop, _ := l.prevDropTime.Load().(time.Duration)
		if prevDrop != 0 {
			// already started drop, return directly
			return drop
		}
		// store start drop time
		l.prevDropTime.Store(now)
	}
	return drop
}

// Stat tasks a snapshot of the bbr limiter.
func (l *Ratelimit) Stat() stat {
	return stat{
		CPU:         l.cpu(),
		MinRt:       l.minRT(),
		MaxPass:     l.maxPASS(),
		MaxInFlight: l.maxInFlight(),
		InFlight:    atomic.LoadInt64(&l.inFlight),
	}
}

// Allow checks all inbound traffic.
// Once overload is detected, it raises limit.ErrLimitExceed error.
func (l *Ratelimit) Allow() (func(), error) {
	if l.shouldDrop() {
		return nil, ErrLimitExceed
	}
	atomic.AddInt64(&l.inFlight, 1)
	start := time.Now().UnixNano()
	ms := float64(time.Millisecond)
	return func() {
		//nolint
		if rt := int64(math.Ceil(float64(time.Now().UnixNano()-start)) / ms); rt > 0 {
			l.rtStat.Append(float64(rt))
		}
		atomic.AddInt64(&l.inFlight, -1)
		l.passStat.Append(float64(1))
	}, nil
}
