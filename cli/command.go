package cli

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/nmaupu/http2back/provider"
	"github.com/nmaupu/http2back/server"
	"os"
)

const (
	AppName    = "http2back"
	AppDesc    = "Push to backends over HTTP"
	AppVersion = "v0.2"
)

var (
	addr *string
	port *int
)

func Process() {
	app := cli.App(AppName, AppDesc)
	app.Version("v version", fmt.Sprintf("%s %s", AppName, AppVersion))

	addr = app.StringOpt("b bind", "127.0.0.1", "Bind address")
	port = app.IntOpt("p port", 8080, "Port to listen connections from")

	app.Command("filesystem fs", "Use filesystem provider", providerFilesystem)
	app.Command("ftp", "Use FTP provider", providerFtp)
	app.Command("dropbox d", "Use Dropbox provider", providerDropbox)
	app.Command("aws-s3 s3", "Use AWS S3 provider", providerAwsS3)

	app.Run(os.Args)
}

func providerFilesystem(cmd *cli.Cmd) {
	dest := cmd.StringOpt("d dest", "/tmp", "Destination directory where to drop files into")

	cmd.Action = func() {
		server.Start(port, addr, func() provider.Provider {
			return provider.Filesystem{DestDir: *dest}
		})
	}
}

func providerFtp(cmd *cli.Cmd) {
	dest := cmd.StringOpt("d dest", "/", "Destination directory where to drop files into")
	ftpAddr := cmd.StringOpt("a addr", "127.0.0.1:21", "FTP host or address with port if need be")

	username := cmd.String(cli.StringOpt{
		Name:   "u username",
		Value:  "anonymous",
		Desc:   "Username to use for FTP login command",
		EnvVar: "FTP_USERNAME",
	})

	password := cmd.String(cli.StringOpt{
		Name:   "p password",
		Value:  "anonymous",
		Desc:   "Password to use for FTP login command",
		EnvVar: "FTP_PASSWORD",
	})

	cmd.Action = func() {
		server.Start(port, addr, func() provider.Provider {
			return provider.Ftp{
				Addr:     *ftpAddr,
				Username: *username,
				Password: *password,
				DestDir:  *dest,
			}
		})
	}
}

func providerDropbox(cmd *cli.Cmd) {
	dest := cmd.StringOpt("d dest", "", "Destination directory where to drop files into")
	accessToken := cmd.String(cli.StringOpt{
		Name:   "t token access-token",
		Value:  "my-token",
		Desc:   "Dropbox access token for API",
		EnvVar: "DROPBOX_API_TOKEN",
	})

	cmd.Action = func() {
		server.Start(port, addr, func() provider.Provider {
			return provider.Dropbox{
				DestDir:     *dest,
				AccessToken: *accessToken,
			}
		})
	}
}

func providerAwsS3(cmd *cli.Cmd) {
	bucket := cmd.StringOpt("b bucket", "my-bucket", "Destination bucket")
	dest := cmd.StringOpt("d dest", "", "Destination directory where to drop files into")

	region := cmd.String(cli.StringOpt{
		Name:   "r region",
		Value:  "eu-west-1",
		Desc:   "Region corresponding to the provided bucket",
		EnvVar: "AWS_REGION",
	})
	key := cmd.String(cli.StringOpt{
		Name:   "k key aws-access-key",
		Value:  "my-access-key",
		Desc:   "AWS api access key",
		EnvVar: "AWS_ACCESS_KEY_ID",
	})
	secret := cmd.String(cli.StringOpt{
		Name:   "s secret secret-access-key",
		Value:  "my-secret",
		Desc:   "AWS api secret key",
		EnvVar: "AWS_SECRET_ACCESS_KEY",
	})
	token := cmd.String(cli.StringOpt{
		Name:   "t token session-token",
		Value:  "",
		Desc:   "AWS token (only for temporary authentication - optional)",
		EnvVar: "AWS_SESSION_TOKEN",
	})

	cmd.Action = func() {
		server.Start(port, addr, func() provider.Provider {
			return provider.AwsS3{
				Bucket:             *bucket,
				DestDir:            *dest,
				Region:             *region,
				AwsAccessKeyId:     *key,
				AwsSecretAccessKey: *secret,
				Token:              *token,
			}
		})
	}
}
