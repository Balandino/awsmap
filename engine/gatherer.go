package engine

func GatherData(account string) DiagramRoot {
	yamlTree := StartDiagram()

	if account == "1" {

		vpcId := "VPC1"
		yamlTree.GetResource("vpc", vpcId, "VPC 1")
		yamlTree.AddParent(AWSCloudID, vpcId)

		subId := "PublicSubnet"
		yamlTree.GetResource("public_subnet", subId, "Public Subnet")
		yamlTree.AddParent(vpcId, subId)

		webServerId := "WebServer"
		yamlTree.GetResource("ec2", webServerId, "Web Server")
		yamlTree.AddParent(subId, webServerId)

		options := make(map[string]string)
		options["sourceLeft"] = "HTTP:80"
		options["targetRight"] = "HTTP:80"

		webServer2Id := "WebServer2"
		yamlTree.AddLink(webServerId, "E", webServer2Id, "W", options)
	}

	if account == "2" {

		vpc2Id := "VPC2"
		yamlTree.GetResource("vpc", vpc2Id, "VPC 2")
		yamlTree.AddParent(AWSCloudID, vpc2Id)

		sub2Id := "PublicSubnet2"
		yamlTree.GetResource("public_subnet", sub2Id, "Public Subnet")
		yamlTree.AddParent(vpc2Id, sub2Id)

		webServer2Id := "WebServer2"
		yamlTree.GetResource("ec2", webServer2Id, "Web Server 2")
		yamlTree.AddParent(sub2Id, webServer2Id)
	}

	return yamlTree
}
