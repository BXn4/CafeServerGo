package objects

type Player struct {
	ID                  int
	Cash                int
	Gold                int
	XP                  int
	InstantCookings     int
	OpenJobs            int
	PlayedWheel         bool
	AllowFriendRequests bool
	AllowEmails         bool
	EmailVerified       bool
	NewGifts            int
	Username            string
	Avatar              Avatar
	Mastery             string // TODO: Create proper mastery
}
