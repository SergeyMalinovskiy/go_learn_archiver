package cmd

import (
	"awesome-archiver/lib/vlc"
	"errors"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var vlcCmd = &cobra.Command{
	Use:   "vlc",
	Short: "Pack file using variable-length code",
	Run:   pack,
}

const packedExtension = "vlc"

var ErrEmptyPath = errors.New("path to file is not specified")

func pack(_ *cobra.Command, args []string) {
	if len(args) == 0 || args[0] == "" {
		handleErr(ErrEmptyPath)
	}

	filePath := args[0]

	data, err := os.ReadFile(filePath)
	if err != nil {
		handleErr(err)
	}

	packed := vlc.Encode(string(data))

	err = os.WriteFile(packedFileName(filePath), []byte(packed), 0644)
	if err != nil {
		handleErr(err)
	}
}

func packedFileName(path string) string {
	fileName := filepath.Base(path)

	return strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "." + packedExtension
}

func init() {
	packCmd.AddCommand(vlcCmd)
}
