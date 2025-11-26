package player

import (
	"cafego/internal/models/balancing"
	"cafego/internal/models/cafe"
	"cafego/internal/models/object"
	"time"
)

type PlayerJob struct {
	StartedAt     time.Time
	FinishesAt    time.Time
	Started       bool
	Location      *cafe.Cafe
	ActionsDone   int
	Offers        []int
	DishID        int
	DishStatus    int
	ReservedChair *object.Object
}

func (p *Player) StartJob(cafe *cafe.Cafe) {
	p.Job.Started = true
	p.Job.StartedAt = time.Now().UTC()
	p.WorkTimeLeft = balancing.BalancingConstants.WorkTimeLeft

	p.SetJobLocation(cafe)

	workTimeEnd := time.Now().UTC().Add(time.Duration(balancing.BalancingConstants.WorkTimeLeft) * time.Second)

	p.Job.FinishesAt = workTimeEnd

	go func() {
		<-time.After(time.Until(workTimeEnd))
		p.FinishJob()
	}()
}

func (p *Player) SetJobLocation(cafe *cafe.Cafe) {
	p.Job.Location = cafe
}

func (p *Player) GetJobStarted() bool {
	return p.Job.Started
}

func (p *Player) GetJobStartedAt() time.Time {
	return p.Job.StartedAt
}

func (p *Player) GetWorkTimeLeft() int {
	if p.GetJobStarted() {
		return max(0, (balancing.BalancingConstants.WorkTimeLeft - int(time.Since(p.GetJobStartedAt()).Seconds())))
	}
	return 0
}

func (p *Player) FinishJob() {

}

func (p *Player) DoAction() {
	p.Job.ActionsDone++
}

func (p *Player) AddJobOffer(playerID int) {
	p.Job.Offers = append(p.Job.Offers, playerID)
}

func (p *Player) RemoveJobOffer(playerID int) {
	for i, id := range p.Job.Offers {
		if id == playerID {
			p.Job.Offers = append(p.Job.Offers[:i], p.Job.Offers[i+1:]...)
			return
		}
	}
}

func (p *Player) ClearOffers() {
	p.Job.Offers = nil
}
