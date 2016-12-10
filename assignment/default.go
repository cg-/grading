package assignment

func GetDefaultAssignment() *Assignment {
	defaultCE1 := CommonError{
		Error:     "Test Error",
		Deduction: 10,
	}

	defaultCE2 := CommonError{
		Error:     "Test Error2",
		Deduction: 20,
	}

	CEs := make([]CommonError, 0)
	CEs = append(CEs, defaultCE1)
	CEs = append(CEs, defaultCE2)

	defaultQ1 := Question{
		Question:     "Test Question",
		Answer:       "Test Answer",
		Value:        20,
		CommonErrors: CEs,
	}
	defaultQ2 := Question{
		Question:     "Test Question2",
		Value:        20,
		CommonErrors: CEs,
	}

	Qs := make([]Question, 0)
	Qs = append(Qs, defaultQ1)
	Qs = append(Qs, defaultQ2)

	return &Assignment{
		Name:      "Assignment Name",
		Questions: Qs,
	}
}
