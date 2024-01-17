package internal

type ImageInfo struct {
	ID          string            `json:"id"`
	RepoTags    []string          `json:"repo_tags"`
	Created     string            `json:"created"`
	Size        string            `json:"size"`
	VirtualSize string            `json:"virtual_size"`
	ParentID    string            `json:"parent_id"`
	RepoDigests []string          `json:"repo_digests"`
	Labels      map[string]string `json:"labels"`
}
