package access_test

// securityConfigMock is a simple mock for SecurityConfig
type securityConfigMock struct {
	accessKey string
}

func (sc *securityConfigMock) RefreshKey() string {
	return ""
}

func (sc *securityConfigMock) AccessKey() string {
	return sc.accessKey
}

func (sc *securityConfigMock) RefreshExp() byte {
	return 0
}

func (sc *securityConfigMock) AccessExp() byte {
	return 0
}