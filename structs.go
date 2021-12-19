package nextcloudapi

type User struct {
	Userid      string
	Password    string
	DisplayName string
	Email       string
	Groups      []string
	Subadmin    []string
	Quota       string
	Language    string
}
