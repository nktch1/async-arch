package task

type Task struct {
	Description string `json:"description"`
	JiraID      string `json:"jira_id"`
	IsOpen      bool   `json:"is_open"`
	PopugID     string `json:"popug_id"`
}
