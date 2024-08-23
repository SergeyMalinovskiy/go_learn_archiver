package cmd

import (
	"awesome-archiver/lib/compression"
	"awesome-archiver/lib/compression/vlc"
	"awesome-archiver/lib/compression/vlc/table/shannon_fano"
	"errors"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "Pack file",
	Run:   pack,
}

const packedExtension = "vlc"

var ErrEmptyPath = errors.New("path to file is not specified")

func pack(cmd *cobra.Command, args []string) {
	var encoder compression.Encoder

	if len(args) == 0 || args[0] == "" {
		handleErr(ErrEmptyPath)
	}

	method := cmd.Flag("method").Value.String()

	switch method {
	case "vlc":
		encoder = vlc.New(shannon_fano.NewGenerator())
	default:
		cmd.PrintErrln("unknown method")
	}

	filePath := args[0]

	data, err := os.ReadFile(filePath)
	if err != nil {
		handleErr(err)
	}

	packed := encoder.Encode(string(data))

	err = os.WriteFile(packedFileName(filePath), packed, 0644)
	if err != nil {
		handleErr(err)
	}
}

func packedFileName(path string) string {
	fileName := filepath.Base(path)

	return strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "." + packedExtension
}

func init() {
	rootCmd.AddCommand(packCmd)

	packCmd.Flags().StringP("method", "m", "", "compression method: vlc")
	if err := packCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
}
