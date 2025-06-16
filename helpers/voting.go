package helpers

func CalculateScore(goal, assist, save int64) int64 {
	var score int64 = 0

	if goal != 0 {
		score += (goal * 4)
	}
	if assist != 0 {
		score += (assist * 3)
	}
	if save != 0 {
		score += (save * 5)
	}

	return score
}
