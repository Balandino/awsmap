package main

import (
	"sync"

	"awsmap/engine"
)

func main() {
	// yamlTree := engine.GatherData()

	accounts := []string{
		"1", // Account A
		"2", // Account B
	}

	resultsChan := make(chan engine.YamlResults, len(accounts))
	var wg sync.WaitGroup

	// 1. Launch a goroutine for each account
	for _, account := range accounts {
		wg.Add(1)
		go func(acc string) {
			defer wg.Done()

			resourceTree := engine.GatherData(account)
			resultsChan <- engine.YamlResults{
				Account:  account,
				YamlTree: resourceTree,
			}
		}(account)
	}

	// 2. Wait for all goroutines to finish and close the channel
	wg.Wait()
	close(resultsChan)

	var trees []engine.DiagramRoot
	var coreTree engine.DiagramRoot

	// 4. Aggregate collected account data into the main diagram tree
	for res := range resultsChan {
		if res.Account != "1" {
			trees = append(trees, res.YamlTree)
			continue
		}

		coreTree = res.YamlTree
	}

	for _, yamlTree := range trees {
		coreTree.MergeTree(yamlTree)
	}

	// 5. Write out the final combined YAML
	coreTree.WriteYAML()
}
