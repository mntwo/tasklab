package zaplog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option interface {
	apply(cfg *config)
}

type ExtraKey string

type option func(cfg *config)

func (fn option) apply(cfg *config) {
	fn(cfg)
}

type CoreConfig struct {
	Enc zapcore.Encoder
	Ws  zapcore.WriteSyncer
	Lvl zapcore.LevelEnabler
}

type config struct {
	extraKeys   []ExtraKey
	coreConfigs []CoreConfig
	zapOpts     []zap.Option
}

// defaultCoreConfig default zapcore config: json encoder, atomic level, stdout write syncer
func defaultCoreConfig() *CoreConfig {
	// default log encoder
	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	// default log level
	lvl := zap.NewAtomicLevelAt(zap.InfoLevel)
	// default write syncer stdout
	ws := zapcore.AddSync(os.Stdout)

	return &CoreConfig{
		Enc: enc,
		Ws:  ws,
		Lvl: lvl,
	}
}

// defaultConfig default config
func defaultConfig() *config {
	return &config{
		coreConfigs: []CoreConfig{*defaultCoreConfig()},
		zapOpts:     []zap.Option{},
	}
}

// WithCoreEnc zapcore encoder
func WithCoreEnc(enc zapcore.Encoder) Option {
	return option(func(cfg *config) {
		cfg.coreConfigs[0].Enc = enc
	})
}

// WithCoreWs zapcore write syncer
func WithCoreWs(ws zapcore.WriteSyncer) Option {
	return option(func(cfg *config) {
		cfg.coreConfigs[0].Ws = ws
	})
}

// WithCoreLevel zapcore log level
func WithCoreLevel(lvl zap.AtomicLevel) Option {
	return option(func(cfg *config) {
		cfg.coreConfigs[0].Lvl = lvl
	})
}

// WithCores zapcore
func WithCores(coreConfigs ...CoreConfig) Option {
	return option(func(cfg *config) {
		cfg.coreConfigs = coreConfigs
	})
}

// WithZapOptions add origin zap option
func WithZapOptions(opts ...zap.Option) Option {
	return option(func(cfg *config) {
		cfg.zapOpts = append(cfg.zapOpts, opts...)
	})
}

// WithExtraKeys allow you log extra values from context
func WithExtraKeys(keys []ExtraKey) Option {
	return option(func(cfg *config) {
		for _, k := range keys {
			if !inArray(k, cfg.extraKeys) {
				cfg.extraKeys = append(cfg.extraKeys, k)
			}
		}
	})
}

// inArray check if a string in a slice
func inArray(key ExtraKey, arr []ExtraKey) bool {
	for _, k := range arr {
		if k == key {
			return true
		}
	}
	return false
}
