package attack


type AttackOptions struct {
	Host string
	Port string
	Script string
	UsernameFile string
	PasswordFile string
	Timeout int
	Threads int
	MaxAttempts int
	RateLimit int

}
