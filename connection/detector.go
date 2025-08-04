package connection

type ContainerInfo struct {
	Name     string
	Type     string
	Port     string
	Database string
	Status   string
}

func DetectLocalContainers()([]ContainerInfo, error) {
	// This function should implement logic to detect local containers
	// For now, we return an empty slice and nil error
	return []ContainerInfo{}, nil
}
