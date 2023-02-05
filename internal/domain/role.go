package domain

type Roles struct {
	IsAdmin  bool
	IsEditor bool
}

func GetRolesFromUser(user User) Roles {
	return Roles{
		IsAdmin:  user.IsAdmin,
		IsEditor: user.IsEditor,
	}
}
