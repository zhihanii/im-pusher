package biz

type NodeInstance struct {
	Region   string
	Zone     string
	Env      string
	AppID    string
	Hostname string
	Addrs    []string
	Version  string
	LastTs   int64
	Metadata map[string]string
}

type Online struct {
	Server    string
	RoomCount map[string]int32
	Updated   int64
}
