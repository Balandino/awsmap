// Package engine provides the means to build resources representing the YAML file
package engine

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

const AWSCloudID = "AWSCloud"

var topLevelResources = []string{AWSCloudID, "Canvas"}

// var yamlTree DiagramRoot // Singleton

func StartDiagram() DiagramRoot {
	yamlTree := DiagramRoot{
		Diagram: DiagramResources{
			DefinitionFiles: []DefinitionFile{
				{
					Type:      "LocalFile",
					LocalFile: "../imgs/dark_imgs.yaml",
				},
			},
			Links:     []Link{},
			Resources: make(map[string]*Resource),
		},
	}

	yamlTree.Diagram.Resources["Canvas"] = &Resource{
		Type:      "AWS::Diagram::Canvas",
		Children:  []string{"AWSCloud"},
		FillColor: "rgba(0, 0, 0, 255)",
	}

	yamlTree.Diagram.Resources[AWSCloudID] = &Resource{
		Type:   "AWS::Diagram::Cloud",
		Preset: "AWSCloudNoLogo",
	}

	return yamlTree
}

func (yamlTree *DiagramRoot) GetResource(resourceType, id, title string) *Resource {
	var ptrResource *Resource

	switch resourceType {
	case "vpc":
		ptrResource = &Resource{
			Type:     "AWS::EC2::VPC",
			Title:    title,
			Children: []string{},
		}

	case "public_subnet":
		ptrResource = &Resource{
			Type:     "AWS::EC2::Subnet",
			Title:    title,
			Children: []string{},
		}
	case "ec2":
		ptrResource = &Resource{
			Type:     "AWS::EC2::Instance",
			Title:    title,
			Children: []string{},
		}
	}

	yamlTree.Diagram.Resources[id] = ptrResource
	return ptrResource
}

func (yamlTree *DiagramRoot) AddParent(parentID, thisId string) {
	yamlTree.Diagram.Resources[parentID].Children = append(yamlTree.Diagram.Resources[parentID].Children, thisId)
}

func (yamlTree *DiagramRoot) AddLink(source, sourcePosition, target, targetPosition string, options map[string]string) {
	var link Link

	link.Source = source
	link.Target = target
	link.SourcePosition = sourcePosition
	link.TargetPosition = targetPosition
	link.LineColor = "rgba(255, 255, 255, 255)"

	_, hasSourceLeft := options["sourceLeft"]
	_, hasSourceRight := options["sourceRight"]

	if !hasSourceLeft && !hasSourceRight {
		yamlTree.Diagram.Links = append(yamlTree.Diagram.Links, link)
		return
	}

	var labels Labels

	if val, ok := options["sourceLeft"]; ok {
		var labelDetail LabelDetail
		labelDetail.Title = val
		labels.SourceLeft = &labelDetail
	}

	if val, ok := options["targetRight"]; ok {
		var labelDetail LabelDetail
		labelDetail.Title = val
		labels.TargetRight = &labelDetail
	}

	link.Labels = &labels
	yamlTree.Diagram.Links = append(yamlTree.Diagram.Links, link)
}

func (yamlTree *DiagramRoot) WriteYAML() {
	// Marshal struct into YAML bytes
	yamlData, err := yaml.Marshal(&yamlTree)
	if err != nil {
		log.Fatalf("Error marshalling to YAML: %v", err)
	}

	// Write the marshalled YAML to disk
	outputFilePath := "./outputs/diagram.yaml"
	err = os.WriteFile(outputFilePath, yamlData, 0o644)
	if err != nil {
		log.Fatalf("Error writing YAML to file: %v", err)
	}

	cmd := exec.Command("./yamlfix/bin/yamlfix", outputFilePath)

	// Direct output to the terminal (stdout and stderr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run executes and waits for completion
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run yamlfix: %v", err)
	}

	fmt.Printf("Successfully generated %s!\n\n", outputFilePath)
	// fmt.Println("File Contents:")
	// fmt.Println(string(yamlData))
}

func (ptrCoreTree *DiagramRoot) MergeTree(tree DiagramRoot) {
	mergeResources(ptrCoreTree, tree)
	mergeLinks(ptrCoreTree, tree)
	mergeParents(ptrCoreTree, tree)
}

func mergeParents(ptrCoreTree *DiagramRoot, tree DiagramRoot) {
	for _, resource := range topLevelResources {
		prtCoreAWSCloudChildren := &ptrCoreTree.Diagram.Resources[resource].Children
		treeAWSCloudChildren := tree.Diagram.Resources[resource].Children
		MergeUniqueChildren(prtCoreAWSCloudChildren, treeAWSCloudChildren)
	}
}

func mergeLinks(ptrCoreTree *DiagramRoot, tree DiagramRoot) {
	// Build list of links already in coreTree to avoid duplication
	var linkBuilder strings.Builder
	for _, link := range ptrCoreTree.Diagram.Links {
		linkBuilder.WriteString(fmt.Sprintf("%s%s|", link.Source, link.Target))
	}

	coreLinks := linkBuilder.String()

	// Add Links
	for _, link := range tree.Diagram.Links {
		linkId := fmt.Sprintf("%s%s|", link.Source, link.Target)

		if !strings.Contains(coreLinks, linkId) {
			ptrCoreTree.Diagram.Links = append(ptrCoreTree.Diagram.Links, link)
		}
	}
}

func mergeResources(ptrCoreTree *DiagramRoot, tree DiagramRoot) {
	// Iterate over each tree, adding resources into core tree if they are not already there
	for id, resource := range tree.Diagram.Resources {
		if _, ok := ptrCoreTree.Diagram.Resources[id]; !ok {
			ptrCoreTree.Diagram.Resources[id] = resource
			continue
		}

		ptrCoreChildren := &ptrCoreTree.Diagram.Resources[id].Children
		treeChildren := tree.Diagram.Resources[id].Children
		MergeUniqueChildren(ptrCoreChildren, treeChildren)
	}
}

func MergeUniqueChildren(coreChildren *[]string, treeChildren []string) {
	// 1. Both inputs have no elements (either nil or length 0)
	if len(*coreChildren) == 0 && len(treeChildren) == 0 {
		return
	}

	// 2. core is empty, tree has items -> copy tree into core
	if len(*coreChildren) == 0 {
		*coreChildren = append([]string(nil), treeChildren...)
		return
	}

	// 3. tree is empty, core has items -> nothing to merge
	if len(treeChildren) == 0 {
		return
	}

	// 4. Both have items -> merge unique entries using map[string]bool
	seen := make(map[string]bool, len(*coreChildren))
	for _, child := range *coreChildren {
		seen[child] = true
	}

	for _, child := range treeChildren {
		if !seen[child] {
			*coreChildren = append(*coreChildren, child)
			seen[child] = true // Prevent duplicates within treeChildren itself
		}
	}
}

// // Chunk 1 Builder
// func addVPCChunk(resources map[string]*Resource) {
// 	resources["VPC"] = &Resource{
// 		Type:      "AWS::EC2::VPC",
// 		Direction: "vertical",
// 		Children:  []string{"VPCPublicStack", "ALB"},
// 		BorderChildren: []BorderChild{
// 			{Position: "S", Resource: "IGW"},
// 		},
// 	}
// }
//
// // Chunk 2 Builder
// func addPublicStackChunk(resources map[string]*Resource) {
// 	resources["VPCPublicSubnet1Instance"] = &Resource{Type: "AWS::EC2::Instance"}
// 	resources["VPCPublicSubnet2Instance"] = &Resource{Type: "AWS::EC2::Instance"}
//
// 	resources["VPCPublicStack"] = &Resource{
// 		Type: "AWS::AutoScaling::AutoScalingGroup",
// 		Children: []string{
// 			"VPCPublicSubnet1Instance",
// 			"VPCPublicSubnet2Instance",
// 		},
// 	}
// }
//
// // Chunk 3 Builder
// func addEdgeChunk(resources map[string]*Resource) {
// 	resources["ALB"] = &Resource{
// 		Type:   "AWS::ElasticLoadBalancingV2::LoadBalancer",
// 		Preset: "Application Load Balancer",
// 	}
//
// 	resources["IGW"] = &Resource{
// 		Type: "AWS::EC2::InternetGateway",
// 		IconFill: &IconFill{
// 			Type:  "rect",
// 			Color: "rgba(0, 0, 0, 255)",
// 		},
// 	}
// }
//
// // Chunk 4 Builder
// func addUserChunk(resources map[string]*Resource) {
// 	resources["User"] = &Resource{
// 		Type:   "AWS::Diagram::Resource",
// 		Preset: "User",
// 	}
// }
//
