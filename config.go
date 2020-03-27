package logger

type Config struct {
	File   FileConfig
	Format FormatConfig
}
type FileConfig struct {
	PrintToFile  bool
	Name         string
	Path         string
	RotationType string
	RotationSize string
	RotationTime string
}
type FormatConfig struct {
	Level      string
	TimeLayout string
}
