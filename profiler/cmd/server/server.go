package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xhermitx/gitpulse-01/profiler/types"
	"github.com/xhermitx/gitpulse-01/profiler/utils"
)

const (
	QUEUE__JOB_STATUS     = "JOB_STATUS_QUEUE"
	UNPARSED_CACHE_PREFIX = "UNPARSED: "
	PARSED_CACHE_PREFIX   = "PARSED: "
)

type Server struct {
	store types.CandidateStore
	git   types.GitService
	queue types.Queue
	cache types.Cache
}

func NewServer(s types.CandidateStore, g types.GitService, q types.Queue, c types.Cache) Server {
	return Server{
		store: s,
		git:   g,
		queue: q,
		cache: c,
	}
}

func (s Server) Run() error {
	msgs, err := s.queue.Subscribe(QUEUE__JOB_STATUS)
	if err != nil {
		return err
	}

	go s.handleQueueData(msgs)

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	select {}
}

func (s Server) handleQueueData(msgs <-chan amqp.Delivery) {
	if msgs == nil {
		log.Println("empty body from queue")
	}

	for d := range msgs {
		var jobQueue types.JobQueue
		if err := json.Unmarshal(d.Body, &jobQueue); err != nil {
			utils.LogError(err, "failed to parse candidate data")
			if err := s.cache.Append(context.Background(), UNPARSED_CACHE_PREFIX+jobQueue.JobId, jobQueue.Filename+" "); err != nil {
				fmt.Printf("\nfailed to cache %s: %v", jobQueue.Filename, err)
			}
		}

		fmt.Println("Queue Msg: ", jobQueue)

		// Fetch user details from git
		for _, id := range jobQueue.GithubIDs {

			res, err := s.git.FetchUserDetails(id)
			if err != nil {
				// FIXME: Handle this Better
				fmt.Printf("\nError fetching details for %s: %v", id, err)
				continue
			}

			fmt.Println("Git User: ", res.Name)
			// Store Git response in DB
			if err := s.handleGitData(jobQueue.JobId, res); err != nil {
				// FIXME: Handle this Better
				fmt.Printf("\nError saving details for %s: %v", id, err)
				continue
			}

			tmp, err := s.cache.Get(context.Background(), PARSED_CACHE_PREFIX+jobQueue.JobId)
			if err != nil {
				log.Println("Failed to get PARSED CACHE", err)
			}

			n, err := strconv.Atoi(tmp)
			if err != nil {
				log.Println("failed to convert the cache value to int", err)
				continue
			}
			if err := s.cache.Set(context.Background(), PARSED_CACHE_PREFIX+id, n+1, 0); err != nil {
				log.Println("failed to update PARSED CACHE")
			}
		}
	}
}

func (s Server) handleGitData(jobId string, u *types.GitUser) error {
	candidate := types.Candidate{
		CandidateId:        uuid.NewString(),
		CandidateMeta:      u.CandidateMeta,
		TotalContributions: u.Contributions.ContributionCalendar.TotalContributions,
		TotalFollowers:     u.Followers.TotalCount,
		JobId:              jobId,
	}

	if len(u.TopRepo.Nodes) > 0 {
		candidate.TopRepo = u.TopRepo.Nodes[0].Name
		candidate.TopRepoStars = u.TopRepo.Nodes[0].Stargazers.TotalCount
		candidate.Languages += len(u.TopRepo.Nodes[0].Languages.Nodes)
		candidate.Topics += len(u.TopRepo.Nodes[0].Topics.Nodes)
	}

	if len(u.TopContributedRepo.Nodes) > 0 {
		candidate.TopContributedRepo = u.TopContributedRepo.Nodes[0].Name
		candidate.TopContributedRepoStars = u.TopContributedRepo.Nodes[0].Stargazers.TotalCount
		candidate.Languages += len(u.TopContributedRepo.Nodes[0].Languages.Nodes)
		candidate.Topics += len(u.TopContributedRepo.Nodes[0].Topics.Nodes)
	}

	// Store in DB
	if err := s.store.SaveCandidate(&candidate); err != nil {
		return err
	}
	return nil
}
