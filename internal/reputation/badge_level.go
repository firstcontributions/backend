package reputation

type BadgeLevel int

const (
	Nil BadgeLevel = iota
	Bronze
	BronzeX2
	BronzeX3
	Silver
	SilverX2
	SilverX3
	Gold
	GoldX2
	GoldX3
	Diamond
)

/**
  	i.      1X Bronze  if > 64 lines of code
	ii.     2X Bronze  if > 256 lines of code
	iii.    3X Bronze  if > 1k lines of code
	iv.     1X Silver 4K
	v.      2X Silver 16K
	vi.     3X Silver 32K
	vii.    1X Gold 64K
	viii.   2X Gold 128K
	ix.     3X Gold 256K
	x.      1X Diamond 512K
*/

var levelTargetMap = map[BadgeLevel]int{
	Nil:      0,
	Bronze:   64,
	BronzeX2: 265,
	BronzeX3: 1000,
	Silver:   4000,
	SilverX2: 16000,
	SilverX3: 32000,
	Gold:     64000,
	GoldX2:   128000,
	GoldX3:   256000,
	Diamond:  512000,
}

func GetLevelTarget(level BadgeLevel) int {
	return levelTargetMap[level]
}
func GetLevelFromPoints(points int) BadgeLevel {
	for level := Nil; level < Diamond; level++ {
		if points < GetLevelTarget(level) {
			return level
		}
	}
	return Diamond
}

func GetProgressPercentageToNextLevel(points int) int64 {
	currLevel := GetLevelFromPoints(points)
	prevTarget := GetLevelTarget(currLevel)
	nextTarget := GetLevelTarget(currLevel + 1)

	return int64((float64(points-prevTarget) / float64(nextTarget-prevTarget)) * 100)
}
