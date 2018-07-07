package ruyi

type Request struct {
	Q            string `json:"q"`
	AppKey       string `json:"app_key"`
	UserID       string `json:"user_id"`
	ResetSession bool   `json:"reset_session"`
}

type Response struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result struct {
		Text    string   `json:"_text"`
		MsgID   string   `json:"msg_id"`
		Intents []Intent `json:"intents"`
	} `json:"result"`
}

type Intent struct {
	Action     string                 `json:"action"`
	Name       string                 `json:"name"`
	Parameters map[string]string      `json:"parameters"`
	Result     map[string]interface{} `json:"result"`
	Outputs    []Output               `json:"outputs"`
	Emotion    string                 `json:"emotion"`
}

type Output struct {
	Type     string            `json:"type"`
	Property map[string]string `json:"property"`
}
