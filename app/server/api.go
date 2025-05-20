package server

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khaledhikmat/tr-extractor/service/config"
	"github.com/khaledhikmat/tr-extractor/service/data"
	"github.com/khaledhikmat/tr-extractor/service/storage"
	"github.com/khaledhikmat/tr-extractor/service/trello"

	"github.com/khaledhikmat/tr-extractor/job"
	jobattachments "github.com/khaledhikmat/tr-extractor/job/attachments"
	jobprops "github.com/khaledhikmat/tr-extractor/job/properties"
)

const (
	version = "1.0.0"
)

var jobProcs = map[data.JobType]job.Processor{
	data.JobTypeProperties:  jobprops.Processor,
	data.JobTypeAttachments: jobattachments.Processor,
}

func apiRoutes(ctx context.Context,
	r *gin.Engine,
	errorStream chan error,
	cfgsvc config.IService,
	datasvc data.IService,
	trsvc trello.IService,
	storagesvc storage.IService) {

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("version: %s - env: %s", version, cfgsvc.GetRuntimeEnvironment()),
		})
	})

	r.POST("/admins/reset", func(c *gin.Context) {
		_ = isPermitted(c, datasvc)
		if 1 == 1 {
			c.JSON(403, gin.H{
				"message": "Invalid or missing API key or not alllowed",
			})
			return
		}

		err := datasvc.ResetFactory()
		if err != nil {
			c.JSON(400, gin.H{
				"message": fmt.Sprintf("reset factory produced %s", err.Error()),
			})
			return
		}

		c.JSON(200, gin.H{
			"data": nil,
		})
	})

	r.GET("/properties", func(c *gin.Context) {
		isPermitted := isPermitted(c, datasvc)
		if !isPermitted {
			c.JSON(403, gin.H{
				"message": "Invalid or missing API key",
			})
			return
		}

		page, e := strconv.Atoi(c.Query("p"))
		if e != nil {
			page = 1
		}

		pageSize, e := strconv.Atoi(c.Query("s"))
		if e != nil {
			pageSize = 50
		}

		order := c.Query("o")
		if order == "" {
			order = "updated_at"
		}

		dir := c.Query("d")
		if dir == "" {
			dir = "desc"
		}

		//jobType := c.Query("t")
		props, err := datasvc.RetrieveProperties(page, pageSize, order, dir)
		if err != nil {
			c.JSON(400, gin.H{
				"message": fmt.Sprintf("retrieve properties produced %s", err.Error()),
			})
			return
		}

		c.JSON(200, gin.H{
			"data": props,
		})
	})

	r.GET("/jobs", func(c *gin.Context) {
		isPermitted := isPermitted(c, datasvc)
		if !isPermitted {
			c.JSON(403, gin.H{
				"message": "Invalid or missing API key",
			})
			return
		}

		jobID := c.Query("i")
		if jobID == "" {
			c.JSON(400, gin.H{
				"message": "job ID is required",
			})
			return
		}

		id, e := strconv.Atoi(jobID)
		if e != nil {
			c.JSON(400, gin.H{
				"message": "job ID could not be parsed",
			})
			return
		}

		job, err := datasvc.RetrieveJobByID(int64(id))
		if err != nil {
			c.JSON(400, gin.H{
				"message": fmt.Sprintf("retrieve job produced %s", err.Error()),
			})
			return
		}

		c.JSON(200, gin.H{
			"data": job,
		})
	})

	r.POST("/jobs", func(c *gin.Context) {
		isPermitted := isPermitted(c, datasvc)
		if !isPermitted {
			c.JSON(403, gin.H{
				"message": "Invalid or missing API key",
			})
			return
		}

		var job data.Job
		if err := c.ShouldBindJSON(&job); err != nil {
			c.JSON(400, gin.H{
				"message": fmt.Sprintf("invalid job: %s", err.Error()),
			})
			return
		}

		if job.Type == "" {
			c.JSON(400, gin.H{
				"message": "job type is required",
			})
			return
		}

		pageSize, e := strconv.Atoi(c.Query("s"))
		if e != nil {
			pageSize = 50
		}

		id, err := ProcessJob(ctx, job, pageSize, true, errorStream, cfgsvc, datasvc, trsvc, storagesvc)
		if err != nil {
			c.JSON(400, gin.H{
				"message": fmt.Sprintf("process job produced %s", err.Error()),
			})
			return
		}

		c.JSON(200, gin.H{
			"data": id, // TODO: Return a status URL
		})
	})

	r.POST("/errors", func(c *gin.Context) {
		isPermitted := isPermitted(c, datasvc)
		if !isPermitted {
			c.JSON(403, gin.H{
				"message": "Invalid or missing API key",
			})
			return
		}

		var thisError data.Error
		if err := c.ShouldBindJSON(&thisError); err != nil {
			c.JSON(400, gin.H{
				"message": fmt.Sprintf("invalid error: %s", err.Error()),
			})
			return
		}

		// Force an initial state
		err := datasvc.NewError(thisError.Source, thisError.Body)
		if err != nil {
			c.JSON(400, gin.H{
				"message": fmt.Sprintf("new error produced %s", err.Error()),
			})
			return
		}

		c.JSON(200, gin.H{
			"data": nil,
		})
	})
}

func ProcessJob(ctx context.Context,
	job data.Job,
	pageSize int,
	async bool,
	errorStream chan error,
	cfgsvc config.IService,
	datasvc data.IService,
	trsvc trello.IService,
	storagesvc storage.IService) (int64, error) {
	fmt.Printf("Processing job %s\n", job.Type)
	// Validate there is a processor for the job type
	proc, ok := jobProcs[job.Type]
	if !ok {
		return -1, fmt.Errorf("job type %s does not have a processor", job.Type)
	}

	// Check to make sure there is no existing job for the same type and channel
	isPending, err := datasvc.IsPendingJobsByType(job.Type)
	if err != nil {
		return -1, fmt.Errorf("is pending jobs by type and channel produced %s", err.Error())
	}
	if isPending {
		return -1, fmt.Errorf("job type %s is already pending", job.Type)
	}

	// Force an initial state
	job.State = data.JobStateQueued
	job.StartedAt = time.Now()
	id, err := datasvc.NewJob(job)
	if err != nil {
		return -1, fmt.Errorf("new job produced %s", err.Error())
	}

	if async {
		// Start the job processor asynchronously
		go proc(ctx, id, pageSize, errorStream, cfgsvc, datasvc, trsvc, storagesvc)
	} else {
		// Start the job processor synchronously
		proc(ctx, id, pageSize, errorStream, cfgsvc, datasvc, trsvc, storagesvc)
	}

	return id, nil
}

func isPermitted(c *gin.Context, datasvc data.IService) bool {
	apiKey := c.GetHeader("api-key")
	if apiKey == "" {
		return false
	}

	isvalid, err := datasvc.IsAPIKeyValid(apiKey)
	if err != nil {
		return false
	}

	if !isvalid {
		return false
	}

	return true
}
