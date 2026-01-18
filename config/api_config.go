package config

func apiConfig() (*APIConfig, error) {
	grpcPort, err := getEnvInt("API_GRPC_PORT")
	if err != nil {
		return nil, err
	}
	apiGatewayPort, err := getEnvInt("API_GATEWAY_PORT")
	if err != nil {
		return nil, err
	}
	apiConfig := &APIConfig{
		GrpcPort:       uint16(grpcPort),
		ApiGatewayPort: uint16(apiGatewayPort),
	}
	return apiConfig, nil
}