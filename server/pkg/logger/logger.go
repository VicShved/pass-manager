package logger

import "go.uber.org/zap"

var Log *zap.Logger = zap.NewNop()

func InitLogger(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	log, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = log
	defer Log.Sync()
	return nil
}
