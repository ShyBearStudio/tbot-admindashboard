package main

const (
	configEnvVarName string = "admindashboardconfig"
)

type configuration struct {
	Address         string
	StaticResources string
	Db              struct {
		Driver           string
		ConnectionString string
	}
	Log struct {
		Dir string
	}
	Tbots map[string]string
}
