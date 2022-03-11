// Package kitsu provides methods for Kitsu task management software
package kitsu

import (
	"app/src/utils"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Task struct {
	Assignees       []string    `json:"assignees,omitempty"`
	ID              string      `json:"id,omitempty"`
	CreatedAt       string      `json:"created_at,omitempty"`
	UpdatedAt       string      `json:"updated_at,omitempty"`
	Name            string      `json:"name,omitempty"`
	LastCommentDate string      `json:"last_comment_date,omitempty"`
	Data            interface{} `json:"data,omitempty"`
	ProjectID       string      `json:"project_id,omitempty"`
	TaskTypeID      string      `json:"task_type_id,omitempty"`
	TaskStatusID    string      `json:"task_status_id,omitempty"`
	EntityID        string      `json:"entity_id,omitempty"`
	AssignerID      string      `json:"assigner_id,omitempty"`
	Type            string      `json:"type,omitempty"`
}
type Tasks struct {
	Each []Task
}

type Person struct {
	ID                        string `json:"id,omitempty"`
	CreatedAt                 string `json:"created_at,omitempty"`
	UpdatedAt                 string `json:"updated_at,omitempty"`
	FirstName                 string `json:"first_name,omitempty"`
	LastName                  string `json:"last_name,omitempty"`
	Email                     string `json:"email,omitempty"`
	Phone                     string `json:"phone,omitempty"`
	Active                    bool   `json:"active,omitempty"`
	LastPresence              string `json:"last_presence,omitempty"`
	DesktopLogin              string `json:"desktop_login,omitempty"`
	ShotgunID                 string `json:"shotgun_id,omitempty"`
	Timezone                  string `json:"timezone,omitempty"`
	Locale                    string `json:"locale,omitempty"`
	Data                      string `json:"data,omitempty"`
	Role                      string `json:"role,omitempty"`
	HasAvatar                 bool   `json:"has_avatar,omitempty"`
	NotificationsEnabled      bool   `json:"notifications_enabled,omitempty"`
	NotificationsSlackEnabled bool   `json:"notifications_slack_enabled,omitempty"`
	NotificationsSlackUserid  string `json:"notifications_slack_userid,omitempty"`
	Type                      string `json:"type,omitempty"`
	FullName                  string `json:"full_name,omitempty"`
}

type Persons struct {
	Each []Person
}

type Entity struct {
	EntitiesOut     []interface{} `json:"entities_out,omitempty"`
	InstanceCasting []interface{} `json:"instance_casting,omitempty"`
	CreatedAt       string        `json:"created_at,omitempty"`
	UpdatedAt       string        `json:"updated_at,omitempty"`
	ID              string        `json:"id,omitempty"`
	Name            string        `json:"name,omitempty"`
	Code            interface{}   `json:"code,omitempty"`
	Description     interface{}   `json:"description,omitempty"`
	ShotgunID       interface{}   `json:"shotgun_id,omitempty"`
	Canceled        bool          `json:"canceled,omitempty"`
	NbFrames        interface{}   `json:"nb_frames,omitempty"`
	ProjectID       string        `json:"project_id,omitempty"`
	EntityTypeID    string        `json:"entity_type_id,omitempty"`
	ParentID        string        `json:"parent_id,omitempty"`
	SourceID        interface{}   `json:"source_id,omitempty"`
	PreviewFileID   interface{}   `json:"preview_file_id,omitempty"`
	Data            interface{}   `json:"data,omitempty"`
	EntitiesIn      []interface{} `json:"entities_in,omitempty"`
	Type            string        `json:"type,omitempty"`
}

type Entities struct {
	Each []Entity
}

type EntityType struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type EntityTypes struct {
	Each []EntityType
}

type TaskStatuses struct {
	Each []TaskStatus
}

type TaskStatus struct {
	ID              string      `json:"id,omitempty"`
	CreatedAt       string      `json:"created_at,omitempty"`
	UpdatedAt       string      `json:"updated_at,omitempty"`
	Name            string      `json:"name,omitempty"`
	ShortName       string      `json:"short_name,omitempty"`
	Color           string      `json:"color,omitempty"`
	IsDone          bool        `json:"is_done,omitempty"`
	IsArtistAllowed bool        `json:"is_artist_allowed,omitempty"`
	IsClientAllowed bool        `json:"is_client_allowed,omitempty"`
	IsRetake        bool        `json:"is_retake,omitempty"`
	ShotgunID       interface{} `json:"shotgun_id,omitempty"`
	IsReviewable    bool        `json:"is_reviewable,omitempty"`
	Type            string      `json:"type,omitempty"`
}

type Comment struct {
	ID        string      `json:"id,omitempty"`
	CreatedAt string      `json:"created_at,omitempty"`
	UpdatedAt string      `json:"updated_at,omitempty"`
	ShotgunID interface{} `json:"shotgun_id,omitempty"`
	ObjectID  string      `json:"object_id,omitempty"`
	PersonID  string      `json:"person_id,omitempty"`
	Text      string      `json:"text,omitempty"`
}

type Comments struct {
	Each []Comment
}

type TaskType struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	ShortName string `json:"short_name,omitempty"`
}

type TaskTypes struct {
	Each []TaskType
}

type Project struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	ProjectStatusID string `json:"project_status_id,omitempty"`
}

type Projects struct {
	Each []Project
}

type ProjectStatus struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ProjectStatuses struct {
	Each []ProjectStatus
}

type Attachment struct {
	ID        string `json:"id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Name      string `json:"name,omitempty"`
	Size      int    `json:"size,omitempty"`
	Extension string `json:"extension,omitempty"`
	Mimetype  string `json:"mimetype,omitempty"`
	CommentID string `json:"comment_id,omitempty"`
	Comment   struct {
		ObjectID   string `json:"object_id,omitempty"`
		ObjectType string `json:"object_type,omitempty"`
	}
}

type Attachments struct {
	Each []Attachment
}

func GetComments() Comments {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/comments"
	response := Comments{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response.Each)

	return response
}

func GetComment(objectID string) Comments {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/comments?object_id=" + objectID
	response := Comments{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response.Each)

	return response
}

func GetTasks() Tasks {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/tasks?relations=true"
	response := Tasks{}
	println(os.Getenv("KitsuJWTToken"))
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response.Each)

	return response
}

func GetTask(taskID string) Task {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/tasks/" + taskID
	response := Task{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

func GetPerson(personID string) Person {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/persons/" + personID
	response := Person{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

func GetPersons() Persons {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/persons/"
	response := Persons{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response.Each)

	return response
}

func GetEntities() Entities {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/entities/"
	response := Entities{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response.Each)

	return response
}

func GetEntity(EntityID string) Entity {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/entities/" + EntityID
	response := Entity{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

func GetEntityTypes() EntityTypes {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/entity-types/"
	response := EntityTypes{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response.Each)

	return response
}

func GetEntityType(entityTypeID string) EntityType {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/entity-types/" + entityTypeID
	response := EntityType{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

func GetTaskStatuses() TaskStatuses {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/task-status/"
	response := TaskStatuses{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response.Each)

	return response
}

func GetTaskStatus(taskStatusID string) TaskStatus {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/task-status/" + taskStatusID
	response := TaskStatus{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

func GetTaskType(taskID string) TaskType {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/task-types/" + taskID
	response := TaskType{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

func GetTaskTypes() TaskTypes {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/task-types/"
	response := TaskTypes{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response.Each)

	return response
}

func GetProject(projectID string) Project {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/projects/" + projectID
	response := Project{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

func GetProjects() Projects {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/projects/"
	response := Projects{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response.Each)

	return response
}

func GetProjectStatus(projectStatusID string) ProjectStatus {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/project-status/" + projectStatusID
	response := ProjectStatus{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

func GetAttachments() Attachments {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/attachment-files/"
	response := Attachments{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response.Each)
	return response
}

func GetAttachment(AttachmentID string) Attachment {
	path := utils.ConfRead().Kitsu.Hostname + "api/data/attachment-files/" + AttachmentID
	response := Attachment{}
	utils.Request(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)
	return response
}

func DownloadAttachment(localPath, id, filename string, conf utils.Config) (int64, error) {
	// Create dir
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		err := os.Mkdir(localPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	// Create the file
	out, err := os.Create(localPath + "/" + filename)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Make request
	path := utils.ConfRead().Kitsu.Hostname + "api/data/attachment-files/" + id + "/file/" + filename
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		panic(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("KitsuJWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		//return fmt.Errorf("bad status: %s", resp.Status)
		//panic("bad status:" + resp.Status)
		return 0, fmt.Errorf(resp.Status)
	}

	// Writer the body to file
	size, err := io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

	return size, nil
}
