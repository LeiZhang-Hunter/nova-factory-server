package main

type exception struct{}

func newException() *exception {
	return &exception{}
}
