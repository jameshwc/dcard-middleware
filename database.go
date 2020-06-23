package main

type Database interface {
	// Check if the IP is in database and whether it's forbidden or not
	Find(string) (bool, bool)
	// return X-RateLimit-Remaining and X-RateLimit-Reset
	GetKey(string) (string, string, error)
	// If IP is not found in database, then create one
	SetKey(string) error
	// Increment the visit counter of the IP
	IncrementVisitByIP(string) error
}
