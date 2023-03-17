package types

import "encoding/json"

// MessageNatsIn - сообщение от Nats
type MessageNatsIn struct {
	Head struct {
		DestVer string `json:"destVer"`
		Sender  string `json:"sender"`
		NetID   string `json:"netID"`
		Created string `json:"created"`
	} `json:"head"`
	Body struct {
		Command string `json:"command"`
		Params  struct {
			TemplateID string  `json:"template_id"`
			Template   JsonMap `json:"template"`
		} `json:"params"`
	} `json:"body"`
}

type JsonMap struct {
	ParamsMap map[string]string
}

func (j *JsonMap) UnmarshalJSON(data []byte) error {
	var err error

	//MapNew := make(map[string]interface{}, 0)
	err = json.Unmarshal(data, &j.ParamsMap)

	//j.ParamsMap = MapNew

	return err
}
