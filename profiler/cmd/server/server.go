package server

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xhermitx/gitpulse-01/profiler/types"
	"github.com/xhermitx/gitpulse-01/profiler/utils"
)

const (
	QUEUE__JOB_STATUS = "JOB_STATUS_QUEUE"
)

type Server struct {
	Store types.CandidateStore
	Git   types.GitService
	Queue types.Queue
	Cache types.Cache
}

func NewServer(s types.CandidateStore, g types.GitService, q types.Queue, c types.Cache) Server {
	return Server{
		Store: s,
		Git:   g,
		Queue: q,
		Cache: c,
	}
}

func (s Server) Run() error {
	msgs, err := s.Queue.Subscribe(QUEUE__JOB_STATUS)
	if err != nil {
		return err
	}

	go s.handleQueueData(msgs)

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-msgs
	return nil
}

func (s Server) handleQueueData(msgs <-chan amqp.Delivery) {
	if msgs == nil {
		log.Println("empty body from queue")
	}

	for d := range msgs {
		var jobQueue types.JobQueue
		if err := json.Unmarshal(d.Body, &jobQueue); err != nil {
			utils.LogError(err, "failed to parse candidate data")
			// FIXME: save the candidate name as unparsed
		}
		// Fetch user details from git
		for _, id := range jobQueue.GithubIDs {
			res, err := s.Git.FetchUserDetails(id)
			if err != nil {
				// FIXME: Handle error
				_ = err
			}

			// Store Git response in DB
			if err := s.handleGitData(jobQueue.JobId, res); err != nil {
				// FIXME: Handle error
				_ = err
			}
		}
	}
}

func (s Server) handleGitData(jobId string, u *types.GitUser) error {
	languages := getLanguages(u.TopRepo.Nodes[0].Languages)
	languages = append(languages, getLanguages(u.TopContributedRepo.Nodes[0].Languages)...)

	topics := getTopics(u.TopRepo.Nodes[0].Topics)
	topics = append(topics, getTopics(u.TopContributedRepo.Nodes[0].Topics)...)

	candidate := types.Candidate{
		CandidateId:             uuid.NewString(),
		CandidateMeta:           u.CandidateMeta,
		TotalContributions:      u.Contributions.ContributionCalendar.TotalContributions,
		TotalFollowers:          u.Followers.TotalCount,
		TopRepo:                 u.TopRepo.Nodes[0].Name,
		TopRepoStars:            u.TopRepo.Nodes[0].Stargazers.TotalCount,
		TopContributedRepo:      u.TopContributedRepo.Nodes[0].Name,
		TopContributedRepoStars: u.TopContributedRepo.Nodes[0].Stargazers.TotalCount,
		Languages:               languages,
		Topics:                  topics,
		JobId:                   jobId,
	}

	// Store in DB
	if err := s.Store.SaveCandidate(&candidate); err != nil {
		return err
	}

	return nil
}

// fetch top languages used in the repositories
func getLanguages(v types.LanguageConnection) []string {
	var languages []string
	for _, l := range v.Nodes {
		languages = append(languages, l.Name)
	}
	return languages
}

// fetch the topics used in the repositories
func getTopics(v types.RepositoryTopics) []string {
	var topics []string
	for _, t := range v.Nodes {
		topics = append(topics, t.Topic.Name)
	}
	return topics
}
