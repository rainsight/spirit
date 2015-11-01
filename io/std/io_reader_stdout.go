package std

import (
	"io"
	"sync"

	"github.com/gogap/spirit"
)

const (
	stdReaderURN = "urn:spirit:reader:std"
)

var _ io.ReadCloser = new(Stdout)

type Stdout struct {
	onceInit sync.Once
	conf     StdIOConfig
	delim    string

	proc *STDProcess
}

func init() {
	spirit.RegisterReader(stdReaderURN, NewStdout)
}

func NewStdout(options spirit.Options) (w io.ReadCloser, err error) {
	conf := StdIOConfig{}
	options.ToObject(&conf)

	w = &Stdout{
		conf: conf,
	}
	return
}

func (p *Stdout) Read(data []byte) (n int, err error) {
	if p.proc == nil {
		p.onceInit.Do(func() {
			if proc, e := takeSTDIO(_Input, p.conf); e != nil {
				err = e
			} else {
				p.proc = proc
				p.proc.Start()
			}
		})
	}

	if p.proc == nil || err != nil {
		return
	}

	return p.proc.Read(data)
}

func (p *Stdout) Close() (err error) {
	if p.proc == nil {
		return
	}

	p.proc.Stop()

	return
}