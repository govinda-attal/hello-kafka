package example

type Greetings interface {
	Hello(rq HelloRq) (HelloRs, error)
}

type HelloRq struct {
	Name string `json:"name"`
}

type HelloRs struct {
	Greetings string `json:"greetings"`
}
