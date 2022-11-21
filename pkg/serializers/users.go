package serializers

import "github.com/Dhruv9449/mou/pkg/models"

func UserSerializer(user models.User) map[string]interface{} {
	return map[string]interface{}{
		"id":      user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"picture": user.Picture,
		"role":    user.Role,
	}
}

func UserLoginSerializer(user models.User, token string) map[string]interface{} {
	return map[string]interface{}{
		"id":      user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"picture": user.Picture,
		"token":   token,
	}
}

func UserBlockSerializer(user models.User) map[string]interface{} {
	return map[string]interface{}{
		"id":      user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"picture": user.Picture,
	}
}

func UserListSerializer(users []models.User) []map[string]interface{} {
	var userList []map[string]interface{}
	for _, user := range users {
		userList = append(userList, UserBlockSerializer(user))
	}
	return userList
}
