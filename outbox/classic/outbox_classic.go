package classic

import (
	"github.com/gogap/spirit"
	"sync"
	"time"
)

const (
	outboxURN = "urn:spirit:outbox:classic"
)

var _ spirit.Outbox = new(ClassicOutbox)

type ClassicOutboxConfig struct {
	Size       int           `json:"size"`
	GetTimeout int           `json:"get_timeout"`
	Labels     spirit.Labels `json:"labels"`
}

type ClassicOutbox struct {
	name string

	statusLocker sync.Mutex
	status       spirit.Status

	deliveriesChan chan []spirit.Delivery

	conf ClassicOutboxConfig
}

func init() {
	spirit.RegisterOutbox(outboxURN, NewClassicOutbox)
}

func NewClassicOutbox(name string, options spirit.Map) (box spirit.Outbox, err error) {
	conf := ClassicOutboxConfig{}

	if err = options.ToObject(&conf); err != nil {
		return
	}

	box = &ClassicOutbox{
		name: name,
		conf: conf,
	}
	return

}

func (p *ClassicOutbox) Name() string {
	return p.name
}

func (p *ClassicOutbox) URN() string {
	return outboxURN
}

func (p *ClassicOutbox) Start() (err error) {
	p.statusLocker.Lock()
	defer p.statusLocker.Unlock()

	spirit.Logger().WithField("actor", spirit.ActorOutbox).
		WithField("urn", outboxURN).
		WithField("name", p.name).
		WithField("event", "start").
		Debugln("enter start")

	if p.status == spirit.StatusRunning {
		return
	}

	p.deliveriesChan = make(chan []spirit.Delivery, p.conf.Size)

	p.status = spirit.StatusRunning

	spirit.Logger().WithField("actor", spirit.ActorOutbox).
		WithField("urn", outboxURN).
		WithField("name", p.name).
		WithField("event", "start").
		Infoln("started")

	return
}

func (p *ClassicOutbox) Stop() (err error) {
	p.statusLocker.Lock()
	defer p.statusLocker.Unlock()

	spirit.Logger().WithField("actor", spirit.ActorOutbox).
		WithField("urn", outboxURN).
		WithField("name", p.name).
		WithField("event", "stop").
		Debugln("enter stop")

	if p.status == spirit.StatusStopped {
		return
	}

	p.status = spirit.StatusStopped

	close(p.deliveriesChan)

	spirit.Logger().WithField("actor", spirit.ActorOutbox).
		WithField("urn", outboxURN).
		WithField("name", p.name).
		WithField("event", "stop").
		Infoln("stopped")

	return
}

func (p *ClassicOutbox) Status() (status spirit.Status) {
	return p.status
}

func (p *ClassicOutbox) Labels() (labels spirit.Labels) {
	labels = p.conf.Labels
	return
}

func (p *ClassicOutbox) Put(deliveries []spirit.Delivery) (err error) {
	p.deliveriesChan <- deliveries
	return
}

func (p *ClassicOutbox) Get() (deliveries []spirit.Delivery, err error) {
	if p.conf.GetTimeout < 0 {
		deliveries = <-p.deliveriesChan
	} else {
		more := false
		select {
		case deliveries, more = <-p.deliveriesChan:
			{
				if more {
					spirit.Logger().WithField("actor", spirit.ActorOutbox).
						WithField("urn", outboxURN).
						WithField("name", p.name).
						WithField("event", "get deliveries").
						WithField("length", len(deliveries)).
						Debugln("deliveries received from deliveries chan")
				}
			}
		case <-time.After(time.Duration(p.conf.GetTimeout) * time.Millisecond):
			{
				spirit.Logger().WithField("actor", spirit.ActorOutbox).
					WithField("urn", outboxURN).
					WithField("name", p.name).
					WithField("event", "get deliveries").
					Debugln("get deliveries timeout")

				return
			}
		}
	}

	return
}
