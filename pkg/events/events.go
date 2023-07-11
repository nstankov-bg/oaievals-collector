
package events


type Data map[string]interface{}


type Event struct {
	RunID     string `json:"run_id"`
	Type      string `json:"type"`
	EventID   int    `json:"event_id"`
	Data      Data   `json:"data"`
	SampleID  string `json:"sample_id"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
}
