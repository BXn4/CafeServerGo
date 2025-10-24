package player

import (
	"cafego/internal/utils"
	"math"
)

func (p *Player) GetXP() int {
	return p.XP
}

func (p *Player) AddXP(amount int) {
	if p.IsLevelUp(amount) {
		// println("LEVEL UP!")
		nextLevel := p.GetLevel() + 1
		if utils.GetLevelCashReward(nextLevel) > 0 {
			p.AddCash(utils.GetLevelCashReward(nextLevel))
		}

		if utils.GetLevelGoldReward(nextLevel) > 0 {
			p.AddGold(utils.GetLevelGoldReward(nextLevel))
		}
	}

	p.XP += amount
}

func (p *Player) SetXP(amount int) {
	p.XP = amount
}

func (p *Player) GetLevel() int {
	if p.GetXP() < 90 {
		return int(math.Floor(math.Pow(math.Floor(float64(p.GetXP())/10.0), 1.0/2.0)))
	}
	return int(math.Floor(math.Pow(math.Floor(float64(p.GetXP())/5.0), 1.0/3.72)))
}

func (p *Player) IsLevelUp(givenXP int) bool {
	return GetLevelByXP(p.GetXP()+givenXP) > GetLevelByXP(p.GetXP())
}

func GetLevelByXP(xp int) int {
	if xp < 90 {
		return int(math.Floor(math.Pow(math.Floor(float64(xp)/10.0), 1.0/2.0)))
	}
	return int(math.Floor(math.Pow(math.Floor(float64(xp)/5.0), 1.0/3.72)))
}

func GetXPByLevel(level int) int {
	if level <= 2 {
		return int(math.Floor(math.Pow(float64(level), 2)+0.99) * 10)
	}
	return int(math.Floor(math.Pow(float64(level), 3.72)+0.99) * 5)
}
