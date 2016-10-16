package indicator

type mockStringWriter struct {
	output []string
}

func (m *mockStringWriter) Write(str string) {
	m.output = append(m.output, str)
}
