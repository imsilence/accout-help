package services

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/imsilence/account-help/domains"

	"github.com/imroc/req"
)

type gitlabService struct {
}

var GitlabService = new(gitlabService)

func (s *gitlabService) header() req.Header {
	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(
		"%s:%s",
		os.Getenv("account.help.github.user"),
		os.Getenv("account.help.github.password"),
	)))

	return req.Header{
		"Accept":        "application/vnd.github.v3+json",
		"Authorization": fmt.Sprintf("Basic %s", token),
	}
}

func (s *gitlabService) AddTeamMember(org, team, member string) string {
	path := fmt.Sprintf("https://api.github.com/orgs/%s/teams/%s/memberships/%s", org, team, member)

	resp, err := req.Put(path, s.header())
	if err != nil {
		return err.Error()
	}
	statusCode := resp.Response().StatusCode
	log.Println("AddTeamMember", path, statusCode)
	switch statusCode {
	case http.StatusOK:
		var rs struct {
			State string
			Role  string
		}
		resp.ToJSON(&rs)
		message := ""
		switch rs.State {
		case "active":
			message = "已加入组织, 请尝试对仓库代码修改"
		case "pending":
			message = "请查找Github注册邮箱中邮件并确认加入组织"
		default:
			message = "未知"
		}
		return fmt.Sprintf("%s. 当前状态为: %s, 角色为: %s", message, rs.State, rs.Role)
	case http.StatusNotFound:
		return "账号信息不存在"
	}

	return "未知错误信息, 请联系管理员"
}

func (s *gitlabService) Members(org, team string) ([]*domains.Member, error) {
	path := fmt.Sprintf("https://api.github.com/orgs/%s/teams/%s/members", org, team)

	page, pre_page := 1, 100
	members := make([]*domains.Member, 0, pre_page)

	for {
		resp, err := req.Get(path, req.QueryParam{"per_page": pre_page, "page": page}, s.header())
		if err != nil {
			break
		}
		statusCode := resp.Response().StatusCode
		log.Println("Members", path, statusCode)
		if statusCode != http.StatusOK {
			break
		}
		var rs []struct {
			ID   int64  `json:"id"`
			Name string `json:"login"`
			Link string `json:"html_url"`
		}
		if err := resp.ToJSON(&rs); err != nil {
			log.Println("Error Members:", err)
			break
		}
		if len(rs) == 0 {
			break
		}

		for _, r := range rs {
			members = append(members, &domains.Member{ID: r.ID, Name: r.Name, Link: r.Link})
		}
		page++
	}

	return members, nil
}

func (s *gitlabService) Invitations(org, team string) ([]*domains.Member, error) {
	path := fmt.Sprintf("https://api.github.com/orgs/%s/teams/%s/invitations", org, team)

	page, pre_page := 1, 100
	members := make([]*domains.Member, 0, pre_page)

	for {
		resp, err := req.Get(path, req.QueryParam{"per_page": pre_page, "page": page}, s.header())
		if err != nil {
			break
		}
		statusCode := resp.Response().StatusCode
		log.Println("Invitations", path, statusCode)
		if statusCode != http.StatusOK {
			break
		}
		var rs []struct {
			ID        int64      `json:"id"`
			Name      string     `json:"login"`
			Link      string     `json:"html_url"`
			Email     string     `json:"email"`
			CreatedAt *time.Time `json:"created_at"`
		}
		if err := resp.ToJSON(&rs); err != nil {
			log.Println("Error Invitations:", err)
			break
		}
		if len(rs) == 0 {
			break
		}

		for _, r := range rs {
			members = append(members, &domains.Member{ID: r.ID, Name: r.Name, Link: r.Link, Email: r.Email, CreatedAt: r.CreatedAt})
		}
		page++
	}

	return members, nil
}

func (s *gitlabService) CancelInvitation(org, id string) error {
	path := fmt.Sprintf("https://api.github.com/orgs/%s/invitations/%s", org, id)

	resp, err := req.Delete(path, s.header())
	if err != nil {
		return err
	}
	statusCode := resp.Response().StatusCode
	log.Println("CancelInvitation", path, statusCode)
	return nil
}

func (s *gitlabService) Repos(org, team string) ([]*domains.Repo, error) {
	path := fmt.Sprintf("https://api.github.com/orgs/%s/teams/%s/repos", org, team)

	page, pre_page := 1, 100

	resp, err := req.Get(path, req.QueryParam{"per_page": pre_page, "page": page}, s.header())
	if err != nil {
		return nil, err
	}
	statusCode := resp.Response().StatusCode
	log.Println("Repos", path, statusCode)
	switch statusCode {
	case http.StatusOK:
		var rs []struct {
			ID        int64      `json:"id"`
			Name      string     `json:"name"`
			Link      string     `json:"html_url"`
			CreatedAt *time.Time `json:"created_at"`
		}
		resp.ToJSON(&rs)
		repos := make([]*domains.Repo, len(rs))
		for i, r := range rs {
			repos[i] = &domains.Repo{ID: r.ID, Name: r.Name, Link: r.Link, CreatedAt: r.CreatedAt}
		}
		return repos, nil
	}

	return nil, fmt.Errorf("未知错误, %d", statusCode)
}
