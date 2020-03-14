package jwt

import (
	"fmt"
	fundation "test_redis"
	"testing"
)

func TestJWT_New(t *testing.T) {
	var jwt JWT
	var u fundation.User
	u.Username="2"
	jwt.New(u)
	jwt.Token="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJyZWRyb2NrIiwiZXhwIjoiMTU4NDA4MDU3NyIsImlhdCI6IjE1ODQwNjk3NzciLCJ1c2VybmFtZSI6IiJ9.K9Qdv33RMPKWd6oD+Dni4pdQDD5mlMwsn6F57QAZQgE="
	var k JWT
	k.Check(jwt.Token)
	fmt.Println(k)

}
