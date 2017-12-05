package disneyland

type User struct {
	Username      string
	ProjectAccess string
	KindAccess    string
}

func (u *User) IsUser() bool {
	// if user
	if u.ProjectAccess != "ANY" {
		return true
	}
	return false
}
func (u *User) IsWorker() bool {
	// if worker
	if u.ProjectAccess == "ANY" && u.KindAccess != "ANY" {
		return true
	}
	return false
}

func (u *User) CanAccessJob(in *Job) bool {
	// if worker
	if u.IsWorker() && in.Kind != u.KindAccess {
		return false
	}
	// if user
	if u.IsUser() && in.Project != u.ProjectAccess {
		return false
	}
	return true
}
