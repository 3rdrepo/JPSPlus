package jpsplus

type JumpDistancesAndGoalBounds struct {
	blockedDirectionBitfield int
	jumpDistance             [8]int
	bounds                   [8][4]int
}
type JumpDistancesAndGoalBoundsMap [][]JumpDistancesAndGoalBounds

var DefaultJumpDistancesAndGoalBounds *JumpDistancesAndGoalBoundsMap
