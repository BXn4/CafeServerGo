package player

import (
	"cafego/internal/models/simple"
	"cafego/internal/utils"
)

func (p *Player) GetAchivements() simple.IntMap {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.Achievements
}

// PLEASE DO NOT MODIFY IT!!!
// How achievements works?
// The game gives rewards by the prevoius level, and we need to use current levels:
// <achievements a="0" w="2001" l="0" r="25000" x="100" ch="500" g="0" />
// idk, just dont touch it.
func (p *Player) MakeAchievementCurrentLevels() {
	if p.achievementsLevel == nil {
		// Set the achievement levels by the progess
		p.achievementsLevel = make(map[int]int)
		for id := 2001; id <= 2030; id++ {
			progress := p.GetAchievementProgress(id)
			level := p.GetAchievementLevelByProgress(id, progress)
			p.achievementsLevel[id] = level
			// println("SET ACHIEVEMENT LEVELS FOR: ID:", id, "Progress:", progress, "Current level:", level)
		}
	}
}

func (p *Player) GetAchievementProgress(achievementID int) int {
	return p.achievementsLevel[achievementID]
}

func (p *Player) SetAchievementLevel(achievementID, level int) {
	p.achievementsLevel[achievementID] = level
}

func (p *Player) GetAchievementLevel(achievementID int) int {
	return p.achievementsLevel[achievementID]
}

// These functions are wrappers for updating achivements
func (p *Player) SetAchievement(achivementID, proggression int) {
	if _, ok := p.Achievements[achivementID]; ok {
		p.achievementsLevel[achivementID] = proggression
	}
}

// PLEASE DO NOT MODIFY IT!!!
// How achievements works?
// The game gives rewards by the prevoius level, and we need to use current levels:
// <achievements a="0" w="2001" l="0" r="25000" x="100" ch="500" g="0" />
// idk, just dont touch it.
func (p *Player) CheckProgress(achievementID int) {
	progress := p.GetAchievementProgress(achievementID)
	currentLevel := p.GetAchievementLevel(achievementID)
	levelProgressTarget := utils.GetAchievementTarget(achievementID, currentLevel)
	maxLevelForAchievement := utils.GetAchievementMAXLevel(achievementID)

	/* println(progress)
	println(currentLevel)
	println(levelProgressTarget) */

	// time.Sleep(5 * time.Second)
	// nested loops happened here. fuck im out. PLEASE KEEP THIS

	if progress >= levelProgressTarget && currentLevel <= maxLevelForAchievement {
		// println("LEVEL UP!")
		p.SetAchievementLevel(achievementID, currentLevel+1)
		p.OnAchievementEarned(achievementID, currentLevel+1, p)
	} else {
		// println("NO LEVEL UP!")
		return
	}

}

// PLEASE DO NOT MODIFY IT!!!
// How achievements works?
// The game gives rewards by the prevoius level, and we need to use current levels:
// <achievements a="0" w="2001" l="0" r="25000" x="100" ch="500" g="0" />
// idk, just dont touch it.
func (p *Player) GetAchievementLevelByProgress(achievementID, progress int) int {
	achievements := utils.GetAchievementList(achievementID)
	level := 0

	for i := range achievements {
		achievement := achievements[i]
		if progress >= achievement.Target {
			level = achievement.Level + 1 // need to add + 1, because the list starts from 0 level
		} else {
			break
		}
	}

	return level
}

func (p *Player) UpdateAchievement(achievementID int, value int) {
	p.Achievements[achievementID] += value
	p.CheckProgress(achievementID)
}

// <wod id="2001" n="Basic" g="Achievement" t="Earnedchips" />
func (p *Player) UpdateAchivementEarnedChips(chips int) {
	p.UpdateAchievement(2001, chips)
}

// <wod id="2002" n="Basic" g="Achievement" t="Spentchips" />
func (p *Player) UpdateAchivementSpendChips(chips int) {
	p.UpdateAchievement(2002, chips)
}

// <wod id="2003" n="Basic" g="Achievement" t="Spentgold" />
func (p *Player) UpdateAchivementSpendGold(gold int) {
	p.UpdateAchievement(2003, gold)
}

// <wod id="2004" n="Basic" g="Achievement" t="Boughtdeco" />
func (p *Player) UpdateAchivementBoughtDecoration() {
	p.UpdateAchievement(2004, 1)
}

// <wod id="2005" n="Basic" g="Achievement" t="Boughtingredients" />
func (p *Player) UpdateAchivementBoughtIngredients() {
	p.UpdateAchievement(2005, 1)
}

// <wod id="2006" n="Basic" g="Achievement" t="Solditems" />
func (p *Player) UpdateAchivementSoldItems(value int) {
	p.UpdateAchievement(2006, value)
}

// <wod id="2007" n="Basic" g="Achievement" t="Servingscount" />
func (p *Player) UpdateAchivementServingsCount(amount int) {
	p.UpdateAchievement(2007, amount)
}

// <wod id="2008" n="Basic" g="Achievement" t="Overcookedfoods" />
func (p *Player) UpdateAchivementOvercookedFoods() {
	p.UpdateAchievement(2008, 1)
}

// <wod id="2009" n="Basic" g="Achievement" t="Firenpc" />
func (p *Player) UpdateAchivementFireNPC() {
	p.UpdateAchievement(2009, 1)
}

// <wod id="2010" n="Social" g="Achievement" t="Friendscount" />
func (p *Player) SetAchivementFriendsCount(v int) {
	p.UpdateAchievement(2010, v)
}

// <wod id="2011" n="Basic" g="Achievement" t="Curiercount" />
func (p *Player) UpdateAchivementCurierCount(amount int) {
	p.UpdateAchievement(2011, amount)
}

// <wod id="2012" n="Basic" g="Achievement" t="Coopcount" />
func (p *Player) UpdateAchivementCoopCount() {
	p.UpdateAchievement(2012, 1)
}

// <wod id="2013" n="Basic" g="Achievement" t="Coopgoldcount" />
func (p *Player) UpdateAchivementCoopGoldCount() {
	p.UpdateAchievement(2013, 1)
}

// <wod id="2014" n="Basic" g="Achievement" t="Servingcountsweets" />
func (p *Player) UpdateAchivementServingCountSweets() {
	p.UpdateAchievement(2014, 1)
}

// <wod id="2015" n="Basic" g="Achievement" t="Servingcountmeals" />
func (p *Player) UpdateAchivementServingCountMeals() {
	p.UpdateAchievement(2015, 1)
}

// <wod id="2016" n="Basic" g="Achievement" t="Servingcountsoups" />
func (p *Player) UpdateAchivementServingCountSoups() {
	p.UpdateAchievement(2016, 1)
}

// <wod id="2017" n="Basic" g="Achievement" t="Servingcountsalads" />
func (p *Player) UpdateAchivementServingCountSalads() {
	p.UpdateAchievement(2017, 1)
}

// <wod id="2018" n="Basic" g="Achievement" t="Servingcountvegans" />
func (p *Player) UpdateAchivementServingCountVegans() {
	p.UpdateAchievement(2018, 1)
}

// <wod id="2019" n="Basic" g="Achievement" t="Servingcountsnacks" />
func (p *Player) UpdateAchivementServingCountSnacks() {
	p.UpdateAchievement(2019, 1)
}

// <wod id="2020" n="Basic" g="Achievement" t="Cookingcount" />
func (p *Player) UpdateAchivementCookingCount() {
	p.UpdateAchievement(2020, 1)
}

// <wod id="2021" n="Basic" g="Achievement" t="Differentdishes" />
func (p *Player) UpdateAchivementDifferentDishes() {
	p.UpdateAchievement(2021, 1)
}

// <wod id="2022" n="Basic" g="Achievement" t="Masteies" />
func (p *Player) UpdateAchivementMasteries() {
	p.UpdateAchievement(2022, 1)
}

// <wod id="2023" n="Basic" g="Achievement" t="Masteriesgold" />
func (p *Player) UpdateAchivementMasteriesGold() {
	p.UpdateAchievement(2023, 1)
}

// <wod id="2024" n="Basic" g="Achievement" t="Fancycount" />
func (p *Player) UpdateAchivementFancyCount() {
	p.UpdateAchievement(2024, 1)
}

// <wod id="2025" n="Basic" g="Achievement" t="Instantcount" />
func (p *Player) UpdateAchivementInstantCount() {
	p.UpdateAchievement(2025, 1)
}

// <wod id="2026" n="Basic" g="Achievement" t="Jobserve" />
func (p *Player) UpdateAchivementJobServe() {
	p.UpdateAchievement(2026, 1)
}

// <wod id="2027" n="Basic" g="Achievement" t="Jobclean" />
func (p *Player) UpdateAchivementJobClean() {
	p.UpdateAchievement(2027, 1)
}

// <wod id="2028" n="Basic" g="Achievement" t="Muffinmancash" />
func (p *Player) UpdateAchivementMuffinmanCash(amount int) {
	p.UpdateAchievement(2028, amount)
}

// <wod id="2029" n="Basic" g="Achievement" t="Muffinmangold" />
func (p *Player) UpdateAchivementMuffinmanGold(amount int) {
	p.UpdateAchievement(2029, amount)
}

// <wod id="2030" n="Basic" g="Achievement" t="Wheeloffortune" />
func (p *Player) UpdateAchivementWheelOfFortune() {
	p.UpdateAchievement(2030, 1)
}
