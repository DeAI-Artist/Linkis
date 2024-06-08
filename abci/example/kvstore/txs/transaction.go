package txs

type Transaction struct {
	Msg       Message `json:"msg"`
	Signature []byte  `json:"signature"`
}

type Message struct {
	Type    uint8  `json:"type"`
	Content []byte `json:"content"`
}
