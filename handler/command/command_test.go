package command

type mockOutputter struct {
	output []string
}

func (m *mockOutputter) Write(str string) {
	m.output = append(m.output, str)
}
