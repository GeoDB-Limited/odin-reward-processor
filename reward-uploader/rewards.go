package main

import (
	"io/ioutil"
	"strconv"
)

func RewardsParser() (map[string]float64, error) {
	rewardsMap := make(map[string]float64)

	files, err := ioutil.ReadDir(GetDir())
	if err != nil {
		return rewardsMap, err
	}

	for _, f := range files {
		account := GetCosmosAccount(f.Name())

		var d RewardDoc
		d.getDocument(f.Name())
		if len(d.Total) > 0 {
			rewardAmount, err := strconv.ParseFloat(d.Total[0].Amount, 64)
			if err != nil {
				return rewardsMap, err
			}
			rewardsMap[account] = rewardAmount
		}
	}

	return rewardsMap, nil
}
