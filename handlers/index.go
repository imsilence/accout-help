package handlers

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/imsilence/account-help/domains"
	"github.com/imsilence/account-help/services"
)

var (
	org  = "htgolang"
	team = "htgolang"
)

//go:embed views/*.html
var view embed.FS

func Index(w http.ResponseWriter, r *http.Request) {
	type Url struct {
		Path string
		Text string
	}
	urls := []Url{
		{"/member/", "加入组织"},
		{"/invitations/", "待确认成员列表"},
		{"/members/", "已加入成员列表"},
		{"/repos/", "仓库列表"},
	}

	template.Must(template.ParseFS(view, "views/index.html")).Execute(w, struct {
		Urls []Url
	}{urls})

}

func Members(w http.ResponseWriter, r *http.Request) {
	members, err := services.GitlabService.Members(org, team)
	template.Must(template.ParseFS(view, "views/members.html")).Execute(w, struct {
		Error   error
		Members []*domains.Member
	}{err, members})
}

func Invitations(w http.ResponseWriter, r *http.Request) {
	members, err := services.GitlabService.Invitations(org, team)
	template.Must(template.ParseFS(view, "views/invitations.html")).Execute(w, struct {
		Error   error
		Members []*domains.Member
	}{err, members})
}

func CancelInvitation(w http.ResponseWriter, r *http.Request) {
	services.GitlabService.CancelInvitation(org, r.FormValue("id"))
	http.Redirect(w, r, "/invitations/", http.StatusFound)
}

func Member(w http.ResponseWriter, r *http.Request) {
	member, message := "", ""
	if r.Method == http.MethodPost {
		member = r.PostFormValue("member")
		message = services.GitlabService.AddTeamMember(org, team, member)
	}

	template.Must(template.ParseFS(view, "views/member.html")).Execute(w, struct {
		Message string
		Member  string
	}{message, member})
}

func Repos(w http.ResponseWriter, r *http.Request) {
	repos, err := services.GitlabService.Repos(org, team)
	template.Must(template.ParseFS(view, "views/repos.html")).Execute(w, struct {
		Error error
		Repos []*domains.Repo
	}{err, repos})
}
