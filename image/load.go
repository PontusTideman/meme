package image

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/PontusTideman/meme/cli"
	"github.com/PontusTideman/meme/data"
	"github.com/PontusTideman/meme/image/stream"
	"github.com/PontusTideman/meme/output"
)

var (
	imageMap = make(map[string]string)
)

// Initialise the package.
func init() {
	images, err := data.Files.ReadDir(data.ImagePath)
	output.OnError(err, "Could not read embedded images")

	for _, image := range images {
		id := strings.TrimSuffix(filepath.Base(image.Name()), data.ImageExtension)
		imageMap[id] = path.Join(data.ImagePath, image.Name())
	}
}

// Load an image from the passed string or stdin.
// The string will be a embedded asset id, an image URL or a local file.
func Load(opt cli.Options) stream.Stream {
	var s io.Reader

	if isURL(opt.Image) {
		s = downloadURL(opt.Image)

	} else if isStdin(opt.Image) {
		s = readStdin()

	} else if isAsset(opt.Image) {
		s = loadAsset(opt.Image)

	} else if isLocalFile(opt.Image) {
		s = readFile(opt.Image)

	} else {
		output.Error("Image not recognised")
	}

	return stream.NewStream(s)
}

// Return true if the passed string is an embedded asset id, false if not.
func isAsset(id string) bool {
	_, ok := imageMap[id]
	return ok
}

// Load and return an embedded asset (image) by id.
// The id is assumed to exist.
func loadAsset(id string) io.Reader {
	image, _ := imageMap[id]

	st, err := data.Files.ReadFile(image)
	output.OnError(err, "Could not read embedded image")

	return bytes.NewReader(st)
}

// Return true if the passed string is an image URL, false if not.
func isURL(url string) bool {
	return strings.HasPrefix(url, "http")
}

// Download the image located at the passed image URL, decode and return it.
func downloadURL(url string) io.Reader {
	res, err := http.Get(url)
	output.OnError(err, "Request error")
	defer res.Body.Close()

	if res.StatusCode != 200 {
		output.Error("Could not access URL")
	}

	st, err := ioutil.ReadAll(res.Body)
	output.OnError(err, "Could not read response body")

	return bytes.NewReader(st)
}

// Return true if the passed string is a file that exists on the local
// filesystem, false if not.
func isLocalFile(path string) bool {
	path, err := homedir.Expand(path)
	output.OnError(err, "Could not expand path")

	_, err = os.Stat(path)
	return err == nil
}

// Read and return a file on the local filesystem.
// The file is assumed to exist.
func readFile(path string) io.Reader {
	path, err := homedir.Expand(path)
	output.OnError(err, "Could not expand path")

	st, err := ioutil.ReadFile(path)
	output.OnError(err, "Could not read local file")

	return bytes.NewReader(st)
}

// return true if the passed string is '-' meaning we should read the image
// from stdin.
func isStdin(path string) bool {
	return path == "-"
}

// Read the image from stdin.
func readStdin() io.Reader {
	st, err := ioutil.ReadAll(os.Stdin)
	output.OnError(err, "Could not read stdin")

	return bytes.NewReader(st)
}

// Decal returns the named decal as a stream.
func Decal(name string) stream.Stream {
	st, err := data.Files.ReadFile(name)
	output.OnError(err, "Could not read embedded decal")

	return stream.NewStream(bytes.NewReader(st))
}
