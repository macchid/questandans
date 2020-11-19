package questions

import "errors"

type Question interface {
	Read() string
	Answer(who, text string) error
	Change(who, text string) error
	Blame() (string, error)
}

type Questionaire interface {
	Add(who, text string) (Question, error)
	Change(question Question, who, newtext string) (Question, error)
	Delete(id string) (Question, error)
	Read() ([]Question, error)
	Peek(id string) (Question, error)
}

type questionaire []Question

func (qs *questionaire) Add(who, text string) (Question, error) {
	return nil, errors.New("Not implemented yet")
}

func (qs *questionaire) Change(question Question, who, newtext string) (Question, error) {
	return nil, errors.New("Not implemented yet")
}

func (qs *questionaire) Delete(id string) (Question, error) {
	return nil, errors.New("Not implemented yet")
}

func (qs *questionaire) Read() ([]Question, error) {
	return nil, errors.New("Not implemented yet")
}

func (qs *questionaire) Peek(id string) (Question, error) {
	return nil, errors.New("Not implemented yet")
}

type question struct {
	id	string, 
	who	string,	
	question string,
	answer string
}

func (q *question) Read() string {
	return ""
}

func (q *question) Answer(who, text string) error {
	return errors.New("Not implemented yet")
}

func (q *question) Change(who, text string) error {
	return errors.New("Not implemented yet")
}

func (q *question) Blame() string {
	return ""
}