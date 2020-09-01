package user

// UserService provides a function for creating a user from the given provider and providerUserID.
// The created user contains the global unique user ID that is used within your environment.
// user/hashuserservice is the most basic way to create a unique user by simply
// hashing both the provider and providerUserID. It would also be possible to load
// the user id from a DB using the provider and providerUserID as a key or just
// concatenate both strings.
type UserService interface {
	UniqueUser(provider string, providerUserID string) (string, error)
}
