package components

import (
	"strings"

	"github.com/RyannKim327/constellation/config/utils"
)

type navItem struct {
	Name string
	URL  string
	Key  string
}

// Header renders the shared site navigation with the current page highlighted.
func Header(currentPage string) string {
	items := []navItem{
		{Name: "Home", URL: "/", Key: ""},
		{Name: "Compose", URL: "/compose", Key: "compose"},
		{Name: "Read", URL: "/read", Key: "read"},
		{Name: "About", URL: "/about", Key: "about"},
	}

	if len(currentPage) > 0 && currentPage[0] == '/' {
		currentPage = currentPage[1:]
	}

	navigation := ""
	for _, item := range items {
		active := ""
		isActive := false
		if item.Key == "" {
			isActive = currentPage == ""
		} else {
			isActive = currentPage == item.Key || strings.HasPrefix(currentPage, item.Key+"/")
		}
		if isActive {
			active = " active"
		}
		navigation += utils.HTParser(`
			<a class="nav-link{active}" href="{url}">{name}</a>
		`, utils.PARAMS{
			"active": active,
			"url":    item.URL,
			"name":   item.Name,
		})
	}

	return utils.HTParser(`
		<header class="site-header">
			<a class="brand" href="/">Clio<span class="brand-dot">.</span></a>
			<nav class="site-nav">{navigation}</nav>
		</header>
	`, utils.PARAMS{
		"navigation": navigation,
	})
}
