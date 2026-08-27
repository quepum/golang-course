package dto

type RepoInfoResponse struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int32  `json:"stars"`
	Forks       int32  `json:"forks"`
	CreatedAt   string `json:"created_at"`
}
