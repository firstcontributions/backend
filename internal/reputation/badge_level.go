package reputation

type BadgeLevel int

const (
	Nil BadgeLevel = iota
	Level1
	Level2
	Level3
	Level4
	Level5
	Level6
	Level7
	Level8
	Level9
	Level10
)

/**
  	i.      Level1  if > 64 lines of code
	ii.     Level2  if > 256 lines of code
	iii.    Level3  if > 1k lines of code
	iv.     Level4  4K
	v.      Level5  16K
	vi.     Level6  32K
	vii.    Level7  64K
	viii.   Level8  128K
	ix.     Level9  256K
	x.      Level10 512K
*/

var levelTargetMap = map[BadgeLevel]int{
	Nil:      0,
	Level1: 64,
	Level2: 256,
	Level3: 1000,
	Level4: 4000,
	Level5: 16000,
	Level6: 32000,
	Level7: 64000,
	Level8: 128000,
	Level9: 256000,
	Level10: 512000,
}

func GetLevelTarget(level BadgeLevel) int {
	return levelTargetMap[level]
}
func GetLevelFromPoints(points int) BadgeLevel {
	for level := Nil; level < Level10; level++ {
		if points < GetLevelTarget(level) {
			return level
		}
	}
	return Level10
}

func GetProgressPercentageToNextLevel(points int) int64 {
	currLevel := GetLevelFromPoints(points)
	prevTarget := GetLevelTarget(currLevel - 1)
	nextTarget := GetLevelTarget(currLevel)

	return int64((float64(points-prevTarget) / float64(nextTarget-prevTarget)) * 100)
}

func GetLinesOfCodeToNextLevel(points int) int64 {
	currLevel := GetLevelFromPoints(points)
	nextTarget := GetLevelTarget(currLevel)

	return int64(nextTarget - points)
}
