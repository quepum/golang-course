package dto

type AddSubscriptionRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type AddSubscriptionResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type SubscriptionItem struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type ListSubscriptionsResponse struct {
	Subscriptions []SubscriptionItem `json:"subscriptions"`
}

type RemoveSubscriptionResponse struct {
	Status bool `json:"status"`
}

type RepoStatItem struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int32  `json:"stars"`
	Forks       int32  `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

type SubscriptionsInfoResponse struct {
	Stats []RepoStatItem `json:"stats"`
}
