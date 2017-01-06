package echobot

const (
	configEnvVarName string = "echobotconfig"
)

type configuration struct {
	Token string
	Db    struct {
		Driver           string
		ConnectionString string
	}
	Log struct {
		Dir string
	}
}

var config configuration = configuration{}
