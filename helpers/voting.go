package helpers

func CalculateScore(scorePoint PerformancePoint, candidatePerformanceCount CandidatePerformanceCount) int64 {
	var score int64 = 0

	if candidatePerformanceCount.Goal != 0 {
		score += (scorePoint.Goal * candidatePerformanceCount.Goal)
	}
	if candidatePerformanceCount.Assist != 0 {
		score += (scorePoint.Assist * candidatePerformanceCount.Assist)
	}
	if candidatePerformanceCount.Save != 0 {
		score += (scorePoint.Save * candidatePerformanceCount.Save)
	}

	return score
}

type PerformancePoint struct {
	Goal   int64
	Assist int64
	Save   int64
}

type CandidatePerformanceCount struct {
	Goal   int64
	Assist int64
	Save   int64
}
