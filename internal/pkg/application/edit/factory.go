package edit

import (
	"context"
	"fmt"
	"strings"
	"time"

	ghub "github.com/google/go-github/github"
	"github.com/valeriobelli/gh-milestones/internal/pkg/domain/constants"
	"github.com/valeriobelli/gh-milestones/internal/pkg/infrastructure/gh"
	"github.com/valeriobelli/gh-milestones/internal/pkg/infrastructure/github"
	"github.com/valeriobelli/gh-milestones/internal/pkg/infrastructure/http"
)

type EditMilestoneConfig struct {
	DueOn       *time.Time
	State       *string
	Title       *string
	Description *string
	Verbose     bool
}

type EditMilestone struct {
	config EditMilestoneConfig
}

func NewEditMilestone(config EditMilestoneConfig) *EditMilestone {
	return &EditMilestone{config: config}
}

func (em EditMilestone) mapTitle(milestone *ghub.Milestone) string {
	if *milestone.Title == "" {
		return "<No Title>"
	}

	return *milestone.Title
}

func (em EditMilestone) mapDescription(milestone *ghub.Milestone) string {
	if *milestone.Description == "" {
		return "<No description>"
	}

	return *milestone.Description
}

func (em EditMilestone) mapDueOn(milestone *ghub.Milestone) string {
	if milestone.DueOn == nil {
		return "<No due date>"
	}

	return milestone.DueOn.Format(constants.DateFormat)
}

func (em EditMilestone) Execute(number int) {
	repoInfo, err := gh.RetrieveRepoInformation()

	if err != nil {
		fmt.Println(err.Error())

		return
	}

	client := github.NewRestClient(http.NewClient())

	milestone, _, err := client.Issues.EditMilestone(
		context.Background(),
		repoInfo.Owner,
		repoInfo.Name,
		number,
		&ghub.Milestone{
			Description: em.config.Description,
			DueOn:       em.config.DueOn,
			State:       em.config.State,
			Title:       em.config.Title,
		},
	)

	if err != nil {
		fmt.Println(err.Error())

		return
	}

	if em.config.Verbose {
		fmt.Println("Milestone has been edited. New status is:")
		fmt.Printf("  Title: %s\n", em.mapTitle(milestone))
		fmt.Printf("  Description: %s\n", em.mapDescription(milestone))
		fmt.Printf("  State: %s\n", strings.ToLower(*milestone.State))
		fmt.Printf("  Due On: %s\n", em.mapDueOn(milestone))
	}
}