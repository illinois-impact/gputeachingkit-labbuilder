package build

type questionAnswer struct {
	Question string
	Answer   string
}

type doc struct {
	Module          int
	FileName        string
	Name            string
	Description     string
	QuestionAnswers []questionAnswer
	CodeTemplate    string
	CodeSolution    string
}

type resource struct {
	fileName string
	content  string
}
