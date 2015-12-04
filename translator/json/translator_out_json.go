package json

import (
	"encoding/json"
	"io"

	"github.com/gogap/spirit"
)

const (
	jsonTranslatorOutURN = "urn:spirit:translator:out:json"
)

var _ spirit.OutputTranslator = new(JSONOutputTranslator)

type JSONOutputTranslatorConfig struct {
	DataOnly bool `json:"data_only"`
}

type JSONOutputTranslator struct {
	name string
	conf JSONOutputTranslatorConfig
}

func init() {
	spirit.RegisterOutputTranslator(jsonTranslatorOutURN, NewJSONOutputTranslator)
}

func NewJSONOutputTranslator(name string, options spirit.Map) (translator spirit.OutputTranslator, err error) {
	conf := JSONOutputTranslatorConfig{}

	if err = options.ToObject(&conf); err != nil {
		return
	}

	translator = &JSONOutputTranslator{
		name: name,
		conf: conf,
	}
	return
}

func (p *JSONOutputTranslator) Name() string {
	return p.name
}

func (p *JSONOutputTranslator) URN() string {
	return jsonTranslatorOutURN
}

func (p *JSONOutputTranslator) outDataOnly(w io.Writer, delivery spirit.Delivery) (err error) {

	var data interface{}
	if data, err = delivery.Payload().GetData(); err != nil {
		return
	}

	switch d := data.(type) {
	case []byte:
		{
			_, err = w.Write(d)
		}
	case string:
		{
			_, err = w.Write([]byte(d))
		}
	default:
		encode := json.NewEncoder(w)
		err = encode.Encode(data)
	}

	return
}

func (p *JSONOutputTranslator) outDeliveryData(w io.Writer, delivery spirit.Delivery) (err error) {

	var data interface{}
	if data, err = delivery.Payload().GetData(); err != nil {
		return
	}

	payload := _JSONPayload{
		Id:      delivery.Payload().Id(),
		Data:    data,
		Errors:  delivery.Payload().Errors(),
		Context: delivery.Payload().Context(),
	}

	jd := _JSONDelivery{
		Id:        delivery.Id(),
		URN:       delivery.URN(),
		SessionId: delivery.SessionId(),
		Payload:   payload,
		Timestamp: delivery.Timestamp(),
		Metadata:  delivery.Metadata(),
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(jd)

	return
}

func (p *JSONOutputTranslator) Out(w io.Writer, delivery spirit.Delivery) (err error) {

	if p.conf.DataOnly {
		return p.outDataOnly(w, delivery)
	}

	return p.outDeliveryData(w, delivery)
}
