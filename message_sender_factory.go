package spirit

import (
	"fmt"
	"reflect"

	"github.com/gogap/errors"
)

type MessageSenderFactory interface {
	RegisterMessageSenders(senders ...MessageSender)
	IsExist(senderType string) bool
	NewSender(senderType string) (sender MessageSender, err error)
}

type DefaultMessageSenderFactory struct {
	senderDrivers map[string]reflect.Type

	instanceCache map[string]MessageSender
}

func NewDefaultMessageSenderFactory() MessageSenderFactory {
	fact := new(DefaultMessageSenderFactory)
	fact.senderDrivers = make(map[string]reflect.Type)
	fact.instanceCache = make(map[string]MessageSender)
	return fact
}

func (p *DefaultMessageSenderFactory) RegisterMessageSenders(senders ...MessageSender) {
	if senders == nil {
		panic("senders is nil")
	}

	if len(senders) == 0 {
		return
	}

	for _, sender := range senders {
		if _, exist := p.senderDrivers[sender.Type()]; exist {
			panic(fmt.Errorf("sender driver of [%s] already exist", sender.Type()))
		}

		vof := reflect.ValueOf(sender)
		vType := vof.Type()
		if vof.Kind() == reflect.Ptr {
			vType = vof.Elem().Type()
		}

		p.senderDrivers[sender.Type()] = vType
	}

	return
}

func (p *DefaultMessageSenderFactory) IsExist(senderType string) bool {
	if _, exist := p.senderDrivers[senderType]; exist {
		return true
	}
	return false
}

func (p *DefaultMessageSenderFactory) NewSender(senderType string) (sender MessageSender, err error) {
	if senderInstance, exist := p.instanceCache[senderType]; exist && senderInstance != nil {
		return senderInstance, nil
	}

	if senderDriver, exist := p.senderDrivers[senderType]; !exist {
		err = ERR_SENDER_DRIVER_NOT_EXIST.New(errors.Params{"type": senderType})
		return
	} else {
		if vOfMessageSender := reflect.New(senderDriver); vOfMessageSender.CanInterface() {
			iMessageSender := vOfMessageSender.Interface()
			if r, ok := iMessageSender.(MessageSender); ok {
				if err = r.Init(); err != nil {
					return
				}
				sender = r
				p.instanceCache[senderType] = sender
				return
			} else {
				err = ERR_SENDER_CREATE_FAILED.New(errors.Params{"type": senderType})
				return
			}
		}
		err = ERR_SENDER_BAD_DRIVER.New(errors.Params{"type": senderType})
		return
	}
}
