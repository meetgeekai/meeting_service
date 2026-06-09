package utils

import (
	"fmt"

	"github.com/meetgeekai/go-common/utils"
)

func GetAddr() string {
	host := utils.GetEnv[string]("HOST")
	port := utils.GetEnv[int]("PORT")

	return fmt.Sprintf("%s:%d", host, port)
}
