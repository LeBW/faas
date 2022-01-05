package handlers

import (
	"github.com/openfaas/faas/gateway/scaling"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func NewFunctionScheduler(config scaling.ScalingConfig, defaultNamespace string, functionCacher scaling.FunctionCacher) FunctionScheduler {
	return FunctionScheduler{
		Cache:            functionCacher,
		Config:           config,
		DefaultNamespace: defaultNamespace,
		cron:             cron.New(cron.WithSeconds()),
	}
}

type FunctionScheduler struct {
	Cache            scaling.FunctionCacher
	Config           scaling.ScalingConfig
	DefaultNamespace string
	cron             *cron.Cron
}

func (scheduler *FunctionScheduler) AddPredictions(predictions []Prediction) {
	len := len(predictions)
	log.Printf("[AddPredictions] number of predictions: %d\n", len)
	for i := range predictions {
		scheduler.schedule(predictions[i])
	}
}

func (scheduler *FunctionScheduler) schedule(prediction Prediction) {
	// easy schedule
	scheduleTimestamp := time.Now().UnixNano()/1e6 + int64(prediction.PredictTime)
	scheduleCron := time.Unix(0, scheduleTimestamp).Format("05 04 15 ? 01 ?")
	log.Printf("%#v", prediction)
	log.Printf("[schedule] scheduleTimestamp: %d, cron expression: %s\n", scheduleTimestamp, scheduleCron)
	scheduler.cron.AddFunc(scheduleCron, func() {
		log.Printf("Cron job start. Schedule function %s", prediction.FunctionName)
		err := scheduler.Config.ServiceQuery.SetReplicas(prediction.FunctionName, scheduler.DefaultNamespace, uint64(prediction.Probability))
		if err != nil {
			log.Printf("Schedule function %s Failed, %s\n", prediction.FunctionName, err)
		}
	})
}
