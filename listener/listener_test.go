package listener

// Mock indicator

type mockIndicator struct {
	promptShown bool
}

func (m *mockIndicator) ShowPrompt() {
	m.promptShown = true
}

// Mock inputter

type mockInputter struct {
	scanCallback func() bool
	textCallback func() string
}

func (m mockInputter) Scan() bool {
	return m.scanCallback()
}

func (m mockInputter) Text() string {
	return m.textCallback()
}
