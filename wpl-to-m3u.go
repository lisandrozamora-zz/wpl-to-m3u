package main

import "encoding/xml"
import "fmt"
import "flag"
import "io/ioutil"
import "os"
import "path/filepath"
import "strings"

const DIRECTORY string = "directory"
const FILE string = "file"
const NONE string = "none"
const WPL_EXTENSION string = ".wpl"
const NEWLINE string = "\n"
const M3U_FILE_HEADER string = "#EXTM3U" + NEWLINE

type WPLFile struct {
	path string
}

type XMLMeta struct {
	XMLName xml.Name `xml:"meta"`
	Name    string   `xml:"name,attr"`
	Content string   `xml:"content,attr"`
}

type XMLHead struct {
	XMLName xml.Name  `xml:"head"`
	Meta    []XMLMeta `xml:"meta"`
	Author  string    `xml:"author"`
	Title   string    `xml:"title"`
}

type XMLMedia struct {
	XMLName xml.Name `xml:"media"`
	Source  string   `xml:"src,attr"`
	TID     string   `xml:"tid,attr"`
	CID     string   `xml:"cid,attr"`
}

type XMLSequence struct {
	XMLName xml.Name   `xml:"seq"`
	Media   []XMLMedia `xml:"media"`
}

type XMLBody struct {
	XMLName  xml.Name    `xml:"body"`
	Sequence XMLSequence `xml:"seq"`
}

type XMLPlaylist struct {
	XMLName xml.Name `xml:"smil"`
	Head    XMLHead  `xml:"head"`
	Body    XMLBody  `xml:"body"`
}

func getFileOrDirectoryToConvertCommandLineArgument() string {
	flag.Parse()
	return flag.Arg(0)
}

func getFileType(file string) string {
	var fileType string
	fi, _ := os.Stat(file)
	switch {
	case fi.IsDir():
		// it's a directory
		fileType = DIRECTORY
	default:
		// it's not a directory
		fileType = FILE
	}

	return fileType
}

func printUsageMessageAndExit() {
	fmt.Println("Usage: wpl2m3u <path to wpl file or directory containing wpl files>")
	os.Exit(-1)
}

func verifyExistenceAndExitIfNotExists(file string) {
	if len(file) == 0 {
		printUsageMessageAndExit()
	} else {
		_, err := os.Stat(file)
		switch {
		case err != nil:
			// does not exist...exit
			fmt.Println("No such file or directory: ", file)
			os.Exit(-1)
		}
	}
}

func getPlaylistsToConvert(wplFileOrDirectoryToConvert string) []WPLFile {
	var resultList []WPLFile
	playlistArgType := getFileType(wplFileOrDirectoryToConvert)
	if playlistArgType == FILE {
		resultList = append(resultList, WPLFile{path: wplFileOrDirectoryToConvert})
	} else {
		dirEntries, _ := ioutil.ReadDir(wplFileOrDirectoryToConvert)
		for _, f := range dirEntries {
			if strings.HasSuffix(f.Name(), WPL_EXTENSION) {
				resultList = append(resultList, WPLFile{path: wplFileOrDirectoryToConvert + "/" + f.Name()})
			}
		}
	}

	return resultList
}

func getMediaSequenceFromPlaylistFile(playlistPath string) XMLSequence {
	xmlFile, _ := os.Open(playlistPath)
	defer xmlFile.Close()
	XMLData, _ := ioutil.ReadAll(xmlFile)
	var pl XMLPlaylist
	xml.Unmarshal(XMLData, &pl)
	return pl.Body.Sequence
}

func getTrackInformationLine(track string) string {
	// must replace Windows path separators with Unix ones so the filepath.Split() works when run on Linux
	modifiedTrack := strings.Replace(track, "\\", "/", 20)
	_, filename := filepath.Split(modifiedTrack)
	return "#EXTINF:0," + filename + "\n"
}

func writeOutM3uFile(currentWplFile string, media []XMLMedia) {
	newM3uFile := strings.Replace(currentWplFile, ".wpl", ".m3u", 1)
	fmt.Println("File to write out:", newM3uFile)
	f, _ := os.Create(newM3uFile)
	_, _ = f.WriteString(M3U_FILE_HEADER)
	for _, element := range media {
		trackInformationLine := getTrackInformationLine(element.Source)
		_, _ = f.WriteString(trackInformationLine)
		_, _ = f.WriteString(element.Source)
		_, _ = f.WriteString(NEWLINE)
		_, _ = f.WriteString(NEWLINE)
	}
}

func convertPlaylist(playlist WPLFile) {
	fmt.Println("Currently converting : ", playlist.path)
	mediaSequence := getMediaSequenceFromPlaylistFile(playlist.path)
	writeOutM3uFile(playlist.path, mediaSequence.Media)
}

func main() {

	fileOrDirectoryToConvert := getFileOrDirectoryToConvertCommandLineArgument()
	verifyExistenceAndExitIfNotExists(fileOrDirectoryToConvert)
	listOfPlaylists := getPlaylistsToConvert(fileOrDirectoryToConvert)
	for _, currentPlaylist := range listOfPlaylists {
		convertPlaylist(currentPlaylist)
	}
	fmt.Println("Done.")
}
