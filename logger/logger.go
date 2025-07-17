package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger encapsule une instance logrus configurée
type Logger struct {
	*logrus.Logger
	config *Config
}

var (
	// Instance globale du logger
	defaultLogger *Logger
)

// New crée un nouveau logger avec la configuration fournie
func New(config *Config) *Logger {
	if config == nil {
		config = DefaultConfig()
	}

	// Valider la configuration
	config.Validate()

	logger := logrus.New()

	// Configuration du niveau
	logger.SetLevel(config.ParseLevel())

	// Configuration du format
	if config.IsJSON() {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		formatter := &logrus.TextFormatter{
			FullTimestamp: true,
			DisableColors: config.NoColor,
		}
		logger.SetFormatter(formatter)
	}

	// Par défaut, on écrit sur stdout
	logger.SetOutput(os.Stdout)

	return &Logger{
		Logger: logger,
		config: config,
	}
}

// Init initialise le logger global avec la configuration fournie
func Init(config *Config) {
	defaultLogger = New(config)
}

// GetDefault retourne l'instance globale du logger
func GetDefault() *Logger {
	if defaultLogger == nil {
		defaultLogger = New(nil) // Configuration par défaut
	}
	return defaultLogger
}

// SetOutput change la sortie du logger (utile pour les tests)
func (l *Logger) SetOutput(output io.Writer) {
	l.Logger.SetOutput(output)
}

// WithField crée un logger avec un champ supplémentaire
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields crée un logger avec plusieurs champs supplémentaires
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// WithComponent crée un logger avec le champ "component"
func (l *Logger) WithComponent(component string) *logrus.Entry {
	return l.WithField("component", component)
}

// WithRequest crée un logger avec des informations de requête HTTP
func (l *Logger) WithRequest(method, path, remoteAddr string) *logrus.Entry {
	return l.WithFields(logrus.Fields{
		"method":      method,
		"path":        path,
		"remote_addr": remoteAddr,
	})
}

// Fonctions de convenance pour le logger global
func Debug(args ...interface{}) {
	GetDefault().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	GetDefault().Debugf(format, args...)
}

func Info(args ...interface{}) {
	GetDefault().Info(args...)
}

func Infof(format string, args ...interface{}) {
	GetDefault().Infof(format, args...)
}

func Warn(args ...interface{}) {
	GetDefault().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	GetDefault().Warnf(format, args...)
}

func Error(args ...interface{}) {
	GetDefault().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	GetDefault().Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	GetDefault().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	GetDefault().Fatalf(format, args...)
}

func WithComponent(component string) *logrus.Entry {
	return GetDefault().WithComponent(component)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return GetDefault().WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetDefault().WithFields(fields)
}

func WithRequest(method, path, remoteAddr string) *logrus.Entry {
	return GetDefault().WithRequest(method, path, remoteAddr)
}
