package command

type Server struct {
	LogLevel      string
	KubeConfig    string
	Namespace     string
	LabelSelector string
	LogFormat     string
}
