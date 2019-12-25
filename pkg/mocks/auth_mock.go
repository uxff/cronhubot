package mocks

type AuthMock struct {}

func (a *AuthMock) CheckPassword(user, sign string) bool {
	return true
}

func NewAuth() *AuthMock {
	return &AuthMock{}
}
