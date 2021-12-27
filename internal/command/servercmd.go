package command

type Server struct {
	LogLevel      string
	KubeConfig    string
	Namespace     string
	Container     string
	LabelSelector string
	LogFormat     string
	IncludeUname  bool
}
