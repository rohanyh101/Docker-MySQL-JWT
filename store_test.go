package main

// Mocks

type MockStore struct{}

func (m *MockStore) CreateUser(u *User) (*User, error) {
	return &User{}, nil
}

func (m *MockStore) CreateTask(t *Task) (*Task, error) {
	return &Task{}, nil
}

func (m *MockStore) GetTask(id string) (*Task, error) {
	return &Task{}, nil
}

func (m *MockStore) GetUserByID(id string) (*User, error) {
	return &User{}, nil
}

func (m *MockStore) CreateProject(u *Project) (*Project, error) {
	return &Project{}, nil
}

func (m *MockStore) GetProjectByID(id string) (*Project, error) {
	return &Project{}, nil
}

func (m *MockStore) DeleteProject(id string) (int64, error) {
	return 0, nil
}
