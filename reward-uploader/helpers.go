package main

import "fmt"

const DIR = "./rewards-yaml-db"

func GetDir() string { return DIR }

func BuildFilePath(filename string) string { return fmt.Sprintf("%s/%s", GetDir(), filename) }

func GetCosmosAccount(filename string) string { return filename[7 : len(filename)-5] }
