package business

import "fmt"

var AlreadyExistUserErr = fmt.Errorf("user already exist")
var SecretNumNotValid = fmt.Errorf("secret_num is not valid")
var ScoreNotValid = fmt.Errorf("score is not valid")
var UsernameNotValid = fmt.Errorf("username is not valid")
