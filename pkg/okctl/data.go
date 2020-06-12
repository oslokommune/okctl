package okctl

func (o *Okctl) Username() string {
	return o.AppData.User.Username
}

func (o *Okctl) Region() string {
	return o.RepoData.Region
}
