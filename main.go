package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// For testing: "Marvel.Agents.of.S.H.I.E.L.D.S04E15.XD.Garbage.Junk-StupidStuff"

// Spec is used to hold properties to rename and/or move files/folders
type Spec struct {
	OldPath    string
	NewPath    string
	MainFolder string
}

var flagPath string
var flagDebug bool
var flagSplits bool

var specs []Spec
var showFoldersPaths map[string]string

func init() {
	flag.StringVar(&flagPath, "path", "", "path to traverse")
	flag.BoolVar(&flagDebug, "debug", false, "debug")
	flag.BoolVar(&flagSplits, "splits", false, "debug splits")
	flag.Parse()
	flagDebug = flagDebug || flagSplits
	showFoldersPaths = make(map[string]string)
}

func main() {
	checkFlags()
	validatePath()
	organize()
}

func checkFlags() {
	if flagPath == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func validatePath() {
	logDebug(flagDebug, "Path: %s", flagPath)
	flagPath, err := filepath.Abs(flagPath)
	if err != nil {
		logFatal("Unable to get Abs Path: %s", err.Error())
	}
	logDebug(flagDebug, "Abs Path: %s", flagPath)
	if !strings.HasPrefix(flagPath, "/Users/") {
		logRed("Path must start with /Users/")
		os.Exit(1)
	}
	_, err = os.Stat(flagPath)
	if err != nil {
		if os.IsNotExist(err) {
			logRed("%s is not a valid path - does not exist", flagPath)
			os.Exit(1)
		}
	}
	if strings.HasPrefix(filepath.Base(flagPath), ".") {
		logRed("Path cannot point to a hidden folder")
	}
}

func organize() {
	filepath.Walk(flagPath, visit)
	// Since we are renaming paths, and filepath.Walk walks the file tree lexically, we should rename starting from the end of the array

	for i := len(specs) - 1; i >= 0; i-- {
		// Rename the shows
		spec := specs[i]
		logCyan("Old path: %s", spec.OldPath)
		logMagenta("New path: %s", spec.NewPath)
		renameErr := os.Rename(spec.OldPath, spec.NewPath)
		if renameErr != nil {
			logRed("Unable to rename file: %s", renameErr.Error())
		}
		if mainFolderPath, exists := showFoldersPaths[spec.MainFolder]; exists {
			// TODO: Move shit
			logGreen("Main folder: %s", mainFolderPath)
		}
	}

	// for k, v := range showFoldersPaths {
	// 	logDefault("%s:%s", k, v)
	// }
}

func visit(path string, f os.FileInfo, err error) (e error) {
	if err != nil {
		log.Fatal(e.Error())
	}
	if f.IsDir() && !strings.HasPrefix(filepath.Base(path), ".") {
		checkPath(path)
	}
	return
}

func checkPath(path string) {

	dirName := strings.TrimSpace(filepath.Base(path))

	showFoldersRegex := regexp.MustCompile("(.*) Season [0-9]{1,2}")
	showNameRegex := regexp.MustCompile("(.*)(((S|s)[0-9]{1,2}(E|e)[0-9]{1,2})|(Season(.+)[0-9]{1,2}(.+)Episode(.+)[0-9]{1,2}))")
	showFoldersIndex := showFoldersRegex.FindStringIndex(dirName)
	showNameIndex := showNameRegex.FindStringIndex(dirName)

	if showFoldersIndex != nil && showFoldersIndex[1]-showFoldersIndex[0] == len(dirName) {
		logDebug(flagDebug, "%s is a show main folder match", dirName)
		showFoldersPaths[strings.ToLower(dirName)] = path
		return
	}
	if showNameIndex == nil {
		logDebug(flagDebug, "%s is not a match", dirName)
		return
	}
	logDebug(flagDebug, "%s is a match", dirName)

	dirName = dirName[showNameIndex[0]:showNameIndex[1]]
	// Split the directory name by spaces and periods
	dirNameSplit := regexp.MustCompile("(\\.|\\s)").Split(dirName, -1)

	// If no split, then directory doesn't need to be renamed - else it might or might not need to be
	if len(dirNameSplit) == 0 {
		return
	}

	newDirName, mainFolder := getNewDirNameAndMainFolder(dirNameSplit)
	if newDirName == "" {
		return
	}

	// Figure out the old path and new path
	oldPath, _ := filepath.Abs(path)
	newPath := filepath.Dir(oldPath) + "/" + newDirName

	// If this directory needs to be renamed, then add it to the Specs
	if oldPath != newPath {
		specs = append(specs, Spec{OldPath: oldPath, NewPath: newPath, MainFolder: strings.ToLower(mainFolder)})
	}
}

func getNewDirNameAndMainFolder(dirNameSplit []string) (newDirName string, mainFolder string) {

	// Acronym Variable used for building any acroynms such as S.H.I.E.L.D.
	var acronym string
	seasonEpRegex := regexp.MustCompile("(S|s)[0-9]{1,2}(E|e)[0-9]{1,2}")

	for i := 0; i < len(dirNameSplit); i++ {

		if dirNameSplit[i] == "" || dirNameSplit[i] == " " {
			// Ignore nil or empty splits
			continue
		} else if len(dirNameSplit[i]) == 1 {
			// If the length is 1, assume its an acronym - if its not, then we fix it automatically
			acronym += dirNameSplit[i] + "."
			continue
		} else if acronym != "" {
			// If the acronym is not empty after all above cases, then append the acronym to the new directory name
			if len(acronym) == 2 && strings.HasSuffix(acronym, ".") {
				// Remove "." if the acronym is only 1 letter long - assumption is that 1 letter strings aren't acronyms
				acronym = acronym[0:1]
			}
			logDebug(flagSplits, "%s", acronym)
			newDirName += acronym + " "
			acronym = ""
		}

		// If the string matches the SXXEXX regex, then expand it
		if seasonEpRegex.MatchString(dirNameSplit[i]) {
			seasonNumber, episodeNumber := getSeasonAndEpNumber(dirNameSplit[i])
			if seasonNumber == "" || episodeNumber == "" {
				return
			}
			mainFolder = newDirName + "Season " + seasonNumber
			newSeasonEp := "Season " + seasonNumber + " Episode " + episodeNumber
			dirNameSplit[i] = newSeasonEp
		}

		logDebug(flagSplits, "%s", dirNameSplit[i])
		newDirName += dirNameSplit[i] + " "
	}

	// If the name ends with an acronym
	if acronym != "" {
		if len(acronym) == 2 && strings.HasSuffix(acronym, ".") {
			// Remove "." if the acronym is only 1 letter long - assumption is that 1 letter strings aren't acronyms
			acronym = acronym[0:1]
		}
		logDebug(flagSplits, "%s", acronym)
		newDirName += acronym + " "
	}

	// Remove the extra whitespace we added at the end
	newDirName = newDirName[0 : len(newDirName)-1]

	return newDirName, mainFolder
}

func getSeasonAndEpNumber(seasonEpisode string) (sNum string, epNum string) {
	seasonNumRegex := regexp.MustCompile("[0-9]{1,2}")
	numIndex := seasonNumRegex.FindAllStringSubmatchIndex(seasonEpisode, -1)
	if len(numIndex) != 2 {
		return
	}
	seasonNumber := seasonEpisode[numIndex[0][0]:numIndex[0][1]]
	episodeNumber := seasonEpisode[numIndex[1][0]:numIndex[1][1]]
	seasonNumber = stripZero(seasonNumber)
	episodeNumber = stripZero(episodeNumber)
	return seasonNumber, episodeNumber
}

func stripZero(str string) string {
	if strings.HasPrefix(str, "0") {
		if len(str) > 1 {
			str = str[1:2]
		}
	}
	return str
}
