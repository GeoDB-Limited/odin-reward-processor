package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type RewardDoc struct {
	Rewards []RewardRow `yaml:"rewards"`
	Total   []AmountRow `yaml:"total"`
}

type RewardRow struct {
	Reward []AmountRow `yaml:"reward"`
}

type AmountRow struct {
	Amount string `yaml:"amount"`
}

func (d *RewardDoc) getDocument(filename string) *RewardDoc {
	yamlFile, err := ioutil.ReadFile(BuildFilePath(filename))
	if err != nil {
		log.Printf("yamlFile.Get err #%v", err)
	}
	err = yaml.Unmarshal(yamlFile, d)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return d
}
