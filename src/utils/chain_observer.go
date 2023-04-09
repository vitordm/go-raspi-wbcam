package utils

// Define a struct do sujeito
type ChainByteSubject struct {
	observers [](chan []byte)
}

func (subject *ChainByteSubject) Attach(c chan []byte) {
	subject.observers = append(subject.observers, c)
}

func (subject *ChainByteSubject) Detach(c chan []byte) {
	for i, observer := range subject.observers {
		if observer == c {
			subject.observers = append(subject.observers[:i], subject.observers[i+1:]...)
			break
		}
	}
}

func (subject *ChainByteSubject) Notify(data []byte) {
	for _, observer := range subject.observers {
		go func(c chan []byte) {
			c <- data
		}(observer)
	}
}

func (subject *ChainByteSubject) IsEmpty() bool {
	return len(subject.observers) <= 0
}
