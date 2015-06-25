package tests

import "github.com/revel/revel/testing"

type AppTest struct {
	testing.TestSuite
}

func (t *AppTest) Before() {
	println("Set up")
}

func (t *AppTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *AppTest) TestThatCommentBasicWorks() {
	t.Get("/commentbasic/apptest")
t.AssertOk()
t.AssertContentType("text/html; charset=utf-8")
}

func (t *AppTest) After() {
	println("Tear down")
}
