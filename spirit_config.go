package spirit

type ActorConfig struct {
	Name    string  `json:"name"`
	URN     string  `json:"urn"`
	Options Options `json:"options"`
}

type WriterPoolActorConfig struct {
	ActorConfig
	Writer ActorConfig `json:"writer"`
}

type ReaderPoolActorConfig struct {
	ActorConfig
	Reader ActorConfig `json:"reader"`
}

type ComposeReceiverConfig struct {
	Name       string `json:"name"`
	Translator string `json:"translator"`
	ReaderPool string `json:"reader_pool"`
}

type ComposeSenderConfig struct {
	Name       string `json:"name"`
	Translator string `json:"translator"`
	WriterPool string `json:"writer_pool"`
}

type ComposeInboxConfig struct {
	Name      string                  `json:"name"`
	Receivers []ComposeReceiverConfig `json:"receivers"`
}

type ComposeOutboxConfig struct {
	Name    string                `json:"name"`
	Senders []ComposeSenderConfig `json:"senders"`
}

type ComposeLabelMatchConfig struct {
	Component string `json:"component"`
	Outbox    string `json:"outbox"`
}

type ComposeRouterConfig struct {
	Name          string                  `json:"name"`
	LabelMatchers ComposeLabelMatchConfig `json:"label_matchers"`
	Components    []string                `json:"components"`
	Inboxes       []ComposeInboxConfig    `json:"inboxes"`
	Outboxes      []ComposeOutboxConfig   `json:"outboxes"`
}

type SpiritConfig struct {
	ReaderPools       []ReaderPoolActorConfig `json:"reader_pools"`
	WriterPools       []WriterPoolActorConfig `json:"writer_pools"`
	InputTranslators  []ActorConfig           `json:"input_translators"`
	OutputTranslators []ActorConfig           `json:"output_translators"`
	Inboxes           []ActorConfig           `json:"inboxes"`
	Outboxes          []ActorConfig           `json:"outboxes"`
	Receivers         []ActorConfig           `json:"receivers"`
	Senders           []ActorConfig           `json:"senders"`
	Routers           []ActorConfig           `json:"routers"`
	Components        []ActorConfig           `json:"components"`
	LabelMatchers     []ActorConfig           `json:"label_matchers"`

	Compose []ComposeRouterConfig `json:"compose"`
}

func (p *SpiritConfig) Validate() (err error) {
	actorNames := map[string]bool{}

	if p.ReaderPools != nil {
		poolDupCheck := map[string]bool{}
		for _, pool := range p.ReaderPools {
			if _, exist := newReaderPoolFuncs[pool.URN]; !exist {
				err = ErrReaderPoolURNNotExist
				return
			}

			if _, exist := poolDupCheck[pool.Name]; exist {
				err = ErrReaderPoolNameDuplicate
				return
			} else {
				poolDupCheck[pool.Name] = true
			}
			actorNames[actorTypedName(ActorReaderPool, pool.Name)] = true

			readerDupCheck := map[string]bool{}
			if _, exist := newReaderFuncs[pool.Reader.URN]; !exist {
				err = ErrReaderURNNotExist
				return
			}

			if _, exist := readerDupCheck[pool.Reader.Name]; exist {
				err = ErrReaderNameDuplicate
				return
			} else {
				readerDupCheck[pool.Reader.Name] = true
			}

			actorNames[actorTypedName(ActorReader, pool.Reader.Name)] = true
			actorNames[actorTypedName(ActorReaderPool, pool.Name)] = true
		}
	}

	if p.WriterPools != nil {
		poolDupCheck := map[string]bool{}
		for _, pool := range p.WriterPools {
			if _, exist := newWriterPoolFuncs[pool.URN]; !exist {
				err = ErrWriterPoolURNNotExist
				return
			}

			if _, exist := poolDupCheck[pool.Name]; exist {
				err = ErrWriterPoolNameDuplicate
				return
			} else {
				poolDupCheck[pool.Name] = true
			}
			actorNames[actorTypedName(ActorReaderPool, pool.Name)] = true

			writerDupCheck := map[string]bool{}
			if _, exist := newWriterFuncs[pool.Writer.URN]; !exist {
				err = ErrWriterURNNotExist
				return
			}

			if _, exist := writerDupCheck[pool.Writer.Name]; exist {
				err = ErrWriterNameDuplicate
				return
			} else {
				writerDupCheck[pool.Writer.Name] = true
			}

			actorNames[actorTypedName(ActorWriter, pool.Writer.Name)] = true
			actorNames[actorTypedName(ActorWriterPool, pool.Name)] = true
		}
	}

	if p.InputTranslators != nil {
		dupCheck := map[string]bool{}
		for _, actor := range p.InputTranslators {
			if _, exist := newInputTranslatorFuncs[actor.URN]; !exist {
				err = ErrInputTranslatorURNNotExist
				return
			}

			if _, exist := dupCheck[actor.Name]; exist {
				err = ErrInputTranslatorNameDuplicate
				return
			} else {
				dupCheck[actor.Name] = true
			}
			actorNames[actorTypedName(ActorInputTranslator, actor.Name)] = true
		}
	}

	if p.OutputTranslators != nil {
		dupCheck := map[string]bool{}
		for _, actor := range p.OutputTranslators {
			if _, exist := newOutputTranslatorFuncs[actor.URN]; !exist {
				err = ErrOutputTranslatorURNNotExist
				return
			}

			if _, exist := dupCheck[actor.Name]; exist {
				err = ErrOutputTranslatorNameDuplicate
				return
			} else {
				dupCheck[actor.Name] = true
			}
			actorNames[actorTypedName(ActorOutputTranslator, actor.Name)] = true
		}
	}

	if p.Receivers != nil {
		dupCheck := map[string]bool{}
		for _, actor := range p.Receivers {
			if _, exist := newReceiverFuncs[actor.URN]; !exist {
				err = ErrReceiverURNNotExist
				return
			}

			if _, exist := dupCheck[actor.Name]; exist {
				err = ErrReceiverNameDuplicate
				return
			} else {
				dupCheck[actor.Name] = true
			}
			actorNames[actorTypedName(ActorReceiver, actor.Name)] = true
		}
	}

	if p.Senders != nil {
		dupCheck := map[string]bool{}
		for _, actor := range p.Senders {
			if _, exist := newSenderFuncs[actor.URN]; !exist {
				err = ErrSenderURNNotExist
				return
			}

			if _, exist := dupCheck[actor.Name]; exist {
				err = ErrSenderNameDuplicate
				return
			} else {
				dupCheck[actor.Name] = true
			}
			actorNames[actorTypedName(ActorSender, actor.Name)] = true
		}
	}

	if p.Inboxes != nil {
		dupCheck := map[string]bool{}
		for _, actor := range p.Inboxes {
			if _, exist := newInboxFuncs[actor.URN]; !exist {
				err = ErrInboxURNNotExist
				return
			}

			if _, exist := dupCheck[actor.Name]; exist {
				err = ErrInboxNameDuplicate
				return
			} else {
				dupCheck[actor.Name] = true
			}
			actorNames[actorTypedName(ActorInbox, actor.Name)] = true
		}
	}

	if p.Outboxes != nil {
		dupCheck := map[string]bool{}
		for _, actor := range p.Outboxes {
			if _, exist := newOutboxFuncs[actor.URN]; !exist {
				err = ErrOutboxURNNotExist
				return
			}

			if _, exist := dupCheck[actor.Name]; exist {
				err = ErrOutboxNameDuplicate
				return
			} else {
				dupCheck[actor.Name] = true
			}
			actorNames[actorTypedName(ActorOutbox, actor.Name)] = true
		}
	}

	if p.Routers != nil {
		dupCheck := map[string]bool{}
		for _, actor := range p.Routers {
			if _, exist := newRouterFuncs[actor.URN]; !exist {
				err = ErrRouterURNNotExist
				return
			}

			if _, exist := dupCheck[actor.Name]; exist {
				err = ErrRouterNameDuplicate
				return
			} else {
				dupCheck[actor.Name] = true
			}
			actorNames[actorTypedName(ActorRouter, actor.Name)] = true
		}
	}

	if p.Components != nil {
		for _, actor := range p.Components {
			if _, exist := newComponentFuncs[actor.URN]; !exist {
				err = ErrComponentURNNotExist
				return
			}
			actorNames[actorTypedName(ActorComponent, actor.Name)] = true
		}
	}

	if p.LabelMatchers != nil {
		dupCheck := map[string]bool{}
		for _, actor := range p.LabelMatchers {
			if _, exist := newLabelMatcherFuncs[actor.URN]; !exist {
				err = ErrLabelMatcherURNNotExist
				return
			}

			if _, exist := dupCheck[actor.Name]; exist {
				err = ErrLabelMatcherNameDuplicate
				return
			} else {
				dupCheck[actor.Name] = true
			}

			actorNames[actorTypedName(ActorLabelMatcher, actor.Name)] = true
		}
	}

	for _, router := range p.Compose {
		for _, compName := range router.Components {
			if _, exist := actorNames[actorTypedName(ActorComponent, compName)]; !exist {
				err = ErrActorNotExist
				return
			}
		}

		for _, inbox := range router.Inboxes {
			if _, exist := actorNames[actorTypedName(ActorInbox, inbox.Name)]; !exist {
				err = ErrActorNotExist
				return
			}

			for _, receiver := range inbox.Receivers {
				if _, exist := actorNames[actorTypedName(ActorReceiver, receiver.Name)]; !exist {
					err = ErrActorNotExist
					return
				}

				if _, exist := actorNames[actorTypedName(ActorInputTranslator, receiver.Translator)]; !exist {
					err = ErrActorNotExist
					return
				}

				if _, exist := actorNames[actorTypedName(ActorReaderPool, receiver.ReaderPool)]; !exist {
					err = ErrActorNotExist
					return
				}
			}
		}

		if _, exist := actorNames[actorTypedName(ActorLabelMatcher, router.LabelMatchers.Component)]; !exist {
			err = ErrActorNotExist
			return
		}

		if _, exist := actorNames[actorTypedName(ActorLabelMatcher, router.LabelMatchers.Outbox)]; !exist {
			err = ErrActorNotExist
			return
		}

		for _, outbox := range router.Outboxes {
			if _, exist := actorNames[actorTypedName(ActorOutbox, outbox.Name)]; !exist {
				err = ErrActorNotExist
				return
			}
			for _, sender := range outbox.Senders {
				if _, exist := actorNames[actorTypedName(ActorSender, sender.Name)]; !exist {
					err = ErrActorNotExist
					return
				}
				if _, exist := actorNames[actorTypedName(ActorOutputTranslator, sender.Translator)]; !exist {
					err = ErrActorNotExist
					return
				}

				if _, exist := actorNames[actorTypedName(ActorWriterPool, sender.WriterPool)]; !exist {
					err = ErrActorNotExist
					return
				}
			}
		}
	}

	return
}

func actorTypedName(actorType ActorType, name string) string {
	return string(actorType) + ":" + name
}
