package balancing

import (
	"cafego/internal/utils"
	"strconv"
	"strings"
	"time"
)

type BalancingConfig struct {
	CleanCostCash          int
	StaffPrice             int
	SellFactorCash         int
	SellFactorGold         int
	TimeFactor             int
	RatingGuestHappy       int
	RatingGuestUnhappy     int
	CourierPrice           int
	MaxCourierSize         int
	InstantCookHourPerGold int
	ServerTimestamp        time.Time
	JobsPerDay             int
	JobRefillGold          int
	WorkTimeLeft           int
	CoopExpansionHours     int
	CoopExpansionGold      int
	CoopTimeToGold         float64
	CoopTimeToSilver       float64
	CoopRewardFactorGold   int
	CoopRewardFactorSilver int
	RefreshFoodCost        int
	MasteryDaysLv1         int
	MasteryDaysLv2         int
	MasteryDaysLv3         int
	MasteryStoveCount      int
	MasteryBonusServing    float64
	MasteryBonusXP         float64
	MasteryBonusTime       float64
	EmailVerificationGold  int
}

var BalancingConstants *BalancingConfig

func (balancing *BalancingConfig) AsResponse() string {
	args := []string{
		strconv.Itoa(balancing.CleanCostCash),
		strconv.Itoa(balancing.StaffPrice),
		strconv.Itoa(balancing.SellFactorCash),
		strconv.Itoa(balancing.SellFactorGold),
		strconv.Itoa(balancing.TimeFactor),
		strconv.Itoa(balancing.RatingGuestHappy),
		strconv.Itoa(balancing.RatingGuestUnhappy),
		strconv.Itoa(balancing.CourierPrice),
		strconv.Itoa(balancing.MaxCourierSize),
		strconv.Itoa(balancing.InstantCookHourPerGold),
		balancing.GetServerTimestamp(),
		strconv.Itoa(balancing.JobsPerDay),
		strconv.Itoa(balancing.JobRefillGold),
		strconv.Itoa(balancing.WorkTimeLeft),
		strconv.Itoa(balancing.CoopExpansionHours),
		strconv.Itoa(balancing.CoopExpansionGold),
		strconv.FormatFloat(balancing.CoopTimeToGold, 'f', -1, 64),
		strconv.FormatFloat(balancing.CoopTimeToSilver, 'f', -1, 64),
		strconv.Itoa(balancing.CoopRewardFactorGold),
		strconv.Itoa(balancing.CoopRewardFactorSilver),
		strconv.Itoa(balancing.RefreshFoodCost),
		strconv.Itoa(balancing.MasteryDaysLv1),
		strconv.Itoa(balancing.MasteryDaysLv2),
		strconv.Itoa(balancing.MasteryDaysLv3),
		strconv.Itoa(balancing.MasteryStoveCount),
		strconv.FormatFloat(balancing.MasteryBonusServing, 'f', -1, 64),
		strconv.FormatFloat(balancing.MasteryBonusXP, 'f', -1, 64),
		strconv.FormatFloat(balancing.MasteryBonusTime, 'f', -1, 64),
		strconv.Itoa(balancing.EmailVerificationGold),
	}

	return strings.Join(args, "%")
}

func LoadBalancing(hasConfig bool, envFile map[string]string) {
	cleanCostCash, _ := strconv.Atoi(utils.If(hasConfig, envFile["CLEAN_COST_CASH"], "15"))
	staffPrice, _ := strconv.Atoi(utils.If(hasConfig, envFile["STAFF_PRICE"], "2"))
	sellFactorCash, _ := strconv.Atoi(utils.If(hasConfig, envFile["SELL_FACTOR_CASH"], "20"))
	sellFactorGold, _ := strconv.Atoi(utils.If(hasConfig, envFile["SELL_FACTOR_GOLD"], "20"))
	timeFactor, _ := strconv.Atoi(utils.If(hasConfig, envFile["TIME_FACTOR"], "1"))
	ratingGuestHappy, _ := strconv.Atoi(utils.If(hasConfig, envFile["RATING_GUEST_HAPPY"], "1"))
	ratingGuestUnhappy, _ := strconv.Atoi(utils.If(hasConfig, envFile["RATING_GUEST_UNHAPPY"], "-2"))
	courierPrice, _ := strconv.Atoi(utils.If(hasConfig, envFile["COURIER_PRICE"], "1"))
	maxCourierSize, _ := strconv.Atoi(utils.If(hasConfig, envFile["MAX_COURIER_SIZE"], "8"))
	instantCookHourPerGold, _ := strconv.Atoi(utils.If(hasConfig, envFile["INSTANT_COOK_HOUR_PER_GOLD"], "1"))
	jobsPerDay, _ := strconv.Atoi(utils.If(hasConfig, envFile["JOBS_PER_DAY"], "5"))
	jobRefillGold, _ := strconv.Atoi(utils.If(hasConfig, envFile["JOB_REFILL_GOLD"], "1"))
	workTimeLeft, _ := strconv.Atoi(utils.If(hasConfig, envFile["WORK_TIME_LEFT"], "0"))
	coopExpansionHours, _ := strconv.Atoi(utils.If(hasConfig, envFile["COOP_EXPANSION_HOURS"], "20"))
	coopExpansionGold, _ := strconv.Atoi(utils.If(hasConfig, envFile["COOP_EXPANSION_GOLD"], "1"))
	coopTimeToGold, _ := strconv.ParseFloat(utils.If(hasConfig, envFile["COOP_TIME_TO_GOLD"], "0.5"), 64)
	coopTimeToSilver, _ := strconv.ParseFloat(utils.If(hasConfig, envFile["COOP_TIME_TO_SILVER"], "0.75"), 64)
	coopRewardFactorGold, _ := strconv.Atoi(utils.If(hasConfig, envFile["COOP_REWARD_FACTOR_GOLD"], "4"))
	coopRewardFactorSilver, _ := strconv.Atoi(utils.If(hasConfig, envFile["COOP_REWARD_FACTOR_SILVER"], "2"))
	refreshFoodCost, _ := strconv.Atoi(utils.If(hasConfig, envFile["REFRESH_FOOD_COST"], "1"))
	masteryDaysLv1, _ := strconv.Atoi(utils.If(hasConfig, envFile["MASTERY_DAYS_LV1"], "1"))
	masteryDaysLv2, _ := strconv.Atoi(utils.If(hasConfig, envFile["MASTERY_DAYS_LV2"], "5"))
	masteryDaysLv3, _ := strconv.Atoi(utils.If(hasConfig, envFile["MASTERY_DAYS_LV3"], "13"))
	masteryStoveCount, _ := strconv.Atoi(utils.If(hasConfig, envFile["MASTERY_STOVE_COUNT"], "4"))
	masteryBonusServing, _ := strconv.ParseFloat(utils.If(hasConfig, envFile["MASTERY_BONUS_SERVING"], "1.05"), 64)
	masteryBonusXP, _ := strconv.ParseFloat(utils.If(hasConfig, envFile["MASTERY_BONUS_XP"], "1.05"), 64)
	masteryBonusTime, _ := strconv.ParseFloat(utils.If(hasConfig, envFile["MASTERY_BONUS_TIME"], "0.95"), 64)
	emailVerificationGold, _ := strconv.Atoi(utils.If(hasConfig, envFile["EMAIL_VERIFICATION_GOLD"], "5"))

	BalancingConstants = &BalancingConfig{
		CleanCostCash:          cleanCostCash,
		StaffPrice:             staffPrice,
		SellFactorCash:         sellFactorCash,
		SellFactorGold:         sellFactorGold,
		TimeFactor:             timeFactor,
		RatingGuestHappy:       ratingGuestHappy,
		RatingGuestUnhappy:     ratingGuestUnhappy,
		CourierPrice:           courierPrice,
		MaxCourierSize:         maxCourierSize,
		InstantCookHourPerGold: instantCookHourPerGold,
		JobsPerDay:             jobsPerDay,
		JobRefillGold:          jobRefillGold,
		WorkTimeLeft:           workTimeLeft,
		CoopExpansionHours:     coopExpansionHours,
		CoopExpansionGold:      coopExpansionGold,
		CoopTimeToGold:         coopTimeToGold,
		CoopTimeToSilver:       coopTimeToSilver,
		CoopRewardFactorGold:   coopRewardFactorGold,
		CoopRewardFactorSilver: coopRewardFactorSilver,
		RefreshFoodCost:        refreshFoodCost,
		MasteryDaysLv1:         masteryDaysLv1,
		MasteryDaysLv2:         masteryDaysLv2,
		MasteryDaysLv3:         masteryDaysLv3,
		MasteryStoveCount:      masteryStoveCount,
		MasteryBonusServing:    masteryBonusServing,
		MasteryBonusXP:         masteryBonusXP,
		MasteryBonusTime:       masteryBonusTime,
		EmailVerificationGold:  emailVerificationGold,
	}
}

func (balancing *BalancingConfig) GetCleanCostCash() int      { return balancing.CleanCostCash }
func (balancing *BalancingConfig) GetStaffPrice() int         { return balancing.StaffPrice }
func (balancing *BalancingConfig) GetSellFactorCash() int     { return balancing.SellFactorCash }
func (balancing *BalancingConfig) GetSellFactorGold() int     { return balancing.SellFactorGold }
func (balancing *BalancingConfig) GetTimeFactor() int         { return balancing.TimeFactor }
func (balancing *BalancingConfig) GetRatingGuestHappy() int   { return balancing.RatingGuestHappy }
func (balancing *BalancingConfig) GetRatingGuestUnhappy() int { return balancing.RatingGuestUnhappy }
func (balancing *BalancingConfig) GetCourierPrice() int       { return balancing.CourierPrice }
func (balancing *BalancingConfig) GetMaxCourierSize() int     { return balancing.MaxCourierSize }
func (balancing *BalancingConfig) GetInstantCookHourPerGold() int {
	return balancing.InstantCookHourPerGold
}

func (balancing *BalancingConfig) GetServerTimestamp() string {
	return strconv.FormatInt(time.Now().UTC().Unix(), 10)
}

func (balancing *BalancingConfig) GetJobsPerDay() int           { return balancing.JobsPerDay }
func (balancing *BalancingConfig) GetJobRefillGold() int        { return balancing.JobRefillGold }
func (balancing *BalancingConfig) GetCoopExpansionHours() int   { return balancing.CoopExpansionHours }
func (balancing *BalancingConfig) GetCoopExpansionGold() int    { return balancing.CoopExpansionGold }
func (balancing *BalancingConfig) GetCoopTimeToGold() float64   { return balancing.CoopTimeToGold }
func (balancing *BalancingConfig) GetCoopTimeToSilver() float64 { return balancing.CoopTimeToSilver }
func (balancing *BalancingConfig) GetCoopRewardFactorGold() int {
	return balancing.CoopRewardFactorGold
}
func (balancing *BalancingConfig) GetCoopRewardFactorSilver() int {
	return balancing.CoopRewardFactorSilver
}
func (balancing *BalancingConfig) GetRefreshFoodCost() int   { return balancing.RefreshFoodCost }
func (balancing *BalancingConfig) GetMasteryDaysLv1() int    { return balancing.MasteryDaysLv1 }
func (balancing *BalancingConfig) GetMasteryDaysLv2() int    { return balancing.MasteryDaysLv2 }
func (balancing *BalancingConfig) GetMasteryDaysLv3() int    { return balancing.MasteryDaysLv3 }
func (balancing *BalancingConfig) GetMasteryStoveCount() int { return balancing.MasteryStoveCount }
func (balancing *BalancingConfig) GetMasteryBonusServing() float64 {
	return balancing.MasteryBonusServing
}
func (balancing *BalancingConfig) GetMasteryBonusXP() float64   { return balancing.MasteryBonusXP }
func (balancing *BalancingConfig) GetMasteryBonusTime() float64 { return balancing.MasteryBonusTime }
func (balancing *BalancingConfig) GetEmailVerificationGold() int {
	return balancing.EmailVerificationGold
}
