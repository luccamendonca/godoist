package godoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

type Task struct {
	Content   string `json:"content"`
	Id        int    `json:"id,string"`
	ParentId  int    `json:"parent_id,string"`
	ProjectId int    `json:"project_id,string"`
	Url       string `json:"url"`
	Due       *struct {
		Date      string `json:"date"`
		String    string `json:"string"`
		Datetime  string `json:"datetime"`
		Recurring bool   `json:"recurring"`
	} `json:"due"`
}

type TaskRequest struct {
	Content   string `json:"content,omitempty"`
	DueString string `json:"due_string,omitempty"`
	ProjectId int    `json:"project_id,omitempty"`
}

type Project struct {
	Id             int    `json:"id,string"`
	IsInboxProject bool   `json:"is_inbox_project"`
	Name           string `json:"name"`
	Url            string `json:"url"`
}

func FetchProjectsByName(projectName string) ([]Project, error) {
	respProjects, err := makeRequest("GET", "/projects", nil)
	if err != nil {
		return nil, err
	}

	var projects []Project
	err = unmarshalHttpResponse(respProjects, &projects)
	if err != nil {
		return projects, err
	}

	var projectsByName []Project
	for _, project := range projects {
		if project.Name == projectName || projectName == "" {
			projectsByName = append(projectsByName, project)
		}
	}

	return projectsByName, nil
}

func FetchTasksByProjectName(projectName string) ([]Task, error) {
	var err error
	var tasks []Task
	projects, err := FetchProjectsByName(projectName)
	if err != nil || len(projects) == 0 {
		return tasks, err
	}
	project := projects[0]

	tasks, err = fetchTasks()
	if err != nil {
		return tasks, err
	}
	var tasksByProject []Task
	for _, task := range tasks {
		if task.ProjectId == project.Id {
			tasksByProject = append(tasksByProject, task)
		}
	}
	return tasksByProject, err
}

func CreateTask(taskName string) (Task, error) {
	taskRequest := TaskRequest{
		Content:   taskName,
		DueString: "today",
	}
	reqBodyBuf, err := json.Marshal(taskRequest)
	var task Task
	if err != nil {
		return task, err
	}
	resp, err := makeRequest("POST", "/tasks", bytes.NewBuffer(reqBodyBuf))
	if err != nil {
		return task, err
	}
	err = unmarshalHttpResponse(resp, &task)
	return task, err
}

func CreateTaskInProject(taskName string, projectName string) (Task, error) {
	return CreateTaskInProjectWithDue(taskName, projectName, "today")
}

func CreateTaskInProjectWithDue(taskName string, projectName string, dueString string) (Task, error) {
	projects, err := FetchProjectsByName(projectName)
	if err != nil {
		return Task{}, err
	}
	if len(projects) == 0 {
		return Task{}, fmt.Errorf("project '%s' not found", projectName)
	}

	project := projects[0]
	taskRequest := TaskRequest{
		Content:   taskName,
		DueString: dueString,
		ProjectId: project.Id,
	}
	reqBodyBuf, err := json.Marshal(taskRequest)
	var task Task
	if err != nil {
		return task, err
	}
	resp, err := makeRequest("POST", "/tasks", bytes.NewBuffer(reqBodyBuf))
	if err != nil {
		return task, err
	}
	err = unmarshalHttpResponse(resp, &task)
	return task, err
}

func fetchTasks() ([]Task, error) {
	var err error
	var tasks []Task
	resp, err := makeRequest("GET", "/tasks", nil)
	if err != nil {
		return tasks, err
	}
	err = unmarshalHttpResponse(resp, &tasks)
	return tasks, err
}

func unmarshalHttpResponse(resp *http.Response, model interface{}) error {
	defer resp.Body.Close()
	bytesProjects, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytesProjects, model)
	if err != nil {
		return err
	}
	return nil
}

type ParsedTask struct {
	Content   string
	DueString string
}

func getCustomWhenRules() []rules.Rule {
	resp := []rules.Rule{}
	resp = append(resp, &rules.F{
		RegExp: regexp.MustCompile("(?i)" +
			"(?:\\W|^)" +
			"(?:\\s*(?:(every(.+other)?)))?" +
			"(?:\\W|$)",
		),
		Applier: func(m *rules.Match, c *rules.Context, o *rules.Options, ref time.Time) (bool, error) {
			return true, nil
		},
	})

	return resp
}

func ParseTaskWithDate(input string) ParsedTask {
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)
	w.Add(getCustomWhenRules()...)

	result, err := w.Parse(input, time.Now())
	if err != nil || result == nil {
		return ParsedTask{
			Content:   strings.TrimSpace(input),
			DueString: "today",
		}
	}

	datePhrase := input[result.Index : result.Index+len(result.Text)]
	cleanContent := strings.Replace(input, result.Text, "", 1)
	cleanContent = strings.TrimSpace(cleanContent)
	cleanContent = strings.Join(strings.Fields(cleanContent), " ")

	return ParsedTask{
		Content:   cleanContent,
		DueString: datePhrase,
	}
}

func makeRequest(method string, path string, reqBody *bytes.Buffer) (*http.Response, error) {
	if reqBody == nil {
		reqBody = bytes.NewBuffer(nil)
	}
	baseURL := "https://api.todoist.com/rest/v2"
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", baseURL, path), reqBody)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", CONFIG.Todoist.ApiKey))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Request-Id", uuid.NewString())

	return http.DefaultClient.Do(req)
}
