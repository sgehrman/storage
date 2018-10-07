package main

import (
	"github.com/cockroachdb/cockroach/pkg/picolo"
	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
	"sync"
)

var crdbInstWaitGroup sync.WaitGroup // keeps track of running crdb instances

func main() {
	log.Info("A walrus appears")

	log.AddHook(filename.NewHook())

	// create data dir
	picolo.CreateDataDir()

	// self updater auto updates the binary when a new version is available
	go picolo.ScheduleSelfUpdater()

	// initialize discovery service
	picolo.InitAppWithServiceAccount()

	// init picoloNode
	picoloNode := picolo.InitNode()

	// register picoloNode with discovery service
	picolo.Register(picoloNode)

	picolo.ThrowFlare(picoloNode)

	// spawn a crdb instance
	picolo.SpawnCrdbInst(picoloNode, &crdbInstWaitGroup)

	//init crdb cluster
	picolo.InitCrdbCluster(picoloNode, &crdbInstWaitGroup)

	log.Info("printed")

	crdbInstWaitGroup.Wait()
}