package objects

// These functions do simple operations while updating achivements
// They return if the achivement is completed

// <wod id="2001" n="Basic" g="Achievement" t="Earnedchips" />
func (p *Player) EarnedChips(chips int) bool {
	p.Cash += chips

	p.Achievement[2001] += chips

	return p.Achievement[2001] >= 25000
}

// <wod id="2002" n="Basic" g="Achievement" t="Spentchips" />
func (p *Player) SpendChips(chips int) bool {
	p.Cash -= chips

	p.Achievement[2002] += chips

	return p.Achievement[2002] >= 25000
}

// <wod id="2003" n="Basic" g="Achievement" t="Spentgold" />
func (p *Player) SpendGold(gold int) bool {
	p.Gold -= gold

	p.Achievement[2003] += gold

	return p.Achievement[2003] >= 30
}

// <wod id="2004" n="Basic" g="Achievement" t="Boughtdeco" />
func (p *Player) BoughtDecoration() bool {

	p.Achievement[2004]++

	return p.Achievement[2004] >= 30
}

// <wod id="2005" n="Basic" g="Achievement" t="Boughtingredients" />
// <wod id="2006" n="Basic" g="Achievement" t="Solditems" />
// <wod id="2007" n="Basic" g="Achievement" t="Servingscount" />
// <wod id="2008" n="Basic" g="Achievement" t="Overcookedfoods" />
// <wod id="2009" n="Basic" g="Achievement" t="Firenpc" />
// <wod id="2010" n="Social" g="Achievement" t="Friendscount" />
// <wod id="2011" n="Basic" g="Achievement" t="Curiercount" />
// <wod id="2012" n="Basic" g="Achievement" t="Coopcount" />
// <wod id="2013" n="Basic" g="Achievement" t="Coopgoldcount" />
// <wod id="2014" n="Basic" g="Achievement" t="Servingcountsweets" />
// <wod id="2015" n="Basic" g="Achievement" t="Servingcountmeals" />
// <wod id="2016" n="Basic" g="Achievement" t="Servingcountsoups" />
// <wod id="2017" n="Basic" g="Achievement" t="Servingcountsalads" />
// <wod id="2018" n="Basic" g="Achievement" t="Servingcountvegans" />
// <wod id="2019" n="Basic" g="Achievement" t="Servingcountsnacks" />
// <wod id="2020" n="Basic" g="Achievement" t="Cookingcount" />
// <wod id="2021" n="Basic" g="Achievement" t="Differentdishes" />
// <wod id="2022" n="Basic" g="Achievement" t="Masteies" />
// <wod id="2023" n="Basic" g="Achievement" t="Masteriesgold" />
// <wod id="2024" n="Basic" g="Achievement" t="Fancycount" />
// <wod id="2025" n="Basic" g="Achievement" t="Instantcount" />
// <wod id="2026" n="Basic" g="Achievement" t="Jobserve" />
// <wod id="2027" n="Basic" g="Achievement" t="Jobclean" />
// <wod id="2028" n="Basic" g="Achievement" t="Muffinmancash" />
// <wod id="2029" n="Basic" g="Achievement" t="Muffinmangold" />
// <wod id="2030" n="Basic" g="Achievement" t="Wheeloffortune" />
