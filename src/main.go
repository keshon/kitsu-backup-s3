package main

import (
	"app/src/api/kitsu"
	"app/src/api/s3"
	"app/src/model"
	"app/src/utils"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// Create logger with rotation
	setupLogger()

	// Read config
	conf := utils.ConfRead()
	log.Info("[main.go][main] Config read successfully")

	// Auth to Kitsu to get JWT token
	JWTToken := utils.AuthForJWTToken(conf.Kitsu.Hostname+"api/auth/login", conf.Kitsu.Email, conf.Kitsu.Password)
	os.Setenv("KitsuJWTToken", JWTToken)
	log.Info("[main.go][main] JWT token acquired")

	// Connect to DB
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		log.Error("[main.go][main] Failed to connect database")
		os.Exit(1)
	}
	db.AutoMigrate(&model.Attachment{})

	// Setup CRON on schedule
	c := cron.New(cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	))

	// Update tray icon
	log.Info("[main.go][main] Parse all attachments on first run")
	utils.EmptyDir(conf.Backup.LocalStorage)
	parseAllAttachments(conf, db)
	c.AddFunc("@every "+strconv.Itoa(conf.Backup.PollDuration)+"m", func() {
		log.Info("[main.go][main] Parse all attachments on CRON job")
		utils.EmptyDir(conf.Backup.LocalStorage)
		parseAllAttachments(conf, db)

	})
	log.Info("[main.go][main] Run CRON")
	c.Run()
}

func parseAllAttachments(conf utils.Config, db *gorm.DB) {
	log.Info("[main.go][parseAllAttachments] Started parsing all attachments")

	// Get all Attachments
	array := kitsu.GetAttachments()

	if len(array.Each) <= 0 {
		return
	}

	// Concurent threads from conf
	threads := conf.Backup.Threads

	var count int

	if threads < 0 {
		// Async
		var wg sync.WaitGroup
		wg.Add(len(array.Each))

		for _, elem := range array.Each {
			go func(elem kitsu.Attachment) {
				defer wg.Done()
				resp := parseSingleAttachment(conf, db, elem)
				if resp {
					count++
				}
			}(elem)
		}
		wg.Wait()

	} else if threads == 0 {
		// Sync
		for _, elem := range array.Each {
			resp := parseSingleAttachment(conf, db, elem)
			if resp {
				count++
			}
		}

	} else if threads > 0 {
		// Semafore async
		var sem = make(chan int, threads)

		for _, elem := range array.Each {
			sem <- 1
			go func() {

				resp := parseSingleAttachment(conf, db, elem)
				if resp {
					count++
				}
				<-sem
			}()
		}

	}

	if count > 0 {

	}
	log.Info("[main.go][parseAllAttachments] Finished parsing all attachments")

}

func parseSingleAttachment(conf utils.Config, db *gorm.DB, attachment kitsu.Attachment) bool {
	log.Info("[main.go][parseSingleAttachment] Started backing up '" + attachment.Name + "'")

	// Ignore attachents with missing IDs
	if attachment.ID == "" {
		return false
	}

	// Ignore attachments with extenstions from ignore list
	for _, elem := range conf.Backup.IgnoreExtension {
		if attachment.Extension == elem {
			log.Info("[main.go][parseSingleAttachment] Skipping ignored extension: " + elem + "\n")
			return false
		}
	}

	// Parse DB and ignore DONE unchanged attachments
	result := model.FindAttachment(db, attachment.ID)
	if len(result.AttachmentID) > 0 {
		if result.AttachmentStatus == "done" && result.AttachmentUpdatedAt == attachment.UpdatedAt {
			return false
		}
	}

	// Prepare local path
	localPath := conf.Backup.LocalStorage + attachment.ID

	// Prepare attachment name
	attachmentName := ""
	if attachment.Name != "" {
		attachmentName = utils.SanitizeString(attachment.Name)
	} else {
		return false
	}

	s3Path := ""

	if attachment.Comment.ObjectID != "" {

		task := kitsu.GetTask(attachment.Comment.ObjectID)

		/*
			## Kitsu help sheet
			Entity - an actual task e.g. shot01, or prop_chair
			Entity Type - category where Entity belongs e.g.: sho01 is a Shot
			Task Type - Assets/Shots's categories (column name in Kitsu UI)
			Task - actual task's sub-task that fits into the column
		*/

		// Get entity name (Top Task)
		log.Info("[main.go][parseSingleAttachment] ** Entity: **")
		entity := kitsu.GetEntity(task.EntityID)
		log.Info(entity)
		entityName := ""
		if entity.Name != "" {
			entityName = utils.SanitizeString(entity.Name) + "/"
		} else {
			return false
		}

		// Get Sequence Name
		sequenceName := ""
		episodeName := ""
		if entity.ParentID != "" {
			sequence := kitsu.GetEntity(entity.ParentID)
			sequenceName = sequence.Name + "/"

			// Get Episode Name
			if sequence.ParentID != "" {
				episode := kitsu.GetEntity(sequence.ParentID)
				episodeName = episode.Name + "/"
			}

		}

		// Get entity type
		log.Info("[main.go][parseSingleAttachment] ** Entity Type: **")
		entityType := kitsu.GetEntityType(entity.EntityTypeID)
		log.Info(entityType)
		entityTypeName := ""
		if entityType.Name == "" {
			entityTypeName = "_Unsorted" + "/"
		} else {
			// Make more verbose divide between Shots and Assets type
			if utils.SanitizeString(entityType.Name) == "Shot" {
				entityTypeName = "shots/"
			} else {
				entityTypeName = "assets/" + utils.SanitizeString(entityType.Name) + "/"
			}
		}

		// Get task type (Sub Task)
		log.Info("[main.go][parseSingleAttachment] ** Task Type: **")
		taskType := kitsu.GetTaskType(task.TaskTypeID)
		log.Info(taskType)
		taskTypeName := ""
		if taskType.Name != "" {
			taskTypeName = utils.SanitizeString(taskType.Name) + "/"
		}

		// Get Project
		log.Info("[main.go][parseSingleAttachment] ** Project: **")
		project := kitsu.GetProject(task.ProjectID)
		log.Info(project)
		projectName := ""
		if project.Name != "" {
			projectName = utils.SanitizeString(project.Name) + "/"
		}
		//projectStatus := kitsu.GetProjectStatus(project.ProjectStatusID)

		s3Path = conf.Backup.S3.RootFolderName + "/" + projectName + episodeName + entityTypeName + sequenceName + entityName + taskTypeName + attachmentName
	} else {
		s3Path = conf.Backup.S3.RootFolderName + "/" + "LOST.FILES" + "/" + attachment.ID + "/" + attachmentName
	}

	// Alter file path to add timestamp postfix
	datetime := strings.ReplaceAll(attachment.CreatedAt, ":", "-")
	lastInd := strings.LastIndex(s3Path, ".")
	if lastInd > 0 {
		//fmt.Println(filename[:lastInd])   // o/p: a_ab_daqe_sd
		//fmt.Println(filename[lastInd+1:]) // o/p: ew
		s3Path = s3Path[:lastInd] + "_" + datetime + "." + s3Path[lastInd+1:]
	} else {
		s3Path = s3Path + "_" + datetime
	}

	log.Info("[main.go][parseSingleAttachment] Formed path is: " + s3Path)

	if len(result.AttachmentID) > 0 {
		// check if status is different or last comment date don't match
		if result.AttachmentStatus != "done" || result.AttachmentUpdatedAt != attachment.UpdatedAt {
			// update
			model.UpdateAttachment(db, attachment.ID, attachment.UpdatedAt, "new")
			kitsu.DownloadAttachment(localPath, attachment.ID, attachmentName, conf)

			// Read file from local dir
			content, err := ioutil.ReadFile(localPath + "/" + attachmentName)
			if err != nil {
				panic(err)
			}

			// Upload file to S3 storage
			s3.UploadFile(s3Path, string(content), conf)
			model.UpdateAttachment(db, attachment.ID, attachment.UpdatedAt, "done")
		} else {
			log.Info("[main.go][parseSingleAttachment] Skipping existing attachment: " + attachmentName)
			return false
		}

	} else {
		// create
		// Download file from Kitsu
		model.CreateAttachment(db, attachment.ID, attachment.UpdatedAt, "new")
		kitsu.DownloadAttachment(localPath, attachment.ID, attachmentName, conf)

		// Read file from local dir
		content, err := ioutil.ReadFile(localPath + "/" + attachmentName)
		if err != nil {
			panic(err)
		}
		// Upload file to S3 storage
		s3.UploadFile(s3Path, string(content), conf)
		model.UpdateAttachment(db, attachment.ID, attachment.UpdatedAt, "done")
	}

	// Cleaning
	os.RemoveAll(localPath)
	log.Info("[main.go][parseSingleAttachment] Finished with '" + localPath + "'\n")
	return true
}

func setupLogger() {
	lumberjackLogger := &lumberjack.Logger{
		// Log file abbsolute path, os agnostic
		Filename:   filepath.ToSlash("./logs/log.txt"),
		MaxSize:    5, // MB
		MaxBackups: 10,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	// Fork writing into two outputs
	multiWriter := io.MultiWriter(os.Stderr, lumberjackLogger)

	logFormatter := new(log.TextFormatter)
	logFormatter.TimestampFormat = time.RFC1123Z // or RFC3339
	logFormatter.FullTimestamp = true

	log.SetFormatter(logFormatter)
	log.SetLevel(log.InfoLevel)
	log.SetOutput(multiWriter)
}
