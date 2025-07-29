package servers

import (
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// Server represents a server in a GPU cluster
type Server struct {
	ID             string                `json:"id"`
	ImageID        *string               `json:"image_id"`
	TaskID         *string               `json:"task_id"`
	Flavor         string                `json:"flavor"`
	KeypairID      *string               `json:"keypair_id"`
	Name           string                `json:"name"`
	Status         string                `json:"status"`
	IPAddresses    []string              `json:"ip_addresses"`
	SecurityGroups []gcorecloud.ItemName `json:"security_groups"`
	Tags           []Tag                 `json:"tags"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
}

// Tag represents a server tag
type Tag struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	ReadOnly bool   `json:"read_only"`
}
