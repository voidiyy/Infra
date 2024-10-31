package internal

import (
	"math/rand"
	"net/url"
	"time"
)

type BalanceStrategy interface {
	PickServer(tg *Target) (*url.URL, error)
	ShowStrat() string
}

func ChooseStrat(strat string) BalanceStrategy {
	switch strat {
	case "roundRobin":

		return &RoundRobinStrategy{}
	case "random":
		return &RandomStrategy{}
	default:
		return &RoundRobinStrategy{}
	}
}

/////////////////////////////////////

type RoundRobinStrategy struct {
	index int
}

func (r *RoundRobinStrategy) PickServer(tg *Target) (*url.URL, error) {

	server := tg.URLs[r.index]
	r.index = (r.index + 1) % len(tg.URLs)

	u, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *RoundRobinStrategy) ShowStrat() string {
	return "roundRobin"
}

////////////////////////////////////////

type RandomStrategy struct{}

func (r *RandomStrategy) PickServer(tg *Target) (*url.URL, error) {
	rand.Seed(time.Now().UnixNano())
	s := tg.URLs[rand.Intn(len(tg.URLs))]

	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *RandomStrategy) ShowStrat() string {
	return "random"
}

////////////////////////////////
