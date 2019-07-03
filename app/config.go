package app


type config struct {
	FTP ftp `json:"ftp"`
}

type ftp struct {
	Host string `json:"host"`
	Password string `json:"password"`
	Login string `json:"login"`
	StartPath string `json:"startPath"`
	DestinationPath string `json:"destinationPath"`
}


func Config() config {
	return config{FTP:ftp{Host: "ftp.radius.ru:8021",
		Password: "GS",
		Login: "GS",
		StartPath: "PHOTO_HH_MSK",
		DestinationPath: "result",
	}}
}

