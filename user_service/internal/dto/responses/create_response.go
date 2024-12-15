package responses

import "github.com/google/uuid"

type CreateResponse struct {
	Id uuid.UUID `json:"id"`
}
