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
	p.job.Started = true
	p.job.StartedAt = time.Now().UTC()
	p.workTimeLeft = balancing.BalancingConstants.WorkTimeLeft

	p.SetJobLocation(cafe)

	workTimeEnd := time.Now().UTC().Add(time.Duration(balancing.BalancingConstants.WorkTimeLeft) * time.Second)

	p.job.FinishesAt = workTimeEnd

	go func() {
		<-time.After(time.Until(workTimeEnd))
		p.FinishJob()
	}()
}

func (p *Player) SetJobLocation(cafe *cafe.Cafe) {
	p.job.Location = cafe
}

func (p *Player) GetJobStarted() bool {
	return p.job.Started
}

func (p *Player) GetJobStartedAt() time.Time {
	return p.job.StartedAt
}

func (p *Player) GetWorkTimeLeft() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.workTimeLeft
}

func (p *Player) FinishJob() {

}

func (p *Player) DoAction() {
	p.job.ActionsDone++
}

func (p *Player) AddJobOffer(playerID int) {
	p.job.Offers = append(p.job.Offers, playerID)
}

func (p *Player) RemoveJobOffer(playerID int) {
	for i, id := range p.job.Offers {
		if id == playerID {
			p.job.Offers = append(p.job.Offers[:i], p.job.Offers[i+1:]...)
			return
		}
	}
}

func (p *Player) ClearOffers() {
	p.job.Offers = nil
}
