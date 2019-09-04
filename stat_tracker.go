package main

import "math"

// StatTracker is a struct which keeps track of arbitrary intger metrics and provides some statistical helpers as well
type StatTracker struct {
	Total     int64
	StdDev    float64
	Mean      float64
	Median    int64
	Min       int64
	Max       int64
	Frames    []int64
	maxFrames int
}

// NewStatTracker returns a new instance of a StatTracker with a provided maximum number of frames
func NewStatTracker(maxFrames int) StatTracker {
	if maxFrames < 10 {
		maxFrames = 144
	}
	return StatTracker{
		Min:       9223372036854775807,
		Frames:    []int64{},
		maxFrames: maxFrames,
	}
}

// Update Calculates various statistical figures based off of the data currently in its frames
func (st StatTracker) Update(frames []int64) {
	st.Frames = append(frames, st.Frames...)
	newLen := len(st.Frames)
	if newLen > st.maxFrames {
		st.Frames = st.Frames[:st.maxFrames]
	}
	st.Total = sum(st.Frames)
	min, max := minMax(st.Frames)
	st.Min = min
	st.Max = max
	st.calcMean()
	st.calcDeviation()
}

func (st StatTracker) calcMean() {
	mean := float64(st.Total) / float64(len(st.Frames))
	st.Mean = mean
}

func (st StatTracker) calcDeviation() {
	mean := st.Mean
	var sumSQDiffs float64
	for _, frame := range st.Frames {
		diff := float64(frame) - mean
		sumSQDiffs += diff * diff
	}
	stdDev := math.Sqrt((sumSQDiffs / float64(len(st.Frames))))
	st.StdDev = stdDev
}

func sum(vec []int64) int64 {
	var sum int64
	for _, frame := range vec {
		sum += frame
	}
	return sum
}

func minMax(frames []int64) (int64, int64) {
	var min int64
	var max int64
	min = 9223372036854775807
	for _, frame := range frames {
		if frame < min {
			min = frame
		}
		if frame > max {
			max = frame
		}
	}
	return min, max
}
