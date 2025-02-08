package player

import (
	"cafego/internal/models/simple"
	"math"
)

func (p *Player) GetAchivements() simple.IntMap {
	return p.Achievement
}

// These functions are wrappers for updating achivements
func (p *Player) SetAchievement(achivementID, proggression int) {
	if _, ok := p.Achievement[achivementID]; ok {
		p.Achievement[achivementID] = proggression
	}
}

// <wod id="2001" n="Basic" g="Achievement" t="Earnedchips" />
func (p *Player) UpdateAchivementEarnedChips(chips int) {
	p.Achievement[2001] += chips
}

// <wod id="2002" n="Basic" g="Achievement" t="Spentchips" />
func (p *Player) UpdateAchivementSpendChips(chips int) {
	p.Achievement[2002] += int(math.Abs(float64(chips)))
}

// <wod id="2003" n="Basic" g="Achievement" t="Spentgold" />
func (p *Player) UpdateAchivementSpendGold(gold int) {
	p.Achievement[2003] += int(math.Abs(float64(gold)))
}

// <wod id="2004" n="Basic" g="Achievement" t="Boughtdeco" />
func (p *Player) AchivementBoughtDecoration() {
	p.Achievement[2004]++
}

// <wod id="2005" n="Basic" g="Achievement" t="Boughtingredients" />
func (p *Player) UpdateAchivementBoughtIngredients() {
	p.Achievement[2005]++
}

// <wod id="2006" n="Basic" g="Achievement" t="Solditems" />
func (p *Player) UpdateAchivementSoldItems() {
	p.Achievement[2006]++
}

// <wod id="2007" n="Basic" g="Achievement" t="Servingscount" />
func (p *Player) UpdateAchivementServingsCount() {
	p.Achievement[2007]++
}

// <wod id="2008" n="Basic" g="Achievement" t="Overcookedfoods" />
func (p *Player) UpdateAchivementOvercookedFoods() {
	p.Achievement[2008]++
}

// <wod id="2009" n="Basic" g="Achievement" t="Firenpc" />
func (p *Player) UpdateAchivementFireNPC() {
	p.Achievement[2009]++
}

// <wod id="2010" n="Social" g="Achievement" t="Friendscount" />
func (p *Player) SetAchivementFriendsCount() {
	p.Achievement[2010] = len(p.Friends)
}

// <wod id="2011" n="Basic" g="Achievement" t="Curiercount" />
func (p *Player) UpdateAchivementCurierCount() {
	p.Achievement[2011]++
}

// <wod id="2012" n="Basic" g="Achievement" t="Coopcount" />
func (p *Player) UpdateAchivementCoopCount() {
	p.Achievement[2012]++
}

// <wod id="2013" n="Basic" g="Achievement" t="Coopgoldcount" />
func (p *Player) UpdateAchivementCoopGoldCount() {
	p.Achievement[2013]++
}

// <wod id="2014" n="Basic" g="Achievement" t="Servingcountsweets" />
func (p *Player) UpdateAchivementServingCountSweets() {
	p.Achievement[2014]++
}

// <wod id="2015" n="Basic" g="Achievement" t="Servingcountmeals" />
func (p *Player) UpdateAchivementServingCountMeals() {
	p.Achievement[2015]++
}

// <wod id="2016" n="Basic" g="Achievement" t="Servingcountsoups" />
func (p *Player) UpdateAchivementServingCountSoups() {
	p.Achievement[2016]++
}

// <wod id="2017" n="Basic" g="Achievement" t="Servingcountsalads" />
func (p *Player) UpdateAchivementServingCountSalads() {
	p.Achievement[2017]++
}

// <wod id="2018" n="Basic" g="Achievement" t="Servingcountvegans" />
func (p *Player) UpdateAchivementServingCountVegans() {
	p.Achievement[2018]++
}

// <wod id="2019" n="Basic" g="Achievement" t="Servingcountsnacks" />
func (p *Player) UpdateAchivementServingCountSnacks() {
	p.Achievement[2019]++
}

// <wod id="2020" n="Basic" g="Achievement" t="Cookingcount" />
func (p *Player) UpdateAchivementCookingCount() {
	p.Achievement[2020]++
}

// <wod id="2021" n="Basic" g="Achievement" t="Differentdishes" />
func (p *Player) UpdateAchivementDifferentDishes() {
	p.Achievement[2021]++
}

// <wod id="2022" n="Basic" g="Achievement" t="Masteies" />
func (p *Player) UpdateAchivementMasteries() {
	p.Achievement[2022]++
}

// <wod id="2023" n="Basic" g="Achievement" t="Masteriesgold" />
func (p *Player) UpdateAchivementMasteriesGold() {
	p.Achievement[2023]++
}

// <wod id="2024" n="Basic" g="Achievement" t="Fancycount" />
func (p *Player) UpdateAchivementFancyCount() {
	p.Achievement[2024]++
}

// <wod id="2025" n="Basic" g="Achievement" t="Instantcount" />
func (p *Player) UpdateAchivementInstantCount() {
	p.Achievement[2025]++
}

// <wod id="2026" n="Basic" g="Achievement" t="Jobserve" />
func (p *Player) UpdateAchivementJobServe() {
	p.Achievement[2026]++
}

// <wod id="2027" n="Basic" g="Achievement" t="Jobclean" />
func (p *Player) UpdateAchivementJobClean() {
	p.Achievement[2027]++
}

// <wod id="2028" n="Basic" g="Achievement" t="Muffinmancash" />
func (p *Player) UpdateAchivementMuffinmanCash(amount int) {
	p.Achievement[2028] += amount
}

// <wod id="2029" n="Basic" g="Achievement" t="Muffinmangold" />
func (p *Player) UpdateAchivementMuffinmanGold(amount int) {
	p.Achievement[2029] += amount
}

// <wod id="2030" n="Basic" g="Achievement" t="Wheeloffortune" />
func (p *Player) UpdateAchivementWheelOfFortune() {
	p.Achievement[2030]++
}
