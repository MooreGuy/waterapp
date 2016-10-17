package main

func Auth(username string, secret string) (authenticated bool) {
	return username == "test" && secret == "test"
}
