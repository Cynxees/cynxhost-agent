package param

type ParamAwsNodeId struct {
	AmiId           string `json:"ami_id"`
	SecurityGroupId string `json:"security_group_id"`
	KeyPairId       string `json:"keypair_id"`
	KeyPairName     string `json:"keypair_name"`
	VpcId           string `json:"vpc_id"`
}

type ParamAwsNodeScript struct {
	InitialLaunchScript string `json:"initial_launch_script"`
	LaunchScript        string `json:"launch_script"`
}
