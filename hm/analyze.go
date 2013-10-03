package hm

import (
	"github.com/cloudfoundry/hm9000/analyzer"
	"github.com/cloudfoundry/hm9000/config"
	"github.com/cloudfoundry/hm9000/helpers/logger"
	"github.com/cloudfoundry/hm9000/helpers/outbox"
	"github.com/cloudfoundry/hm9000/store"
	"github.com/cloudfoundry/hm9000/storeadapter"

	"os"
)

func Analyze(l logger.Logger, conf config.Config, poll bool) {
	etcdStoreAdapter := connectToETCDStoreAdapter(l, conf)

	if poll {
		l.Info("Starting Analyze Daemon...")
		err := Daemonize(func() error {
			return analyze(l, conf, etcdStoreAdapter)
		}, conf.AnalyzerPollingInterval(), conf.AnalyzerTimeout(), l)
		if err != nil {
			l.Error("Analyze Daemon Errored", err)
		}
		l.Info("Analyze Daemon is Down")
	} else {
		err := analyze(l, conf, etcdStoreAdapter)
		if err != nil {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}
}

func analyze(l logger.Logger, conf config.Config, etcdStoreAdapter storeadapter.StoreAdapter) error {
	store := store.NewStore(conf, etcdStoreAdapter)
	outbox := outbox.New(store, l)

	l.Info("Analyzing...")

	analyzer := analyzer.New(store, outbox, buildTimeProvider(l), l, conf)
	err := analyzer.Analyze()

	if err != nil {
		l.Error("Analyzer failed with error", err)
		return err
	} else {
		l.Info("Analyzer completed succesfully")
		return nil
	}
}
