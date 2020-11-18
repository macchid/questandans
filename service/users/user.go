package users

type User interface {
	Ask(text string) error
	Read() (questions.Questionaire, error)
	Answer(text string) (question Question, error)	 
}
