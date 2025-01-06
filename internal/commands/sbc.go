package commands

import (
  "cafego/internal/types/requests"
  "cafego/internal/client"
  "cafego/internal/managers"
)

// sbc - SendBalancingConstant
func SendBalancingConstant(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {

  c.SendExtensionResponse("sbc", "1", "0",
        "15",          // cleanCostCash
        "2",           // staffPrice
        "20",          // sellFactorCash
        "20",          // sellFactorGold
        "1",           // timeFactor
        "1",           // rating_guest_happy
        "-2",          // rating_guest_unhappy
        "0",           // courierPrice (NOT READY)
        "0",           // maxCourierSize (NOT READY)
        "0",           // instantCookHourPerGold (NOT READY)
        "1704067200",  // serverTimeStamp
        "5",           // jobsPerDay
        "1",           // jobRefillGold
        "0",           // workTimeLeft (NOT READY)
        "0",           // coopExpansionHoures (NOT READY)
        "0",           // coopExpansionGold (NOT READY)
        "0",           // coopTimeToGold (NOT READY)
        "0",           // coopTimeToSilver (NOT READY)
        "0",           // coopRewardFactorGold (NOT READY)
        "0",           // coopRewardFactorSilver (NOT READY)
        "1",           // refreshFoodCost (NOT READY)
        "1",           // masteryDaysLV1
        "5",           // masteryDaysLV2
        "13",          // masteryDaysLV3
        "4",           // masteryStoveCount
        "1.05",        // masteryBonusServing
        "1.05",        // masteryBonusXP
        "0.95",        // masteryBonusTime
        "5",            // emailVerificationGold
  )

  return nil
}


