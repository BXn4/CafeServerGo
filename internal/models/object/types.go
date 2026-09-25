/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package object

type CafeObjectKind int

const (
	STOVE   CafeObjectKind = 0
	COUNTER CafeObjectKind = 1
	CHAIR   CafeObjectKind = 2
	TABLE   CafeObjectKind = 3
	VENDING CafeObjectKind = 4
	OTHER   CafeObjectKind = 5
)

type CafeObjectRotation int

const (
	Up CafeObjectRotation = iota
	Left
	Down
	Right
)
