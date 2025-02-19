package answer

type Answer struct {
	RecievedMessage string
	AnswerMessage   string
}

func CreateAnswer(recievedMessage string) Answer {
	answer := Answer{
		RecievedMessage: recievedMessage,
		AnswerMessage:   recievedMessage,
	}

	return answer
}
