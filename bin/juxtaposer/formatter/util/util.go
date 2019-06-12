package util

import (
	"math"
	"time"

	formatter_api "github.com/cyberark/secretless-broker/bin/juxtaposer/formatter/api"
)

func GetSuccessPercentage(timingInfo *formatter_api.BackendTiming) float64 {
	successfulRounds := timingInfo.Count - len(timingInfo.Errors)
	return (float64(successfulRounds) / float64(timingInfo.Count)) * 100
}

func GetAverageDuration(timingInfo *formatter_api.BackendTiming) time.Duration {
	averageDuration := 0 * time.Second
	successfulRounds := timingInfo.Count - len(timingInfo.Errors)
	if successfulRounds > 0 {
		averageDuration = time.Duration(int64(timingInfo.Duration) /
			int64(successfulRounds))
	}

	return averageDuration
}

func getMappedDataPointCount(mappedCounts *map[int]int) int {
	countOfDataPoints := 0
	for _, mappedCount := range *mappedCounts {
		countOfDataPoints += mappedCount
	}

	return countOfDataPoints
}

func GetMean(mappedCounts *map[int]int) float64 {
	if len(*mappedCounts) == 0 {
		return 0.0
	}

	total := 0.0
	for valueAmount, occurrences := range *mappedCounts {
		total += float64(valueAmount * occurrences)
	}

	return total / float64(getMappedDataPointCount(mappedCounts))
}

func GetStandardDeviation(mappedCounts *map[int]int) float64 {
	if len(*mappedCounts) == 0 {
		return 0.0
	}

	mean := GetMean(mappedCounts)
	totalDeviation := 0.0
	for valueAmount, occurrences := range *mappedCounts {
		deviation := (float64(valueAmount) - mean) * (float64(valueAmount) - mean)
		totalDeviation += deviation * float64(occurrences)
	}
	standardDeviation := math.Pow(totalDeviation/float64(getMappedDataPointCount(mappedCounts)), 0.5)

	return standardDeviation
}

func GetConfidenceInterval90(mappedCounts *map[int]int) (lowerBound float64, upperBound float64) {
	// http://mathworld.wolfram.com/ConfidenceInterval.html
	confidenceDeviation := 1.64485 // 90% confidence for the mean

	mean := GetMean(mappedCounts)
	deviation := GetStandardDeviation(mappedCounts) / math.Sqrt(float64(getMappedDataPointCount(mappedCounts)))
	return mean - deviation*confidenceDeviation, mean + deviation*confidenceDeviation
}

func GetThresholdBreachedPercent(mappedCounts *map[int]int, thresholdPercent int) (percent float64) {
	count := 0
	for valueAmount, occurrences := range *mappedCounts {
		if valueAmount >= thresholdPercent {
			count += occurrences
		}
	}
	return float64(count) / float64(getMappedDataPointCount(mappedCounts)) * 100.0
}