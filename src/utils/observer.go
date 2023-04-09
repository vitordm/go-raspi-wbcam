package utils

import "fmt"

// Define a interface para o observador
type Observer interface {
	Update()
}

// Define a struct do sujeito
type Subject struct {
	observers []Observer
	state     string
}

// Adiciona um observador à lista
func (s *Subject) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

// Remove um observador da lista
func (s *Subject) Detach(o Observer) {
	for i, observer := range s.observers {
		if observer == o {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

// Notifica todos os observadores sobre mudanças de estado
func (s *Subject) Notify() {
	for _, observer := range s.observers {
		observer.Update()
	}
}

// Define o estado do sujeito
func (s *Subject) SetState(state string) {
	s.state = state
	s.Notify()
}

// Define a struct do observador
type ConcreteObserver struct{}

// Implementa a interface do observador
func (o ConcreteObserver) Update() {
	fmt.Println("Observador notificado!")
}
