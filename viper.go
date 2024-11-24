package singleton

import (
	"github.com/spf13/viper"
)

type ViperOptions struct {
	Path     string // path to look for the config file in, absolute path start with / and relative path start with .
	FileName string // name of config file (without extension)
	FileType string // REQUIRED if the config file does not have the extension in the name
	EnvName  string // Env to read from the config file
	Target   any
}

func WithViper(options ViperOptions) Option {
	return func(s *Singleton) (err error) {
		v := viper.New()
		v.SetConfigType(options.FileType)
		v.SetConfigName(options.FileName)
		v.AddConfigPath(options.Path)
		err = v.ReadInConfig()
		if err != nil {
			// global.LOGGER.Error(fmt.Sprintf("fatal error config file: %s", err))
			return
		}
		// v.WatchConfig() // watch config change, hot reload
		// v.OnConfigChange(func(e fsnotify.Event) {
		// 	log.Println("Config file changed:", e.Name)
		// 	// singleton.Logger.Info(fmt.Sprintf("config file changed: %s", e.Name))
		// 	// handleResizerWorkersChange(context.TODO(), v)
		// })
		if err = v.UnmarshalKey(options.EnvName, &options.Target); err != nil {
			return
		}
		s.Viper = v
		return
	}
}
